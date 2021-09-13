package main

import (
	"Aurora/aurora"
	"Aurora/config"
	"Aurora/request/get"
	"Aurora/request/post"
)

type Body struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	config.RegisterResource("js", "js", "test")
	config.RegisterInterceptor(MyInterceptor1{})
	post.Mapping("/", func(ctx *aurora.Context) interface{} {
		var body Body
		ctx.PostBody(&body)
		return body
	})

	get.Mapping("/abc", func(ctx *aurora.Context) interface{} {

		return "/abc"
	})
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/abc"
	})
	config.RegisterPathInterceptor("/abc", MyInterceptor2{})

	config.RegisterPathInterceptor("/", MyInterceptor3{})

	aurora.RunApplication("8080")
}
