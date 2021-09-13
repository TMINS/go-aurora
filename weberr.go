package main

import (
	"Aurora/aurora"
)

type TestErr struct {
	error
}

func (t TestErr) ErrorHandler(ctx *aurora.Context) interface{} {
	//对error 进行指定处理，选择输出

	return "/html/index.html"
}
