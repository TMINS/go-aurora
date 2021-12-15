package aurora

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

/*
	静态资源处理
	浏览器请求静态资源的方式是Get，html中引入的资源会被一起发送资源请求服务，请求路径则是，导入资源的路径比如

	假定有以下静态资源目录结构
	/			root目录
	/static
		/js/    存放js文件
		/css/   存放css文件
		/html/  存放html文件
	在html文件中 正确的映入方式应是   ../js/xxx.js 或者  ../css/xxx.css

	浏览器如何对服务器资源进行请求
		1.当一个请求返回的是一个服务器html页面，浏览器接到响应解析请求头，会根据之前发送请求的url对返回的页面进行构建一个和服务器内部静态资源存储相同目录结构
		2.根据浏览器生成的目录结构，浏览器解析到html上面有导入资源，会自动携带url这个信息去查找服务器上的资源
		3.服务器需要解析结构得到正确存储路径，才能够响应给请求者
		4.得到静态资源，会把该资源放到构建好的目录中，以便html能够正确引入资源

	如此我们需要配置一个专门处理这一类请求的服务
	Golang Web 默认处理静态资源 是通过写入的方式
*/

const contentType = "Content-Type"
const favicon = "favicon.ico"

// Views 是整个服务器对视图渲染的核心接口,开发者实现改接口对需要展示的页面进行自定义处理
type Views interface {
	View(*Ctx, string)
}

type ViewFunc func(*Ctx, string)

func (vf ViewFunc) View(c *Ctx, p string) {
	vf(c, p)
}

// ResourceFun w 响应体，path 资源真实路径，rt资源类型
// 根据rt资源类型去找到对应的resourceMapType 存储的响应头，进行发送资源
func (a *Aurora) resourceFun(w http.ResponseWriter, mapping string, path string, rt string) {
	data := a.readResource(a.projectRoot + "/" + a.resource + path)
	if data != nil {
		h := w.Header()
		if h.Get(contentType) == "" {
			h.Set(contentType, a.resourceMapType[rt])
		}
		h.Add(contentType, a.resourceMapType[rt])
		sendResource(w, data)
	}
	//http.NotFound(w,req)  //没有读取到对应静态资源，发送404
	http.Error(w, newErrorResponse(mapping, "The server static resource does not exist", 500), 500)
}

// SendResource 发送静态资源
func sendResource(w http.ResponseWriter, data []byte) {
	if data == nil {
		return
	}
	_, err := w.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// readResource 读取成功则返回结果，失败则返回nil
func (a *Aurora) readResource(path string) []byte {
	if f, err := ioutil.ReadFile(path); err == nil {
		return f
	} else {
		if os.IsNotExist(err) {
			fmt.Println(err.Error())
		}
	}
	return nil
}

// RegisterResourceType 加载静态资源路径，静态资源读取路径，服务器处理静态资源策略改为ServeHTTP处判别，最终静态资源的处理取决于 resource 根的设置
//存在不同的图片类型需要多次调用设置对应的存储路径（图片类型存在不同，待解决）
func (a *Aurora) registerResourceType(t string, paths ...string) {
	if a.resourceMappings == nil {
		a.resourceMappings = make(map[string][]string)
	}
	for i := 0; i < len(paths); i++ {
		pl := len(paths[i])
		if paths[i][:1] != "/" {
			paths[i] = "/" + paths[i]
		}
		if paths[i][pl-1:] != "/" {
			paths[i] = paths[i] + "/"
		}
	}
	a.resourceMappings[t] = paths
}

func loadResourceHead(a *Aurora) {
	a.resourceMapType["js"] = "text/javascript"
	a.resourceMapType["css"] = "text/css"
	a.resourceMapType["html"] = "text/html"
	a.resourceMapType["encoding"] = "charset=utf-8"
	a.resourceMapType["gif"] = "image/gif"
	a.resourceMapType["png"] = "image/png"
	a.resourceMapType["svg"] = "image/svg+xml"
	a.resourceMapType["webp"] = "image/webp"
	a.resourceMapType["ico"] = "image/x-icon"
}
