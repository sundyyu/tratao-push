package main

import (
	"flag"
	"os"
	"xcurrency-push/config"
	"xcurrency-push/server/svrcheck/check"
	"xcurrency-push/util"
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

	config.LoadConfig(path)

	forever := make(chan int, 1)

	// 汇率Check
	var exrate check.Exrate = &check.ExrateYahoo{}
	exrate.Update()
	exrate.Loop()

	// 预警信息Check
	var checkAlarm check.Check = &check.CheckAlarm{UExrate: exrate}
	checkAlarm.Loop()

	// 推送信息Check
	var checkPushMsg check.Check = &check.CheckMessage{}
	checkPushMsg.Loop()

	<-forever
}
