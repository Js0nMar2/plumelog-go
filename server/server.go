package main

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"plumelog/config"
	plumeEs "plumelog/elastic"
	"plumelog/kafka"
	"plumelog/log"
	"plumelog/rpc"
)

type RpcServer struct {
	ch       chan *rpc.PlumelogInfo
	isOpenWs bool
}

// SendPlumelog 收集日志
func (p *RpcServer) SendPlumelog(ctx context.Context, req *rpc.PlumelogRequest) (*rpc.PlumelogResponse, error) {
	var plumelogInfos []string
	err := json.Unmarshal([]byte(req.Message), &plumelogInfos)
	if err != nil {
		log.Error(err.Error())
	}
	if req.Topic == "plumelog" {
		p.isOpenWs = false
		info := <-p.ch
		return &rpc.PlumelogResponse{Message: info.String()}, nil
	}
	go func() {
		if !p.isOpenWs {
			<-p.ch
		}
	}()
	if len(plumelogInfos) > 0 {
		for _, info := range plumelogInfos {
			plumelogInfo := rpc.PlumelogInfo{}
			json.Unmarshal([]byte(info), &plumelogInfo)
			if plumelogInfo.AppName == "" {
				continue
			}
			if plumelogInfo.TraceId == "" {
				plumelogInfo.TraceId = "0"
			}
			kafka.ProducerMap[req.Topic].SendMessage(&plumelogInfo, req.Topic, plumelogInfo.TraceId)
			p.ch <- &plumelogInfo
		}
	}
	return &rpc.PlumelogResponse{Message: "ok"}, nil
}

// PutService 添加服务设置topic生产者
func (*RpcServer) PutService(ctx context.Context, req *rpc.PlumelogServiceRequest) (*rpc.PlumelogResponse, error) {
	m := make(map[string]string, 100)
	m["serviceName"] = req.ServiceName
	jsonStr, _ := json.Marshal(m)
	query := elastic.NewTermQuery("serviceName", req.ServiceName)
	do, err := plumeEs.Client.Search("plume_log_services").Query(query).Do(ctx)
	if err != nil {
		return &rpc.PlumelogResponse{Message: err.Error()}, err
	}
	if do.TotalHits() > 0 {
		return &rpc.PlumelogResponse{Message: req.ServiceName + " is already exists"}, err
	}
	producer := kafka.NewKafkaProducer([]string{config.Conf.Kafka.Host}, req.ServiceName)
	kafka.ProducerMap[req.ServiceName] = producer
	response, err := plumeEs.Client.Index().Index("plume_log_services").BodyJson(string(jsonStr)).Do(ctx)
	cm := &kafka.ConsumerGroup{
		GroupId: req.ServiceName,
		Topics:  []string{req.ServiceName},
	}
	kafka.ConsumerMap[req.ServiceName] = cm
	return &rpc.PlumelogResponse{Message: response.Result}, err
}

func (x *RpcServer) WritePlumeLog(ctx context.Context, req *rpc.PlumelogInfo) (*rpc.PlumelogInfo, error) {
	if req.AppName != "" {
		x.isOpenWs = true
	}
	return <-x.ch, nil
}
