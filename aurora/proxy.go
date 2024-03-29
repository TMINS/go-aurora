package aurora

import (
	"fmt"
	"github.com/awensir/go-aurora/aurora/is"
	"net/http"
	"os"
	"strings"
	"sync"
)

// HandlerProxy 代理 路由处理，负责生成上下文变量和调用具体处理函数
type proxy struct {
	rw           sync.RWMutex
	rew          http.ResponseWriter
	req          *http.Request
	ServeHandler                        //处理函数
	args         map[string]interface{} //REST API 参数解析
	ctx          *Ctx                   //上下文
	result       interface{}            //业务结果
	view         Views                  //支持自定义视图渲染机制
	ar           *Aurora

	index   int          //全局插件索引
	plugins []PluginFunc //全局插件

	Interceptor  bool //是否放行拦截器
	AInterceptor bool

	ExecuteStack, AfterStack *interceptorStack // ExecuteStack,AfterStack 全局拦截器

	TreeInter                                  []Interceptor //通配拦截器集合
	TreeExecuteInterStack, TreeAfterInterStack *interceptorStack

	InterceptorList        []Interceptor     //局部拦截器
	ExecutePart, AfterPart *interceptorStack //ExecutePart,AfterPart
}

// Start 路由查询入口
func (sp *proxy) start() {
	//初始化 ctx
	sp.initCtx()
	sp.AInterceptor = is.Yes //初始化放行所有全局拦截器
	sp.Interceptor = is.Yes  //初始化放行所有路径拦截器
	defer func(sp *proxy) {
		//用于捕捉 plugin或上一级拦截器 执行期间的 panic
		if i := recover(); i != nil {
			switch i.(type) {
			case string:

				sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(string)))
			case error:
				sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(error).Error()))
			}
			return
		}
		//用于捕捉 拦截器发生 的panic
		defer func(sp *proxy) {
			if i := recover(); i != nil {
				switch i.(type) {
				case string:
					sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(string)))
				case error:
					sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(error).Error()))
				}
				return
			}
		}(sp)

		//全局拦截器的AfterCompletion 修改在其他拦截器返回 拦截的情况下 全局拦截器无法 完全执行的bug,此处的实现被提升到上一层函数中 start()中处理全局拦截器的执行,
		//业务处理之后的调用位置暂时不改东,这里使用全局拦截器就会出现,如果路径拦截器阻断会 导致全局拦截器的PostHandle Ctx 上下文中可能没有你想要处理的数据,进而AfterCompletion 中的逻辑可能出错
		//使用全局拦截器尽可能避免直接接触业务逻辑.
		if len(sp.ar.interceptorList) > 0 {
			for {
				if f := sp.AfterStack.Pull(); f != nil {
					f.AfterCompletion(sp.ctx)
				} else {
					break
				}
			}
		}
	}(sp)

	//全局拦截器 运行
	if len(sp.ar.interceptorList) > 0 {
		for _, v := range sp.ar.interceptorList { //依次执行注册过的 拦截器
			if sp.AInterceptor = v.PreHandle(sp.ctx); !sp.AInterceptor { //如果返回false 则终止
				//清空拦截器栈，释放资源
				break //拦截器不放行,后续拦截器也不再执行
			}
			//入栈，为下面的执行周期做出栈的准备
			if sp.ExecuteStack == nil && sp.AfterStack == nil {
				sp.ExecuteStack = &interceptorStack{}
				sp.AfterStack = &interceptorStack{}
			}
			sp.ExecuteStack.Push(v)
			sp.AfterStack.Push(v)
		}
	}

	//执行插件，插件优先级 高于路径拦截器，低于全局拦截器
	for _, p := range sp.plugins {
		if b := p(sp.ctx); !b {
			//插件执行 返回false 则中断 该请求的后续执行 退出 本次请求处理，在退出服务处理之前 应该在插件层面 对 应用用户做出一个发送消息的动作 以表示在插件处理过程中出现问题

			//对中断插件进行消息处理
			message := sp.ctx.GetMessage(plugin)
			if message == nil {
				//如果没有拿到 Error，需要给出一个错误提示，这个消息提取不到会影响到整个框架运行逻辑因此会 painc 或者结束服务器程序
				sp.ar.auroraLog.Error("Plugin Error Message not find")
				os.Exit(1)
			}
			//正确拿到消息后，对客户端进行响应并且发出 500 错误
			http.Error(sp.rew, message.(string), 500)
			//通过goto跳转掉下面的执行
			goto PluginsEnd
		}
	}

	if sp.AInterceptor { //判断全局 拦截器是否放行 ，如果plugin处发生了，panic 后续业务将无法执行下去
		sp.before()
		if sp.Interceptor { //拦截器不放行的情况下是不走业务
			sp.execute()
			sp.after()
		}
	}
PluginsEnd: //结束 插调用链的执行，此处不走结果处理器
}

