package aurora

import (
	"net/http"
	"strings"
)

/*
	路由分组
*/

type Group struct {
	a      *Aurora
	prefix string
}

// GET 请求
func (g *Group) GET(path string, servlet Servlet) {
	g.a.Register(http.MethodGet, g.prefix+path, servlet)
}

// POST 请求
func (g *Group) POST(path string, servlet Servlet) {
	g.a.Register(http.MethodPost, g.prefix+path, servlet)
}

// PUT 请求
func (g *Group) PUT(path string, servlet Servlet) {
	g.a.Register(http.MethodPut, g.prefix+path, servlet)
}

// DELETE 请求
func (g *Group) DELETE(path string, servlet Servlet) {
	g.a.Register(http.MethodDelete, g.prefix+path, servlet)
}

// HEAD 请求
func (g *Group) HEAD(path string, servlet Servlet) {
	g.a.Register(http.MethodHead, g.prefix+path, servlet)
}

// Group 路由分组  必须以 “/” 开头分组
func (g *Group) Group(path string) *Group {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return &Group{
		prefix: g.prefix + path,
		a:      g.a,
	}
}
