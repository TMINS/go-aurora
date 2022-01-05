package aurora

import "github.com/go-redis/redis/v8"

//初步完成
func (a *Aurora) loadGoRedis() {
	if a.cnf == nil {
		//如果配置文件没有加载成功，将不做任何事情
		return
	}
	configs := a.cnf.Get("aurora.redis.go-redis")
	if configs == nil {
		return
	}
	a.auroraLog.Info("start loading go-redis configuration")
	config, b := configs.([]map[string]interface{})
	if !b {
		return
	}
	for _, v := range config {
		c := redis.NewClient(&redis.Options{
			Addr:     v["addr"].(string),
			Password: v["password"].(string), // no password set
			DB:       v["db"].(int),          // use default DB
		})
		a.goredis = append(a.goredis, c)
	}
}
