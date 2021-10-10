package main

import (
	"github.com/awensir/Aurora/aurora"

	"reflect"
)

func BuilderAPI(service interface{}) interface{} {
	of := reflect.ValueOf(service)
	if of.NumMethod() < 0 {
		return nil
	}
	api := make(map[string]aurora.Servlet)
	for i := 0; i < of.NumMethod(); i++ {
		method := of.Type().Method(i)
		FunName := method.Name
		Fun := of.Method(i).Interface().(aurora.Servlet)
		api[FunName] = Fun
	}
	return nil
}
