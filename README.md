# Aurora

## 简介

​		框架特点，简单易配置。对请求返回值的处理也更加灵活，约定自动返回json格式，也可以进行视图解析加载给浏览器。自定义的错误机制可以让业务逻辑中经可能减少err的判空，也可以对专属错误进行可控的提示响应给浏览器。官方交流群:836414068，备注hub，部分功能不完善使用过程中有任何疑问可以添加QQ：1219449282，微信：Saber__o，备注hub，即可。

## Demo

### 拉取依赖

```
import "github.com/awensir/go-aurora/aurora"
```

### 入门示例

```go
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		//结果响应给 调用者
		return "hello web"
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
out：
hello web
```



## 路由注册

​			aurora支持REST API，在方便的同时以意味着也有一定的缺陷，使用了rest api的路由注册容易和其它路径产生冲突。

aurora路由设计规则如下：

```tex
路由存储规则参考HttpRouter
基于查询树的路由器
路由器规则:
   1.无法存储相同的路径
      1)形同路径的判定：校验参数相同，并且节点函数不为nil，节点函数为nil的节点说明，这个路径是未注册过，被提取为公共根
      2)第一条规则进行修改，注册相同路径处理函数，默认覆盖前面相同的处理函数。
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
   8. 浏览器访问接口，不能带有可编码符号，特别是{} ，{}是框架解析rest ful 参数的标识，接收到带有{},比如/a/b/{sss}/c,带有{或}的请求都视为非法url
```

## 路由分组

路由分组用于注册服务处理时候带来的超长的路径重复复制粘贴问题。分组对象只能进行路径注册等简单操作，对web的配置上还是以基础实例为主

```
核心方法
// Group 路由分组  必须以 “/” 开头分组
func (a *Aurora) Group(path string) *group

// Group 路由分组  必须以 “/” 开头分组，支持多级分组
func (g *group) Group(path string) *group

```



## 映射管理

服务处理函数签名：

```go
	type Servlet func(c *Ctx) interface{}
```

支持常用的注册方式：

```go
type Routes interface {
	GET(string, Servlet) interface{}
	POST(string, Servlet) interface{}
	PUT(string, Servlet) interface{}
	DELETE(string, Servlet) interface{}
	HEAD(string, Servlet) interface{}
}
```

## 响应处理

### String

```
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		//直接返回字符串
		return "hello web"
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
##############################################################################################################################
out：
hello web
```



### Json

```go
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		s := struct {
			Name string
			Age  int
		}{Name: "test", Age: 20}
		//直接返回结构体，自动编码json
		return s
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
##############################################################################################################################
out:
{"Name":"test", "Age":20}
```



### Error

```go
package main

import (
	"errors"
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return errors.New("is error")
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
##############################################################################################################################
out: status：500
error:is error
```



### WebError

```go
package main

import (
	"errors"
	"github.com/awensir/go-aurora/aurora"
)

// 绑定 ErrorHandler(c *aurora.Context) interface{} 函数即可
type AgeErr struct {
	err error
}

func (receiver *AgeErr) ErrorHandler(c *aurora.Ctx) interface{} {
	/*
		对同一类型的错误 统一处理
	*/
	return receiver.err
}

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		s := struct {
			Name string
			Age  int
		}{Name: "test", Age: 20}
		if s.Age == 20 {
			return &AgeErr{err: errors.New("is error")}
		}
		return s
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
##############################################################################################################################
out: status：500
error:is error
```

使用WebError 进行统一错误处理的方式，一定要避免返回其它错误机制，必须处理错误给出确切响应，继续返回错误机制会造成服务器资源崩溃。

### Nil

```go
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return nil
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
##############################################################################################################################
out:

```

返回nil ，框架不会对 nil 结果做出任何响应，等于本次访问什么都没发生

### 自定义错误处理：

实现错误处理方法既可以自定义错误的处理，错误处理和服务处理参数虽然相同，但是不会走全局拦截器，只负责对产生的错误进行包装处理，然后给浏览器做出需要的响应。

