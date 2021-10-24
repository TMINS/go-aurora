package aurora

import (
	"encoding/json"
	"io"
)

// JsonBody 读取Post请求体Json或表单数据数据解析到body中,
func (c *Ctx) JsonBody(body interface{}) error {
	data, err := io.ReadAll(c.Request.Body)
	defer func(Body io.ReadCloser, ctx *Ctx) {
		err := Body.Close()
		if err != nil {
			c.monitor.En(executeInfo(err))
			c.ar.runtime <- c.monitor
		}
	}(c.Request.Body, c)
	if err == nil {
		err := json.Unmarshal(data, &body)
		if err != nil {
			c.monitor.En(executeInfo(err))
			c.ar.runtime <- c.monitor
			return err
		}
		return nil
	}
	c.monitor.En(executeInfo(err))
	c.ar.runtime <- c.monitor
	return err
}
