package main

import (
	"github.com/awensir/Aurora/aurora"
	"github.com/awensir/Aurora/aurora/start"
	"github.com/awensir/Aurora/request/get"
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
	p := []string{
		"/",
		"/a/b",
		"/a/c",
		"/a/d",
		"/a/e",
		"/a/f",
		"/ag/mac",
		"/a/bc",
		"/a/bbb",
		"/a/cab",
		"/a/ccc",
		"/a/mack",
		"/a/old",
		"/index",
		"/index/home",
		"/index/user",
		"/index/use",
		"/index/user/login",
		"/index/user/lout",
		"/index/user/update",
		"/aa",
		"/aaa",
	}
	for _, v := range p {
		get.Mapping(v, Servlet)
	}

	start.Running("8080")
}
