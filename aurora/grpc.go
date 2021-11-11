package aurora

import (
	"fmt"
	"github.com/awensir/Aurora/aurora/option"
)

const ()

func (a *Aurora) GrpcConfig(opt Opt) {
	o := opt()
	s := o[option.GRPC_SERVER]
	l := o[option.GRPC_LISTEN]
	fmt.Println(s, l)
}
