package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
)

type Service struct {
	ServiceName string
	ServiceAddr string
	ServicePort string
	ApiNum      int
}

type RpcService struct{}

func (as *RpcService) Pay(ctx *aurora.Context) interface{} {
	fmt.Println("Pay")
	return nil
}

func (as *RpcService) GetPort(ctx *aurora.Context) interface{} {
	fmt.Println("GetPort")
	return nil
}

func (as *RpcService) CreatOrder(ctx *aurora.Context) interface{} {
	fmt.Println("CreatOrder")
	return nil
}

func (as *RpcService) DeleteOrder(ctx *aurora.Context) interface{} {
	fmt.Println("DeleteOrder")
	return nil
}
