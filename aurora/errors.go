package aurora

/*
	错误定义
*/

// ArgsError 获取参数错误,用于解析请求参数或调用请求参数有问题时候产生这个错误
type ArgsError struct {
	Type    string //类型
	Message string
}

func (e ArgsError) Error() string {
	return "ArgsError : " + "Type:" + e.Type + " Message:" + e.Message
}

type UrlPathError struct {
	Type     string
	Method   string
	Message  string
	NodePath string
	Path     string
}

func (e UrlPathError) Error() string {
	return e.Method + ":UrlPathError : " + "Type:" + e.Type + " " + e.Path + "   " + e.NodePath + ". Message:" + e.Message
}

// WebResponseError 业务处理期间，的特定错误
type WebResponseError interface {
	ErrorHandler(ctx *Context) interface{} //ErrorHandler 处理对应的错误
}

type WebErr func(ctx *Context) interface{}