```go
 ErrorHandler(ctx *Ctx) interface{}


//编写一个结构体 实现  ErrorHandler(ctx *Ctx) interface{} 方法即可
type TestErr struct {
	 error
}
// 绑定的方式 使用结构体方式即可，对于指针的支持后续进行改进
func (t TestErr) ErrorHandler(ctx *aurora.Ctx) interface{} {
	//对error 进行指定处理，选择输出

	return "error"
}

```

***错误处理，中需要避免再次返回错误处理类型，会造成死循环，无限递归最终栈溢出。请求处理的返回值同样适用于错误处理的返回值***

## 静态资源

静态资源的解析，路由器默认的解析文件夹是static，所有静态资源需要放到static下面进行使用。这个默认路径是读取资源文件的主要路径，若设置为 ""空字符串则默认为项目跟目录，根据浏览器在html页面上请求资源的格式，静态资源所在的文件夹必须添加路由。html和静态资源应该都放在static文件夹下。

默认静态资源根路径为static

配置静态资源

```go
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	//设置静态资源根路径
	a.StaticRoot("static")
	// ResourceMapping 资源映射
	//添加静态资源配置，t资源类型必须以置源后缀命名，
	//paths为t类型资源的子路径(一级子路径：static/xxx/aaa,xxx为第一级)，可以一次性设置多个。
	//每个资源类型最调用一次设置方法否则覆盖原有设置
	a.ResourceMapping("js", "js", "jsfiles")
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		//直接返回静态资源中的 页面作为响应数据
		return "/html/index.html"
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
```



## REST FUL

```go
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/{name}", func(c *aurora.Ctx) interface{} {
		
		return c.Args["name"]
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
```



## 上下文对象

### 请求转发

支持返回值方式指定转发服务，真实处理交给转发者返回给浏览器，或者返回值转发服务，两者转发的格式遵循路由注册的url方式，理论上是可以支持RAST API转发，但是不推荐。

```go
package main

import (
	"github.com/awensir/go-aurora/aurora"
)

func main() {
	//获取 aurora 路由实例
	a := aurora.New()
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return "forward:/abc"
	})

	a.GET("/abc", func(c *aurora.Ctx) interface{} {

		return "/abc"
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
```

上下文对象对原始的web api封装不是很，但是提供了暴露了请求体和响应体。给开发者使用，后续会逐渐完善。web的请求处理可以完全按照go web的方式进行，把页面响应交给框架即可，同时对浏览器响应信息可能导致冲突，但是响应头的设置不会有所影响。

### Get请求参数获取

```go
// Get 获取一个字符串参数
func (c *Context) Get(Args string) (string, error)

// GetInt 获取一个整数参数
func (c *Context) GetInt(Args string) (int, error)

// GetFloat64 获取一个64位浮点参数
func (c *Context) GetFloat64(Args string) (float64, error) 

// GetSlice 获取切片类型参数
func (c *Context) GetSlice(Args string) ([]string, error) 

// GetIntSlice 整数切片
func (c *Context) GetIntSlice(Args string) ([]int, error) 

// GetFloat64Slice 浮点切片
func (c *Context) GetFloat64Slice(Args string) ([]float64, error) 

```

### Post 请求体绑定

```go
// JsonBody 读取Post请求体Json或表单数据数据解析到body中,
func (c *Ctx) JsonBody(body interface{}) error 
```



### 文件上传

参考gin api



### 日志调用

日志相关api，Aurora 默认使用的是logrus 日志框架作为封装提供上下文变量进行替代fmt.Print 输出消息信息，一下是相关的api

```go
// INFO 打印 info 日志信息
func (c *Ctx) INFO(info ...interface{})

// WARN 打印 警告信息
func (c *Ctx) WARN(warning ...interface{}) 

// ERROR 打印错误信息
func (c *Ctx) ERROR(error ...interface{}) 

// PANIC 打印信息并且结束程序
func (c *Ctx) PANIC(panic ...interface{}) 
```



## 拦截器

