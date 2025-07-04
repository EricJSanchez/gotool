package sys

import (
	"encoding/json"
	"github.com/EricJSanchez/gotool/environment"
	"gorm.io/gorm"
)

var gormManager = NewGormClientManager()

func Gorm(names ...string) (client *gorm.DB) {
	var name = Cfg("app").GetString("default_db")

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
	client = gormManager.Get(name, config)
	return
}
