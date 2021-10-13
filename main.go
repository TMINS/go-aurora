package main

import (
	"github.com/awensir/Aurora/aurora"
	"github.com/awensir/Aurora/aurora/start"
	"github.com/awensir/Aurora/request/get"
)

func main() {
	get.Mapping("/", func(ctx *aurora.Context) interface{} {
		v, err := ctx.GetIntSlice("arr1")
		if err != nil {
			return err
		}
		return v
	})
	start.Running("8080")
}
