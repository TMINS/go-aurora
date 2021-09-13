package aurora

import (
	"Aurora/logs"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

// ServletProxy 代理 路由处理，负责生成上下文变量和调用具体处理函数
type ServletProxy struct {
	rw                       sync.RWMutex
	rew                      http.ResponseWriter
	req                      *http.Request
	ServletHandler                             //处理函数
	Args                     map[string]string //REST API 参数解析
	ctx                      *Context          //上下文
	result                   interface{}       //业务结果
	Interceptor              bool

	ExecuteStack, AfterStack *InterceptorStack // ExecuteStack,AfterStack 全局拦截器

	TreeInter 				 []Interceptor	   //通配拦截器集合
	TreeExecuteInterStack,TreeAfterInterStack           *InterceptorStack

	InterceptorList          []Interceptor
	ExecutePart, AfterPart   *InterceptorStack //ExecutePart,AfterPart局部
}

// Start 路由查询入口
func (sp *ServletProxy) Start() {
	sp.Init()
	sp.Before()
	if sp.Interceptor {
		sp.Execute()
		sp.After()
	}
}

// Before 服务处理之前
func (sp *ServletProxy) Before() {
	sp.Interceptor = true //初始化放行所有拦截器
	defer func(ctx *Context, sp *ServletProxy) {
		//处理全局拦截器和局部拦截器之前，临时构造一个拦截器执行序列

		if len(InterceptorList) > 0 { //全局拦截器
			for _, v := range InterceptorList { //依次执行注册过的 拦截器
				if sp.Interceptor = v.PreHandle(ctx); !sp.Interceptor { //如果返回false 则终止
					//清空拦截器栈，释放资源
					break //拦截器不放行,后续拦截器也不再执行
				}
				if sp.ExecuteStack == nil && sp.AfterStack == nil {
					sp.ExecuteStack = &InterceptorStack{}
					sp.AfterStack = &InterceptorStack{}
				}
				sp.ExecuteStack.Push(v)
				sp.AfterStack.Push(v)
			}
		}

		if sp.TreeInter != nil && len(sp.TreeInter) > 0 && sp.Interceptor { //通配拦截器
			for _, v := range sp.TreeInter {
				if sp.Interceptor = v.PreHandle(ctx); !sp.Interceptor {
					break
				}
				if sp.TreeExecuteInterStack == nil && sp.TreeAfterInterStack == nil {
					sp.TreeExecuteInterStack = &InterceptorStack{}
					sp.TreeAfterInterStack = &InterceptorStack{}
				}
				sp.TreeExecuteInterStack.Push(v)
				sp.TreeAfterInterStack.Push(v)
			}
		}

		if sp.InterceptorList != nil && len(sp.InterceptorList) > 0 && sp.Interceptor { //局部拦截器
			for _, v := range sp.InterceptorList {
				if sp.Interceptor = v.PreHandle(ctx); !sp.Interceptor {
					break
				}
				if sp.ExecutePart == nil && sp.AfterPart == nil {
					sp.ExecutePart = &InterceptorStack{}
					sp.AfterPart = &InterceptorStack{}
				}
				sp.ExecutePart.Push(v)
				sp.AfterPart.Push(v)
			}
		}
	}(sp.ctx, sp)
}

// Execute 执行业务
func (sp *ServletProxy) Execute() {

	defer func(ctx *Context, sp *ServletProxy) {
		if len(InterceptorList) > 0 { //全局拦截器
			for {
				if f := sp.ExecuteStack.Pull(); f != nil {
					f.PostHandle(ctx)
				} else {
					break
				}
			}
		}
		if sp.TreeInter != nil { //通配
			for {
				if f := sp.TreeExecuteInterStack.Pull(); f != nil {
					f.PostHandle(ctx)
				} else {
					break
				}
			}
		}
		if sp.InterceptorList != nil { //局部
			for {
				if f := sp.ExecutePart.Pull(); f != nil {
					f.PostHandle(ctx)
				} else {
					break
				}
			}
		}

	}(sp.ctx, sp)
	sp.result = sp.ServletHandler.ServletHandler(sp.ctx)
}

// After 服务处理之后，主要处理业务结果
func (sp *ServletProxy) After() {

	defer func(ctx *Context, sp *ServletProxy) {
		if len(InterceptorList) > 0 { //全局拦截器
			for {
				if f := sp.AfterStack.Pull(); f != nil {
					f.AfterCompletion(ctx)
				} else {
					break
				}
			}
		}
		if sp.TreeInter != nil { //通配
			for {
				if f := sp.TreeAfterInterStack.Pull(); f != nil {
					f.PostHandle(ctx)
				} else {
					break
				}
			}
		}
		if sp.InterceptorList != nil { //局部
			for {
				if f := sp.AfterPart.Pull(); f != nil {
					f.PostHandle(ctx)
				} else {
					break
				}
			}
		}
	}(sp.ctx, sp)
	sp.ResultHandler()
}

// Init 初始化 Context变量
func (sp *ServletProxy) Init() {
	if sp.ctx == nil {
		sp.ctx = &Context{}
		sp.ctx.Request = sp.req
		sp.ctx.Response = sp.rew
		if sp.Args != nil {
			sp.ctx.Args = sp.Args
		}
	}
}

func (sp *ServletProxy) ResultHandler() {
	switch sp.result.(type) {
	case string:
		if strings.HasSuffix(sp.result.(string), ".html") {
			SendResource(sp.rew, readResource(sp.result.(string))) //直接响应 html 页面
		} else {
			sp.ctx.JSON(sp.result)
		}
	case WebResponseError:
		v := reflect.ValueOf(sp.result)
		method := v.MethodByName("ErrorHandler")
		value := method.Call([]reflect.Value{
			reflect.ValueOf(sp.ctx),
		})
		if len(value) != 1 {
			panic("Call return failed")
		}
		r := value[0].Interface()
		sp.result = r      //更新递归变量
		sp.ResultHandler() //递归处理错误输出
	case error:
		logs.WebRequestError(sp.result.(error).Error())
	default:
		sp.ctx.JSON(sp.result)
	}
}
