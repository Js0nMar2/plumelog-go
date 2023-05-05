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
	result, _ := plumelogEs.Client.Search("plume_log_services").Do(context.Background())
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
	go cg.Consume(cg)
	ConsumerMap[config.Conf.GroupId] = cg
}

type ConsumerGroup struct {
	sarama.ConsumerGroup
	GroupId     string
	Topics      []string
	PlumeInfoCh chan model.PlumelogInfo
	WsChMap     map[string]chan model.PlumelogInfo
	HealthMap   map[string]chan []byte
}

type ConsumerGroupHandler interface {
	Consume(handler sarama.ConsumerGroupHandler)
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
	log.Info.Println("kafka version:", config.Version, "init address is:", addr, "topics is:", topics)
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
func (cg *ConsumerGroup) Consume(handler sarama.ConsumerGroupHandler) {
	ctx := context.Background()
	for {
		err := cg.ConsumerGroup.Consume(ctx, cg.Topics, handler)
		if err != nil {
			log.Error.Println(err)
			time.Sleep(1 * time.Second)
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
func (cg *ConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	for {
		if sess == nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	rwLock := new(sync.RWMutex)
	if claim.Topic() == "app-index" {
		messages := claim.Messages()
		for msg := range messages {
			err := putService(string(msg.Value))
			if err != nil {
				log.Error.Println("set elastic index err:", err)
				return err
			}
			sess.MarkMessage(msg, "")
		}
		return nil
	}
	go func() {
		bulkService := plumelogEs.Client.Bulk()
		for {
			select {
			case plumelogInfo := <-cg.PlumeInfoCh:
				timeStr := time.Now().Format("20060102")
				request := elastic.NewBulkIndexRequest().Index("plume_log_" + timeStr).Doc(plumelogInfo)
				bulkService.Add(request)
				_, err := bulkService.Do(context.Background())
				if err != nil {
					log.Error.Println(err)
				}
			}
		}
	}()
	//读取kafka数据记录到分片,使用定时批量读取分片内数据
	messages := claim.Messages()
	for msg := range messages {
		//log.Info.Println(string(msg.Value))
		b, err := putServiceStatus(msg.Value)
		if b || err != nil {
			if err == nil {
				for _, guardCh := range cg.HealthMap {
					guardCh <- msg.Value
				}
			}
			sess.MarkMessage(msg, "")
			continue
		}
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
					cg.PlumeInfoCh <- plumelogInfo
				}
			}()
			for key, ch := range cg.WsChMap {
				if model.ConnMap[key] == nil {
					delete(cg.WsChMap, key)
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
	query := elastic.NewBoolQuery().Filter(elastic.NewTermQuery("serviceName", serviceName)).Filter(elastic.NewIdsQuery().Ids(serviceName))
	do, err := plumelogEs.Client.Search("plume_log_services").Query(query).Do(context.Background())
	if err != nil {
		return err
	}
	if do.TotalHits() > 0 {
		return nil
	}
	m := make(map[string]string)
	m["serviceName"] = serviceName
	jsonStr, _ := json.Marshal(m)
	res, err := plumelogEs.Client.Index().Index("plume_log_services").BodyJson(string(jsonStr)).Id(serviceName).Do(context.Background())
	if res.Result == "created" {
		ConsumerMap[config.Conf.GroupId].Topics = append(ConsumerMap[config.Conf.GroupId].Topics, serviceName)
		log.Info.Println("put new service:", serviceName)
	}
	return err
}

func putServiceStatus(bytes []byte) (bool, error) {
	healthGuard := model.HealthGuard{}
	json.Unmarshal(bytes, &healthGuard)
	if healthGuard.Type == "" || healthGuard.Type != "healthGuard" {
		return false, nil
	}
	if healthGuard.Status == "DOWN" {
		query := elastic.NewBoolQuery().Filter(elastic.NewTermQuery("appName", healthGuard.AppName)).
			Filter(elastic.NewTermQuery("ip", healthGuard.Ip)).Filter(elastic.NewTermQuery("env", healthGuard.Env))
		_, err := plumelogEs.Client.DeleteByQuery("plume_log_services_status").Query(query).Do(context.Background())
		if err != nil {
			return true, err
		}
		return true, nil
	}
	m := make(map[string]string)
	m["appName"] = healthGuard.AppName
	m["ip"] = healthGuard.Ip
	m["env"] = healthGuard.Env
	m["time"] = healthGuard.Time
	jsonStr, _ := json.Marshal(m)
	_, err := plumelogEs.Client.Index().Index("plume_log_services_status").BodyJson(string(jsonStr)).Id(healthGuard.Env + "-" + healthGuard.AppName + "-" + healthGuard.Ip).Do(context.Background())
	return true, err
}
