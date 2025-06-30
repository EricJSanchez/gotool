package sys

import (
	"encoding/json"
	"github.com/EricJSanchez/gotool/environment"
	"github.com/redis/go-redis/v9"
)

var redisManager = NewRedisClientManager()

func Redis(names ...string) (client *redis.Client) {
	var name = Cfg("app").GetString("default_redis")

	if len(names) > 0 {
		name = names[0]
	}
	var config map[string]interface{}
	if environment.Is(environment.Development) {
		config = Cfg("db").GetStringMap(name)
		if len(config) == 0 {
			config = Nacos("database.toml").GetStringMap(name)
		}
	} else {
		config = Nacos("database.toml").GetStringMap(name)
	}
	connectUniq, _ := json.Marshal(config)
	name = name + Md5(string(connectUniq))
	pool := redisManager.Get(name, config)
	client = pool.pool
	return
}
