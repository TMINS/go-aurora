package main

import (
	"Aurora/aurora"
	"Aurora/aurora/start"
	"Aurora/config"
	"Aurora/request/get"
)

type Body struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	config.RegisterResource("js", "js", "test")
	//config.RegisterInterceptor(MyInterceptor1{})

	get.Mapping("/abc", func(ctx *aurora.Context) interface{} {

		return "/abc"
	})
	get.Mapping("/abc/bbc", func(ctx *aurora.Context) interface{} {

		return "/abc/bbc"
	})
	get.Mapping("/abc/bbc/asd", func(ctx *aurora.Context) interface{} {

		return "/abc/bbc/asd"
	})

	get.Mapping("/abc/bbc/aaa", func(ctx *aurora.Context) interface{} {

		return "/abc/bbc/aaa"
	})
	get.Mapping("/abc/qaq", func(ctx *aurora.Context) interface{} {

		return "/abc/qaq"
	})
	get.Mapping("/abc/qaq/csdn", func(ctx *aurora.Context) interface{} {

		return "/abc/qaq/csdn"
	})
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/"
	})
	config.RegisterPathInterceptor("/abc/*", MyInterceptor2{})

	config.RegisterPathInterceptor("/", MyInterceptor3{})

	config.RegisterPathInterceptor("/abc/bbc/aaa", MyInterceptor4{})

	start.Running("8080")
}
