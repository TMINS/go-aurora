package main

import (
	"context"
	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/awensir/go-aurora/aurora/option"
	"github.com/go-redis/redis/v8"
)

func main() {

	//获取 aurora 路由实例
	a := aurora.New()

	//加载 redis 配置
	a.GoRedisConfig(func() map[string]interface{} {
		return map[string]interface{}{
			option.GOREDIS_CONFIG: &redis.Options{
				Addr:     "82.157.160.117:6379",
				Password: "duanzhiwen",
				DB:       0,
			},
		}
	})
	// GET 方法注册 web get请求
	a.GET("/set", func(c *aurora.Ctx) interface{} {
		client := a.Get(frame.GO_REDIS).(*redis.Client)
		if err := client.Set(context.TODO(), "name", "test", 0).Err(); err != nil {
			return err
		}
		return "ok!"
	})

	a.GET("/get", func(c *aurora.Ctx) interface{} {
		client := a.Get(frame.GO_REDIS).(*redis.Client)
		result, err := client.Get(context.TODO(), "name").Result()
		if err != nil {
			c.ERROR(err.Error())
		}
		return result
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
