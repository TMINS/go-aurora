package aurora

import (
	"fmt"
	"github.com/awensir/go-aurora/aurora/option"
	"google.golang.org/grpc"
)

func (a *Aurora) GrpcConfig(opt Opt) {
	o := opt()
	s := o[option.GRPC_SERVER]
	l := o[option.GRPC_LISTEN]
	fmt.Println(s, l)
}

// GetClient 获取grpc 远程调用客户端
func (a *Aurora) GetClient(raddr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	return nil, nil
}
