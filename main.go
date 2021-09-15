package main

import (
	"github.com/awensir/Aurora/aurora"
	"github.com/awensir/Aurora/aurora/start"
	"github.com/awensir/Aurora/config"
	"github.com/awensir/Aurora/request/get"
)

func main() {
	config.Resource("js", "js", "test")

	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})
	start.Running("8080")
}
