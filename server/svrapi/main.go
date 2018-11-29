package main

import (
	"flag"
	"github.com/astaxie/beego"
	"os"
	"xcurrency-push/config"
	"xcurrency-push/server/svrapi/api"
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

	// cfg := config.NewConfig("../../config/cfg.yaml")
	cfg := config.LoadConfig(path)
	alarmPath := cfg.GetString("http.alarmPath")

	beego.BConfig.Listen.HTTPPort = cfg.GetInt("http.listenAddr")
	beego.BConfig.CopyRequestBody = true

	beego.Router(alarmPath+"/add", &api.AlarmController{}, "post:AddAlarm")
	beego.Router(alarmPath+"/update", &api.AlarmController{}, "put:UpdateAlarm")
	beego.Router(alarmPath+"/list/:account", &api.AlarmController{}, "get:ListAlarm")
	beego.Router(alarmPath+"/del", &api.AlarmController{}, "delete:DelAlarm")
	beego.Router(alarmPath+"/updatedev", &api.AlarmController{}, "put:UpdateDevice")

	beego.Router(alarmPath+"/push", &api.AlarmController{}, "post:AddPushMsg")
	beego.Run()
}
