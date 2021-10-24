package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
)

func main() {

	//获取 aurora 路由实例
	a := aurora.New()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		get, err := c.Get("obj")
		if err != nil {
			return err
		}
		return get
	})

	a.POST("/", func(c *aurora.Ctx) interface{} {
		var post interface{}
		err := c.JsonBody(&post)
		if err != nil {
			return err
		}
		fmt.Println(post)
		return post
	})

	a.POST("/upload", func(c *aurora.Ctx) interface{} {
		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
