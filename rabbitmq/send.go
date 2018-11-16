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
	queue, err := ch.QueueDeclare(
		"test_queue", // name
		true,         // durable 持久化
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		util.LogError(queue.Name, err)
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

func DoSend(ch *amqp.Channel, body []byte, send chan int) {
	if send != nil {
		defer util.ReadChan(send)
	}
	if ch == nil {
		return
	}

	DoPublish(ch, body, nil)
}

// 发布消息
func DoPublish(ch *amqp.Channel, body []byte, pub chan int) {
	if pub != nil {
		defer util.ReadChan(pub)
	}

	err := ch.Publish(
		"",           // exchange
		"test_queue", // routing key
		false,        // mandatory
		false,        // immediate
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
