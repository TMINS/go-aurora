package main

import (
	"Aurora/aurora"
	"Aurora/logs"
)

type MyInterceptor1 struct {
}

func (de MyInterceptor1) PreHandle(ctx *aurora.Context) bool {
	logs.Info("MyPreHandle111")
	return true
}

func (de MyInterceptor1) PostHandle(ctx *aurora.Context) {
	logs.Info("MyPostHandle111")
}

func (de MyInterceptor1) AfterCompletion(ctx *aurora.Context)  {
	logs.Info("MyAfterCompletion111")
}


type MyInterceptor2 struct {
}

func (de MyInterceptor2) PreHandle(ctx *aurora.Context) bool {
	logs.Info("MyPreHandle222")
	return true
}

func (de MyInterceptor2) PostHandle(ctx *aurora.Context) {
	logs.Info("MyPostHandle222")
}

func (de MyInterceptor2) AfterCompletion(ctx *aurora.Context)  {
	logs.Info("MyAfterCompletion222")
}

type MyInterceptor3 struct {
}

func (de MyInterceptor3) PreHandle(ctx *aurora.Context) bool {
	logs.Info("MyPreHandle333")
	return true
}

func (de MyInterceptor3) PostHandle(ctx *aurora.Context) {
	logs.Info("MyPostHandle333")
}

func (de MyInterceptor3) AfterCompletion(ctx *aurora.Context)  {
	logs.Info("MyAfterCompletion333")
}

type MyInterceptor4 struct {
}

func (de MyInterceptor4) PreHandle(ctx *aurora.Context) bool {
	logs.Info("MyPreHandle444")
	return true
}

func (de MyInterceptor4) PostHandle(ctx *aurora.Context) {
	logs.Info("MyPostHandle444")
}

func (de MyInterceptor4) AfterCompletion(ctx *aurora.Context)  {
	logs.Info("MyAfterCompletion444")
}


