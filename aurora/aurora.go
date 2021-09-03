package aurora

import (
	"Aurora/logs"
	"Aurora/message"
	"Aurora/uuid"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

/*
	aurora 路由管理器
	1.请求映射管理
		-封装基本元素（请求类型，url，具体处理）(***)
	2.路由规则
		- 路径参数
	3.静态资源处理
		- 响应解析
	4.监听处理
		- 请求监听(***)
*/
func init() {
	logs.LoadWebLog(&logs.WebLogs{logs.Log{Head: "Aurora"}})                    //初始化日志
	startLoading()    //开启路由加载监听
}

//
var  ctx, cancel =context.WithCancel(context.TODO())

type CtxListenerKey string

//全局管理
var sessionMap=make(map[string]*Session)
var InterceptorList = make([]Interceptor, 0)
var sessionIdCreater=uuid.NewWorker(1,1)
// 全局路由器
var aurora = &Aurora{
	Router: ServerRouter{},
	resource: "static",   //设定资源默认存储路径
	Ctx: ctx,
	Cancel: cancel,
	InitError: make(chan error),
	StartInfo: make(chan message.Message),
}
type Aurora struct {
	 rw  sync.RWMutex
	 Port string
	 Router ServerRouter  		            //服务管理
	 resource string				        //静态资源管理 默认为 root 目录
	 resourceMapping map[string][]string	//静态资源映射路径标识
	 InitError   chan error		            //路由器级别错误通道 一旦初始化出错，则结束服务，检查配置
	 StartInfo chan message.Message         //输出启动信息
	 Ctx   context.Context                  //服务器顶级上下文，通过此上下文可以跳过web 上下文去开启纯净的子go程
	 Cancel func()
}

// RunApplication 启动服务器
func RunApplication(port string) {
	if port[0:1]!=":"{
		port=":"+port
	}
	server := &http.Server{
		Addr: port,
		Handler: aurora,
		BaseContext: CreateConText,
	}
	aurora.Port=port
	aurora.Router.OptimizeTree()        //路由树节点排序
	err := server.ListenAndServe()      //启动服务器
	if err != nil {
		aurora.InitError<-err
		return
	}
}

// CreateConText 提供web自定义父级上下文
func CreateConText(listener net.Listener) context.Context{
	key:=CtxListenerKey("Listener")
	p:=context.TODO()
	vCtx:=context.WithValue(p,key,listener)
	aurora.Ctx,aurora.Cancel=context.WithCancel(vCtx)       //重新封装上下文，把连接对象保存在上下文中，在次之前使用aurora.Ctx 将可能无法释放资源
	aurora.StartInfo<-message.StartSuccessful{Port: aurora.Port}
	return aurora.Ctx
}

// SetResourceRoot 设置静态资源根路径
func SetResourceRoot(root string) {
	rl:=len(root)
	if root[:1]=="/"{
		root=root[1:]
	}
	if root[rl-1:]=="/"{
		root=root[:rl-1]
	}
	aurora.resource=root
}


// startLoading 启动加载
func startLoading()  {
	
	//启动日志
	go func(aurora *Aurora) {
		open, err := ioutil.ReadFile("aurora/start.txt")
		if err != nil {
			logs.WebErrorLogger(err.Error())
			return
		}
		fmt.Printf("%s \n\r",string(open))
		for true {
			select {
			case msg:=<-aurora.StartInfo:    //启动日志
				logs.WebLogger(msg.ToString())
			
			case err:=<-aurora.InitError:    //启动初始化错误处理
				logs.WebErrorLogger(err.Error())
				os.Exit(-1)            //结束程序
			}
		}
	}(aurora)
	
	//session 生命周期检查，定时任务
	go func (aurora *Aurora) {
		Ticker:=time.NewTicker(time.Second*5)  //每隔 60秒执行一次 session 清理
		defer Ticker.Stop()
		for true {
			select {
			case t:=<-Ticker.C:
				if len(sessionMap)>0{
					aurora.rw.Lock()
					for k,_:=range sessionMap{
						s:=sessionMap[k]
						//now:=time.Now()
						if t.After(s.MaxAge){
							//session过期 删除
							delete(sessionMap,k)
							logs.Info("销毁session ：",k)
						}
					}
					aurora.rw.Unlock()
				}
			case <-aurora.InitError:  //初始化错误 结束线程
				return
			}
		}
	}(aurora)
}




