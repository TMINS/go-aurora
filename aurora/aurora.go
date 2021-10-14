package aurora

import (
	"context"
	"fmt"
	"github.com/awensir/Aurora/logs"
	"net/http"
	"os"
	"strings"
	"sync"
)

var log = logs.NewLog()
var rlog = logs.NewRouteLog()

type CtxListenerKey string

type Aurora struct {
	rw               sync.RWMutex
	Port             string        //服务端口号
	Router           *ServerRouter //路由服务管理
	server           *http.Server
	Resource         string              //静态资源管理 默认为 root 目录
	ResourceMappings map[string][]string //静态资源映射路径标识
	InitError        chan error          //路由器级别错误通道 一旦初始化出错，则结束服务，检查配置
	StartInfo        chan string         //输出启动信息
	Ctx              context.Context     //服务器顶级上下文，通过此上下文可以跳过web 上下文去开启纯净的子go程
	Cancel           func()
	ProjectRoot      string        //项目根路径
	InterceptorList  []Interceptor //全局拦截器
	ResourceMapType  map[string]string
	Api              chan string
}

func Default() *Aurora {
	a := New()
	if a == nil {
		return nil
	}
	return a
}

// New :最基础的 Aurora 实例
func New() *Aurora {
	a := &Aurora{
		Port:            "8080", //默认端口号
		Router:          &ServerRouter{},
		server:          &http.Server{},
		Resource:        "static", //设定资源默认存储路径
		InitError:       make(chan error),
		StartInfo:       make(chan string),
		ResourceMapType: make(map[string]string),
		Api:             make(chan string),
	}
	projectRoot, _ := os.Getwd()
	a.ProjectRoot = projectRoot
	a.Router.View = a.DefaultView //使用默认视图解析
	a.Router.AR = a
	a.InterceptorList = []Interceptor{
		0: &DefaultInterceptor{},
	}
	LoadResourceHead(a)
	startLoading(a)
	return a
}

// Guide 启动 Aurora 服务器
func (a *Aurora) Guide(port ...string) {
	a.run(port...)
}

func (a *Aurora) run(port ...string) {
	if port != nil && len(port) > 1 {
		panic("too mach port")
	}
	if port == nil {
		a.server.Addr = ":" + a.Port
	} else {
		p := port[0]
		if p[0:1] != ":" {
			p = ":" + p
		}
		a.server.Addr = p
		a.Port = p
	}
	a.server.Handler = a
	err := a.server.ListenAndServe() //启动服务器
	if err != nil {
		a.InitError <- err
		return
	}
}

// ResourceMapping 资源映射
//添加静态资源配置，t资源类型必须以置源后缀命名，
//paths为t类型资源的子路径，可以一次性设置多个。
//每个资源类型最调用一次设置方法否则覆盖原有设置
func (a *Aurora) ResourceMapping(Type string, Paths ...string) {
	a.RegisterResourceType(Type, Paths...)
}

// StaticRoot 设置静态资源根路径
func (a *Aurora) StaticRoot(root string) {
	if root == "" {
		panic(" static resource paths cannot be empty! ")
	}
	if strings.HasPrefix(root, "/") {
		root = root[1:]
	}
	if strings.HasSuffix(root, "/") {
		root = root[:len(root)-1]
	}
	a.Resource = root
}

// RouteIntercept path路径上添加一个或者多个路由拦截器
func (a *Aurora) RouteIntercept(path string, interceptor ...Interceptor) {
	a.Router.RegisterInterceptor(path, interceptor...)
}

// DefaultInterceptor 配置默认顶级拦截器
func (a *Aurora) DefaultInterceptor(interceptor Interceptor) {
	a.InterceptorList[0] = interceptor
	l := fmt.Sprintf("Web Default Global Rout Interceptor successds")
	a.StartInfo <- l
}

// AddInterceptor 追加全局拦截器
func (a *Aurora) AddInterceptor(interceptor ...Interceptor) {
	//追加全局拦截器
	for _, v := range interceptor {
		a.InterceptorList = append(a.InterceptorList, v)
		l := fmt.Sprintf("Web Global Rout Interceptor successds")
		a.StartInfo <- l
	}
}

// GET 请求
func (a *Aurora) GET(path string, servlet Servlet) {
	a.Register(http.MethodGet, path, servlet)
}

// POST 请求
func (a *Aurora) POST(path string, servlet Servlet) {
	a.Register(http.MethodPost, path, servlet)
}

// PUT 请求
func (a *Aurora) PUT(path string, servlet Servlet) {
	a.Register(http.MethodPut, path, servlet)
}

// DELETE 请求
func (a *Aurora) DELETE(path string, servlet Servlet) {
	a.Register(http.MethodDelete, path, servlet)
}

// HEAD 请求
func (a *Aurora) HEAD(path string, servlet Servlet) {
	a.Register(http.MethodHead, path, servlet)
}

// Group 路由分组  必须以 “/” 开头分组
func (a *Aurora) Group(path string) *Group {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return &Group{
		prefix: path,
		a:      a,
	}
}

// startLoading 启动加载
func startLoading(a *Aurora) {
	//启动日志
	go func(a *Aurora) {

		s := "    /\\\n   /  \\  _   _ _ __ ___  _ __ __ _\n  / /\\ \\| | | | '__/ _ \\| '__/ _` |\n / ____ \\ |_| | | | (_) | | | (_| |\n/_/    \\_\\__,_|_|  \\___/|_|  \\__,_|\n:: Aurora ::   (v0.0.8.RELEASE)"
		/*
		       /\
		      /  \  _   _ _ __ ___  _ __ __ _
		     / /\ \| | | | '__/ _ \| '__/ _` |
		    / ____ \ |_| | | | (_) | | | (_| |
		   /_/    \_\__,_|_|  \___/|_|  \__,_|
		   :: Aurora ::   (v0.0.1.RELEASE)

		*/
		fmt.Println(s)
		for true {
			select {
			case msg := <-a.StartInfo: //
				log.Info(msg)
			case api := <-a.Api:
				rlog.Info(api)
			case err := <-a.InitError: //启动初始化错误处理
				log.Error(err.Error())
				os.Exit(-1) //结束程序
			}
		}
		fmt.Println(11111)
	}(a)
}
