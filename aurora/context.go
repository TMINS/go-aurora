package aurora

import (
	"Aurora/logs"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Attribute interface {
	SetAttribute(key string,value interface{})
	GetAttribute(key string) interface{}
}

type Context struct {
	http.ResponseWriter
	*http.Request
	rw sync.RWMutex
	Args                map[string]string       //REST API 参数
	QueryArgs           map[string]string       //普通k/v
	Attribute           map[string]interface{}  //Context属性
	sessionV            *Session                //session变量
}

// GetSession 服务器获取Session，存在session则返回，不再存在则创建
func (c *Context) GetSession() *Session {
	var session *Session
	if c.sessionV!=nil{
		return c.sessionV
	}
	cookie, err := c.Request.Cookie("SESSIONID")
	if err != nil {
		aurora.rw.Lock()
		//未查询到session 则创建一个session
		if session==nil{
			id,err:=sessionIdCreater.NextID()
			if err!=nil {
			
			}
			IdValue:=strconv.FormatUint(id,10)
			session=NewSession(IdValue)
			sessionMap[IdValue]=session
		}
		aurora.rw.Unlock()
		c.ResponseWriter.Header().Set("Set-Cookie",session.cookie.String())
		if c.sessionV==nil{
			c.sessionV=session  //初始化请求上下文 session变量
		}
		return session
	}else {
		aurora.rw.RLock()
		session,_:=sessionMap[cookie.Value]
		session.Keep()
		aurora.rw.RUnlock()
		h:=c.ResponseWriter.Header()
		if h.Get("Set-Cookie")!=""{
			c.ResponseWriter.Header().Set("Set-Cookie",session.cookie.String())
		}
		return session
	}
}

func (c *Context) SetAttribute(key string,value interface{})  {
	c.rw.Lock()
	defer c.rw.Unlock()
	if c.Attribute==nil{
		c.Attribute=make(map[string]interface{})
	}
	c.Attribute[key]=value
}

func (c *Context) GetAttribute(key string) interface{}  {
	if c.Attribute==nil{
		return nil
	}
	c.rw.RLock()
	defer c.rw.RUnlock()
	if v,ok:=c.Attribute[key];ok{
		return v
	}
	return nil
}

func (c *Context) NewCookie(name , value string) *Cookie {
	return &Cookie{&http.Cookie{
		Name: name,
		Value: value,
		MaxAge: 60,
	}}
}

func (c *Context) AddCookie(cooke Cookie) {
	  h:=c.ResponseWriter.Header()
	  h.Add("Set-Cookie",cooke.Cookie.String())
}

// RequestUrl 获取请求接口url
func (c *Context) RequestUrl() string {
	return c.RequestURI
}

// GetArgs 获取REST API 参数
func (c *Context) GetArgs(key string) string {
	if c.Args!=nil{
		if v,ok:=c.Args[key];ok{
			return v
		}else {
			//查询了一个不存在的key
			return ""
		}
	}
	return ""
}



func (c *Context) Get(key string) string {
	return c.URL.Query().Get(key)
}

func (c *Context) JSON(data interface{}) {
	marshal, err := json.Marshal(data)
	if err != nil {
		logs.WebRequestError(err.Error())
		return
	}
	_, err = c.ResponseWriter.Write(marshal)
	if err != nil {
		logs.WebRequestError(err.Error())
	}
}
