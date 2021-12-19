package aurora

import "google.golang.org/grpc"

/*
	待定
*/

// GrpcConfig 整合grpc 配置
func (a *Aurora) GrpcConfig(s *grpc.Server) {
	a.grpc = s
}
