package test

import (
	"github.com/EricJSanchez/gotool/pkg/environment"
	"github.com/EricJSanchez/gotool/service"
	"github.com/EricJSanchez/gotool/sys"
	"github.com/sirupsen/logrus"
)

var (
	serv *service.List
	log  *logrus.Logger
)

func init() {
	// 初始化环境
	sys.InitEnv(environment.Development)
	serv = service.Register()
	//初始化 config
	sys.InitConfig("../configs/")

	log = sys.InitLog()
	//nacos初始化
	_ = serv.Nacos.InitClient()
	//nacos 注册服务
	//defer service.Factory.Nacos.DeRegister()
	serv.Nacos.Register()

}
