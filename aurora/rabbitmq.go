package aurora

import (
	"github.com/awensir/Aurora/aurora/frame"
	"github.com/streadway/amqp"
	"log"
)

// RabbitMqConfig 链接RabbitMQ address 链接地址
//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
func (a *Aurora) RabbitMqConfig(address string) {
	conn, err := amqp.Dial(address)
	failOnError(err, "Failed to connect to RabbitMQ")
	a.container.store(frame.RABBITMQ, conn)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
