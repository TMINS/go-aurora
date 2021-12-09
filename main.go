package main

import (
	"github.com/awensir/go-aurora/aurora"
	"log"
)

type Mm1 struct {
}

func (m *Mm1) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm1")
	return true
}

func (m *Mm1) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm1")
}

func (m *Mm1) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm1")
}

type Mm2 struct {
}

func (m *Mm2) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm2")
	return false
}

func (m *Mm2) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm2")
}

func (m *Mm2) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm2")
}

type Mm3 struct {
}

func (m *Mm3) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm3")
	return true
}

func (m *Mm3) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm3")
}

func (m *Mm3) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm3")
}

func main() {

	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		log.Println("service..")
		return nil
	})

	a.RouteIntercept("/", &Mm1{}, &Mm2{}, &Mm3{})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
