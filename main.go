package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {

	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/s")

		return nil
	})

	a.GET("/a", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/s")

		return nil
	})

	a.GET("/b", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/s")

		return nil
	})
	a.GET("/b/{name}/v", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/s")

		return nil
	})
	a.GET("/b/{name}/age", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/s")

		return nil
	})

	a.GET("/ab/s", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/s")

		return nil
	})
	a.GET("/ab/c/a", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/c/a")

		return nil
	})

	a.GET("/ab", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab")

		return nil
	})

	a.GET("/ab/c/a", func(c *aurora.Ctx) interface{} {
		c.INFO("/ab/c/a  2")

		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
