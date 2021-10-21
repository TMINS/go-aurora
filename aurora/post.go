package aurora

import (
	"encoding/json"
	"fmt"
	"io"
)

// PostBody 读取Post请求体数据解析到body中
func (c *Ctx) PostBody(body interface{}) error {
	data, err := io.ReadAll(c.Request.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err.Error())
		}
	}(c.Request.Body)
	if err == nil {
		err := json.Unmarshal(data, &body)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	}
	return err
}
