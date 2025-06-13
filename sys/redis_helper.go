package sys

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"sync"
	"time"
)

// 封装成一个redis资源池
type RedisConn struct {
	pool      *redis.Client
	showDebug bool
}

// 设置是否打印操作日志
func (rds *RedisConn) ShowDebug(b bool) {
	rds.showDebug = b
}

type RedisClientManager struct {
	rw      *sync.RWMutex
	clients map[string]*RedisConn
}

func NewRedisClientManager() *RedisClientManager {
	return &RedisClientManager{
		rw:      &sync.RWMutex{},
		clients: make(map[string]*RedisConn),
	}
}

// 获取给定名称的 Gorm 客户端实例（如果客户端不存在则返回 nil）
func (m *RedisClientManager) Get(name string, config map[string]interface{}) *RedisConn {
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
func (m *RedisClientManager) NewInstance(config map[string]interface{}) (client *RedisConn) {
	ro := &redis.Options{
		Addr:            config["addr"].(string) + ":" + cast.ToString(config["port"].(int64)),
		Password:        config["password"].(string),
		DB:              int(config["database"].(int64)),
		PoolSize:        int(config["pool_size"].(int64)),
		MinIdleConns:    int(config["min_idle_conns"].(int64)),
		ConnMaxIdleTime: time.Duration(config["conn_max_idle_time"].(int64)) * time.Second,
		ConnMaxLifetime: time.Duration(config["conn_max_lifetime"].(int64)) * time.Second,
		MaxRetries:      int(config["max_retries"].(int64)),
		MinRetryBackoff: time.Duration(config["min_retry_backoff"].(int64)) * time.Millisecond,
	}
	pool := redis.NewClient(ro)

	pong, err := pool.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis connect fail", err)
		//Log().Error(err)
		//Log().Info("redis连接失败")
		//Log().Panic(err)
	}
	fmt.Println("redis连接成功:", pong)
	client = &RedisConn{
		pool: pool,
	}

	//非线上环境开启debug
	//if !environment.Is(environment.Production) {
	//	client.ShowDebug(true)
	//}

	return
}
