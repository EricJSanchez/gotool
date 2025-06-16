package sys

import (
	"github.com/olivere/elastic/v7"
	"reflect"
)

var esManager = NewEsClientManager()

func Elastic(names ...string) (client *elastic.Client) {
	var name = Cfg("app").GetString("default_es")

	if len(names) > 0 {
		name = names[0]
	}
	config := Nacos("database.toml").GetStringMap(name)
	client = esManager.Get(name, config)

	return
}

func EsToStruct[T any](result *elastic.SearchResult) (ret []T, err error) {
	var typ *T
	//遍历命中的数据，对数据进行类型断言，获取数据
	for _, item := range result.Each(reflect.TypeOf(typ)) {
		ret = append(ret, *item.(*T))
	}
	return
}
