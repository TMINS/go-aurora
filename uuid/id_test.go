package uuid

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestUuid(t *testing.T) {
	w := NewWorker(1,1)
	id,err:=w.NextID()
	if err!=nil{
		t.Error(err.Error())
	}
	t.Log(id)
	t.Log(strconv.FormatUint(id,10))
}

func TestTime(t *testing.T) {
	c:=time.NewTicker(time.Second*2)
	
	for true {
		select {
		case <-c.C:
			fmt.Printf("c")
		}
	}
}
