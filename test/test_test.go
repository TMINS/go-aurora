package test

import (
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	st := reflect.StructOf([]reflect.StructField{
		{Name: "Name", Type: reflect.TypeOf(string("")), Tag: "name"},
		{Name: "Age", Type: reflect.TypeOf(int(0)), Tag: "age"},
	})
	sv := reflect.New(st).Elem()
	sv.Field(0).SetString("saber")
	sv.Field(1).SetInt(19)
	t.Log(st.String())
	d := sv.Interface()
	t.Log(d)
}
