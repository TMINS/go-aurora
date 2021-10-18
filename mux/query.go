package mux

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
func (c *Ctx) Get(Args string) (string, error) {
	c.monitor.En(ExecuteInfo(nil))
	return c.get(Args, c.Request.URL.Query())
}

// GetInt 获取一个整数参数
func (c *Ctx) GetInt(Args string) (int, error) {
	c.monitor.En(ExecuteInfo(nil))
	return c.getInt(Args, c.Request.URL.Query())
}

// GetFloat64 获取一个64位浮点参数
func (c *Ctx) GetFloat64(Args string) (float64, error) {
	c.monitor.En(ExecuteInfo(nil))
	return c.getFloat64(Args, c.Request.URL.Query())
}

// GetSlice 获取切片类型参数
func (c *Ctx) GetSlice(Args string) ([]string, error) {
	c.monitor.En(ExecuteInfo(nil))
	return c.getSlice(Args, c.Request.URL.Query())
}

// GetIntSlice 整数切片
func (c *Ctx) GetIntSlice(Args string) ([]int, error) {
	c.monitor.En(ExecuteInfo(nil))
	return c.getIntSlice(Args, c.Request.URL.Query())
}

// GetFloat64Slice 浮点切片
func (c *Ctx) GetFloat64Slice(Args string) ([]float64, error) {
	c.monitor.En(ExecuteInfo(nil))
	return c.getFloat64Slice(Args, c.Request.URL.Query())
}

func (c *Ctx) getFloat64Slice(Args string, values url.Values) ([]float64, error) {
	arr, err := c.getValues(Args, values)
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

func (c *Ctx) getIntSlice(Args string, values url.Values) ([]int, error) {
	arr, err := c.getValues(Args, values)
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

func (c *Ctx) getSlice(Args string, values url.Values) ([]string, error) {
	return c.getValues(Args, values)
}

func (c *Ctx) get(Args string, values url.Values) (string, error) {
	arr, err := c.getValue(Args, values)
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return "", err
	}
	return arr[0], nil
}

func (c *Ctx) getFloat64(Args string, values url.Values) (float64, error) {
	arr, err := c.getValue(Args, values)
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return 0, err
	}
	f, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return 0, err
	}
	return f, nil
}

func (c *Ctx) getInt(Args string, values url.Values) (int, error) {
	arr, err := c.getValue(Args, values)
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return 0, err
	}
	a, err := strconv.Atoi(arr[0])
	if err != nil {
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return 0, err
	}
	return a, nil
}

func (c *Ctx) getValue(Args string, values url.Values) ([]string, error) {
	arr, b := values[Args]
	if !b {
		err := errors.New("Query Param Not Exist return 0 ")
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return nil, err
	}
	if len(arr) != 1 {
		err := errors.New("Query Param Not Exist return 0 ")
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return nil, err
	}
	return arr, nil
}

func (c *Ctx) getValues(Args string, values url.Values) ([]string, error) {
	arr, b := values[Args]
	if !b {
		err := errors.New("Query Param Not Exist return 0 ")
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return nil, err
	}
	if len(arr) != 1 {
		err := errors.New("Query Param Not Exist return 0 ")
		c.monitor.En(ExecuteInfo(err))
		c.ar.runtime <- c.monitor
		return nil, err
	}
	split := strings.Split(arr[0], ",")
	return split, nil
}
