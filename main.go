package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
	"github.com/awensir/Aurora/aurora/start"
	"github.com/awensir/Aurora/config"
	"github.com/awensir/Aurora/request/get"
	"sync"
)

func Test(ctx *aurora.Context) interface{} {
	return "test"
}

func P(group *sync.WaitGroup, ch chan int) {
	for v := range ch {
		fmt.Println(v)
		group.Done()
	}
}
func main() {
	config.Resource("js", "js", "test")

	get.Mapping("/", func(next aurora.Servlet) aurora.Servlet {
		return func(ctx *aurora.Context) interface{} {
			fmt.Println("before")
			v := next(ctx)
			fmt.Println("after")
			return v
		}
	}(Test))

	start.Running("8080")
}
