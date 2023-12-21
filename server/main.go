package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"plumelog/log"
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
	//putService("wf-authorization-center-service")
	for i := 0; i < 100; i++ {
		file, err := UploadFile("https://569bf5f41bb4af8ea658c6ed698d12f9.r2.cloudflarestorage.com/us168-image-p/fat/test/c3895b5eb4cb49d9bfe6636d691e37c41658398833643618304.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20230516T090748Z&X-Amz-SignedHeaders=host&X-Amz-Expires=3600&X-Amz-Credential=02f4094c2ffaa2164abb50d9c4a51a26%2F20230516%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=ccd6cd99514936c6c5c49ebc3beb54c06b5f88403935aad0122ca19bbc74b05e", "D:\\\\googleDownload\\\\img\\\\微信图片_20211026090842.png")
		if err != nil {
			log.Error(err.Error())
			//return
		}
		log.Info(string(file))
	}

}

//func putService(serviceName string) error {
//	m := make(map[string]string, 100)
//	m["serviceName"] = serviceName
//	jsonStr, _ := json.Marshal(m)
//	query := elastic.NewTermQuery("serviceName", serviceName)
//	do, err := plumelogEs.Client.Search("plume_log_services").Query(query).Do(context.Background())
//	if err != nil {
//		return err
//	}
//	if do.TotalHits() > 0 {
//		return nil
//	}
//	_, err = plumelogEs.Client.Index().Index("plume_log_services").BodyJson(string(jsonStr)).Do(context.Background())
//	return err
//}

func UploadFile(url string, fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "image/png")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
