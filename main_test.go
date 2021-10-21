package main

import (
	"fmt"
	"testing"
)

func Ttt(a ...interface{}) {
	for i, _ := range a {
		fmt.Println(a[i])
	}
}
func TestLoading(t *testing.T) {
	a := []int{1, 2, 3}
	Ttt(a)
}
