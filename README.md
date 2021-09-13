Aurora

[TOC]



## 简介

​		框架特点，简单易配置，通过包名对服务器中的各种属性进行设置，甚至不需要初始化任何变量，开发者可以专注的处理业。对请求返回值的处理也更加灵活，约定自动返回json格式，也可以进行视图解析加载给浏览器。自定义的错误机制可以让业务逻辑中经可能减少err的判空，也可以对专属错误进行可控的提示响应给浏览器。官方交流群:836414068，备注hub

## 路由注册

aurora支持REST API，在方便的同时以意味着也有一定的缺陷，使用了rest api的路由注册容易和其它路径产生冲突。

aurora路由设计规则如下：

```tex
路由存储规则参考HttpRouter
基于查询树的路由器
路由器规则:
   1.无法存储相同的路径
      1)形同路径的判定：校验参数相同，并且节点函数不为nil，节点函数为nil的节点说明，这个路径是未注册过，被提取为公共根
   2.路径查找按照逐层检索
   3.路由树上面存储者当前路径匹配的服务处理函数
   4.注册路径必须以 / 开头
   5.发生公共根
      1)节点和被添加路径产生公共根，提取公共根后，若公共根未注册，服务处理函数将为nil
      2)若节点恰好是公共根，则设置函数
   6.REST 风格注册
      1)同一个根路径下只能有一个REST 子路径
      2)REST 作为根路径也只能拥有一个REST 子路径
      3)REST 路径会和其它非REST同级路径发生冲突
   7.注册路径不能以/结尾（bug未修复，/user /user/ 产生 /user 的公共根 使用切割解析路径方式，解析子路径，拼接剩余子路径会存在bug ,注册路径的时候强制无法注册 / 结尾的 url）
```



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

请求返回值是一个interface{}类型，意味着你可以返回任何类型，aurora约定返回三种类型，结构体，页面，错误，自定义错误处理

- struct

  struct：返回一个任意的结构体，将被解析为json发送给浏览器

- path

  path：是一个页面路径，约定必须以 / 开头，path一定是静态资源目录下开始的路径，一般是返回html页面

- error

  error：返回一个错误(对错误的处理暂定直接发送浏览器json串)

  页面响应对模板的支持还在设计中，预计知识对golang的模板语法简单的封装一下，还是尽可能的以json方式为主。



### 自定义错误处理：

实现错误处理方法既可以自定义错误的处理，错误处理和服务处理参数虽然相同，但是不会走全局拦截器，只负责对产生的错误进行包装处理，然后给浏览器做出需要的响应。

```go
 ErrorHandler(ctx *Context) interface{}


//编写一个结构体 实现  ErrorHandler(ctx *Context) interface{} 方法即可
type TestErr struct {
	 error
}
// 绑定的方式 使用结构体方式即可，对于指针的支持后续进行改进
func (t TestErr) ErrorHandler(ctx *aurora.Context) interface{} {
	//对error 进行指定处理，选择输出

	return "error"
}

func main() {
	config.RegisterResource("js", "js","test")
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return TestErr{fmt.Errorf("err")}
	})
	aurora.RunApplication("8080")
}

```

错误处理，中需要避免再次返回处理者本身的类型，会造成死循环，无限递归最终栈溢出。请求处理的返回值同样适用于错误处理的返回值

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

## 上下文对象

rest api 参数获取

请求转发



## 拦截器

拦截器机制，当前支持全局拦截器，和局部拦截器两种，局部拦截器支持path子路径匹配（/*,的形式），造成了一个缺陷需要匹配大多数路径而非常麻烦，后续改进。

拦截器机制，和springboot类似的效果：

```go
// Interceptor 拦截器统一接口，实现这个接口就可以向服务器注册一个全局或者指定路径的处理拦截，此处的Context是aurora包下的上下文变量
type Interceptor interface {
	PreHandle(ctx *Context) bool
	PostHandle(ctx *Context)
	AfterCompletion(ctx *Context)
}
```

服务器内置了一个默认的全局拦截器

```go
// DefaultInterceptor 实现全局请求处理前后环绕
type DefaultInterceptor struct {
             
}

func (de DefaultInterceptor) PreHandle(ctx *Context) bool {
	//处理服务之前必须经过的函数，返回true表示放行服务继续执行后面的拦截器或者业务，false则会终止服务继续执行，将不会执行到对应业务
	return true
}

func (de DefaultInterceptor) PostHandle(ctx *Context) {
	//业务完成处理后执行此函数，此刻还没有想浏览器发送视图信息
}

func (de DefaultInterceptor) AfterCompletion(ctx *Context)  {
	//试图解析发送完成后执行，此处表示服务业务已经完全处理完毕
}
```

拦截器的注册（config包中调用方法）

```go
func main() {
    //RegisterInterceptor添加全局拦截器
	config.RegisterInterceptor(MyInterceptor1{})
	post.Mapping("/", func(ctx *aurora.Context) interface{} {
		var body Body
		ctx.PostBody(&body)
		return body
	})

	get.Mapping("/abc", func(ctx *aurora.Context) interface{} {

		return "/abc"
	})
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/abc"
	})
	config.RegisterPathInterceptor("/abc",MyInterceptor2{})

	config.RegisterPathInterceptor("/",MyInterceptor3{})

	aurora.RunApplication("8080")
}
```

通配符拦截器

```go
func main() {

	get.Mapping("/abc/bbc", func(ctx *aurora.Context) interface{} {

		return "/abc"
	})
	get.Mapping("/abc/bbc/asd", func(ctx *aurora.Context) interface{} {

		return "/abc/bbc/asd"
	})

	get.Mapping("/abc/bbc/aaa", func(ctx *aurora.Context) interface{} {

		return "/abc/bbc/aaa"
	})
	get.Mapping("/abc/qaq", func(ctx *aurora.Context) interface{} {

		return "/abc/qaq"
	})
	get.Mapping("/abc/qaq/csdn", func(ctx *aurora.Context) interface{} {

		return "/abc/qaq/csdn"
	})
	get.Mapping("/", func(ctx *aurora.Context) interface{} {

		return "/"
	})
    
    // /abc/* 只能匹配以/abc/结尾的父路径，即/abc 这个路径是不能被 MyInterceptor2所拦截的，如果有需要只能在加一个
    // 例如多加一个config.RegisterPathInterceptor("/abc", MyInterceptor2{}) 即可
	config.RegisterPathInterceptor("/abc/*", MyInterceptor2{})

	config.RegisterPathInterceptor("/", MyInterceptor3{})

	config.RegisterPathInterceptor("/abc/bbc/aaa", MyInterceptor4{})

	aurora.RunApplication("8080")
}
```

注意事项：通配符配置拦截器，只支持 /xxx/* ,xxx是具体完整的父路径，不支持使用*来切割子路径进行匹配

局部拦截器，依赖于路由树，所以注册局部拦截器时候必须等待路由注册完毕才能正常注册成功，全局则不需要依赖于路由树。

## session机制

