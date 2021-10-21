package aurora

/*
	错误处理
*/

// WebError 业务处理期间，的特定错误
type WebError interface {
	ErrorHandler(c *Ctx) interface{} //ErrorHandler 处理对应的错误
}

type WebErr func(c *Ctx) interface{}
