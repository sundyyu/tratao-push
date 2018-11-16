package rabbitmq

import (
	"github.com/streadway/amqp"
	"strconv"
	"tratao-push/config"
	"tratao-push/util"
)

func DoReceive(call CallBack) {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "recover DoReceive error.")
		}
	}()

	// cfg := config.LoadConfig("../../config/cfg.yaml")
	cfg := config.GetConfig()
	user := cfg.GetString("rabbitmq.user")
	pass := cfg.GetString("rabbitmq.password")
	ip := cfg.GetString("rabbitmq.ip")
	port := cfg.GetInt("rabbitmq.port")

	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@" + ip + ":" + strconv.Itoa(port) + "/")
	util.LogError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	util.LogError(err)
	defer ch.Close()

	// 消息队列
	q, err := ch.QueueDeclare(
		"test_queue", // name
		true,         // durable 持久化
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	util.LogError(err)

	// 公平派遣
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	util.LogError(err)

	// 获取消息
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack  false为不自动应答
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	util.LogError(err)

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			call.Call(m) // 回调处理消息
		}
	}()

	util.LogInfo("Waiting for messages. To exit press CTRL+C")
	<-forever
}
