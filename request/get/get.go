package get

import (
	"github.com/awensir/Aurora/aurora"
	"net/http"
)

// Mapping GET请求注册器
func Mapping(url string, fun aurora.Servlet) {
	aurora.RegisterServlet(http.MethodGet, url, fun)
}
