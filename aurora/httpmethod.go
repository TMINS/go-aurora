package aurora

import (
	"net/http"
	"strings"
)

// GET 请求
func (a *Aurora) GET(path string, servlet Serve) {
	a.register(http.MethodGet, path, servlet)
}

// POST 请求
func (a *Aurora) POST(path string, servlet Serve) {
	a.register(http.MethodPost, path, servlet)
}

// PUT 请求
func (a *Aurora) PUT(path string, servlet Serve) {
	a.register(http.MethodPut, path, servlet)
}

// DELETE 请求
func (a *Aurora) DELETE(path string, servlet Serve) {
	a.register(http.MethodDelete, path, servlet)
}

// HEAD 请求
func (a *Aurora) HEAD(path string, servlet Serve) {
	a.register(http.MethodHead, path, servlet)
}

// register 通用注册器
func (a *Aurora) register(method string, mapping string, fun Serve) {

	a.router.addRoute(method, mapping, fun)
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
