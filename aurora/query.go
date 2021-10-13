package aurora

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

/*
	get 请求参数查询
*/

// Get 获取一个字符串参数
func (c *Context) Get(Args string) (string, error) {
	return get(Args, c.Request.URL.Query())
}

// GetInt 获取一个整数参数
func (c *Context) GetInt(Args string) (int, error) {
	return getInt(Args, c.Request.URL.Query())
}

// GetFloat64 获取一个64位浮点参数
func (c *Context) GetFloat64(Args string) (float64, error) {
	return getFloat64(Args, c.Request.URL.Query())
}

// GetSlice 获取切片类型参数
func (c *Context) GetSlice(Args string) ([]string, error) {
	return getSlice(Args, c.Request.URL.Query())
}

// GetIntSlice 整数切片
func (c *Context) GetIntSlice(Args string) ([]int, error) {
	return getIntSlice(Args, c.Request.URL.Query())
}

// GetFloat64Slice 浮点切片
func (c *Context) GetFloat64Slice(Args string) ([]float64, error) {
	return getFloat64Slice(Args, c.Request.URL.Query())
}

func getFloat64Slice(Args string, values url.Values) ([]float64, error) {
	arr, err := getValues(Args, values)
	if err != nil {
		return nil, err
	}
	a := make([]float64, 0)
	for _, v := range arr {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		a = append(a, f)
	}
	return a, err
}

func getIntSlice(Args string, values url.Values) ([]int, error) {
	arr, err := getValues(Args, values)
	if err != nil {
		return nil, err
	}
	a := make([]int, 0)
	for _, v := range arr {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		a = append(a, atoi)
	}
	return a, err
}

func getSlice(Args string, values url.Values) ([]string, error) {
	return getValues(Args, values)
}

func get(Args string, values url.Values) (string, error) {
	arr, err := getValue(Args, values)
	if err != nil {
		return "", err
	}
	return arr[0], nil
}

func getFloat64(Args string, values url.Values) (float64, error) {
	arr, err := getValue(Args, values)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func getInt(Args string, values url.Values) (int, error) {
	arr, err := getValue(Args, values)
	if err != nil {
		return 0, err
	}
	a, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, err
	}
	return a, nil
}

func getValue(Args string, values url.Values) ([]string, error) {
	arr, b := values[Args]
	if !b {
		return nil, errors.New("Query Param Not Exist return 0 ")
	}
	if len(arr) != 1 {
		return nil, errors.New("Query Param Not Exist return 0 ")
	}
	return arr, nil
}

func getValues(Args string, values url.Values) ([]string, error) {
	arr, b := values[Args]
	if !b {
		return nil, errors.New("Query Param Not Exist return 0 ")
	}
	if len(arr) != 1 {
		return nil, errors.New("Query Param Not Exist return 0 ")
	}
	split := strings.Split(arr[0], ",")
	return split, nil
}
