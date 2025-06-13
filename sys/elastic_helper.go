package sys

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type EsClientManager struct {
	rw      *sync.RWMutex
	clients map[string]*elastic.Client
}

func NewEsClientManager() *EsClientManager {
	return &EsClientManager{
		rw:      &sync.RWMutex{},
		clients: make(map[string]*elastic.Client),
	}
}

// 获取给定名称的 Es 客户端实例（如果客户端不存在则返回 nil）
func (m *EsClientManager) Get(name string, config map[string]interface{}) *elastic.Client {
	// 1、获取连接实例
	m.rw.RLock()
	if client, exists := m.clients[name]; exists {
		m.rw.RUnlock()
		return client
	}
	m.rw.RUnlock()

	// 获取读写锁
	m.rw.Lock()
	defer m.rw.Unlock()

	if client, exists := m.clients[name]; exists {
		return client
	}
	// 2、添加连接实例
	client := m.NewInstance(config)
	m.clients[name] = client
	return client
}

// 创建连接实例
func (m *EsClientManager) NewInstance(config map[string]interface{}) (client *elastic.Client) {
	address := strings.Split(config["addresses"].(string), ",")
	client, err := elastic.NewClient(

		elastic.SetHealthcheck(false),
		elastic.SetHealthcheckTimeoutStartup(5*time.Second),
		elastic.SetHealthcheckTimeout(1*time.Second),
		elastic.SetHealthcheckInterval(100*time.Second),

		elastic.SetURL(address...),
		elastic.SetBasicAuth(
			config["username"].(string),
			config["password"].(string),
		),
		elastic.SetSniff(false),
		//elastic.SetHealthcheckInterval(10*time.Second),
		//elastic.SetRetrier(NewCustomRetrier()),
		//elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		//elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetHeaders(http.Header{
			"X-Caller-Id": []string{"..."},
		}),
	)
	if err != nil {
		fmt.Println("es实例化出错", err)
	}
	return
}

// 清空 Es 客户端实例
func (m *EsClientManager) Clear() {
	m.rw.Lock()
	defer m.rw.Unlock()
	for k := range m.clients {
		delete(m.clients, k)
	}
}
