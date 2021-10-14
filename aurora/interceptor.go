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
	PreHandle(ctx *Context) bool
	PostHandle(ctx *Context)
	AfterCompletion(ctx *Context)
}

// DefaultInterceptor 实现全局请求处理前后环绕
type DefaultInterceptor struct {
	t time.Time
}

func (de *DefaultInterceptor) PreHandle(ctx *Context) bool {
	de.t = time.Now()
	return true
}

func (de *DefaultInterceptor) PostHandle(ctx *Context) {

}

func (de *DefaultInterceptor) AfterCompletion(ctx *Context) {
	times := time.Now().Sub(de.t)
	re := ctx.Request
	radd := re.RemoteAddr
	if radd[0:5] == "[::1]" {
		ip := "172.0.0.1"
		radd = radd[5:]
		radd = ip + radd
	}
	l := fmt.Sprintf(" %s → %s | %s %s | %s", radd, re.URL.Host, re.Method, re.URL.Path, times)
	ctx.AR.Api <- l
}

// InterceptorData 实现拦截器压栈出栈功能
type InterceptorData struct {
	imp Interceptor
	pre *InterceptorData
}

type InterceptorStack struct {
	top   *InterceptorData
	stack *InterceptorData
}

func (s *InterceptorStack) Push(i Interceptor) {
	if s.stack == nil && s.top == nil {
		s.stack = &InterceptorData{imp: i}
		s.top = s.stack
		return
	}
	t := &InterceptorData{imp: i, pre: s.top}
	//更新栈顶
	s.top = t
}

// Pull 栈为空时 返回为nil
func (s *InterceptorStack) Pull() Interceptor {
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
