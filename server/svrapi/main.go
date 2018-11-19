package main

import (
	"flag"
	"github.com/astaxie/beego"
	"os"
	"tratao-push/config"
	"tratao-push/server/svrapi/controller"
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

	// cfg := config.NewConfig("../../config/cfg.yaml")
	cfg := config.LoadConfig(path)
	alarmPath := cfg.GetString("http.alarmPath")
	pushMsgPath := cfg.GetString("http.pushMsgPath")

	beego.BConfig.Listen.HTTPPort = cfg.GetInt("http.listenAddr")
	beego.BConfig.CopyRequestBody = true

	beego.Router(alarmPath+"/add", &controller.AlarmController{}, "post:AddAlarm")
	beego.Router(alarmPath+"/update", &controller.AlarmController{}, "put:UpdateAlarm")
	beego.Router(alarmPath+"/list/:account", &controller.AlarmController{}, "get:ListAlarm")
	beego.Router(alarmPath+"/del", &controller.AlarmController{}, "delete:DelAlarm")
	beego.Router(alarmPath+"/updatedev", &controller.AlarmController{}, "put:UpdateDevice")

	beego.Router(pushMsgPath+"/add", &controller.AlarmController{}, "post:AddPushMsg")
	beego.Run()
}
