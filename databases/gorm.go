package databases

import (
	"errors"
	"github.com/awensir/Aurora/manage"
	"github.com/awensir/Aurora/manage/frame"
	"gorm.io/gorm"
)

const (
	DBT    = "database"
	CONFIG = "config"
)

/*
	整合gorm 框架
	默认使用 v2版本
	提供配置项 初始化默认gorm变量
	需要连接多个库，存放在容器中，实现 manage.Variable 接口 Clone() Variable 方法即可存入容器
*/

type GORM struct {
	*gorm.DB
}

func (g *GORM) Clone() manage.Variable {
	return g
}

//GormConfig 整合gorm
func GormConfig(opt map[string]interface{}) {
	//读取配置项
	dil, b := opt[DBT].(gorm.Dialector)
	if !b {
		panic(errors.New("gorm config option gorm.Dialector type error！"))
	}

	config, b := opt[CONFIG].(gorm.Option)
	if !b {
		panic(errors.New("gorm config option gorm.Option type error！"))
	}
	db, err := gorm.Open(dil, config)
	if err != nil {
		panic(err.Error())
	}
	gormDB := &GORM{db}
	manage.Container.Store(frame.GORM, gormDB)
}
