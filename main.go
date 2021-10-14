package main

import (
	"github.com/awensir/Aurora/aurora"
)

func main() {
	a := aurora.Default()
	a.ResourceMapping("js", "js")
	a.GET("/a/${name}", func(c *aurora.Context) interface{} {
		args := c.GetArgs()
		i := args["name"]
		return i
	})
	a.Guide()
}
