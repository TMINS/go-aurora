package aurora

import (
	"Aurora/logs"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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
	Args      map[string]string      //REST API 参数
	QueryArgs map[string]string      //普通k/v
	Attribute map[string]interface{} //Context属性
	sessionV  *Session               //session变量
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
			logs.Info(err.Error())
			return false
		}
		return true
	}
	return false
}

// GetSession 服务器获取Session，存在session则返回，不再存在则创建
func (c *Context) GetSession() *Session {
	var session *Session
	if c.sessionV != nil { //查看上下文中是否已经获取了session
		return c.sessionV
	}
	cookie, err := c.Request.Cookie("SESSIONID") //请求中读取session
	if err != nil {                              //读取不到session 则创建一个
		aurora.rw.Lock() //创建存储在全局
		//未查询到session 则创建一个session
		if session == nil { //避免二次创建
			id, err := sessionIdCreater.NextID()
			if err != nil {

			}
			IdValue := strconv.FormatUint(id, 10)
			session = NewSession(IdValue)
			sessionMap[IdValue] = session
		}
		aurora.rw.Unlock()
		c.Response.Header().Set("Set-Cookie", session.cookie.String()) //设置即将响应的响应头，发送给浏览器
		if c.sessionV == nil {
			c.sessionV = session //初始化请求上下文 session变量
		}
		return session
	} else {
		aurora.rw.RLock()
		session, _ := sessionMap[cookie.Value] //可能存在bug 如果伪造session这里就会出现问题***待解决
		session.Keep()                         //重置session存活时间以保持连接
		aurora.rw.RUnlock()
		h := c.Response.Header()
		if h.Get("Set-Cookie") != "" { //避免设置两次session cookie
			c.Response.Header().Set("Set-Cookie", session.cookie.String())
		}
		return session
	}
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

// NewCookie 创建一个cookie 默认存活60秒
func (c *Context) NewCookie(name, value string) *Cookie {
	return &Cookie{&http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: 60,
	}}
}

// AddCookie 添加cookie响应头
func (c *Context) AddCookie(cooke Cookie) {
	h := c.Response.Header()
	h.Add("Set-Cookie", cooke.Cookie.String())
}

// RequestUrl 获取请求接口url
func (c *Context) RequestUrl() string {
	return c.Request.RequestURI
}

// GetArgs 获取REST API 参数，查询不存在的key或者不存在REST API 参数则返回""
func (c *Context) GetArgs(key string) string {
	if c.Args != nil {
		if v, ok := c.Args[key]; ok {
			return v
		} else {
			//查询了一个不存在的key
			return ""
		}
	}
	return ""
}

// JSON 向浏览器输出json数据
func (c *Context) JSON(data interface{}) {
	s, b := data.(string) //返回值如果是json字符串或者直接是字符串，将不再转码,json 二次转码对原有的json格式会进行二次转义
	if b {
		_, err := c.Response.Write([]byte(s))
		if err != nil {
			logs.WebRequestError(err.Error())
		}
		return
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		logs.WebRequestError(err.Error())
		return
	}
	_, err = c.Response.Write(marshal)
	if err != nil {
		logs.WebRequestError(err.Error())
	}
}

// RequestForward 服务转发
func (c *Context) RequestForward(path string)  {
	aurora.Router.SearchPath(c.Request.Method,path,c.Response,c.Request,c)
}
