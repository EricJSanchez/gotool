package test

import (
	"fmt"
	"testing"
)

func init() {
	fmt.Println("init")
}

func TestLib(t *testing.T) {
	gotool.InitConfig("../configs/", "development")
	nacosAddr := php2go.Cfg("app").GetString("nacos.addr")
	php2go.Pr(nacosAddr)
	return
}
