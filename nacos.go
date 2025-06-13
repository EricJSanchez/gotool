package gotool

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
		fmt.Println("read nacos config file err 0")
		return nil
	}
	config.Lock()
	defer config.Unlock()
	if cfg, ok := config.vipers[file]; ok {
		return cfg
	}
	// 读取基础配置
	baseConfig := viper.New()
	baseConfig.SetConfigType("toml")
	err := baseConfig.ReadConfig(bytes.NewBuffer([]byte(NacosConfig[file])))
	if err != nil {
		fmt.Println("read nacos config err 1")
		return nil
	} else {
		config.vipers[file] = baseConfig
		return baseConfig
	}
}
