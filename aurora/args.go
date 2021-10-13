package aurora

// GetArgs 获取REST API 参数，查询不存在的key或者不存在REST API 参数则返回""
func (c *Context) GetArgs(key string) interface{} {
	return c.args
}
