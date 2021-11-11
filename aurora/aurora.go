package aurora

import (
	"context"
	"fmt"
	"github.com/awensir/go-aurora/logs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

func init() {
	s := "    /\\\n   /  \\  _   _ _ __ ___  _ __ __ _\n  / /\\ \\| | | | '__/ _ \\| '__/ _` |\n / ____ \\ |_| | | | (_) | | | (_| |\n/_/    \\_\\__,_|_|  \\___/|_|  \\__,_|\n:: aurora ::   (v0.1.2.RELEASE)"
	/*
	       /\
	      /  \  _   _ _ __ ___  _ __ __ _
	     / /\ \| | | | '__/ _ \| '__/ _` |
	    / ____ \ |_| | | | (_) | | | (_| |
	   /_/    \_\__,_|_|  \___/|_|  \__,_|
	   :: aurora ::   (v0.0.1.RELEASE)

	*/
	fmt.Println(s)
}

type Aurora struct {
	rw               *sync.RWMutex
	ctx              context.Context     //服务器顶级上下文，通过此上下文可以跳过 go web 自带的子上下文去开启纯净的子go程，结束此上下文 web服务也将结束
	cancel           func()              //取消上下文
	port             string              //服务端口号
	router           *route              //路由服务管理
	projectRoot      string              //项目根路径
	resource         string              //静态资源管理 默认为 root 目录
	resourceMappings map[string][]string //静态资源映射路径标识
	resourceMapType  map[string]string   //常用的静态资源头
	load             chan struct{}
	message          chan string        //启动自带的日志信息
	initError        chan error         //路由器级别错误通道 一旦初始化出错，则结束服务，检查配置
	runtime          chan *localMonitor //单体服务运行时错误时候的链路调用日志
	serviceInfo      chan string        //业务 info日志
	serviceWarning   chan string        //业务 警告日志
	serviceError     chan string        //业务 错误日志
	servicePanic     chan string        //业务 panic日志
	routeInterceptor []interceptorArgs  //拦截器初始华切片
	interceptorList  []Interceptor      //全局拦截器
	container        *containers        //第三方配置整合容器,原型模式
	log              *logrus.Logger     // Aurora 实例日志变量
	serviceLog       *logrus.Logger     // 业务实例日志
	cnf              *viper.Viper       // 配置实例
	Server           *http.Server       // web服务器
	Ln               net.Listener       // web服务器监听
}

// New :最基础的 Aurora 实例
func New() *Aurora {
	a := &Aurora{
		rw:   &sync.RWMutex{},
		port: "8080", //默认端口号
		router: &route{
			mx: &sync.Mutex{},
		},
		Server:          &http.Server{},
		resource:        "static", //设定资源默认存储路径
		initError:       make(chan error),
		resourceMapType: make(map[string]string),
		load:            make(chan struct{}),
		message:         make(chan string),
		runtime:         make(chan *localMonitor),
		serviceInfo:     make(chan string),
		serviceWarning:  make(chan string),
		servicePanic:    make(chan string),
		serviceError:    make(chan string),
		container: &containers{
			rw:         &sync.RWMutex{},
			prototypes: make(map[string]interface{}),
		},
		log:        logs.NewLog(),
		serviceLog: logs.NewServiceLog(),
	}
	startLoading(a)
	loadResourceHead(a)
	projectRoot, _ := os.Getwd()
	a.projectRoot = projectRoot
	a.router.defaultView = a //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
	a.router.AR = a
	a.interceptorList = []Interceptor{
		0: &defaultInterceptor{},
	}
	a.message <- fmt.Sprintf("Golang Version :%1s", runtime.Version())
	a.message <- fmt.Sprintf("Project Path:%1s", a.projectRoot)
	a.message <- fmt.Sprintf("Default Server Port :%1s", a.port)

	//加载 cnf 配置实例
	//a.ViperConfig()
	a.message <- fmt.Sprintf("Default Static Resource Path:%1s", a.resource)
	a.Server.BaseContext = a.baseContext //配置 上下文对象属性
	return a
}

// Guide 启动 Aurora 服务器，默认端口号8080
func (a *Aurora) Guide(port ...string) {
	a.run(port...)
}

func (a *Aurora) run(port ...string) {
	if port != nil && len(port) > 1 {
		panic("too mach port")
	}
	if port == nil {
		a.Server.Addr = ":" + a.port
	} else {
		p := port[0]
		if p[0:1] != ":" {
			p = ":" + p
		}
		a.Server.Addr = p
		a.port = p
	}
	a.Server.Handler = a
	err := a.Server.ListenAndServe() //启动服务器
	if err != nil {
		a.initError <- err
	}
}

// ResourceMapping 资源映射
//添加静态资源配置，t资源类型必须以置源后缀命名，
//paths为t类型资源的子路径，可以一次性设置多个。
//每个资源类型最调用一次设置方法否则覆盖原有设置
func (a *Aurora) ResourceMapping(Type string, Paths ...string) {
	a.registerResourceType(Type, Paths...)
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
	a.resource = root
}

// RouteIntercept path路径上添加一个或者多个路由拦截器
func (a *Aurora) RouteIntercept(path string, interceptor ...Interceptor) {
	if a.routeInterceptor == nil {
		a.routeInterceptor = make([]interceptorArgs, 0)
	}
	r := interceptorArgs{path: path, list: interceptor}
	a.routeInterceptor = append(a.routeInterceptor, r)
	//a.router.RegisterInterceptor(path, &LocalMonitor{mx: &sync.Mutex{}}, interceptor...)
}

// DefaultInterceptor 配置默认顶级拦截器
func (a *Aurora) DefaultInterceptor(interceptor Interceptor) {
	a.interceptorList[0] = interceptor
	l := fmt.Sprintf("Web Default Global Rout Interceptor successds")
	a.message <- l
}

// AddInterceptor 追加全局拦截器
func (a *Aurora) AddInterceptor(interceptor ...Interceptor) {
	//追加全局拦截器
	for _, v := range interceptor {
		a.interceptorList = append(a.interceptorList, v)
		l := fmt.Sprintf("Web Global Rout Interceptor successds")
		a.message <- l
	}
}

// GET 请求
func (a *Aurora) GET(path string, servlet Servlet) {
	a.register(http.MethodGet, path, servlet)
}

// POST 请求
func (a *Aurora) POST(path string, servlet Servlet) {
	a.register(http.MethodPost, path, servlet)
}

// PUT 请求
func (a *Aurora) PUT(path string, servlet Servlet) {
	a.register(http.MethodPut, path, servlet)
}

// DELETE 请求
func (a *Aurora) DELETE(path string, servlet Servlet) {
	a.register(http.MethodDelete, path, servlet)
}

// HEAD 请求
func (a *Aurora) HEAD(path string, servlet Servlet) {
	a.register(http.MethodHead, path, servlet)
}

// register 通用注册器
func (a *Aurora) register(method string, mapping string, fun Servlet) {
	list := &localMonitor{mx: &sync.Mutex{}}
	list.En(executeInfo(nil))
	a.router.addRoute(method, mapping, fun, list)
}

// Group 路由分组  必须以 “/” 开头分组
func (a *Aurora) Group(path string) *group {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return &group{
		prefix: path,
		a:      a,
	}
}

func (a *Aurora) ViewHandle(views Views) {
	a.router.defaultView = views
}

// View 默认视图解析
func (a *Aurora) View(ctx *Ctx, html string) {
	parseFiles, err := template.ParseFiles(a.projectRoot + "/" + a.resource + html)
	if err != nil {
		log.Fatal("ParseFiles" + err.Error())
		return
	}
	err = parseFiles.Execute(ctx.Response, ctx.Attribute)
	if err != nil {
		log.Fatal("Execute" + err.Error())
		return
	}
}

func (a *Aurora) loadingInterceptor() {
	if a.routeInterceptor != nil {
		for i := 0; i < len(a.routeInterceptor); i++ {
			e := a.routeInterceptor[i]
			a.router.RegisterInterceptor(e.path, &localMonitor{mx: &sync.Mutex{}}, e.list...)
		}
	}
}

func (a *Aurora) baseContext(ln net.Listener) context.Context {
	//初始化 Aurora net.Listener 变量，用于整合grpc
	a.Ln = ln
	a.loadingInterceptor() //加载 拦截器
	l := fmt.Sprintf("The server successfully runs on port %s", a.port)
	c, f := context.WithCancel(context.TODO())
	a.ctx = c
	a.cancel = f
	a.message <- l
	return c
}

// Get 获取加载
func (a *Aurora) Get(name string) interface{} {
	return a.container.get(name)
}

// Store 加载
func (a *Aurora) Store(name string, variable interface{}) {
	a.container.store(name, variable)
}

// startLoading 启动加载
func startLoading(a *Aurora) {
	//启动日志
	go func(a *Aurora) {
		for true {
			select {

			case info := <-a.message:
				a.log.Info(info)

			case msg := <-a.runtime:
				a.log.Error(msg.Message())

			case info := <-a.serviceInfo:
				a.serviceLog.Info(info)

			case info := <-a.serviceWarning:
				a.serviceLog.Warning(info)

			case info := <-a.serviceError:
				a.serviceLog.Error(info)

			case info := <-a.servicePanic:
				a.serviceLog.Error(info)
				os.Exit(-2) //结束程序

			case err := <-a.initError: //启动初始化错误处理
				a.log.Error(err.Error())
				os.Exit(-1) //结束程序
			}
		}
	}(a)
}
