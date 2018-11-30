package rabbitmq

import (
	"github.com/streadway/amqp"
)

type CallBack interface {
	Call(msg amqp.Delivery)
}
