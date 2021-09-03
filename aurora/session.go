package aurora

import (
	"net/http"
	"sync"
	"time"
)

type Session struct {
	cookie *Cookie
	rw sync.RWMutex
	MaxAge  time.Time
	Attribute  map[string]interface{}  //session 域
}

// NewSession 默认生成一个60秒的session
func NewSession(value string) *Session {
	t:=sessionAge()
	return &Session{cookie:&Cookie{&http.Cookie{
		Name: "SESSIONID",
		Value: value,
		MaxAge: 60,
		HttpOnly: true,
		Path: "/",
		Domain: "localhost",
	}},MaxAge: t}
}

func (s *Session) Keep()  {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.MaxAge=sessionAge()
}

func (s *Session) GetSessionId() string {
	return s.cookie.Value
}

// SetAttribute 添加属性，同属性直接覆盖
func (s *Session) SetAttribute(key string,value interface{})  {
	 s.rw.Lock()
	 defer s.rw.Unlock()
	 if s.Attribute==nil{
	 	s.Attribute=make(map[string]interface{})
	 }
	 s.Attribute[key]=value
}

// GetAttribute 读取属性，没有存储属性返回nil，不存在的key 返回nil
func (s *Session) GetAttribute(key string) interface{}  {
	if s.Attribute==nil{
		return nil
	}
	s.rw.RLock()
	defer s.rw.RUnlock()
	if v,ok:=s.Attribute[key];ok{
		return v
	}
	return nil
}

func sessionAge() time.Time {
	return time.Now().Add(time.Second*60)
}



