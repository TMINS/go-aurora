package main

/*
	grpc 整合 consul 注册中心
*/

import (
	"github.com/awensir/go-aurora/aurora/consul_test/services"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
)

func main() {
	server := grpc.NewServer()
	services.RegisterTestServiceServer(server, &services.TestService{})
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	listen, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	check := &api.AgentServiceCheck{
		Timeout:                        "5s",
		Interval:                       "4s",
		DeregisterCriticalServiceAfter: "30s",
		GRPC:                           "localhost:8088",
	}
	service := &api.AgentServiceRegistration{}
	service.Name = "grpc"         //Name 属性标识在consul中的服务名称
	service.ID = "grpc"           //ID属性是基于Name属性下面的编号,使用的时候不应该出现重复(准备时间戳+name属性来标识id)
	service.Check = check         //Check 属性用于配置服务健康检查，相对的Checks可以配置多个
	service.Address = "localhost" //设置服务地址信息
	service.Port = 8088           //设置服务端口信息
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	agent := client.Agent()
	err = agent.ServiceRegister(service)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	err = server.Serve(listen)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
