package aurora

import (
	"Aurora/logs"
	"net/http"
	"strings"
	"sync"
)

// ServletProxy 代理 路由处理，负责生成上下文变量和调用具体处理函数
type ServletProxy struct {
	rw sync.RWMutex
	rew http.ResponseWriter
	req *http.Request
	ServletHandler                  //处理函数
	Query string                    //url 查询参数
	Args map[string]string          //REST API 参数
	ctx *Context                    //上下文
	result interface{}              //业务结果
	Interceptor bool
	ExecuteStack,AfterStack,ExecutePart,AfterPart    *InterceptorStack    // ExecuteStack,AfterStack 全局拦截器, ExecutePart,AfterPart局部
}

// Start 路由查询入口
func (sp *ServletProxy) Start() {
	sp.Init()
	sp.Before()
	if sp.Interceptor{
	    sp.Execute()
		sp.After()
	}
}

// Before 服务处理之前
func (sp *ServletProxy) Before() {
	sp.ctx.GetSession()
	sp.Interceptor=true
	defer func(ctx *Context,sp *ServletProxy) {   //全局拦截器
		if len(InterceptorList)>0{
			for _,v:=range InterceptorList{     //依次执行注册过的 拦截器
				if sp.Interceptor=v.PreHandle(ctx);!sp.Interceptor {    //如果返回false 则终止
					//清空拦截器栈
					break   //拦截器不放行,后续拦截器也不再执行
				}
				if sp.ExecuteStack==nil && sp.AfterStack==nil{
					sp.ExecuteStack=&InterceptorStack{}
					sp.AfterStack=&InterceptorStack{}
				}
				sp.ExecuteStack.Push(v)
				sp.AfterStack.Push(v)
			}
		}
	}(sp.ctx,sp)
}

// Execute 执行业务
func (sp *ServletProxy) Execute()  {
	
	defer func(ctx *Context,sp *ServletProxy) {  //全局拦截器
		if len(InterceptorList)>0{
			for {
				if f:=sp.ExecuteStack.Pull();f!=nil{
					f.PostHandle(ctx)
				}else {
					break
				}
			}
		}
	}(sp.ctx,sp)
	sp.result=sp.ServletHandler.ServletHandler(sp.ctx)
}

// After 服务处理之后，主要处理业务结果
func (sp *ServletProxy) After() {
	
	defer func(ctx *Context,sp *ServletProxy){   //全局拦截器
		if len(InterceptorList)>0{
			for {
				if f:=sp.AfterStack.Pull();f!=nil{
					f.AfterCompletion(ctx)
				}else {
					break
				}
			}
		}
	}(sp.ctx,sp)
	
	switch sp.result.(type) {
		case string:
			if strings.HasSuffix(sp.result.(string),".html") {
				SendResource(sp.rew,readResource(sp.result.(string)))  //直接响应 html 页面
			}else {
				sp.ctx.JSON(sp.result)
			}
		case error:
			logs.WebRequestError(sp.result.(error).Error())
		default:
			sp.ctx.JSON(sp.result)
	}
}

// Init 初始化 Context变量
func (sp *ServletProxy) Init() {
	if sp.ctx==nil{
		sp.ctx=&Context{}
		sp.ctx.Request=sp.req
		sp.ctx.ResponseWriter=sp.rew
		if sp.Args!=nil{
			sp.ctx.Args=sp.Args
		}
	}
}

func (sp *ServletProxy) ResultHandler()  {
	switch sp.result.(type) {
	case string:
		if strings.HasSuffix(sp.result.(string),".html") {
			SendResource(sp.rew,readResource(sp.result.(string)))  //直接响应 html 页面
		}else {
			sp.ctx.JSON(sp.result)
		}
	case error:
		logs.WebRequestError(sp.result.(error).Error())
	default:
		sp.ctx.JSON(sp.result)
	}
}
