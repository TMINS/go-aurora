package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
	"testing"
)

type My struct{}

func (de My) PreHandle(ctx *aurora.Context) bool {

	return true
}

func (de My) PostHandle(ctx *aurora.Context) {

}

func (de My) AfterCompletion(ctx *aurora.Context) {

}

func Servlet(ctx *aurora.Context) interface{} {
	arr := []string{
		"aaaa",
		"bbbb",
		"cccc",
	}
	return arr
}

func TestRunning(t *testing.T) {
	sss(1)
}

func sss(a ...int) {
	if a == nil {
		fmt.Println("nil")
	}
}
