package sys

import (
	"github.com/EricJSanchez/gotool/pkg/environment"
	"github.com/redis/go-redis/v9"
)

var redisManager = NewRedisClientManager()

func Redis(names ...string) (client *redis.Client) {
	var name = Cfg("app").GetString("default_redis")

	if len(names) > 0 {
		name = names[0]
	}
	var redisConfig map[string]interface{}
	if environment.Is(environment.Development) {
		redisConfig = Cfg("db").GetStringMap(name)
		if len(redisConfig) == 0 {
			redisConfig = Nacos("database.toml").GetStringMap(name)
		}
	} else {
		redisConfig = Nacos("database.toml").GetStringMap(name)
	}
	pool := redisManager.Get(name, redisConfig)
	client = pool.pool
	return
}
