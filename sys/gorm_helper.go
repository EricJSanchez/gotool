package sys

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"sync"
	"time"
)

type GormClientManager struct {
	rw         *sync.RWMutex
	clients    map[string]*gorm.DB
	DebugLevel int
}

func NewGormClientManager() *GormClientManager {
	fmt.Println("-------------------NewGormClientManager-------------------")
	return &GormClientManager{
		rw:         &sync.RWMutex{},
		clients:    make(map[string]*gorm.DB),
		DebugLevel: 4,
	}
}

// 获取给定名称的 Gorm 客户端实例（如果客户端不存在则返回 nil）
func (m *GormClientManager) Get(name string, config map[string]interface{}) *gorm.DB {
	// 1、获取连接实例
	//Log().Printf("获取数据库链接：%s 读锁", name)
	m.rw.RLock()
	liveDebugLevel := Nacos().GetInt("MysqlDebugLevel")
	if client, exists := m.clients[name]; exists {
		m.rw.RUnlock()
		//log level
		// Silent LogLevel = iota + 1
		// Error
		// Warn
		// Info
		if m.DebugLevel != liveDebugLevel {
			m.clients[name] = client.Session(&gorm.Session{
				Logger: client.Logger.LogMode(logger.LogLevel(liveDebugLevel)),
			})
			m.DebugLevel = liveDebugLevel
		}
		return m.clients[name]
	}
	m.rw.RUnlock()

	// 获取读写锁
	//Log().Printf("获取数据库链接：%s 读写锁", name)
	m.rw.Lock()
	defer m.rw.Unlock()

	if client, exists := m.clients[name]; exists {
		if m.DebugLevel != liveDebugLevel {
			m.clients[name] = client.Session(&gorm.Session{
				Logger: client.Logger.LogMode(logger.LogLevel(liveDebugLevel)),
			})
			m.DebugLevel = liveDebugLevel
		}
		return m.clients[name]
	}
	// 2、添加连接实例
	client := m.NewInstance(config)
	m.clients[name] = client
	//Log().Printf("获取数据库链接：%s 成功，新建链接", name)
	return client
}

func beforeUpdate(client *gorm.DB) {
	sql := client.Dialector.Explain(client.Statement.SQL.String(), client.Statement.Vars...)
	//Log().Printf("设计的sql===", sql)
	sql = strings.Replace(sql, ",", "分", -1)
	sql = strings.Replace(sql, "'", "分", -1)
	sql = strings.Replace(sql, "`", "分", -1)
	sql = strings.Replace(sql, "=", "分", -1)

	//sql= fmt.Sprintf("`%s`", sql)

	//Log().Printf("设计的sql===", sql)

	if strings.Contains(sql, "ww_client_staff_union") && !strings.Contains(sql, "ww_temp_id") {
		//err:=dao.Factory.WwTempIdDao.AddAgent(sql)
		//
		//if err != nil {
		//	println("插入日志失败",err)
		//}else{
		//	println("插入日志成功",err)
		//
		//}
		_ = client.Raw("INSERT INTO ww_temp_id ( unionid ) VALUES (@unionid)", map[string]interface{}{"unionid": sql}).Error
	}

}

// 创建连接实例
func (m *GormClientManager) NewInstance(config map[string]interface{}) (client *gorm.DB) {
	var dialector gorm.Dialector = nil
	switch config["driver"] {
	case "postgres":
		dialector = m.newPostgresDialector(config)
	case "mysql":
		dialector = m.newMysqlDialector(config)
	default:
	}
	if dialector == nil {
		return
	}
	client, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt:          false,
		DisableAutomaticPing: false,
	})
	if err != nil {
		fmt.Println("db_helper NewInstance err: ", err)
		return nil
	}
	// 日志级别交由 nacos 配置
	client = client.Session(&gorm.Session{
		Logger: client.Logger.LogMode(logger.LogLevel(m.DebugLevel)),
	})

	sqlDB, err := client.DB()
	if err != nil {
		fmt.Println("db_helper NewInstance err: ", err)
		return
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(int(config["max_idle_conn"].(int64)))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(int(config["max_conn"].(int64)))
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	return
}

func (m *GormClientManager) newMysqlDialector(config map[string]interface{}) gorm.Dialector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=5s",
		config["username"],
		config["password"],
		config["host"],
		config["port"],
		config["database"],
		config["charset"])
	return mysql.Open(dsn)
}

func (m *GormClientManager) newPostgresDialector(config map[string]interface{}) gorm.Dialector {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d connect_timeout=5 sslmode=disable TimeZone=Asia/Shanghai",
		config["host"],
		config["username"],
		config["password"],
		config["database"],
		config["port"],
	)
	return postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	})
}

// 清空 Gorm 客户端实例
func (m *GormClientManager) Clear() {
	m.rw.Lock()
	defer m.rw.Unlock()
	for k := range m.clients {
		delete(m.clients, k)
	}
}
