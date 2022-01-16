package aurora

import (
	"context"
	"errors"
	"fmt"
	"github.com/awensir/aurora-email/email"
	"github.com/awensir/minilog/mini"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"html/template"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

/*
	<***> 基于稳定模块，无需更改
	<---> 非稳定模块，可能会随着使用的范围，出现问题
	<+++> 进行中，还没投入使用
*/

const (
	Dev = mini.ALL
	Pro = mini.INFO
)

type Aurora struct {
	name            string            //服务名称
	ctx             context.Context   //服务器顶级上下文，通过此上下文可以跳过 go web 自带的子上下文去开启纯净的子go程，结束此上下文 web服务也将结束 <***>
	cancel          func()            //取消上下文 <***>
	host            string            //主机信息
	port            string            //服务端口号 <***>
	auroraLog       *mini.Log         //日志,(代码文件logs.go 单独分离出去为 minilog库 ，将不再使用logs.log)<***>
	router          *route            //路由服务管理 <***>
	projectRoot     string            //项目根路径 <***>
	resource        string            //静态资源管理 默认为 root 目录 <***>
	resourceMapType map[string]string //常用的静态资源头 <***>

	MaxMultipartMemory int64 //文件上传大小配置

	plugins          []PluginFunc       //<--->全局插件处理链，每个请求都会走一次,待完善只实现了对插件统一调用，还未做出对插件中途取消，等操作。plugin 发生panic会阻断待执行的业务处理器，可借助panic进行中断，配合ctx进行消息返回
	routeInterceptor []interceptorArgs  //拦截器初始华切片 <***>
	interceptorList  []Interceptor      //全局拦截器 <***>
	gorms            map[int][]*gorm.DB //存储gorm各种类型的连接实例，默认初始化从配置文件中读取<***>
	goredis          []*redis.Client    //存储go-redis 配置实例
	email            *email.Client

	cnf          *viper.Viper // 配置实例，读取配置文件 <***>
	remoteConfig func() *viper.Viper
	cnfLock      *sync.RWMutex //分布式配置中心处理动态刷新web 服务配置的读写锁

	Server *http.Server // web服务器 <***>
	grpc   *grpc.Server // 用于接入grpc支持,该整合意义在于让grpc服务和http服务公用一个ip和端口号,仅支持tls通讯情况下的整合
	Ln     net.Listener // web服务器监听,启动服务器时候初始化

	consuls []*api.Client //<+++>
	consul  consulConfig  //集成consul模块<+++>
}

// New :最基础的 Aurora 实例
// config:指定加载配置文件
func New(config ...string) *Aurora {
	//初始化基本属性
	a := &Aurora{
		cnfLock: &sync.RWMutex{},
		port:    "8080", //默认端口号
		router: &route{
			mx: &sync.Mutex{},
		},
		auroraLog:       mini.NewLog(Dev),
		Server:          &http.Server{},
		resource:        "", //设定资源默认存储路径，需要连接项目更目录 和解析出来资源的路径，资源路径解析出来是没有前缀 “/” 的作为 resource属性，在其两边加上 斜杠
		consuls:         make([]*api.Client, 0),
		resourceMapType: make(map[string]string),
		gorms:           make(map[int][]*gorm.DB),
	}
	a.auroraLog.Info(fmt.Sprintf("golang version information:%1s", runtime.Version()))
	a.auroraLog.Info(fmt.Sprintf("start loading the application.yml configuration file."))
	a.viperConfig(config...) //加载默认位置的 application.yml,config可配置路径信息
	projectRoot, _ := os.Getwd()
	a.projectRoot = projectRoot //初始化项目路径信息
	a.router.defaultView = a    //初始化使用默认视图解析,aurora的视图解析是一个简单的实现，可以通过修改 a.Router.DefaultView 实现自定义的试图处理，框架最终调用此方法返回页面响应
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
		a.auroraLog.Info(fmt.Sprintf("server static resource root directory:%1s", a.resource))
	}

	name := a.cnf.GetString("aurora.application.name")
	if name != "" {
		a.name = name
		a.auroraLog.Info("the service name is " + a.name)
	}
	//加载默认的全局拦截器
	a.interceptorList = []Interceptor{
		0: &defaultInterceptor{},
	}
	a.auroraLog.Info("initialize the default top-level interceptor")
	a.loadResourceHead()                 //加载静态资源头
	a.loadGormConfig()                   //加载配置文件中的gorm配置项
	a.loadGoRedis()                      //加载go-redis
	a.loadEmail()                        //加载邮件配置
	a.consulConfig()                     //加载consul
	a.Server.BaseContext = a.baseContext //配置 上下文对象属性
	return a
}

// Level 修改系统日志输出级别，系统默认输出任何级别
func (a *Aurora) Level(level int) {
	a.auroraLog.Level(level)
}

func (a *Aurora) RemoteConfig(fun func() *viper.Viper) {
	a.remoteConfig = fun
}

// ServiceName 设置程序服务名称,配置文件信息由于api设置
func (a *Aurora) ServiceName(name string) {
	if name == "" {
		a.name = name
	}
}

// Port 获取端口号 int类型
func (a *Aurora) Port() int {
	atoi, err := strconv.Atoi(a.port)
	if err != nil {
		return 0
	}
	return atoi
}

// ServerIp 获取服务器ip地址信息
func ServerIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
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
			fmt.Println()
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
}

// TopLevelInterceptor 配置默认顶级拦截器
func (a *Aurora) TopLevelInterceptor(interceptor Interceptor) {
	a.interceptorList[0] = interceptor
	a.auroraLog.Info("web default global Rout Interceptor successds")
}

// AddInterceptor 追加全局拦截器
func (a *Aurora) AddInterceptor(interceptor ...Interceptor) {
	//追加全局拦截器
	for _, v := range interceptor {
		a.interceptorList = append(a.interceptorList, v)
		a.auroraLog.Info("web global Rout Interceptor successds")
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
		a.auroraLog.Error(err.Error())
		return
	}
	err = parseFiles.Execute(ctx.Response, ctx.Attribute)
	if err != nil {
		a.auroraLog.Error(err.Error())
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
	c, f := context.WithCancel(context.TODO())
	a.ctx = c
	a.cancel = f
	a.auroraLog.Info("successfully initialized the context instance and cleanup function.")
	a.auroraLog.Info(fmt.Sprintf("the server successfully binds to the port:%s", a.port))
	return c
}

func (a *Aurora) ProjectPath() string {
	return a.projectRoot
}
