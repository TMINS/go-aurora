package main

import (
	"fmt"
	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/minilog/mini"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	a.Level(mini.INFO)
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		fmt.Println("hello")
		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
