package aurora

import (
	"context"

	"github.com/go-redis/redis/v8"
)

//go-redis 当前只完成了单一的客户端配置,支持同时配置多个客户端，go-redis本身对redis的集群等高可用方式提供了实现，目前还没有对这些进行整合

//初步完成   https://redis.uptrace.dev/
func (a *Aurora) loadGoRedis() {
	if a.config == nil {
		//如果配置文件没有加载成功，将不做任何事情
		return
	}
	configs := a.config.Get("aurora.redis.go-redis")
	if configs == nil {
		return
	}
	a.auroraLog.Info("start loading go-redis configuration")
	config, b := configs.([]map[string]interface{})
	if !b {
		return
	}
	for i, v := range config {
		c := redis.NewClient(&redis.Options{
			Addr:     v["addr"].(string),
			Password: v["password"].(string), // no password set
			DB:       v["db"].(int),          // use default DB
		})
		//检测 redis 连接
		ping := c.Ping(context.TODO())
		if err := ping.Err(); err != nil {
			a.auroraLog.Error("index : ", i, ",", err.Error()) //如果第 i个配置出现问题,如何处理这个问题，待解决
		}
		a.goredis = append(a.goredis, c)
	}
}

// GoRedis 获取默认的go redis 客户端
func (a *Aurora) GoRedis() *redis.Client {
	if len(a.goredis) < 1 {
		return nil
	}
	return a.goredis[0]
}

func (a *Aurora) GoRedisList(index int) *redis.Client {
	if len(a.goredis) <= index {
		return nil
	}
	return a.goredis[index]
}
