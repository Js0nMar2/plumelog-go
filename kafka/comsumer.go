package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/olivere/elastic/v7"
	"plumelog/config"
	plumelogEs "plumelog/elastic"
	"plumelog/log"
	"plumelog/model"
	"plumelog/utils"
	"sync"
	"time"
)

var (
	ConsumerMap  map[string]*ConsumerGroup
	ServiceSlice []*model.HealthGuard
)

// init 初始化已知的服务的消费者
func init() {
	ConsumerMap = make(map[string]*ConsumerGroup)
	ServiceSlice = make([]*model.HealthGuard, 0)
	topics := []string{
		"app-index",
	}
	result, _ := plumelogEs.Client.Search("plume_log_services").Size(100).Do(context.Background())
	for _, hit := range result.Hits.Hits {
		marshalJSON, _ := hit.Source.MarshalJSON()
		m := make(map[string]string)
		json.Unmarshal(marshalJSON, &m)
		serviceName := m["serviceName"]
		// 初始化消费者
		if serviceName != "" {
			topics = append(topics, serviceName)
		}
	}
	cg := NewConsumerGroup(topics, []string{config.Conf.Kafka.Host}, config.Conf.GroupId)
	ConsumerMap[config.Conf.GroupId] = cg
	go cg.StartConsumerGroup(cg)
}

type ConsumerGroup struct {
	sarama.ConsumerGroup
	GroupId     string
	Topics      []string
	PlumeInfoCh chan model.PlumelogInfo
	WsChMap     map[string]chan model.PlumelogInfo
	HealthMap   map[string]chan []byte
}

func NewConsumerGroup(topics, addr []string, groupID string) *ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V3_2_3_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(addr, groupID, config)
	if err != nil {
		panic(err.Error())
	}
	log.Info("kafka version: %s init address is: %s topics is: %s", config.Version, addr, topics)
	plumeInfoCh := make(chan model.PlumelogInfo, 10000)
	wsCh := make(map[string]chan model.PlumelogInfo)
	healthMap := make(map[string]chan []byte)
	return &ConsumerGroup{
		consumerGroup,
		groupID,
		topics,
		plumeInfoCh,
		wsCh,
		healthMap,
	}
}

// 启动消费者
func (cg *ConsumerGroup) StartConsumerGroup(handler sarama.ConsumerGroupHandler) {
	ctx := context.Background()
	for {
		err := cg.ConsumerGroup.Consume(ctx, cg.Topics, handler)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(1 * time.Second)
		}
	}
}

func (*ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (*ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (cg *ConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	rwLock := new(sync.RWMutex)
	bulkService := plumelogEs.Client.Bulk()
	//读取kafka数据记录到分片,使用定时批量读取分片内数据
	messages := claim.Messages()
	for msg := range messages {
		if claim.Topic() == "app-index" {
			err := cg.putService(string(msg.Value))
			if err != nil {
				log.Error("set elastic index err:", err)
				return err
			}
			sess.MarkMessage(msg, "")
			return nil
		}
		m := make([]string, 0)
		err := json.Unmarshal(msg.Value, &m)
		if err != nil {
			log.Error(err.Error())
		}
		for _, message := range m {
			if message == "" {
				continue
			}
			mm := make(map[string]string)
			err := json.Unmarshal([]byte(message), &mm)
			if err != nil {
				log.Error(err.Error())
			}
			plumelogInfo := model.PlumelogInfo{}
			err = json.Unmarshal([]byte(mm["message"]), &plumelogInfo)
			if err != nil {
				log.Error(err.Error())
			}
			plumelogInfo.DateTime = utils.ParseTime(plumelogInfo.DateTime)
			rwLock.Lock()
			healthGuard, err := plumelogEs.PutServiceStatus(plumelogInfo)
			if err != nil {
				log.Error(err.Error())
			}
			if err == nil && healthGuard.Type != "" {
				for _, guardCh := range cg.HealthMap {
					marshal, err := json.Marshal(healthGuard)
					if err != nil {
						log.Error(err.Error())
					}
					guardCh <- marshal
				}
			}
			// 批量写入es
			go func() {
				timeStr := time.Now().Format("20060102")
				request := elastic.NewBulkIndexRequest().Index("plume_log_" + timeStr).Doc(plumelogInfo)
				bulkService.Add(request)
			}()
			// ws
			for key, ch := range cg.WsChMap {
				if model.ConnMap[key] == nil {
					delete(cg.WsChMap, key)
					continue
				}
				ch <- plumelogInfo
			}
			rwLock.Unlock()
		}
		if bulkService.NumberOfActions() > 0 {
			_, err = bulkService.Do(context.Background())
			if err != nil {
				log.Error(err.Error())
			}
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (cg *ConsumerGroup) putService(serviceName string) error {
	res, err := plumelogEs.PutService(serviceName)
	if res == "created" {
		cg.Topics = append(cg.Topics, serviceName)
		log.Info("put new service:", serviceName)
		cg.StartConsumerGroup(cg)
	}
	return err
}
