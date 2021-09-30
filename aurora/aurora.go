package aurora

import (
	"context"
	"fmt"
	"github.com/awensir/Aurora/message"
	"github.com/awensir/Aurora/uuid"
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
*/
func init() {
	projectRoot, err := os.Getwd()
	if err != nil {
		return
	}
	aurora.ProjectRoot = projectRoot
	aurora.sessionMap = make(map[string]*Session)
	aurora.interceptorList = []Interceptor{
		0: DefaultInterceptor{},
	}
	aurora.SessionCreate = uuid.NewWorker(1, 1) //sessionId生成器
	startLoading()                              //开启路由加载监听
}

var ctx, cancel = context.WithCancel(context.TODO())

type CtxListenerKey string

var sessionIdCreater = uuid.NewWorker(1, 1) //sessionId生成器
// 全局路由器
var aurora = &Aurora{
	Port:            "8080",
	Router:          ServerRouter{},
	Resource:        "static", //设定资源默认存储路径
	Ctx:             ctx,
	Cancel:          cancel,
	InitError:       make(chan error),
	StartInfo:       make(chan message.Message),
	sort:            Sort{First: make(chan struct{}), Second: make(chan struct{}), Finally: make(chan struct{})},
	resourceMapType: make(map[string]string),
	vw:              DefaultView,
}

type Aurora struct {
	rw              sync.RWMutex
	Port            string               //服务端口号
	Router          ServerRouter         //路由服务管理
	Resource        string               //静态资源管理 默认为 root 目录
	resourceMapping map[string][]string  //静态资源映射路径标识
	InitError       chan error           //路由器级别错误通道 一旦初始化出错，则结束服务，检查配置
	StartInfo       chan message.Message //输出启动信息
	Ctx             context.Context      //服务器顶级上下文，通过此上下文可以跳过web 上下文去开启纯净的子go程
	Cancel          func()
	ProjectRoot     string              //项目根路径
	interceptorList []Interceptor       //全局拦截器
	sessionMap      map[string]*Session //全局session管理
	SessionCreate   *uuid.Worker        //session id 生成器
	resourceMapType map[string]string
	sort            Sort
	vw              ViewFunc //支持自定义视图渲染机制
}

type Sort struct {
	First   chan struct{}
	Second  chan struct{}
	Finally chan struct{}
}

// RunApplication 启动服务器
func RunApplication(port string) {
	if port[0:1] != ":" {
		port = ":" + port
	}
	server := &http.Server{
		Addr:        port,
		Handler:     aurora,
		BaseContext: CreateConText,
	}
	aurora.Port = port
	aurora.Router.OptimizeTree()   //路由树节点排序
	err := server.ListenAndServe() //启动服务器
	if err != nil {
		aurora.InitError <- err
		return
	}
}

func RegisterInterceptorList(interceptor ...Interceptor) {
	//追加全局拦截器
	for _, v := range interceptor {
		aurora.interceptorList = append(aurora.interceptorList, v)
	}
}

// RegisterDefaultInterceptor 提供修改默认顶级拦截器
func RegisterDefaultInterceptor(interceptor Interceptor) {
	aurora.interceptorList[0] = interceptor
}

// CreateConText 提供web自定义父级上下文
func CreateConText(listener net.Listener) context.Context {
	key := CtxListenerKey("Listener")
	p := context.TODO()
	vCtx := context.WithValue(p, key, listener)
	aurora.Ctx, aurora.Cancel = context.WithCancel(vCtx) //重新封装上下文，把连接对象保存在上下文中，在次之前使用aurora.Ctx 将可能无法释放资源

	return aurora.Ctx
}

// SetResourceRoot 设置静态资源根路径
func SetResourceRoot(root string) {
	if root == "" { //不允许设置""
		return
	}
	rl := len(root)
	if root[:1] == "/" {
		root = root[1:]
	}
	if root[rl-1:] == "/" {
		root = root[:rl-1]
	}
	aurora.Resource = root
}

// startLoading 启动加载
func startLoading() {
	//启动日志
	go func(aurora *Aurora) {
		/*
				/\
			   /  \  _   _ _ __ ___  _ __ __ _
			  / /\ \| | | | '__/ _ \| '__/ _` |
			 / ____ \ |_| | | | (_) | | | (_| |
			/_/    \_\__,_|_|  \___/|_|  \__,_|
			  :: Aurora ::   (v0.0.1.RELEASE)
		*/

		for true {
			select {
			case <-aurora.StartInfo: //启动日志，暂时不做处理
			case err := <-aurora.InitError: //启动初始化错误处理
				fmt.Println(err.Error())
				os.Exit(-1) //结束程序
			case <-aurora.sort.First:

			case <-aurora.sort.Second:

			case <-aurora.sort.Finally:

			}
		}
	}(aurora)

	//session 生命周期检查，定时任务
	go func(aurora *Aurora) {
		Ticker := time.NewTicker(time.Second * 65) //每隔 65秒执行一次 session 清理，存在占用资源bug，在没有任何session情况下会做无用的定时任务（待解决）
		defer Ticker.Stop()
		for true {
			select {
			case t := <-Ticker.C:
				if len(aurora.sessionMap) > 0 {
					aurora.rw.Lock()
					for k, _ := range aurora.sessionMap {
						s := aurora.sessionMap[k]
						if t.After(s.MaxAge) {
							//session过期 删除
							delete(aurora.sessionMap, k)
						}
					}
					aurora.rw.Unlock()
				}
			case <-aurora.InitError: //初始化错误 结束线程
				return
			}
		}
	}(aurora)
}
