package test

import (
	"fmt"
	"gotool"
	"php2go"
	"testing"
)

func init() {
	fmt.Println("init")
}

func TestLib(t *testing.T) {
	gotool.InitConfig("../configs/", "development")
	nacosAddr := gotool.Cfg("app").GetString("nacos.addr")
	php2go.Pr(nacosAddr)
	return
}
