package manage

import "sync"

// Container 全局容器
var Container = &Containers{rw: &sync.RWMutex{}}

type Variable interface {
	Clone() Variable
}

type Containers struct {
	rw         *sync.RWMutex
	prototypes map[string]Variable //容器存储的属性
}

// Get 获取容器内指定变量，不存在的key 返回默认零值
func (c *Containers) Get(name string) Variable {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.prototypes[name]
}

// Store 加载 指定变量
func (c *Containers) Store(name string, variable Variable) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.prototypes[name] = variable
}
