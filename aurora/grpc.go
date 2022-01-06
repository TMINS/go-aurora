package aurora

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

/*
	待定
*/

// Grpc 整合grpc 配置
func (a *Aurora) Grpc(server *grpc.Server) {
	a.grpc = server
	//初始化健康检查,使用grpc的默认实现
	grpc_health_v1.RegisterHealthServer(a.grpc, health.NewServer())
}
