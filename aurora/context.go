package aurora

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Attribute interface {
	SetAttribute(key string, value interface{})
	GetAttribute(key string) interface{}
}

type Context struct {
	Response  http.ResponseWriter
	Request   *http.Request
	rw        *sync.RWMutex
	args      map[string]interface{} //REST API 参数
	Attribute map[string]interface{} //Context属性
	AR        *Aurora                // Aurora 引用
}

// GetValue 获取指定key的查询参数,指定key不存在则返回""
func (c *Context) GetValue(key string) string {
	return c.Request.URL.Query().Get(key)
}

// PostBody 读取Post请求体数据解析到body中
func (c *Context) PostBody(body interface{}) bool {
	data, err := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err == nil {
		err := json.Unmarshal(data, &body)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		return true
	}
	return false
}

// SetAttribute 设置属性值
func (c *Context) SetAttribute(key string, value interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if c.Attribute == nil {
		c.Attribute = make(map[string]interface{})
	}
	c.Attribute[key] = value
}

// GetAttribute 获取属性值
func (c *Context) GetAttribute(key string) interface{} {
	if c.Attribute == nil {
		return nil
	}
	c.rw.RLock()
	defer c.rw.RUnlock()
	if v, ok := c.Attribute[key]; ok {
		return v
	}
	return nil
}

// JSON 向浏览器输出json数据
func (c *Context) JSON(data interface{}) {
	s, b := data.(string) //返回值如果是json字符串或者直接是字符串，将不再转码,json 二次转码对原有的json格式会进行二次转义
	if b {
		_, err := c.Response.Write([]byte(s))
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = c.Response.Write(marshal)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// RequestForward 服务转发
func (c *Context) RequestForward(path string) {
	c.AR.Router.SearchPath(c.Request.Method, path, c.Response, c.Request, c)
}
