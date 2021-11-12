package aurora

import (
	"errors"
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/go-redis/redis/v8"
)

const (
	GOREDIS_CNF = "GO_REDIS_CNF" //go-redis 配置key
)

// GoRedisConfig 根据配置项配置 go-redis
func (a *Aurora) GoRedisConfig(opt *redis.Options) {
	if opt == nil {
		panic(errors.New("go-redis config option not find"))
	}
	r := redis.NewClient(opt)
	a.container.store(frame.GO_REDIS, r)
}
