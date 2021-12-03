package aurora

import (
	"encoding/json"
	"time"
)

/*
	错误处理
*/

// WebError 业务处理期间，的特定错误
type WebError interface {
	ErrorHandler(c *Ctx) interface{} //ErrorHandler 处理对应的错误
}

type WebErr func(c *Ctx) interface{}

type ErrorResponse struct {
	UrlPath      string `json:"url"`
	Status       int    `json:"code"`
	ErrorMessage string `json:"error"`
	Time         string `json:"time"`
}

func newErrorResponse(path, message string, status int) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	msg := ErrorResponse{
		UrlPath:      path,
		Status:       status,
		ErrorMessage: message,
		Time:         now,
	}
	marshal, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}
