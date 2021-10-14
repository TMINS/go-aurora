package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
)

func Test(ctx *aurora.Context) interface{} {
	return "test"
}
func main() {
	//获取 aurora 路由实例
	a := aurora.Default()

	// GET 方法注册 web get请求
	a.GET("/", func(next aurora.Servlet) aurora.Servlet {
		return func(ctx *aurora.Context) interface{} {
			fmt.Println("before")
			v := next(ctx)
			fmt.Println("after")
			return v
		}
	}(Test))
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
