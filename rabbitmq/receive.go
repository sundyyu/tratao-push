package rabbitmq

import (
	"github.com/streadway/amqp"
	"tratao-push/util"
)

const ALARM_QUEUE = "alarm_queue"

func DoReceive(ch *amqp.Channel, call CallBack) {
	DoReceiveQueue(ch, call, ALARM_QUEUE)
}

func DoReceiveQueue(ch *amqp.Channel, call CallBack, queue string) {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "recover DoReceive error.")
		}
	}()

	// 公平派遣
	err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	util.LogError(err)

	// 获取消息
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack  false为不自动应答
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	util.LogError(err)

	go func() {
		for m := range msgs {
			call.Call(m) // 回调处理消息
		}
	}()

	util.LogInfo(queue + " Waiting for messages. To exit press CTRL+C")
}