拦截器机制，当前支持全局拦截器，和局部拦截器两种，局部拦截器支持path子路径匹配（/*,的形式），造成了一个缺陷需要匹配大多数路径而非常麻烦，后续改进。

### 拦截器机制

和springboot类似的效果：

```go
// Interceptor 拦截器统一接口，实现这个接口就可以向服务器注册一个全局或者指定路径的处理拦截，此处的Context是aurora包下的上下文变量
type Interceptor interface {
	PreHandle(ctx *Ctx) bool
	PostHandle(ctx *Ctx)
	AfterCompletion(ctx *Ctx)
}
```

服务器内置了一个默认的全局拦截器

```go
// DefaultInterceptor 实现全局请求处理前后环绕
type DefaultInterceptor struct {
             
}

func (de DefaultInterceptor) PreHandle(ctx *Ctx) bool {
	//处理服务之前必须经过的函数，返回true表示放行服务继续执行后面的拦截器或者业务，false则会终止服务继续执行，将不会执行到对应业务
	return true
}

func (de DefaultInterceptor) PostHandle(ctx *Ctx) {
	//业务完成处理后执行此函数，此刻还没有想浏览器发送视图信息
}

func (de DefaultInterceptor) AfterCompletion(ctx *Ctx)  {
	//试图解析发送完成后执行，此处表示服务业务已经完全处理完毕
}
```

默认连接器实现了对服务访问的控制台日志输出，如果想要自定义默认全局拦截器，可以修改默认拦截器如下，然后自己实现拦截器逻辑即可

```go
// DefaultInterceptor 配置默认顶级拦截器
func (a *Aurora) DefaultInterceptor(interceptor Interceptor) 
```

### 拦截器的注册

```go
// AddInterceptor 追加全局拦截器
func (a *Aurora) AddInterceptor(interceptor ...Interceptor)

// RouteIntercept path路径上添加一个或者多个路由拦截器
func (a *Aurora) RouteIntercept(path string, interceptor ...Interceptor)
```

​		通配符拦截器，eg： /abc/* 只能匹配以/abc/结尾的父路径以及 /abc 本身，如果继续单独添加RouteIntercept("/abc", MyInterceptor2{}) 配置访问/abc 会被两个配置分别拦截，两种方式不冲突。

​		注意事项：通配符配置拦截器，只支持 /xxx/* ,xxx是具体完整的父路径，不支持使用*来切割子路径进行匹配

***局部拦截器，依赖于路由树，所以注册局部拦截器时候必须等待路由注册完毕才能正常注册成功(以优化，统一放在最后处理)，全局则不需要依赖于路由树。***

### 准备一下接口

```go
// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a", func(c *aurora.Ctx) interface{} {

		return c.Args
	})
	a.GET("/b", func(c *aurora.Ctx) interface{} {

		return c.Args
	})
	a.GET("/a/b", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a/b/c/{name}", func(c *aurora.Ctx) interface{} {

		return c.Args
	})

	a.GET("/a/b/cc", func(c *aurora.Ctx) interface{} {

		return c.Args
	})
```

### 创建拦截器

```go
package main

import (
	"fmt"
	"github.com/awensir/Aurora/aurora"
)


type MyInterceptor struct {
}

func (de *MyInterceptor) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle")
	return true
}

func (de *MyInterceptor) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle")
}

func (de *MyInterceptor) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion")
}

type MyInterceptor1 struct {
}

func (de *MyInterceptor1) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle1")
	return true
}

func (de *MyInterceptor1) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle1")
}

func (de *MyInterceptor1) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion1")
}


type MyInterceptor2 struct {
}

func (de *MyInterceptor2) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle2")
	return true
}

func (de *MyInterceptor2) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle2")
}

func (de *MyInterceptor2) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion2")
}


type MyInterceptor3 struct {
}

func (de *MyInterceptor3) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle3")
	return true
}

func (de *MyInterceptor3) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle3")
}

func (de *MyInterceptor3) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion3")
}


type MyInterceptor4 struct {
}

func (de *MyInterceptor4) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle4")
	return true
}

func (de *MyInterceptor4) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle4")
}

func (de *MyInterceptor4) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion4")
}


type MyInterceptor5 struct {
}

func (de *MyInterceptor5) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle5")
	return true
}

func (de *MyInterceptor5) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle5")
}

func (de *MyInterceptor5) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion5")
}


type MyInterceptor6 struct {
}

func (de *MyInterceptor6) PreHandle(c *aurora.Ctx) bool {
	fmt.Println("MyPreHandle6")
	return true
}

func (de *MyInterceptor6) PostHandle(c *aurora.Ctx) {
	fmt.Println("MyPostHandle6")
}

func (de *MyInterceptor6) AfterCompletion(c *aurora.Ctx) {
	fmt.Println("MyAfterCompletion6")
}
```



### 添加全局拦截器

```go
// DefaultInterceptor 配置默认顶级拦截器
func (a *Aurora) DefaultInterceptor(interceptor Interceptor)
// AddInterceptor 追加全局拦截器
func (a *Aurora) AddInterceptor(interceptor ...Interceptor)
```

全局拦截器的代表是内置的api 访问信息拦截，每次访问提示访问那个接口。以上api是配置全局拦截器。在此不做演示

### 指定路径添加拦截器

```go
a.RouteIntercept("/",&MyInterceptor{})
a.RouteIntercept("/b",&MyInterceptor{})

访问：localhost:8080/  , localhost:8080/b
控制台输出：
MyPreHandle
MyPostHandle
MyAfterCompletion
2021/10/21 10:50:03 [info  ] ==>  172.0.0.1:59291 →  | GET / | 535.8µs
MyPreHandle
MyPostHandle
MyAfterCompletion
2021/10/21 10:50:19 [info  ] ==>  172.0.0.1:59291 →  | GET /b | 0s
```



### 指定路径添加多个拦截器

```go
a.RouteIntercept("/",&MyInterceptor{},&MyInterceptor1{},&MyInterceptor2{})

访问：localhost:8080/ 
控制台输出：
MyPreHandle
MyPreHandle1
MyPreHandle2
MyPostHandle2
MyPostHandle1
MyPostHandle
MyAfterCompletion2
MyAfterCompletion1
MyAfterCompletion
2021/10/21 10:52:09 [info  ] ==>  172.0.0.1:65145 →  | GET / | 0s
```



### 指定路径添加通配拦截器

```go
a.RouteIntercept("/a/*",&MyInterceptor1{})
访问：localhost:8080/a, localhost:8080/a/b, localhost:8080/a/b/cc , localhost:8080/a/b/c/{name}
控制台输出：
MyPreHandle1
MyPreHandle1
MyPostHandle1
MyPostHandle1
MyAfterCompletion1
MyAfterCompletion1
2021/10/21 10:55:36 [info  ] ==>  172.0.0.1:55582 →  | GET /a | 0s
MyPreHandle1
MyPostHandle1
MyAfterCompletion1
2021/10/21 10:55:50 [info  ] ==>  172.0.0.1:55582 →  | GET /a/b | 0s
MyPreHandle1
MyPostHandle1
MyAfterCompletion1
2021/10/21 10:56:03 [info  ] ==>  172.0.0.1:55582 →  | GET /a/b/cc | 0s
MyPreHandle1
MyPostHandle1
MyAfterCompletion1
2021/10/21 10:56:18 [info  ] ==>  172.0.0.1:55582 →  | GET /a/b/c/test | 271µs
```



### 指定路径添加多个通配拦截器

```go
a.RouteIntercept("/a/*",&MyInterceptor1{},&MyInterceptor2{})

访问：localhost:8080/a, localhost:8080/a/b, localhost:8080/a/b/cc , localhost:8080/a/b/c/{name}
控制台输出：
MyPreHandle1
MyPreHandle2
MyPreHandle1
MyPreHandle2
MyPostHandle2
MyPostHandle1
MyPostHandle2
MyPostHandle1
MyAfterCompletion2
MyAfterCompletion1
MyAfterCompletion2
MyAfterCompletion1
2021/10/21 11:01:08 [info  ] ==>  172.0.0.1:65037 →  | GET /a | 0s
MyPreHandle1
MyPreHandle2
MyPostHandle2
MyPostHandle1
MyAfterCompletion2
MyAfterCompletion1
2021/10/21 11:01:11 [info  ] ==>  172.0.0.1:65037 →  | GET /a/b | 0s
MyPreHandle1
MyPreHandle2
MyPostHandle2
MyPostHandle1
MyAfterCompletion2
MyAfterCompletion1
2021/10/21 11:01:16 [info  ] ==>  172.0.0.1:65037 →  | GET /a/b/cc | 0s
MyPreHandle1
MyPreHandle2
MyPostHandle2
MyPostHandle1
MyAfterCompletion2
MyAfterCompletion1
2021/10/21 11:01:21 [info  ] ==>  172.0.0.1:65037 →  | GET /a/b/c/test | 866.6µs
```



拦截器注册的规则：

1. 拦截器只支持注册已存在的路径
2. 一个普通路径可以注册多个拦截器，多次调用注册拦截器将覆盖前一次的拦截器
3. 路径添加通配路径拦截器，同时路径本身也添加，访问到路径本身会触发2次，一次是通配路径，一次是路径本身，（运行时会触发1次通配和路径本身的拦截器）
4. REST FUL 风格的路径不支持拦截器注册，但支持使用通配的形式进行拦截。
5. 拦截器是按照路径进行注册，不区分方法，对特殊路径，请求方法在拦截器上面有冲突的可以使用中间件的写法来避免。
6. 全局拦截器，任何路径都会触发。
7. 拦截器的触发优先级：全局注册顺序 > 通配顺序 > 局部顺序 (注：中间件的拦截处理需要根据处理函数定义，其执行顺序不在拦截器规则范围内)

## 常用中间件写法

aurora也支持go web 最常用的中间件编写处理

```go
package main

import (
	"fmt"
	"github.com/awensir/go-aurora/aurora"
)
//真实业务处理
func Test(ctx *aurora.Cxt) interface{} {
	return "test"
}
func main() {
	//获取 aurora 路由实例
	a := aurora.Default()

	// GET 方法注册 web get请求
	a.GET("/", func(next aurora.Serve) aurora.Serve {
		return func(ctx *aurora.Cxt) interface{} {
			fmt.Println("before")
			v := next(ctx)
			fmt.Println("after")
			return v
		}
	}(Test))
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
```

## 第三方框架整合

Aurora约定设置了第三方框架或库的相关实例 的key

```go
package frame

/*
	整合第三方框架标准 key
*/
const (
	GORM     = "gorm"     // gorm    容器数据库连接实例key
	GO_REDIS = "go-redis" // go-redis 容器客户端连接实例key
	RABBITMQ = "RabbitMQ" // rabbit mq 容器客户端连接实例key
	DB       = "db"       // db作为原生 db
)
```

第三方库相关 key 都在 frame包中定义

### Opt 配置选项

源码定义如下

```go
// Opt 配置选项参数
type Opt func() map[string]interface{}
```

Opt 是一个无参函数，返回一个存储interface{}的map，Aurora的配置思想，是约定用户给定配置去初始化第三方库或这框架。

### option包

option 包定义了专有 框架默认读取的配置项，现目前定义的有一下：

```go
package option

