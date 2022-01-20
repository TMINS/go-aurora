package main

import (
	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/go-aurora/aurora/grpc_test/grpc_server/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	//获取 aurora 路由实例

	a := aurora.New()
	{
		file, err := credentials.NewServerTLSFromFile("aurora/test_ca/rootcert.pem", "aurora/test_ca/rootkey.pem")
		if err != nil {
			return
		}
		server := grpc.NewServer(grpc.Creds(file))
		services.RegisterTestServiceServer(server, &services.TestService{})
		a.GRPC(server)
	}

	// GET 方法注册 web get请求
	// a.GET("/", func(c *aurora.Ctx) interface{} {
	// 	fmt.Println("hello")
	// 	return nil
	// })

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.GuideTLS("aurora/test_ca/rootcert.pem", "aurora/test_ca/rootkey.pem", "8088")
}
