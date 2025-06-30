package sys

import (
	"fmt"
	"github.com/EricJSanchez/gotool/environment"
	"os"
	"runtime"
	"sync"
)
import "github.com/spf13/viper"

type configuration struct {
	paths  []string
	vipers map[string]*viper.Viper
	sync.Mutex
}

var configLocal *configuration

func InitConfig(configPath ...string) {
	if len(configPath) == 0 {
		configPath = []string{"../../configs", "configs", "../configs"}
	}
	configLocal = &configuration{
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
	if cfg, ok := configLocal.vipers[file]; ok {
		return cfg
	}
	configLocal.Lock()
	defer configLocal.Unlock()
	// 读取基础配置
	baseConfig := viper.New()
	baseConfigFile := configLocal.getConfigFile(file, Env())
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
	envConfigFile := configLocal.getConfigFile(file, Env())
	if envConfigFile != "" {
		envConfig.SetConfigFile(envConfigFile)
		err = envConfig.ReadInConfig()
		if err != nil {
			fmt.Println("Cfg err 2:", err)
			return nil
		}
	}
	configLocal.vipers[file] = envConfig
	return envConfig
}

func ResetCfgKey(file string) {
	configLocal.Lock()
	defer configLocal.Unlock()
	if _, ok := configLocal.vipers[file]; ok {
		delete(configLocal.vipers, file)
	}
}

func GetRunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
