# Aurora

## 系统架构

字典树

服务处理函数

REST API

路由规则

全局拦截器

静态资源处理

Session会话

## 映射管理

管理服务器启动期间注册的Web服务，服务将分别存放到各自类型的服务容器中，Aurora中的路由管理默认提供两种方式，get和post请求。后续会支持更多的请求方式

- Get
- Post

服务处理函数签名：

```go
	type Servlet func(ctx *Context) interface{}
```

接口统一注册签名

```go
	func Mapping(url string,fun aurora.Servlet)
```

通过包名对不同类型的请求进行注册

```go
	//get请求
	get.Mapping("/", func(ctx *aurora.Context) interface{} {
		
		return "/html/index.html"
	})
	//post请求
	post.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})
```

## 请求处理

请求返回值是一个interface{}类型，意味着你可以返回任何类型，aurora约定返回三种类型，结构体，页面，错误

- struct
- path
- error

## 静态资源

静态资源的解析，路由器默认的解析文件夹是static，所有静态资源需要放到static下面进行使用。这个默认路径是读取资源文件的主要路径，若设置为 ""空字符串则默认为项目跟目录，根据浏览器在html页面上请求资源的格式，静态资源所在的文件夹必须添加路由。html和静态资源应该都放在static文件夹下。

默认静态资源根路径为static

配置静态资源

```go
//t 静态资源类型，paths为对应静态资源下的路径，对于图片的资源处理，还在优化
func RegisterResource(t string,paths ...string) 

func main() {
    //设置js文件的静态资源路径
	config.RegisterResource("js","js")
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})

	post.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})
	
	aurora.RunApplication("8080")
}
```

更改默认静态资源根路径

```go
func main() {
	config.RegisterResource("js","js")
	
	//更改默认静态资源根路径为resource
	config.ResourceRoot("resource")
    
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})

	post.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/html/index.html"
	})

	aurora.RunApplication("8080")
}
```



