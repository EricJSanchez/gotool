package sys

import (
	"fmt"
	"gotool/pkg/environment"
	"os"
	"sync"
)
import "github.com/spf13/viper"

type configuration struct {
	paths  []string
	vipers map[string]*viper.Viper
	sync.Mutex
}

var config *configuration

func InitConfig(configPath ...string) {
	if len(configPath) == 0 {
		configPath = []string{"../../configs", "configs", "../configs"}
	}
	config = &configuration{
		paths:  configPath,
		vipers: make(map[string]*viper.Viper),
	}
}

func (c *configuration) getConfigFile(file string, env environment.Env) string {
	configFile := ""
	for _, path := range c.paths {
		tmpConfigFile := fmt.Sprintf("%s/%s/%s.toml", path, env, file)
		if _, err := os.Stat(tmpConfigFile); err == nil {
			configFile = tmpConfigFile
			break
		}
	}
	return configFile
}

func Cfg(file string) *viper.Viper {
	if cfg, ok := config.vipers[file]; ok {
		return cfg
	}
	config.Lock()
	defer config.Unlock()
	// 读取基础配置
	baseConfig := viper.New()
	baseConfigFile := config.getConfigFile(file, Env())
	baseConfig.SetConfigFile(baseConfigFile)
	err := baseConfig.ReadInConfig()
	if err != nil {
		fmt.Println("Cfg err:", err)
		return nil
	}
	// 将基础配置全部以默认配置写入
	envConfig := viper.New()
	for k, v := range baseConfig.AllSettings() {
		envConfig.SetDefault(k, v)
	}
	envConfigFile := config.getConfigFile(file, Env())
	if envConfigFile != "" {
		envConfig.SetConfigFile(envConfigFile)
		err = envConfig.ReadInConfig()
		if err != nil {
			fmt.Println("Cfg err 2:", err)
			return nil
		}
	}
	config.vipers[file] = envConfig
	return envConfig
}

func ResetCfgKey(file string) {
	config.Lock()
	defer config.Unlock()
	if _, ok := config.vipers[file]; ok {
		delete(config.vipers, file)
	}
}