const (
		//go-redis 配置项键 （*redis.Options）
	GOREDIS_CONFIG = "go-redis"

	//gorm 数据库类型配置项键 （gorm.Dialector）
	GORM_TYPE = "database" //gorm 数据库类型

	//gorm 配置项选项键 （gorm.Option）
	GORM_CONFIG = "config" //gorm 配置项

	RABBITMQ_URL = "rabbit-url"

	//添加配置 配置项
	Config_key = "name" //定义配置 名
	Config_fun = "func" //定义配置 函数 (type Configuration func(Opt) interface{})
	Config_opt = "opt"  //定义配置 参数选项 (type Opt func() map[string]interface{})

)
```



### gorm 连接配置

func (a *Aurora) GormConfig(opt map[string]interface{}) 方法 配置gorm 默认使用 v2 版本，若要使用其他版本可以通过func (a *Aurora) Store(name string, variable interface{})自行定义k并存储

```go
/*
	整合gorm 框架
	默认使用 v2版本
	提供配置项 初始化默认gorm变量
	需要连接多个库，存放在容器中，实现 manage.Variable 接口 Clone() Variable 方法即可存入容器
*/

//GormConfig 整合gorm
func (a *Aurora) GormConfig(opt Opt) {
	o := opt()
	//读取配置项
	dil, b := o[option.GORM_TYPE].(gorm.Dialector)
	if !b {
		panic(errors.New("gorm config option gorm.Dialector type error！"))
	}

	config, b := o[option.GORM_CONFIG].(gorm.Option)
	if !b {
		panic(errors.New("gorm config option gorm.Option type error！"))
	}
	db, err := gorm.Open(dil, config)
	if err != nil {
		panic(err.Error())
	}
	a.container.store(frame.GORM, db)
}
```

GormConfig 的opt 配置参数，注意按照 const  常量来指定配置项。

#### gorm 配置示例

```go
func main() {

	//获取 aurora 路由实例
	a := aurora.New()

	// 加载 gorm 配置
	a.GormConfig(func() map[string]interface{}{
		return map[string]interface{}{
			option.GORM_TYPE: mysql.Open("root:duanzhiwen@tcp(82.157.160.117:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"),
			option.GORM_CONFIG: &gorm.Config{},
		}
	})

	// GET 方法注册 web get请求
	a.GET("/query", func(c *aurora.Ctx) interface{} {
		db:= a.Get(frame.GORM).(*gorm.DB)
		var s []Student
		db.Raw("select * from student").Scan(&s)
		return s
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
```



### redis 连接配置

redis 客户端，默认采用go-redis 框架进行配置。

```go
package aurora

import (
	"errors"
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/awensir/go-aurora/aurora/option"
	"github.com/go-redis/redis/v8"
)

// GoRedisConfig 根据配置项配置 go-redis
func (a *Aurora) GoRedisConfig(opt Opt) {
	if opt == nil {
		panic(errors.New("go-redis config option not find"))
	}
	o := opt()
	r := redis.NewClient(o[option.GOREDIS_CONFIG].(*redis.Options))
	a.container.store(frame.GO_REDIS, r)
}
```

#### redis 配置示例

```go
func main() {

	//获取 aurora 路由实例
	a := aurora.New()

	//加载 redis 配置
	a.GoRedisConfig(func() map[string]interface{} {
		return map[string]interface{}{
			option.GOREDIS_CONFIG: &redis.Options{
				Addr:     "xx.xxx.xxx.xxx:6379",
				Password: "xxxxxxx",
				DB:       0,
			},
		}
	})
	// GET 方法注册 web get请求
	a.GET("/set", func(c *aurora.Ctx) interface{} {
		client := a.Get(frame.GO_REDIS).(*redis.Client)
		if err := client.Set(context.TODO(), "name", "test", 0).Err(); err != nil {
			return err
		}
		return "ok!"
	})
	a.GET("/get", func(c *aurora.Ctx) interface{} {
		client := a.Get(frame.GO_REDIS).(*redis.Client)
		result, err := client.Get(context.TODO(), "name").Result()
		if err != nil {
			c.ERROR(err.Error())
		}
		return result
	})
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}
```



### RabbitMQ 连接配置

```go
package aurora

import (
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/awensir/go-aurora/aurora/option"
	"github.com/streadway/amqp"
	"log"
)

// RabbitMqConfig 链接RabbitMQ address 链接地址
//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
func (a *Aurora) RabbitMqConfig(opt Opt) {
	o:=opt()
	conn, err := amqp.Dial(o[option.RABBITMQ_URL].(string))
	failOnError(err, "Failed to connect to RabbitMQ")
	a.container.store(frame.RABBITMQ, conn)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

```

### grpc配置整合

#### 定义测试proto

```protobuf
syntax="proto3";

option go_package="../pb";

message DefaultRequest {
  int32   a=1;
  int64   b=2;
  float   c=3;
  double  g=7;
  bool    d=4;
  uint32  e=5;
  uint64  f=6;
  string  h=8;
  bytes   i=9;
  map<string,string> jmap=10;
  repeated string  k=11;
}


message DefaultResponse {
   string result =1;
}

service Default {
  rpc TestDefault(DefaultRequest) returns (DefaultResponse);
}
```

#### 实现服务定义

```go
package pb

import "context"

type Default struct {
}

func (d Default) TestDefault(ctx context.Context, request *DefaultRequest) (*DefaultResponse, error) {

	return &DefaultResponse{Result: "success"},nil
}

func (d Default) mustEmbedUnimplementedDefaultServer() {
	panic("implement me")
}
```

#### 定义服务端

```go
func main() {
	a := aurora.New()

	//部署grpc
	c, _ := credentials.NewServerTLSFromFile("ca/rootcert.pem", "ca/rootkey.pem")
	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterDefaultServer(server,&pb.Default{})
	a.GrpcServer=server

	a.GET("/server", func(c *aurora.Ctx) interface{} {

		return "test"
	})
	a.GuideTLS("ca/rootcert.pem","ca/rootkey.pem","8089")
}
```

#### 定义客户端

```go
func main() {
	a := aurora.New()

	a.GET("/client", func(c *aurora.Ctx) interface{} {
		file, err := credentials.NewClientTLSFromFile("ca/rootcert.pem","localhost")
		if err != nil {
			return err
		}
		dial, err := grpc.Dial("127.0.0.1:8089", grpc.WithTransportCredentials(file))
		if err != nil {
			return err
		}

		client := pb.NewDefaultClient(dial)
		testDefault, err := client.TestDefault(context.TODO(), &pb.DefaultRequest{})
		if err != nil {
			return err
		}
		return testDefault.Result
	})
	a.GuideTLS("ca/rootcert.pem","ca/rootkey.pem","8088")
}
```

通过把grpc的 server 提交给aurora 进行 配置，达到服务端可以同时使用一个端口号进行https和grpc远程调用通讯。

## 自定义配置

第三方框架的整合，就是 做一个管理，支持简单的配置。管理实现基于map

```go
package aurora

import (
	"github.com/awensir/Aurora/aurora/frame"
	"sync"
)

type containers struct {
	rw         *sync.RWMutex
	prototypes map[string]interface{} //容器存储的属性
}

// Get 获取容器内指定变量，不存在的key 返回默认零值
func (c *containers) get(name string) interface{} {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.prototypes[name]
}

// Store 加载 指定变量
func (c *containers) store(name string, variable interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.prototypes[name] = variable
}

// Delete 提供删除自定义整合数据
func (c *containers) Delete(name string) {
	switch name {
		// 内置整合 不允许删除
		case frame.DB, frame.GORM, frame.GO_REDIS:
			return
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, b := c.prototypes[name]; !b {
		return
	} else {
		delete(c.prototypes, name)
	}
}
```

也可以通过容器管理托管自己的公有变量等....

```go
// Store 加载
func (a *Aurora) Store(name string, variable interface{}) {
	a.container.store(name, variable)
}
```

## Pool

基于sync.Pool 实现一个对全局实例，提供线程安全的使用操作，使用Pool 加载第三方配置实例 和 自定义的配置 是完全不同的方式，自定义加载的实例变量是可以多线程同时使用的，其变量实例自己做了对并发情况下的安全支持可以选择。

当需要使用Pool方式实现多线程共享安全实例可以通过以下配置实现：

```go
/*
	LoadConfiguration 加载自定义配置项，
	Opt 必选配置项：
	Config_key ="name"	定义配置 名，
	Config_fun ="func"	定义配置 函数，
	Config_opt ="opt"	定义配置 参数选项
*/
func (a *Aurora) LoadConfiguration(options Opt) {
	a.lock.Lock()
	defer a.lock.Unlock()
	o := options()
	key, b := o[option.Config_key].(string)
	opt, b := o[option.Config_opt].(Opt)
	fun, b := o[option.Config_fun].(Configuration)
	if !b {
		//配置选项出现问题
		return
	}
	a.options[key] = &Option{
		opt,
		fun,
	}
}
```



```go
// GetPool 获取容器池中的实例，前提是通过 Pool 向池中加载了。 GetPool 和 PutPool 必须成对出现
func (a *Aurora) GetPool(name string) interface{}

// PutPool 把取出来的 实例返回放入到池中，以便下次使用
func (a *Aurora) PutPool(name string, v interface{})
```



## 全局配置文件管理

Aurora 会默认加载项目根目录下的 application.yml 配置文件。文件不存在则不做任何处理，文件存在就读取并初始化 viper 配置变量。

示例（读取默认位置）：

```go
//读取项目根目录位置
a.ViperConfig()
```

application.yml不在默认位置，可以通过传递文件夹参数进行指定。application.yml 配置文件及其类型和名称是严格读取，不能更改的。

示例：

```go
//（基于项目根目录）读取指定位置  static/config/application.yml  
a.ViperConfig("static","config")
```

指定位置不存在仅提示错误信息
