`使用protobuf生成.go文件`
1.创建`rpc/proto`目录，在`proto`目录下编写.proto文件
2.cd到`proto`目录下，执行`protoc --go_out=plugins=grpc:. server.proto`即可生成.go文件