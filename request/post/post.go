package post

import (
	"Aurora/aurora"
	"net/http"
)

func Mapping(url string,fun aurora.ServletHandler)  {
	aurora.RegisterServlet(http.MethodPost,url,fun)
}
