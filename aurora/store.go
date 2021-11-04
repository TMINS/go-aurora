package aurora

import (
	"github.com/awensir/Aurora/aurora/frame"
	"sync"
)

type containers struct {
	rw         *sync.RWMutex
	prototypes map[string]interface{} //容器存储的属性
}

// Get 获取容器内指定变量，不存在的key 返回默认零值
func (c *containers) get(name string) interface{} {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.prototypes[name]
}

// Store 加载 指定变量
func (c *containers) store(name string, variable interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.prototypes[name] = variable
}

// Delete 提供删除自定义整合数据
func (c *containers) Delete(name string) {
	switch name {
	// 内置整合 不允许删除
	case frame.DB, frame.GORM, frame.GO_REDIS:
		return
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, b := c.prototypes[name]; !b {
		return
	} else {
		delete(c.prototypes, name)
	}
}
