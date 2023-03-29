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
	"strings"
	"sync"
	"time"
)

var (
	ConsumerMap map[string]*ConsumerGroup
)

// init 初始化已知的服务的消费者
func init() {
	ConsumerMap = make(map[string]*ConsumerGroup)
	topics := make([]string, 0)
	for _, hit := range plumelogEs.Hits {
		marshalJSON, _ := hit.Source.MarshalJSON()
		m := make(map[string]string)
		json.Unmarshal(marshalJSON, &m)
		serviceName := m["serviceName"]
		// 初始化消费者
		if serviceName != "" {
			topics = append(topics, serviceName)
		}
	}
	cg := &ConsumerGroup{
		GroupId: config.Conf.GroupId,
		Topics:  topics,
	}
	cg.Init()
	ConsumerMap[config.Conf.GroupId] = cg
}

func (och *ConsumerGroup) Init() {
	och.ConsumerGroup = NewMConsumerGroup(&MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, och.Topics,
		[]string{config.Conf.Kafka.Host}, och.GroupId)
	och.PlumeInfoCh = make(chan model.PlumelogInfo, 10000)
	och.WsChMap = make(map[string]chan model.PlumelogInfo)
	go func() {
		och.Consume(context.Background(), och.Topics, och)
	}()
}

type ConsumerGroup struct {
	sarama.ConsumerGroup
	GroupId     string
	Topics      []string
	PlumeInfoCh chan model.PlumelogInfo
	WsChMap     map[string]chan model.PlumelogInfo
}

type MConsumerGroupConfig struct {
	KafkaVersion   sarama.KafkaVersion
	OffsetsInitial int64
	IsReturnErr    bool
}

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, topics, addr []string, groupID string) *ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = consumerConfig.KafkaVersion
	config.Consumer.Offsets.Initial = consumerConfig.OffsetsInitial
	config.Consumer.Return.Errors = consumerConfig.IsReturnErr
	consumerGroup, err := sarama.NewConsumerGroup(addr, groupID, config)
	if err != nil {
		panic(err.Error())
	}
	log.Info.Println("init address is ", addr, "topics is ", topics)
	plumeInfoCh := make(chan model.PlumelogInfo, 10000)
	wsCh := make(map[string]chan model.PlumelogInfo)
	return &ConsumerGroup{
		consumerGroup,
		groupID,
		topics,
		plumeInfoCh,
		wsCh,
	}
}
func (mc *ConsumerGroup) RegisterHandleAndConsumer(handler sarama.ConsumerGroupHandler) {
	ctx := context.Background()
	for {
		err := mc.ConsumerGroup.Consume(ctx, mc.Topics, handler)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

/*
*
读取kafka消息放入到切片
使用协程定时任务获取切片内的数据(批量)
将获取到的数据打包发送到管道msgDistributionCh(1000打成一个元素,最后一个不满足1000的也为一个元素)
*/
func (och *ConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	for {
		if sess == nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	if claim.Topic() == "app-index" {
		messages := claim.Messages()
		for msg := range messages {
			err := putService(string(msg.Value))
			if err != nil {
				log.Error.Println("set elastic index err:", err)
				return err
			}
		}
		return nil
	}
	rwLock := new(sync.RWMutex)
	consumerGroup := ConsumerMap[config.Conf.GroupId]
	go func() {
		bulkService := plumelogEs.Client.Bulk()
		for {
			select {
			case <-time.After(5 * time.Second):
				timeStr := time.Now().Format("20060102")
				for {
					plumelogInfo := <-consumerGroup.PlumeInfoCh
					request := elastic.NewBulkIndexRequest().Index("plume_log_" + timeStr).Doc(plumelogInfo)
					bulkService.Add(request)
					bulkService.Do(context.Background())
				}
			}
		}
	}()
	//读取kafka数据记录到分片,使用定时批量读取分片内数据
	messages := claim.Messages()
	for msg := range messages {
		m := make([]string, 1)
		json.Unmarshal(msg.Value, &m)
		for _, message := range m {
			if message == "" {
				continue
			}
			mm := make(map[string]string)
			err := json.Unmarshal([]byte(message), &mm)
			if err != nil {
				log.Error.Println(err.Error())
			}
			plumelogInfo := model.PlumelogInfo{}
			err = json.Unmarshal([]byte(mm["message"]), &plumelogInfo)
			if err != nil {
				log.Error.Println(err.Error())
			}
			rwLock.Lock()
			go func() {
				if !strings.Contains(config.Conf.LogLevels, plumelogInfo.LogLevel) {
					consumerGroup.PlumeInfoCh <- plumelogInfo
				}
			}()
			for key, ch := range consumerGroup.WsChMap {
				if model.ConnMap[key] == nil {
					delete(consumerGroup.WsChMap, key)
					continue
				}
				ch <- plumelogInfo

			}
			rwLock.Unlock()
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func putService(serviceName string) error {
	m := make(map[string]string, 100)
	m["serviceName"] = serviceName
	jsonStr, _ := json.Marshal(m)
	query := elastic.NewTermQuery("serviceName", serviceName)
	do, err := plumelogEs.Client.Search("plume_log_services").Query(query).Do(context.Background())
	if err != nil {
		return err
	}
	if do.TotalHits() > 0 {
		return nil
	}
	_, err = plumelogEs.Client.Index().Index("plume_log_services").BodyJson(string(jsonStr)).Do(context.Background())
	return err
}
