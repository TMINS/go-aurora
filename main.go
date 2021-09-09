package main

import (
	"Aurora/aurora"
	"Aurora/config"
	"Aurora/request/get"
	"Aurora/request/post"
)

func main() {
	config.RegisterResource("js", "js")

	//config.RegisterInterceptor(MyInterceptor1{})
	config.ResourceRoot("resource")
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})

	post.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})

	aurora.RunApplication("8080")
}
