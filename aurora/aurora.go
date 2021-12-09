package aurora

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"time"

	"html/template"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

/*
	<***> 基于稳定模块，无需更改
	<---> 非稳定模块，可能会随着使用的范围，出现问题
	<+++> 进行中，还没投入使用
*/

const format = "start message : %s \n"

type Aurora struct {
	lock             *sync.RWMutex
	ctx              context.Context     //服务器顶级上下文，通过此上下文可以跳过 go web 自带的子上下文去开启纯净的子go程，结束此上下文 web服务也将结束 <***>
	cancel           func()              //取消上下文 <***>
	port             string              //服务端口号 <***>
	router           *route              //路由服务管理 <***>
	projectRoot      string              //项目根路径 <***>
	resource         string              //静态资源管理 默认为 root 目录 <***>
	resourceMappings map[string][]string //静态资源映射路径标识 <***>
	resourceMapType  map[string]string   //常用的静态资源头 <--->

	MaxMultipartMemory int64 //文件上传大小配置

	load       chan struct{}
	message    chan string //启动自带的日志信息 <***>
	errMessage chan string
	initError  chan error //路由器级别错误通道 一旦初始化出错，则结束服务，检查配置 <***>

	routeInterceptor []interceptorArgs     //拦截器初始华切片 <***>
	interceptorList  []Interceptor         //全局拦截器 <***>
	container        *containers           //第三方配置整合容器,原型模式
	pools            map[string]*sync.Pool // 容器池，用于存储配置实例，保证了在整个服务器运行期间 不会被多个线程同时占用唯一变量	     	<+++>
	options          map[string]*Option    // 配置项，每个第三方库/框架的唯一  	<+++>
	cnf              *viper.Viper          // 配置实例 <***>
	Server           *http.Server          // web服务器 <***>
	GrpcServer       *grpc.Server          //
	Ln               net.Listener          // web服务器监听
}

// New :最基础的 Aurora 实例
func New() *Aurora {

	a := &Aurora{
		lock: &sync.RWMutex{},
		port: "8080", //默认端口号
		router: &route{
			mx: &sync.Mutex{},
		},
		Server:          &http.Server{},
		resource:        "", //设定资源默认存储路径
		initError:       make(chan error),
		resourceMapType: make(map[string]string),
		load:            make(chan struct{}),
		message:         make(chan string),
		container: &containers{
			rw:         &sync.RWMutex{},
			prototypes: make(map[string]interface{}),
		},
		pools:   make(map[string]*sync.Pool),
		options: make(map[string]*Option),
	}
	startLoading(a)
	loadResourceHead(a)
	//fmt.Println(print_aurora())
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

// GuideTLS 启动 Aurora TLS服务器，默认端口号8080
// args[0]	证书路径参数，必选项
// args[1]	私钥路径参数，必选项
// args[2]	选择端口绑定参数，可选项
func (a *Aurora) GuideTLS(args ...string) {
	a.tls(args...)
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

func (a *Aurora) tls(args ...string) {
	if len(args) < 2 {
		panic("Parameter error")
	}
	if len(args) <= 2 {
		a.Server.Addr = ":" + a.port
	} else {
		p := args[2]
		if p[0:1] != ":" {
			p = ":" + p
		}
		a.Server.Addr = p
		a.port = p
	}
	//a.Server.Handler = a
	if a.GrpcServer != nil {
		// 在 Aurora 和 GrpcServer 两个路由器中间 加一个原生路由器 用于 分别提供 http 和 https 服务（来自grpc 官方文档示例 url: https://pkg.go.dev/google.golang.org/grpc#NewServer ）
		a.Server.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.ProtoMajor == 2 && strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
				a.GrpcServer.ServeHTTP(writer, request)
			} else {
				a.ServeHTTP(writer, request)
			}
			return
		})
	} else {
		a.Server.Handler = a
	}
	err := a.Server.ListenAndServeTLS(args[0], args[1]) //启动服务器
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
	//a.router.RegisterInterceptor(path, interceptor...)
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

// ViewHandle 修改默认视图解析接口
func (a *Aurora) ViewHandle(views Views) {
	a.router.defaultView = views
}

// View 默认视图解析
func (a *Aurora) View(ctx *Ctx, html string) {
	parseFiles, err := template.ParseFiles(a.projectRoot + "/" + a.resource + html)
	if err != nil {
		a.errMessage <- err.Error()
		return
	}
	err = parseFiles.Execute(ctx.Response, ctx.Attribute)
	if err != nil {
		a.errMessage <- err.Error()
		return
	}
}

// loadingInterceptor 加载局部拦截器
func (a *Aurora) loadingInterceptor() {
	if a.routeInterceptor != nil {
		for i := 0; i < len(a.routeInterceptor); i++ {
			e := a.routeInterceptor[i]
			a.router.RegisterInterceptor(e.path, e.list...)
		}
	}
}

// baseContext 初始化 Aurora 顶级上下文
func (a *Aurora) baseContext(ln net.Listener) context.Context {
	//初始化 Aurora net.Listener 变量，用于整合grpc
	a.Ln = ln
	a.loadingInterceptor() //加载 拦截器
	if a.GrpcServer != nil {
		go func(ln net.Listener) {

		}(ln)
	}
	l := fmt.Sprintf("The server successfully runs on port %s", a.port)
	c, f := context.WithCancel(context.TODO())
	a.ctx = c
	a.cancel = f
	a.message <- fmt.Sprintf("Initialize the top-level context and clear the function")
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

			//初始化实例日志信息
			case info := <-a.message:
				log.Printf(format, info)

			//服务器内部错误信息
			case e := <-a.errMessage:
				log.Println(e)

			case err := <-a.initError: //启动初始化错误处理
				log.Fatal(err)
			}
		}
	}(a)
}
func print_aurora() string {
	s := "    /\\\n   /  \\  _   _ _ __ ___  _ __ __ _\n  / /\\ \\| | | | '__/ _ \\| '__/ _` |\n / ____ \\ |_| | | | (_) | | | (_| |\n/_/    \\_\\__,_|_|  \\___/|_|  \\__,_|\n:: aurora ::   (v0.1.5.RELEASE)"
	/*
	       /\
	      /  \  _   _ _ __ ___  _ __ __ _
	     / /\ \| | | | '__/ _ \| '__/ _` |
	    / ____ \ |_| | | | (_) | | | (_| |
	   /_/    \_\__,_|_|  \___/|_|  \__,_|
	   :: aurora ::   (v0.0.1.RELEASE)

	*/
	return s
}

func GetTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}
