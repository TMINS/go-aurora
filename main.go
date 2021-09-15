package main

import (
	"Aurora/aurora"
	"Aurora/aurora/start"
	"Aurora/config"
	"Aurora/request/get"
)

func main() {
	config.Resource("js", "js", "test")

	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/"
	})
	start.Running("8080")
}