// Before 服务处理之前
func (sp *proxy) before() {

	defer func(sp *proxy) {

		//用于捕捉 拦截器发生 的panic，此处的拦截器 发生panic后 before 阶段会结束，panic点之后的 拦截器栈不会初始化，可能导致 后续空指针，这里需要改变放行机制
		defer func(sp *proxy) {
			if i := recover(); i != nil {
				sp.Interceptor = false
				switch i.(type) {
				case string:
					sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(string)))
				case error:
					sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(error).Error()))
				}
				return
			}
		}(sp)

		//通配拦截器链
		if sp.TreeInter != nil && len(sp.TreeInter) > 0 && sp.Interceptor { //通配拦截器
			for _, v := range sp.TreeInter {
				if sp.Interceptor = v.PreHandle(sp.ctx); !sp.Interceptor {
					break
				}
				if sp.TreeExecuteInterStack == nil && sp.TreeAfterInterStack == nil {
					sp.TreeExecuteInterStack = &interceptorStack{}
					sp.TreeAfterInterStack = &interceptorStack{}
				}
				sp.TreeExecuteInterStack.Push(v)
				sp.TreeAfterInterStack.Push(v)
			}
		}

		//路径拦截器
		if sp.InterceptorList != nil && len(sp.InterceptorList) > 0 && sp.Interceptor { //局部拦截器
			for _, v := range sp.InterceptorList {
				if sp.Interceptor = v.PreHandle(sp.ctx); !sp.Interceptor {
					break
				}
				if sp.ExecutePart == nil && sp.AfterPart == nil {
					sp.ExecutePart = &interceptorStack{}
					sp.AfterPart = &interceptorStack{}
				}
				sp.ExecutePart.Push(v)
				sp.AfterPart.Push(v)
			}
		}
	}(sp)

}

// Execute 执行业务
func (sp *proxy) execute() {
	sp.result = sp.ServeHandler.Controller(sp.ctx) // 此处的panic 已在执行阶段处理，如果发生panic 被捕捉，处理函数一般直接返回为 nil，后续结果处理的部分也是 按照nil进行处理
}

// After 服务处理之后，主要处理业务结果
func (sp *proxy) after() {

	defer func(sp *proxy) {
		//用于捕捉 外部拦截器发生 的panic
		if i := recover(); i != nil {
			switch i.(type) {
			case string:
				sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(string)))
			case error:
				sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(error).Error()))
			}
			return
		}

		//用于捕捉 内部 拦截器发生 的panic
		defer func() {
			if i := recover(); i != nil {
				switch i.(type) {
				case string:
					sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(string)))
				case error:
					sp.ar.auroraLog.Error(fmt.Sprintf("panic: %s", i.(error).Error()))
				}
				return
			}
		}()

		//通配
		if sp.TreeInter != nil && sp.TreeAfterInterStack != nil {
			for {
				if f := sp.TreeAfterInterStack.Pull(); f != nil {
					f.AfterCompletion(sp.ctx)
				} else {
					break
				}
			}
		}
		//局部
		if sp.InterceptorList != nil && sp.AfterPart != nil {
			for {
				if f := sp.AfterPart.Pull(); f != nil {
					f.AfterCompletion(sp.ctx)
				} else {
					break
				}
			}
		}
	}(sp)

	if len(sp.ar.interceptorList) > 0 { //全局拦截器,此处需要经过业务,可以不用更改调用位置
		for {
			if f := sp.ExecuteStack.Pull(); f != nil {
				f.PostHandle(sp.ctx)
			} else {
				break
			}
		}
	}

	if sp.TreeInter != nil && sp.TreeExecuteInterStack != nil { //通配
		for {
			if f := sp.TreeExecuteInterStack.Pull(); f != nil {
				f.PostHandle(sp.ctx)
			} else {
				break
			}
		}
	}

	if sp.InterceptorList != nil && sp.ExecutePart != nil { //局部
		for {
			if f := sp.ExecutePart.Pull(); f != nil {
				f.PostHandle(sp.ctx)
			} else {
				break
			}
		}
	}
	// 调用结果处理

	sp.resultHandler()
}

func (sp *proxy) resultHandler() {
	switch sp.result.(type) {
	case string:
		path := sp.result.(string)
		//处理普通页面响应
		if strings.HasSuffix(path, ".html") {
			if path[:1] == "/" {
				path = path[1:]
			}
			path = sp.ar.projectRoot + sp.ar.resource + path //拼接文件服务器本地位置目录
			sp.view.View(sp.ctx, path)                       //视图解析 响应 html 页面
			return
		}
		//处理重定向，重定向本质重新走一边路由，找到对应处理的方法
		if strings.HasPrefix(path, "forward:") {
			path = path[8:]
			sp.ctx.forward(path)
			return
		}
		//处理字符串输出
		sp.ctx.json(sp.result)
	case WebError:
		//处理自定义错误处理器
		a := sp.result.(WebError)
		sp.result = a.ErrorHandler(sp.ctx) //更新递归变量
		sp.resultHandler()                 //递归处理错误输出
		return
	case error:
		//直接返回错误处理,让调用者根据错误进行处理
		sp.ctx.SetStatus(500)
		sp.ctx.Response.Header().Set(contentType, sp.ctx.ar.resourceMapType[".json"])
		sp.ctx.json("error:" + sp.result.(error).Error())
		return
	case nil:
		//对结果不做出处理
		return
	case bool:
		b := sp.result.(bool)
		if b {

		} else {

		}
		return
	default:
		//其它类型直接编码json发送
		sp.ctx.json(sp.result)
	}
}

// Init 初始化 Context变量
func (sp *proxy) initCtx() {
	if sp.ctx == nil {
		sp.ctx = &Ctx{}
		sp.ctx.Attribute = &sync.Map{}
		sp.ctx.Request = sp.req
		sp.ctx.Response = sp.rew
		sp.ctx.rw = &sync.RWMutex{}
		sp.ctx.ar = sp.ar
		sp.ctx.Args = sp.args
	}
}

func (sp *proxy) bad() {

}

func (sp *proxy) ok() {

}
