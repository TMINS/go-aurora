package delet

import (
	"github.com/awensir/Aurora/aurora"
	"net/http"
)

func Mapping(url string, fun aurora.Servlet) {
	aurora.RegisterServlet(http.MethodDelete, url, fun)
}
