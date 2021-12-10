package tests

import (
	"github.com/awensir/go-aurora/aurora"
	"testing"
)

type ServeHandler interface {
	Controller(*aurora.Ctx, ...func(c *aurora.Ctx) interface{}) interface{}
}

type ServeFun func(Serve) Serve

type Serve func(c *aurora.Ctx) interface{}

func (s Serve) Controller(c *aurora.Ctx) interface{} {

	return s(c)
}

func A(c *aurora.Ctx) interface{} {
	return "aaa"
}

func B(serve Serve) Serve {
	return func(c *aurora.Ctx) interface{} {
		defer func() {

		}()
		return serve(c)
	}
}

func TestInterceptor(t *testing.T) {

}
