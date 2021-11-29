package main

import (
	"github.com/awensir/go-aurora/aurora"
	"google.golang.org/grpc"
)

type Student struct {
	Id     int    `gorm:primaryKey`
	Name   string `gorm:column:name`
	Gender string `gorm:column:gender`
	Age    int    `gorm:column:age`
}

func main() {

	//获取 aurora 路由实例
	a := aurora.New()
	grpc.NewServer()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
