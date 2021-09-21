package main

import (
	"github.com/awensir/Aurora/aurora"
	"github.com/awensir/Aurora/aurora/start"
	"github.com/awensir/Aurora/config"
	"github.com/awensir/Aurora/request/delet"
	"github.com/awensir/Aurora/request/get"
	"github.com/awensir/Aurora/request/post"
	"github.com/awensir/Aurora/request/put"
)

func main() {
	config.Resource("js", "js", "test")

	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})
	post.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})
	put.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})
	delet.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})

	start.Running("8080")
}
