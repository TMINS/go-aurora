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
	config.Resource("js", "js", "test")

	get.Mapping("/abc", func(ctx *aurora.Context) interface{} {
		ctx.RequestForward("/abc/bbc/asd")
		return nil
	})
	get.Mapping("/abc/bbc", func(ctx *aurora.Context) interface{} {

		return "forward:/abc/bbc/asd"
	})
	get.Mapping("/abc/bbc/asd", func(ctx *aurora.Context) interface{} {

		return "/abc/bbc/asd"
	})

	get.Mapping("/abc/bbc/aaa", func(ctx *aurora.Context) interface{} {
		session:=ctx.GetSession()
		return session.GetSessionId()
	})
	get.Mapping("/abc/qaq", func(ctx *aurora.Context) interface{} {

		return "/abc/qaq"
	})
	get.Mapping("/abc/qaq/csdn", func(ctx *aurora.Context) interface{} {

		return "/abc/qaq/csdn"
	})

	get.Mapping("/user/${name}/${age}", func(ctx *aurora.Context) interface{} {
		n:=ctx.GetArgs("name")
		a:=ctx.GetArgs("age")
		s:= struct {
			Name string
			Age  string
		}{n,a}

		return s
	})
	config.Interceptor(MyInterceptor1{})
	config.PathInterceptor("/abc/*", MyInterceptor2{})
	config.PathInterceptor("/", MyInterceptor3{})
	config.PathInterceptor("/abc/bbc/aaa", MyInterceptor4{})
	start.Running("8080")
}
