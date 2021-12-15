package aurora

import (
	"fmt"
	"time"
)

/*
	拦截器
	如果需要对web服务中的请求进行预处理，只需要实现Interceptor接口并且注册到InterceptorList中即可
	全局或者同一路基的多个拦截器，执行顺序按照注册先后顺序决定
*/

// Interceptor 拦截器统一接口，实现这个接口就可以向服务器注册一个全局或者指定路径的处理拦截
type Interceptor interface {
	PreHandle(c *Ctx) bool
	PostHandle(c *Ctx)
	AfterCompletion(c *Ctx)
}

// DefaultInterceptor 实现全局请求处理前后环绕
type defaultInterceptor struct {
	t time.Time
}

func (de *defaultInterceptor) PreHandle(ctx *Ctx) bool {
	de.t = time.Now()
	return true
}

func (de *defaultInterceptor) PostHandle(ctx *Ctx) {

}

func (de *defaultInterceptor) AfterCompletion(ctx *Ctx) {

	times := time.Now().Sub(de.t)
	re := ctx.Request
	radd := re.RemoteAddr
	if radd[0:5] == "[::1]" {
		ip := "172.0.0.1"
		radd = radd[5:]
		radd = ip + radd
	}
	l := fmt.Sprintf(" %s → %s | %s %s | %s", radd, re.URL.Host, re.Method, re.URL.Path, times)
	ctx.ar.message <- l
}

// InterceptorData 实现拦截器压栈出栈功能
type interceptorData struct {
	imp Interceptor
	pre *interceptorData
}

type interceptorStack struct {
	top   *interceptorData
	stack *interceptorData
}

func (s *interceptorStack) Push(i Interceptor) {
	if s.stack == nil && s.top == nil {
		s.stack = &interceptorData{imp: i}
		s.top = s.stack
		return
	}
	t := &interceptorData{imp: i, pre: s.top}
	//更新栈顶
	s.top = t
}

// Pull 栈为空时 返回为nil
func (s *interceptorStack) Pull() Interceptor {
	if s.stack == nil && s.top == nil {
		return nil
	}
	fun := s.top.imp
	if s.top.pre != nil {
		s.top = s.top.pre
	} else { // 当弹出最后一个元素 时候清空 初始化栈内存
		s.top = nil
		s.stack = nil
	}
	return fun
}

type interceptorArgs struct {
	path string
	list []Interceptor
}
