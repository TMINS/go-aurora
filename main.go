package main

import (
	"fmt"

	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例

	a := aurora.New()

	// GET 方法注册 web get请求

	a.GET("/test", func(params aurora.HttpRequest) interface{} {

		// 视图解析返回的路径必须是基于静态资源的路径之下开始，不需要斜杠开头
		return "html/index.html"
	})

	a.PUT("/test", func(params aurora.HttpRequest) interface{} {
		fmt.Println(params)
		return nil
	})

	a.POST("/test2", func(params aurora.HttpRequest) interface{} {
		fmt.Println(params)
		return nil
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
