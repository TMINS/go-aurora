package main

import (
	"context"

	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/go-aurora/aurora/grpc_test/grpc_client/services"
	"github.com/awensir/minilog/mini"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	log := mini.NewLog(mini.DEBUG)
	// GET 方法注册 web get请求
	a.GET("/", func(request aurora.HttpRequest) interface{} {
		//该测试，用于整合grpc 和 http 共同端口，测试用的ca文件可能存在过期
		tlsFromFile, err := credentials.NewClientTLSFromFile("aurora/test_ca/rootcert.pem", "localhost")
		if err != nil {
			log.Error(err.Error())
			return err
		}
		conn, _ := grpc.Dial("localhost:8088", grpc.WithTransportCredentials(tlsFromFile))
		// error handling omitted
		client := services.NewTestServiceClient(conn)
		prc, err := client.TestPRC(context.Background(), &services.TestRequest{Message: "msg"})
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return prc
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide("8089")
}
