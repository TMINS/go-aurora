package post

import (
	"github.com/awensir/Aurora/aurora"
	"net/http"
)

func Mapping(url string, fun aurora.Servlet) {
	aurora.RegisterServlet(http.MethodPost, url, fun)
}
