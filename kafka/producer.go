package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
	"plumelog/config"
	plumelogEs "plumelog/elastic"
	"plumelog/log"
	"plumelog/utils"
)

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

var ProducerMap map[string]*Producer

func init() {
	ProducerMap = make(map[string]*Producer)
	result, _ := plumelogEs.Client.Search("plume_log_services").Do(context.Background())
	for _, hit := range result.Hits.Hits {
		marshalJSON, _ := hit.Source.MarshalJSON()
		m := make(map[string]string)
		json.Unmarshal(marshalJSON, &m)
		serviceName := m["serviceName"]
		// 初始化生产者
		if serviceName != "" {
			producer := NewKafkaProducer([]string{config.Conf.Kafka.Host}, serviceName)
			ProducerMap[serviceName] = producer
		}
	}
}

func NewKafkaProducer(addr []string, topic string) *Producer {
	p := Producer{}
	p.addr = addr
	p.topic = topic
	p.config = sarama.NewConfig()             //Instantiate a sarama Config
	p.config.Producer.Return.Successes = true //Whether to enable the successes channel to be notified after the message is sent successfully
	p.config.Producer.Return.Errors = true
	p.config.Producer.RequiredAcks = sarama.WaitForAll        //Set producer Message Reply level 0 1 all
	p.config.Producer.Partitioner = sarama.NewHashPartitioner //Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly
	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
	if err != nil {
		panic(err.Error())
		return nil
	}
	p.producer = producer
	return &p
}

func (p *Producer) SendMessage(m proto.Message, topic, key string) (int32, int64, error) {
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = topic
	kMsg.Key = sarama.StringEncoder(key)
	bMsg, err := proto.Marshal(m)
	if err != nil {
		log.Error("proto marshal err = %s", err.Error())
		return -1, -1, err
	}
	if len(bMsg) == 0 {
		log.Error("len(bMsg) == 0 ")
		return 0, 0, errors.New("len(bMsg) == 0 ")
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)
	if kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 {
		log.Error("kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 ", kMsg)
		return -1, -1, errors.New("key or value == 0")
	}
	a, b, c := p.producer.SendMessage(kMsg)
	return a, b, utils.Wrap(c, "")
}
