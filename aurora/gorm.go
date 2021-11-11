package aurora

import (
	"errors"
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/awensir/go-aurora/aurora/option"
	"gorm.io/gorm"
)

const (
	DBT    = "database" //gorm 数据库类型
	CONFIG = "config"   //gorm 配置项
)

/*
	整合gorm 框架
	默认使用 v2版本
	提供配置项 初始化默认gorm变量
	需要连接多个库，存放在容器中，实现 manage.Variable 接口 Clone() Variable 方法即可存入容器
*/

//GormConfig 整合gorm
func (a *Aurora) GormConfig(opt Opt) {
	o := opt()
	//读取配置项
	dil, b := o[option.GORM_TYPE].(gorm.Dialector)
	if !b {
		panic(errors.New("gorm config option gorm.Dialector type error！"))
	}

	config, b := o[option.GORM_CONFIG].(gorm.Option)
	if !b {
		panic(errors.New("gorm config option gorm.Option type error！"))
	}
	db, err := gorm.Open(dil, config)
	if err != nil {
		panic(err.Error())
	}
	a.container.store(frame.GORM, db)
}
