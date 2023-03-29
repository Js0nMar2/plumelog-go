package main

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	plumelogEs "plumelog/elastic"
)

//func main() {
//	// 1. new一个grpc的server
//	rpcServer := grpc.NewServer()
//	// 2. 将刚刚我们新建的ProdService注册进去
//	r := &RpcServer{
//		ch: make(chan *rpc.PlumelogInfo, 100),
//	}
//	rpc.RegisterPlumelogServiceServer(rpcServer, r)
//	// 3. 新建一个listener，以tcp方式监听8899端口
//	listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.Conf.Rpc.Port))
//	if err != nil {
//		log.Fatal("服务监听端口失败", err)
//	}
//	// 4. 运行rpcServer，传入listener
//	rpcServer.Serve(listener)
//}

func main() {
	putService("wf-authorization-center-service")
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
