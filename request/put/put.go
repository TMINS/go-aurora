package put

import (
	"github.com/awensir/Aurora/aurora"
	"net/http"
)

func Mapping(url string, fun aurora.Servlet) {
	aurora.RegisterServlet(http.MethodPut, url, fun)
}
