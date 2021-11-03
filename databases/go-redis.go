package databases

import (
	"errors"
	"github.com/awensir/Aurora/manage"
	"github.com/awensir/Aurora/manage/frame"
	"github.com/go-redis/redis/v8"
)

const (
	GOREDIS_CNF = "GO_REDIS_CNF"
)

type GO_REDIS struct {
	*redis.Client
}

func (gr *GO_REDIS) Clone() manage.Variable {
	return gr
}

// GoRedisConfig 根据配置项配置 go-redis
func GoRedisConfig(opt *redis.Options) {
	if opt == nil {
		panic(errors.New("go-redis config option not find"))
	}
	r := redis.NewClient(opt)
	goredis := &GO_REDIS{r}
	manage.Container.Store(frame.GO_REDIS, goredis)
}
