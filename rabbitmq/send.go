package rabbitmq

import (
	"github.com/streadway/amqp"
	"strconv"
	"sync"
	"tratao-push/config"
	"tratao-push/util"
)

var connLock *sync.Mutex = new(sync.Mutex)

func GetConn() *amqp.Connection {
	connLock.Lock()
	defer connLock.Unlock()

	// cfg := config.LoadConfig("../../config/cfg.yaml")
	cfg := config.GetConfig()
	user := cfg.GetString("rabbitmq.user")
	pass := cfg.GetString("rabbitmq.password")
	ip := cfg.GetString("rabbitmq.ip")
	port := cfg.GetInt("rabbitmq.port")

	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@" + ip + ":" + strconv.Itoa(port) + "/")

	if err != nil {
		util.LogError(err)
		return nil
	}
	return conn
}

var chLock *sync.Mutex = new(sync.Mutex)

func GetChannel(conn *amqp.Connection) *amqp.Channel {
	return GetChannelQueue(conn, ALARM_QUEUE)
}

func GetChannelQueue(conn *amqp.Connection, queue string) *amqp.Channel {
	chLock.Lock()
	defer chLock.Unlock()

	if conn == nil {
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		util.LogError(err)
		return nil
	}

	// 消息队列
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable 持久化
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		util.LogError(q.Name, err)
		return nil
	}
	return ch
}

func GetChannelArr(conn *amqp.Connection, count int) []*amqp.Channel {
	chanArr := make([]*amqp.Channel, 0, count)
	for i := 0; i < count; i++ {
		chanArr = append(chanArr, GetChannel(conn))
	}
	return chanArr
}

func CloseChannelArr(chanArr []*amqp.Channel) {
	for _, ch := range chanArr {
		if ch != nil {
			ch.Close()
		}
	}
}

func GetConnArr(count int) []*amqp.Connection {
	chanArr := make([]*amqp.Connection, 0, count)
	for i := 0; i < count; i++ {
		chanArr = append(chanArr, GetConn())
	}
	return chanArr
}

func CloseConnArr(chanArr []*amqp.Connection) {
	for _, ch := range chanArr {
		if ch != nil {
			ch.Close()
		}
	}
}

func DoSend(ch *amqp.Channel, body []byte) {
	if ch == nil {
		return
	}
	DoPublish(ch, body, ALARM_QUEUE)
}

func DoSendQueue(ch *amqp.Channel, body []byte, queue string) {
	if ch == nil {
		return
	}
	DoPublish(ch, body, queue)
}

// 发布消息
func DoPublish(ch *amqp.Channel, body []byte, queue string) {
	err := ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久化
			ContentType:  "text/plain",
			Body:         body, // []byte(body)
		})

	if err != nil {
		util.LogError(err)
	}
}

func CloseConn(conn *amqp.Connection) {
	if conn != nil {
		conn.Close()
	}
}

func CloseChannel(ch *amqp.Channel) {
	if ch != nil {
		ch.Close()
	}
}
