package aurora

import (
	"google.golang.org/grpc"
)

/*
	待定
*/

// GRPC 整合grpc 配置
func (a *Aurora) GRPC(server *grpc.Server) {
	a.grpc = server
}

// Grpc 获取grpc 实例
func (a *Aurora) Grpc() *grpc.Server {
	return a.grpc
}
