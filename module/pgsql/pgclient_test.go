package pgclient

import (
	"github.com/liamylian/jsontime"
	"testing"
	"time"
	"xcurrency-push/config"
	"xcurrency-push/model"
	"xcurrency-push/util"
)

func TestPgSql(t *testing.T) {
	config.LoadConfig("../config/cfg.yaml")

	pushmsg := model.PushMsg{}
	pushmsg.Account = "test"
	pushmsg.Title = "推送测试"
	pushmsg.Body = "小米第4次推送测试"
	pushmsg.DeviceId = "ddddd11111"
	pushmsg.DeviceOS = "xiaomi"
	pushmsg.DeviceLang = "zh"
	pushmsg.DeviceCountry = "CN"
	pushmsg.CreateTime = time.Now()

	// InsertPushMsg(pushmsg)
	list, _ := QueryPushMsg("test2", 1, 5)

	var json = jsontime.ConfigWithCustomTimeFormat
	byt, _ := json.Marshal(list)
	util.LogInfo(string(byt))

	defer CloseConn()

}
