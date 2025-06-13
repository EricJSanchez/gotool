package gotool

import (
	"fmt"
	"testing"
)

func init() {
	fmt.Println("init")
}

func TestLib(t *testing.T) {
	InitConfig("../configs/", "development")
	nacosAddr := Cfg("app").GetString("nacos.addr")
	Pr(nacosAddr)
	return
}
