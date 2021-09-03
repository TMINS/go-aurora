package message

type Message interface {
	 ToString() string
}

type UrlRegisterInfo struct {
	 Method string
	 Path   string
}

func (i UrlRegisterInfo) ToString() string {
	return i.Method+":"+i.Path+" 服务注册成功~~~~"
}

type ResourceInfo struct {
	 Type string
	 Path string
}

func (r ResourceInfo) ToString() string {
	return r.Type+" "+r.Path+" 静态资源路径映射添加成功！"
}

type StartSuccessful struct {
	 Port string
}

func (s StartSuccessful) ToString() string {
	return "port "+s.Port+" Successful binding!"
}
