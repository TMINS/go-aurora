package aurora

/*
	封装响应信息
*/

func (c *Ctx) SetStatus(code int) {
	c.Response.WriteHeader(code)
}

// SetHeader 设置响应头，响应头存在则追加，不存在则新添加
func (c *Ctx) SetHeader(key, value string) {
	h := c.Response.Header()
	s := h.Get(key)
	if s == "" {
		h.Set(key, value)
	} else {
		h.Add(key, value)
	}
}

// NewHeader key存在会直接覆盖
func (c *Ctx) NewHeader(key, value string) {
	h := c.Response.Header()
	h.Set(key, value)
}
