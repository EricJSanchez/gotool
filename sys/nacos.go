package sys

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

var NacosConfig map[string]string

func Nacos(files ...string) *viper.Viper {
	var file string
	if len(files) == 0 {
		file = Cfg("app").GetString("nacos.defaultDataId")
	} else {
		file = files[0]
	}
	if file == "" {
		fmt.Println("获取 naocs 失败")
		return nil
	}
	if _, ok := NacosConfig[file]; !ok {
		fmt.Println("read nacos configLocal file err 0")
		return nil
	}
	configLocal.Lock()
	defer configLocal.Unlock()
	if cfg, ok := configLocal.vipers[file]; ok {
		return cfg
	}
	// 读取基础配置
	baseConfig := viper.New()
	baseConfig.SetConfigType("toml")
	err := baseConfig.ReadConfig(bytes.NewBuffer([]byte(NacosConfig[file])))
	if err != nil {
		fmt.Println("read nacos configLocal err 1")
		return nil
	} else {
		configLocal.vipers[file] = baseConfig
		return baseConfig
	}
}
