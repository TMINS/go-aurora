package main

import (
	"Aurora/logs"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	c:=time.NewTicker(time.Second)
	
	for true {
		select {
		case t:=<-c.C:
			logs.Info("c",t)
		}
	}
}