package aurora

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

const (
	Mysql = iota
	SqlServer
	SQLite
	Postgresql
)

/*
	自动读取application。yml 配置对应 gorm实例
	读取配置文件中的url列表进行注册实例
*/

func (a *Aurora) loadGormConfig() {
	if a.cnf == nil {
		//如果配置文件没有加载成功，将不做任何事情
		return
	}
	mysqls := a.cnf.GetStringSlice("gorm.mysql.url")
	if mysqls != nil {
		a.auroraLog.Info("load gorm mysql configuration information")
	}
	add(mysqls, &a.gorms, Mysql, a.auroraLog)

	sqlservers := a.cnf.GetStringSlice("gorm.sqlserver.url")
	if sqlservers != nil {
		a.auroraLog.Info("load gorm sqlservers mysql configuration information")
	}
	add(sqlservers, &a.gorms, SqlServer, a.auroraLog)

	postgresql := a.cnf.GetStringSlice("gorm.postgresql.url")
	if postgresql != nil {
		a.auroraLog.Info("load gorm postgresql configuration information")
	}
	add(postgresql, &a.gorms, Postgresql, a.auroraLog)

	if mysqls != nil || sqlservers != nil || postgresql != nil {
		//该日志之前 任何数据库配置错误将导致服务停止
		a.auroraLog.Info("load database configuration is complete")
	}
}

// 添加一个 gorm 配置
func add(urls []string, gorms *map[int][]*gorm.DB, db int, logs *Log) {
	if urls != nil {
		if _, b := (*gorms)[db]; !b {
			(*gorms)[db] = make([]*gorm.DB, 0)
		}
		dbs := (*gorms)[db]
		for _, url := range urls {
			if dbs != nil {
				dbs = append(dbs, DefaultSQLConfig(db, logs, url))
				(*gorms)[db] = dbs
			}
		}
	}
}

// DefaultSQLConfig 返回一个 *gorm.DB 连接
// config
// config[0]:user		用户名
// config[1]:password	密码
// config[2]:ip			地址
// config[3]:port 		端口
// config[4]:db_name	库
func DefaultSQLConfig(database int, logs *Log, config ...string) *gorm.DB {
	dns := ""
	if len(config) != 1 && len(config) != 5 {
		logs.Error("configuration parameter error returns nil")
		return nil
	}
	if len(config) == 1 {
		dns = config[0]
	}
	if len(config) == 5 {
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config[0], config[1], config[2], config[3], config[4])
	}
	if dns == "" {
		return nil
	}
	switch database {
	case Mysql:
		return mysqlDb(dns, logs)
	case SqlServer:
		return sqlServer(dns, logs)
	case SQLite:
		//return sqlIte(dns)
	case Postgresql:
		return postgreSql(dns, logs)
	}
	return nil
}

// 提供的默认日志
func newLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,          // 彩色打印
		},
	)
}

// 提供的gorm默认配置
func defaultConfig() *gorm.Config {
	return &gorm.Config{
		//GORM 配置项

		//日志配置
		Logger: newLogger(),

		//命名策略配置
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	}
}

// mysqlDb 返回mysql连接实例
func mysqlDb(url string, logs *Log) *gorm.DB {
	db, err := gorm.Open(mysql.New(
		mysql.Config{
			//数据库驱动配置项
			DSN:                       url,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), defaultConfig())
	if err != nil {
		logs.Error("mysql:" + err.Error())
		panic(err)
	}
	return db
}

// sqlIte 返回 sqlIte 连接实例
//func sqlIte(url string) *gorm.DB {
//	db, err := gorm.Open(sqlite.Open(url), defaultConfig())
//	if err != nil {
//		panic(err.Error())
//	}
//	return db
//}

// sqlServer 返回 sqlServer 连接实例
func sqlServer(url string, logs *Log) *gorm.DB {
	db, err := gorm.Open(sqlserver.Open(url), defaultConfig())
	if err != nil {
		logs.Error("sqlserver:" + err.Error())
		panic(err)
	}
	return db
}

// postgreSql 返回 postgreSql 连接实例
func postgreSql(url string, logs *Log) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), defaultConfig())
	if err != nil {
		logs.Error("postgresql:" + err.Error())
		panic(err)
	}
	return db
}

// Mysql 获取注册的 默认 mysql
func (a Aurora) Mysql() *gorm.DB {
	return get(a.gorms, Mysql, 0)
}

// SQLite 获取注册的 默认 SQLite
func (a Aurora) SQLite() *gorm.DB {
	return get(a.gorms, SQLite, 0)
}

// PostgreSql 获取注册的 默认 PostgreSql
func (a Aurora) PostgreSql() *gorm.DB {
	return get(a.gorms, Postgresql, 0)
}

// SqlServer 获取注册的 默认 SqlServer
func (a Aurora) SqlServer() *gorm.DB {
	return get(a.gorms, SqlServer, 0)
}

func (a *Aurora) RegisterGorm(dbtype int, db *gorm.DB) int {
	if _, b := a.gorms[dbtype]; !b {
		a.gorms[dbtype] = make([]*gorm.DB, 0)
	}
	dbs := a.gorms[dbtype]
	if dbs != nil {
		dbs = append(dbs, db)
		a.gorms[dbtype] = dbs
		return len(dbs) - 1
	}
	return -1
}

func get(gorms map[int][]*gorm.DB, db int, index int) *gorm.DB {
	if dbs, b := gorms[db]; b {
		if dbs != nil && len(dbs) > 0 && index < len(dbs) {
			return dbs[index]
		}
	}
	return nil
}
