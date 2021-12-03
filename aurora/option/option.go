package option

const (
	//go-redis 配置项键 （*redis.Options）
	GOREDIS_CONFIG = "go-redis"

	//gorm 数据库类型配置项键 （gorm.Dialector）
	GORM_TYPE = "database" //gorm 数据库类型

	//gorm 配置项选项键 （gorm.Option）
	GORM_CONFIG = "config" //gorm 配置项

	RABBITMQ_URL = "rabbit-url"

	//添加配置 配置项
	Config_key = "name" //定义配置 名
	Config_fun = "func" //定义配置 函数 (type Configuration func(Opt) interface{})
	Config_opt = "opt"  //定义配置 参数选项 (type Opt func() map[string]interface{})

)
