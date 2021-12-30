package aurora

/*
	提供一个默认的 aurora 实例进行快捷的使用 web开发

*/

var defaultAurora = New()

func Config(config ...string) {
	defaultAurora = New(config...)
}

func GET(url string, serve Serve) {
	defaultAurora.GET(url, serve)
}

func POST(url string, serve Serve) {
	defaultAurora.POST(url, serve)
}

func Start(port ...string) error {
	return defaultAurora.Guide(port...)
}

func StartTLS(args ...string) error {
	return defaultAurora.GuideTLS(args...)
}
