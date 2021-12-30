package aurora

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"html/template"
	"log"
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

const format = "[信息]:%s \n"

type Aurora struct {
	name             string
	lock             *sync.RWMutex
	ctx              context.Context     //服务器顶级上下文，通过此上下文可以跳过 go web 自带的子上下文去开启纯净的子go程，结束此上下文 web服务也将结束 <***>
	cancel           func()              //取消上下文 <***>
	port             string              //服务端口号 <***>
	router           *route              //路由服务管理 <***>
	projectRoot      string              //项目根路径 <***>
	resource         string              //静态资源管理 默认为 root 目录 <***>
	resourceMappings map[string][]string //静态资源映射路径标识 <***>
	resourceMapType  map[string]string   //常用的静态资源头 <***>

	MaxMultipartMemory int64       //文件上传大小配置
	message            chan string //启动自带的日志信息 <***>
	errMessage         chan string //服务内部api处理错误消息日志<***>
	initError          chan error  //路由器级别错误通道 一旦初始化出错，则结束服务，检查配置 <***>

	plugins          []PluginFunc       //全局插件处理链，每个请求都会走一次,待完善只实现了对插件统一调用，还未做出对插件中途取消，等操作。plugin 发生panic会阻断待执行的业务处理器，可借助panic进行中断，配合ctx进行消息返回<--->
	routeInterceptor []interceptorArgs  //拦截器初始华切片 <***>
	interceptorList  []Interceptor      //全局拦截器 <***>
	gorms            map[int][]*gorm.DB //存储gorm各种类型的连接实例，默认初始化从配置文件中读取<***>
	goredis          []*redis.Client    //存储go-redis 配置实例

	cnf    *viper.Viper // 配置实例，读取配置文件 <***>
	Server *http.Server // web服务器 <***>
	grpc   *grpc.Server // 用于接入grpc支持https服务 <***>,整合 grpc 需要 http 2
	Ln     net.Listener // web服务器监听,启动服务器时候初始化

	consulClient *api.Client
}

// New :最基础的 Aurora 实例
// config:指定加载配置文件
func New(config ...string) *Aurora {
	a := &Aurora{
		lock: &sync.RWMutex{},
		port: "8080", //默认端口号
		router: &route{
			mx: &sync.Mutex{},
		},
		Server:          &http.Server{},
		resource:        "", //设定资源默认存储路径，需要连接项目更目录 和解析出来资源的路径，资源路径解析出来是没有前缀 “/” 的作为 resource属性，在其两边加上 斜杠
		initError:       make(chan error),
		resourceMapType: make(map[string]string),
		gorms:           make(map[int][]*gorm.DB),
		message:         make(chan string),
		errMessage:      make(chan string),
	}
	startLoading(a) //开启日志线程
	a.message <- fmt.Sprintf("Golang 版本信息:%1s", runtime.Version())
	a.message <- fmt.Sprintf("开始加载application.yml配置文件.")
	a.viperConfig(config...) //加载默认位置的 application.yml
	projectRoot, _ := os.Getwd()
	a.projectRoot = projectRoot
	a.router.defaultView = a //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
	a.router.AR = a
	//加载配置文件中定义的 端口号
	port := a.cnf.GetString("aurora.server.port")
	if port != "" {
		a.port = port
	}
	//读取配置路径
	p := a.cnf.GetString("aurora.resource.static")
	if p != "" {
		if p[:1] != "/" {
			p = "/" + p
		}
		if p[len(p)-1:] != "/" {
			p = p + "/"
		}
		a.resource = p
	}
	name := a.cnf.GetString("aurora.application.name")
	if name != "" {
		a.name = name
	}
	//加载默认的全局拦截器
	a.interceptorList = []Interceptor{
		0: &defaultInterceptor{},
	}
	a.message <- fmt.Sprintf("项目根路径信息:%1s", a.projectRoot)
	a.message <- fmt.Sprintf("服务器端口号:%1s", a.port)
	a.message <- fmt.Sprintf("服务器静态资源根目录:%1s", a.resource)
	loadResourceHead(a) //加载静态资源头
	a.loadGormConfig()  //加载配置文件中的gorm配置项
	a.loadGoRedis()     //加载go-redis
	a.consulConfig()
	a.Server.BaseContext = a.baseContext //配置 上下文对象属性
	return a
}

func (a *Aurora) App(name string) {
	if name == "" {
		a.name = name
	}
}

// Guide 启动 Aurora 服务器，默认端口号8080
func (a *Aurora) Guide(port ...string) error {
	return a.run(port...)
}

// GuideTLS 启动 Aurora TLS服务器，默认端口号8080
// args[0]	证书路径参数，必选项
// args[1]	私钥路径参数，必选项
// args[2]	选择端口绑定参数，可选项
func (a *Aurora) GuideTLS(args ...string) error {
	return a.tls(args...)
}

func (a *Aurora) run(port ...string) error {
	if port != nil && len(port) > 1 {
		return errors.New("too mach port")
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
	return a.Server.ListenAndServe() //启动服务器
}

func (a *Aurora) tls(args ...string) error {
	if len(args) < 2 {
		return errors.New("parameter error")
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
	if a.grpc != nil {
		// 在 Aurora 和 GrpcServer 两个路由器中间 加一个原生路由器 用于 分别提供 http 和 https 服务（来自grpc 官方文档示例 url: https://pkg.go.dev/google.golang.org/grpc#NewServer ）
		a.Server.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.ProtoMajor == 2 && strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
				a.grpc.ServeHTTP(writer, request)
			} else {
				a.ServeHTTP(writer, request)
			}
			return
		})
	} else {
		a.Server.Handler = a
	}
	return a.Server.ListenAndServeTLS(args[0], args[1]) //启动服务器
}

// ResourceMapping 资源映射(暂时弃用)
//添加静态资源配置，t资源类型必须以置源后缀命名，
//paths为t类型资源的子路径，可以一次性设置多个。
//每个资源类型最调用一次设置方法否则覆盖原有设置
func (a *Aurora) resourceMapping(Type string, Paths ...string) {
	a.registerResourceType(Type, Paths...)
}

// StaticRoot 设置静态资源根路径，会覆盖配置文件选项
func (a *Aurora) StaticRoot(root string) {
	if root == "" {
		a.resource = "/"
	}
	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
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

	parseFiles, err := template.ParseFiles(html)
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
	l := fmt.Sprintf("服务器成功绑定到端口:%s", a.port)
	c, f := context.WithCancel(context.TODO())
	a.ctx = c
	a.cancel = f
	a.message <- fmt.Sprintf("初始化上下文实例和清除函数.")
	a.message <- l
	return c
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

func (a *Aurora) ProjectPath() string {
	return a.projectRoot
}
