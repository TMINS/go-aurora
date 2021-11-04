package aurora

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Ctx struct {
	rw        *sync.RWMutex
	ar        *Aurora // Aurora 引用
	monitor   *localMonitor
	Response  http.ResponseWriter
	Request   *http.Request
	Args      map[string]interface{} //REST API 参数
	Attribute map[string]interface{} //Context属性
}

// INFO 打印 info 日志信息
func (c *Ctx) INFO(info ...interface{}) {
	c.ar.serviceInfo <- fmt.Sprint(info...)
}

// WARN 打印 警告信息
func (c *Ctx) WARN(warning ...interface{}) {
	c.ar.serviceWarning <- fmt.Sprint(warning...)
}

// ERROR 打印错误信息
func (c *Ctx) ERROR(error ...interface{}) {
	c.ar.serviceError <- fmt.Sprint(error...)
}

// PANIC 打印信息并且结束程序
func (c *Ctx) PANIC(panic ...interface{}) {
	c.ar.servicePanic <- fmt.Sprint(panic...)
}

// Get 获取加载 ,需要转换类型后使用
func (c *Ctx) Get(name string) interface{} {
	return c.ar.Get(name)
}

// JSON 向浏览器输出json数据
func (c *Ctx) json(data interface{}) {
	s, b := data.(string) //返回值如果是json字符串或者直接是字符串，将不再转码,json 二次转码对原有的json格式会进行二次转义
	if b {
		_, err := c.Response.Write([]byte(s))
		if err != nil {
			c.monitor.En(executeInfo(err))
			c.ar.runtime <- c.monitor
		}
		return
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		c.monitor.En(executeInfo(err))
		c.ar.runtime <- c.monitor
		return
	}
	_, err = c.Response.Write(marshal)
	if err != nil {
		c.monitor.En(executeInfo(err))
		c.ar.runtime <- c.monitor
	}
}

// RequestForward 内部路由转发，会携带 Ctx 本身进行转发，转发之后继续持有该 Ctx
func (c *Ctx) forward(path string) {
	c.ar.router.SearchPath(c.Request.Method, path, c.Response, c.Request, c, c.monitor)
}

// Redirect 发送重定向
func (c *Ctx) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}
