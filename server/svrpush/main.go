package main

import (
	"flag"
	"os"
	"tratao-push/config"
	"tratao-push/rabbitmq"
	"tratao-push/server/svrpush/push"
	"tratao-push/util"
)

func main() {
	c := flag.String("config", "", "配置文件参数")
	flag.Parse()

	path := *c
	_, err := os.Stat(path)
	if err == nil {
		util.LogInfoF("config file %s exists", path)
	} else if os.IsNotExist(err) {
		util.LogInfoF("config file %s not exists", path)
		return
	} else {
		util.LogInfoF("config file %s stat error: %v", path, err)
		return
	}

	forever := make(chan int, 1)
	cfg := config.LoadConfig(path)
	max := cfg.GetInt("check.maxReceive") // 并发执行线程数
	conn := rabbitmq.GetConn()
	defer conn.Close()

	// 监听Alarm推送的消息队列
	alarm := push.AlarmReceive{}
	alarm.CallChan = make(chan int, max)
	ch := rabbitmq.GetChannel(conn)
	defer ch.Close()
	rabbitmq.DoReceive(ch, alarm)

	// 监听PushMsg推送的消息队列
	queue := cfg.GetString("check.pushMsgQueue")
	pushMsg := push.PushMsgReceive{}
	pushMsg.CallChan = make(chan int, max)
	pch := rabbitmq.GetChannelQueue(conn, queue)
	defer pch.Close()
	rabbitmq.DoReceiveQueue(pch, pushMsg, queue)

	<-forever
}
