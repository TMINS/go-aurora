package config

import (
	"Aurora/aurora"
)

/*
    提供对其它包中的可配置接口
*/

// RegisterInterceptor 添加拦截器,暂时只支持全局拦截器
func RegisterInterceptor(interceptor... aurora.Interceptor)  {
	aurora.InterceptorList=interceptor
}

// RegisterResource 添加静态资源配置，t资源类型必须以置源后缀命名，
//paths为t类型资源的子路径，可以一次性设置多个。
//每个资源类型最调用一次设置方法否则覆盖原有设置
func RegisterResource(t string,paths ...string)  {
	aurora.RegisterResourceType(t,paths...)
}

// ResourceRoot 设置静态资源根路径
func ResourceRoot(root string) {
	aurora.SetResourceRoot(root)
}
