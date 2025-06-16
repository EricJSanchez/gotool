package test

import (
	"github.com/EricJSanchez/gotool/pkg/environment"
	"github.com/EricJSanchez/gotool/service"
	"github.com/EricJSanchez/gotool/sys"
)

func init() {
	// 初始化环境
	sys.InitEnv(environment.Development)
	service.Register()
	//初始化 config
	sys.InitConfig("../configs/")
	//nacos初始化
	_ = service.Factory.Nacos.InitClient(sys.Cfg("app").GetStringMap("nacos"))
	//nacos 注册服务
	//defer service.Factory.Nacos.DeRegister()
	//service.Factory.Nacos.Register()

}
