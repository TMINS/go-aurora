package aurora

import (
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

type Ctx struct {
	rw        *sync.RWMutex
	ar        *Aurora // Aurora 引用
	Response  http.ResponseWriter
	Request   *http.Request
	Args      map[string]interface{} //REST API 参数
	Attribute map[string]interface{} //Context属性
}

// Viper 获取项目配置实例,未启动配置则返回nil
func (c *Ctx) Viper() *viper.Viper {
	return c.ar.Viper()
}

// GetRoot 获取项目根路径
func (c *Ctx) GetRoot() string {
	return c.ar.projectRoot
}

// Get 获取加载 ,需要转换类型后使用
func (c *Ctx) Get(name string) interface{} {
	return c.ar.Get(name)
}

// JSON 向浏览器输出json数据
func (c *Ctx) json(data interface{}) {
	s, b := data.(string) //返回值如果是json字符串或者直接是字符串，将不再转码,json 二次转码对原有的json格式会进行二次转义
	if b {
		c.Response.WriteHeader(http.StatusOK)
		_, err := c.Response.Write([]byte(s))
		if err != nil {
			c.ar.errMessage <- err.Error()
		}
		return
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		c.ar.errMessage <- err.Error()
		return
	}
	c.Response.WriteHeader(http.StatusOK)
	_, err = c.Response.Write(marshal)
	if err != nil {
		c.ar.errMessage <- err.Error()
	}
}

// RequestForward 内部路由转发，会携带 Ctx 本身进行转发，转发之后继续持有该 Ctx
func (c *Ctx) forward(path string) {
	c.ar.router.SearchPath(c.Request.Method, path, c.Response, c.Request, c)
}

// Redirect 发送重定向
func (c *Ctx) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}
