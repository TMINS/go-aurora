package main

import (
	"Aurora/aurora"
	"Aurora/config"
	"Aurora/request/get"
	"fmt"
)

func main() {
	config.RegisterResource("js", "js", "test")
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		//return "/html/index.html"
		return TestErr{fmt.Errorf("err")}
	})

	aurora.RunApplication("8080")
}
