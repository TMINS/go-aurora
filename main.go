package main

import (
	"github.com/awensir/Aurora/mux"
)

func main() {

	//获取 mux 路由实例
	a := mux.New()
	// GET 方法注册 web get请求
	a.GET("/c/{a}", func(c *mux.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a/d/b", func(c *mux.Ctx) interface{} {

		return c.Args
	})
	a.GET("/", func(c *mux.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a/b", func(c *mux.Ctx) interface{} {

		return c.Args
	})

	a.GET("/d", func(c *mux.Ctx) interface{} {

		return c.Args
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide("8888")
}
