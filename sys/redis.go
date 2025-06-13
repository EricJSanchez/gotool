package sys

import "github.com/redis/go-redis/v9"

var redisManager = NewRedisClientManager()

func Redis(names ...string) (client *redis.Client) {
	var name = Cfg("app").GetString("default_redis")

	if len(names) > 0 {
		name = names[0]
	}

	redisConfig := Nacos("database.toml").GetStringMap(name)
	pool := redisManager.Get(name, redisConfig)
	client = pool.pool
	return
}
