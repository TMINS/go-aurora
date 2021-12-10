package tests

import (
	"github.com/awensir/go-aurora/aurora"
	"testing"
)

type ServeHandler interface {
	Controller(*aurora.Ctx, ...func(c *aurora.Ctx) interface{}) interface{}
}

type ServeFun func(Serve) interface{}

type Serve func(c *aurora.Ctx) interface{}

func (s Serve) Controller(c *aurora.Ctx) interface{} {

	return s(c)
}

func TestInterceptor(t *testing.T) {

}
