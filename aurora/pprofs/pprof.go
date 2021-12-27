package pprofs

import (
	"github.com/awensir/go-aurora/aurora"
	"net/http/pprof"
)

/*
	用于接口执行性能 的测试 处理器
*/

func Index(c *aurora.Ctx) interface{} {
	pprof.Index(c.Response, c.Request)
	return nil
}

func Profile(c *aurora.Ctx) interface{} {
	pprof.Profile(c.Response, c.Request)
	return nil
}
func Cmdline(c *aurora.Ctx) interface{} {
	pprof.Cmdline(c.Response, c.Request)
	return nil
}
func Symbol(c *aurora.Ctx) interface{} {
	pprof.Symbol(c.Response, c.Request)
	return nil
}
func Trace(c *aurora.Ctx) interface{} {
	pprof.Trace(c.Response, c.Request)
	return nil
}
