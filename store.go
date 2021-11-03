package main

import "sync"

const (
	GORM = "gorm" // 全局默认 gorm作为 gorm 实例的key
	DB   = "db"   // db作为原生 db key
)

// Container 全局容器
var Container = &Containers{rw: &sync.RWMutex{}}

type Variable interface {
	Clone() Variable
}

type Containers struct {
	rw         *sync.RWMutex
	prototypes map[string]Variable //容器存储的属性
}

// Get 获取指定 变量
func (c *Containers) Get(name string) Variable {
	return c.prototypes[name]
}

// Store 加载 指定变量
func (c *Containers) Store(name string, variable Variable) {
	c.prototypes[name] = variable
}
