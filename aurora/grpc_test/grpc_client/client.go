package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	//log := mini.NewLog(mini.DEBUG)
	// GET 方法注册 web get请求
	// a.GET("/", func(c *aurora.Ctx) interface{} {
	// 	//不采用身份认证的方式连接grpc服务端
	// 	tlsFromFile, err := credentials.NewClientTLSFromFile("aurora/test_ca/rootcert.pem", "localhost")
	// 	if err != nil {
	// 		log.Error(err.Error())
	// 		return err
	// 	}
	// 	conn, _ := grpc.Dial("localhost:8088", grpc.WithTransportCredentials(tlsFromFile))
	// 	// error handling omitted
	// 	client := services.NewTestServiceClient(conn)
	// 	prc, err := client.TestPRC(context.Background(), &services.TestRequest{Message: "msg"})
	// 	if err != nil {
	// 		log.Error(err.Error())
	// 		return err
	// 	}
	// 	return prc
	// })

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide("8089")
}
