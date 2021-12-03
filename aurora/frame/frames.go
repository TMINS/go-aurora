package frame

/*
	整合第三方框架标准 key
*/
const (
	GORM     = "gorm"     // gorm    容器数据库连接实例key
	GO_REDIS = "go-redis" // go-redis 容器客户端连接实例key
	RABBITMQ = "RabbitMQ" // rabbit mq 容器客户端连接实例key
	DB       = "db"       // db作为原生 db
)
