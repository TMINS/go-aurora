package aurora

import (
	"github.com/awensir/go-aurora/aurora/frame"
	"github.com/awensir/go-aurora/aurora/option"
	"github.com/streadway/amqp"
	"log"
)

// RabbitMqConfig 链接RabbitMQ address 链接地址
//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
func (a *Aurora) RabbitMqConfig(opt Opt) {
	o := opt()
	conn, err := amqp.Dial(o[option.RABBITMQ_URL].(string))
	failOnError(err, "Failed to connect to RabbitMQ")
	a.container.store(frame.RABBITMQ, conn)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
