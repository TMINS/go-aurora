package aurora

import (
	"fmt"
	"github.com/awensir/go-aurora/aurora/option"
)

func (a *Aurora) GrpcConfig(opt Opt) {
	o := opt()
	s := o[option.GRPC_SERVER]
	l := o[option.GRPC_LISTEN]
	fmt.Println(s, l)
}
