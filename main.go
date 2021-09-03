package main

import (
	"Aurora/aurora"
	"Aurora/config"
	"Aurora/logs"
	"Aurora/request/get"
)

func main() {
	config.RegisterResource("js","js")
	
	//config.RegisterInterceptor(MyInterceptor1{})
	
	get.Mapping("/", func(ctx *aurora.Context) interface{} {
		session:=ctx.GetSession()
		logs.Info(session.GetSessionId())
		return "/html/index.html"
	})
	
	
	aurora.RunApplication("8080")
	
}

