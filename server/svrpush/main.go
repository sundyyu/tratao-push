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

	cfg := config.LoadConfig(path)
	// cfg := config.NewConfig("../../config/cfg.yaml")
	max := cfg.GetInt("check.maxReceive") // 并发执行线程数
	receiveService := push.ReceiveService{}
	receiveService.CallChan = make(chan int, max)

	forever := make(chan int, 1)

	// 监听推送的消息队列
	rabbitmq.DoReceive(receiveService)

	<-forever

}
