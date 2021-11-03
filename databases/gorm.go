package databases

import (
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DBT    = "database"
	CONFIG = "config"
)

/*
	整合gorm 框架
	默认使用 v2版本
	默认使用数据库 MySQL
*/

func init() {

}

var Gorm *GORM

type GORM struct {
	*gorm.DB
}

//InitGorm 默认配置
func InitGorm(g *GORM) {
	if g == nil {
		g = &GORM{nil}
		//准备配置

		//日志配置
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second,   // 慢 SQL 阈值
				LogLevel:                  logger.Silent, // 日志级别
				IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,         // 禁用彩色打印
			},
		)
		//命名策略
		namingStrategy := schema.NamingStrategy{
			TablePrefix:   "t_",                              // 表名前缀，`User`表为`t_users`
			SingularTable: true,                              // 使用单数表名，启用该选项后，`User` 表将是`user`
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
		}
		config := &gorm.Config{
			Logger:         newLogger,
			NamingStrategy: namingStrategy,
		}
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
			DefaultStringSize:         256,                                                                        // string 类型字段的默认长度
			DisableDatetimePrecision:  true,                                                                       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,                                                                       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,                                                                       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false,                                                                      // 根据当前 MySQL 版本自动配置
		}), config)

		if err != nil {
			panic(err.Error())
		}
		g.DB = db
	}
}

//GormOptConfig 配置项初始化Gorm
func GormOptConfig(g *GORM, opt map[string]interface{}) {
	if g == nil {
		g = &GORM{nil}
		//准备配置
		_, err := gorm.Open(opt[DBT].(gorm.Dialector), opt[CONFIG].(gorm.Option))
		if err != nil {
			panic(err.Error())
		}
	}
}
