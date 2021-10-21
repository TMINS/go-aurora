package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
)

type MyInterceptor struct {
}

func (de *MyInterceptor) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle")
	return true
}

func (de *MyInterceptor) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle")
}

func (de *MyInterceptor) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion")
}

type MyInterceptor1 struct {
}

func (de *MyInterceptor1) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle1")
	return true
}

func (de *MyInterceptor1) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle1")
}

func (de *MyInterceptor1) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion1")
}

type MyInterceptor2 struct {
}

func (de *MyInterceptor2) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle2")
	return true
}

func (de *MyInterceptor2) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle2")
}

func (de *MyInterceptor2) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion2")
}

type MyInterceptor3 struct {
}

func (de *MyInterceptor3) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle3")
	return true
}

func (de *MyInterceptor3) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle3")
}

func (de *MyInterceptor3) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion3")
}

type MyInterceptor4 struct {
}

func (de *MyInterceptor4) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle4")
	return true
}

func (de *MyInterceptor4) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle4")
}

func (de *MyInterceptor4) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion4")
}

type MyInterceptor5 struct {
}

func (de *MyInterceptor5) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle5")
	return true
}

func (de *MyInterceptor5) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle5")
}

func (de *MyInterceptor5) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion5")
}

type MyInterceptor6 struct {
}

func (de *MyInterceptor6) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle6")
	return true
}

func (de *MyInterceptor6) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle6")
}

func (de *MyInterceptor6) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion6")
}
