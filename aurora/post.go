package aurora

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
)

type PostValue struct {
	Value map[string]interface{}
	File  map[string][]*multipart.FileHeader
}

// Body 读取Post请求体Json或表单数据数据解析到body中,仅限于单个post请求体，不能用于文件上传和请求体并存
func (c *Ctx) Body(body interface{}) error {
	data, err := io.ReadAll(c.Request.Body)
	defer func(Body io.ReadCloser, ctx *Ctx) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(c.Request.Body, c)

	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		return err
	}
	return nil
}

// MapBody 直接获取map格式的请求体,仅限于单个post请求体，不能用于文件上传和请求体并存
func (c *Ctx) MapBody() (map[string]interface{}, error) {
	data, err := io.ReadAll(c.Request.Body)
	defer func(Body io.ReadCloser, ctx *Ctx) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(c.Request.Body, c)
	if err != nil {
		return nil, err
	}
	var body interface{}
	err = json.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		return nil, errors.New("post解析失败")
	}
	return body.(map[string]interface{}), nil
}

// FileBody 主要用于 文件上传和请求体同时存在
func (c *Ctx) FileBody() (*PostValue, error) {
	err := c.Request.ParseMultipartForm(c.ar.MaxMultipartMemory)
	if err != nil {
		return nil, err
	}
	form := c.Request.MultipartForm
	args, err := parse(form.Value)
	if err != nil {
		return nil, err
	}
	p := &PostValue{
		args,
		c.Request.MultipartForm.File,
	}
	return p, nil
}

//解析 post 和 文件上传中的消息体
func parse(value map[string][]string) (map[string]interface{}, error) {
	args := make(map[string]interface{})
	for k, v := range value {
		if len(v) > 1 {
			return nil, errors.New("parameter parsing error, duplicate request parameter name")
		}
		s := v[0]

		//解析 参数是否为数组
		if s[:1] == "[" && s[len(s)-1:] == "]" {
			//去除前后的中括号
			s = s[1 : len(s)-1]
			split := strings.Split(s, ",")
			args[k] = split
			continue
		}
		//解析json参数体
		if s[:1] == "{" && s[len(s)-1:] == "}" {
			var data interface{}
			err := json.Unmarshal([]byte(s), &data)
			if err != nil {
				return nil, err
			}
			args[k] = data
			continue
		}
		args[k] = s
	}
	return args, nil
}

// String 获取string类型参数
func (pv *PostValue) String(name string) (string, error) {
	if _, b := pv.Value[name]; !b {
		return "", errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	s, b := v.(string)
	if !b {
		return "", errors.New("the query parameter is not a string type")
	}
	return s, nil
}

func (pv *PostValue) Int(name string) (int, error) {
	if _, b := pv.Value[name]; !b {
		return 0, errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	s, b := v.(int)
	if !b {
		return 0, errors.New("the query parameter is not a int type")
	}
	return s, nil
}

func (pv *PostValue) Float64(name string) (float64, error) {
	if _, b := pv.Value[name]; !b {
		return 0, errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	s, b := v.(float64)
	if !b {
		return 0, errors.New("the query parameter is not a float64 type")
	}
	return s, nil
}

func (pv *PostValue) Slice(name string) ([]string, error) {
	if _, b := pv.Value[name]; !b {
		return nil, errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	i, b := v.([]string)
	if !b {
		return nil, errors.New("the query parameter is not a []string type")
	}
	return i, nil
}
func (pv *PostValue) IntSlice(name string) ([]int, error) {
	if _, b := pv.Value[name]; !b {
		return nil, errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	i, b := v.([]string)
	if !b {
		return nil, errors.New("the query parameter is not a []int type")
	}
	arr := make([]int, len(i))
	for _, a := range i {
		quote, err := strconv.Atoi(a)
		if err != nil {
			return nil, errors.New("incorrect type conversion format")
		}
		arr = append(arr, quote)
	}
	return arr, nil
}

func (pv *PostValue) Float64Slice(name string) ([]float64, error) {
	if _, b := pv.Value[name]; !b {
		return nil, errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	i, b := v.([]string)
	if !b {
		return nil, errors.New("the query parameter is not a []float64 type")
	}

	arr := make([]float64, len(i))
	for _, a := range i {
		quote, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return nil, errors.New("incorrect type conversion format")
		}
		arr = append(arr, quote)
	}
	return arr, nil
}

// MapBody 主要用于获取传递的json参数
func (pv *PostValue) MapBody(name string) (map[string]interface{}, error) {
	if _, b := pv.Value[name]; !b {
		return nil, errors.New("query parameter does not exist")
	}
	v := pv.Value[name]
	i, b := v.(map[string]interface{})
	if !b {
		return nil, errors.New("the query parameter is not a map[string]interface{} type")
	}
	return i, nil
}
