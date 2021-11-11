package aurora

import (
	"errors"
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/awensir/go-aurora/aurora/option"
	"github.com/go-redis/redis/v8"
)

// GoRedisConfig 根据配置项配置 go-redis
func (a *Aurora) GoRedisConfig(opt Opt) {
	if opt == nil {
		panic(errors.New("go-redis config option not find"))
	}
	o := opt()
	r := redis.NewClient(o[option.GOREDIS_CONFIG].(*redis.Options))
	a.container.store(frame.GO_REDIS, r)
}
