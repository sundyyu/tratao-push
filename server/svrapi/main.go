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
	rootPath := cfg.GetString("http.rootPath")

	beego.BConfig.Listen.HTTPPort = cfg.GetInt("http.listenAddr")
	beego.BConfig.CopyRequestBody = true

	beego.Router(rootPath+"/add", &controller.AlarmController{}, "post:AddAlarm")
	beego.Router(rootPath+"/update", &controller.AlarmController{}, "put:UpdateAlarm")
	beego.Router(rootPath+"/list/:account", &controller.AlarmController{}, "get:ListAlarm")
	beego.Router(rootPath+"/del", &controller.AlarmController{}, "delete:DelAlarm")
	beego.Router(rootPath+"/updatedev", &controller.AlarmController{}, "put:UpdateDevice")
	beego.Run()
}
