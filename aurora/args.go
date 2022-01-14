package aurora

import (
	"errors"
	"strconv"
)

// InitArgs 获取Int类型的 REST FUL 参数，参数不存在情况下返回 0 和 err
func (c *Ctx) InitArgs(name string) (int, error) {
	args, err := c.getArgs(name)
	if err != nil {
		return 0, err
	}
	atoi, err := strconv.Atoi(args.(string))
	if err != nil {
		return 0, err
	}
	return atoi, nil
}

// Float64Args 获取Float64类型的 REST FUL 参数
func (c *Ctx) Float64Args(name string) (float64, error) {
	args, err := c.getArgs(name)
	if err != nil {
		return 0, err
	}
	float, err := strconv.ParseFloat(args.(string), 64)
	if err != nil {
		return 0, err
	}
	return float, nil
}

// StringArgs 获取字符串形式的 REST FUL 参数
func (c *Ctx) StringArgs(name string) (string, error) {
	args, err := c.getArgs(name)
	if err != nil {
		return "", err
	}
	return args.(string), nil
}

// BoolArgs 获取 逻辑 形式的 REST FUL 参数
func (c *Ctx) BoolArgs(name string) (bool, error) {
	args, err := c.getArgs(name)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(args.(string))
}

func (c *Ctx) getArgs(name string) (interface{}, error) {
	v, b := c.Args[name]
	if !b {
		return nil, errors.New("参数不存在")
	}
	return v, nil
}
