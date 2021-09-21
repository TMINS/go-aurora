package config

import (
	"github.com/awensir/Aurora/aurora"
)

/*
   提供对其它包中的可配置接口
*/

// Interceptor 添加全局拦截器
func Interceptor(interceptor ...aurora.Interceptor) {
	aurora.RegisterInterceptorList(interceptor...)
}

// DefaultInterceptor 修改默认拦截器
func DefaultInterceptor(interceptor aurora.Interceptor) {
	aurora.RegisterDefaultInterceptor(interceptor)
}

// PathInterceptor 注册局部拦截器
func PathInterceptor(path string, interceptor ...aurora.Interceptor) {
	aurora.RegisterInterceptor(path, interceptor...)
}

// Resource 添加静态资源配置，t资源类型必须以置源后缀命名，
//paths为t类型资源的子路径，可以一次性设置多个。
//每个资源类型最调用一次设置方法否则覆盖原有设置
func Resource(Type string, Paths ...string) {
	aurora.RegisterResourceType(Type, Paths...)
}

// ResourceRoot 设置静态资源根路径
func ResourceRoot(root string) {
	aurora.SetResourceRoot(root)
}
