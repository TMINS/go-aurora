package mux

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Attribute interface {
	SetAttribute(key string, value interface{})
	GetAttribute(key string) interface{}
}

type Ctx struct {
	rw        *sync.RWMutex
	ar        *Aurora // Aurora 引用
	monitor   *LocalMonitor
	Response  http.ResponseWriter
	Request   *http.Request
	Args      map[string]interface{} //REST API 参数
	Attribute map[string]interface{} //Context属性
}

// JSON 向浏览器输出json数据
func (c *Ctx) json(data interface{}) {
	s, b := data.(string) //返回值如果是json字符串或者直接是字符串，将不再转码,json 二次转码对原有的json格式会进行二次转义
	if b {
		_, err := c.Response.Write([]byte(s))
		if err != nil {
			c.monitor.En(ExecuteInfo(err))
			c.ar.runtime <- c.monitor
		}
		return
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return
	}
	_, err = c.Response.Write(marshal)
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
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
