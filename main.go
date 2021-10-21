package main

import (
	"github.com/awensir/Aurora/aurora"
)

func main() {

	//获取 aurora 路由实例
	a := aurora.New()
	//a.RouteIntercept("",&MyInterceptor{})
	//a.RouteIntercept("/",&MyInterceptor{},&MyInterceptor1{},&MyInterceptor2{})
	//a.RouteIntercept("/b",&MyInterceptor{})
	a.RouteIntercept("/a/*", &MyInterceptor1{}, &MyInterceptor2{})
	//a.RouteIntercept("/a/b/*",&MyInterceptor2{})
	//a.RouteIntercept("/a/b/c*",&MyInterceptor3{})
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a", func(c *aurora.Ctx) interface{} {

		return c.Args
	})
	a.GET("/b", func(c *aurora.Ctx) interface{} {

		return c.Args
	})
	a.GET("/a/b", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a/b/c/{name}", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a/b/cc", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
