package main

import (
	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/go-aurora/aurora/pprofs"
)

func main() {

	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		return nil
	})
	a.GET("/debug/pprof/heap", pprofs.Index)
	a.GET("/debug/pprof/cmdline", pprofs.Cmdline)
	a.GET("/debug/pprof/profile", pprofs.Profile)
	a.GET("/debug/pprof/symbol", pprofs.Symbol)
	a.GET("/debug/pprof/trace", pprofs.Trace)
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
