package main

import (
	"github.com/awensir/Aurora/aurora"
	"github.com/awensir/Aurora/aurora/frame"
	"testing"
)

func TestLoading(t *testing.T) {
	gorms := &aurora.GORM{nil}
	aurora.Container.Store(frame.GORM, gorms)

	g := aurora.Container.Get(frame.GORM)
	if gorms == g {
		a := g.(*aurora.GORM)
		a.String()
	}
}
