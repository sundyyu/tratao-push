package api

import (
	"github.com/astaxie/beego"
	"github.com/liamylian/jsontime"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
	"xcurrency-push/model"
	"xcurrency-push/module/pgsql"
	"xcurrency-push/module/redis"
	"xcurrency-push/util"
)

type AlarmController struct {
	beego.Controller
}

func (this *AlarmController) AddAlarm() {
	alarm := model.Alarm{}
	reqBody := this.Ctx.Input.RequestBody

	// Account
	account := gjson.GetBytes(reqBody, "account").String()
	if len(account) == 0 {
		errInfo := "[alarm_controller.go AddAlarm] fail,  account can not be empty."
		util.LogError(errInfo)

		this.Data["json"] = util.JsonResult(0, nil, errInfo)
		this.ServeJSON()
		return
	}

	alarm.Account = account
	alarm.BaseCur = gjson.GetBytes(reqBody, "basecur").String()
	alarm.TargetCur = gjson.GetBytes(reqBody, "targetcur").String()
	alarm.LbPrice = gjson.GetBytes(reqBody, "lbprice").Float()
	alarm.UbPrice = gjson.GetBytes(reqBody, "ubprice").Float()
	alarm.Enabled = gjson.GetBytes(reqBody, "enabled").Bool()
	alarm.CreateTime = time.Now().Unix()
	alarm.UpdateTime = time.Now().Unix()

	id, err := redis.AddAlarm(alarm)
	if err != nil {
		util.LogError("[alarm_controller.go AddAlarm] fail, account: ", account, err)

		this.Data["json"] = util.JsonResult(0, nil, err.Error())
		this.ServeJSON()
		return
	}
	util.LogInfo("[alarm_controller.go AddAlarm] success, account: ", account)

	this.Data["json"] = util.JsonResult(1, id, "add alarm success.")
	this.ServeJSON()
}

func (this *AlarmController) UpdateAlarm() {
	alarm := model.Alarm{}
	reqBody := this.Ctx.Input.RequestBody

	// Id
	id := gjson.GetBytes(reqBody, "id").Int()
	if id <= 0 {
		util.LogError("[alarm_controller.go UpdateAlarm] fail, id: ", id, "id is empty or error.")

		this.Data["json"] = util.JsonResult(0, nil, "id error. id is empty or error.")
		this.ServeJSON()
		return
	}
	alarm.Id = id

	// Account
	// account := gjson.GetBytes(reqBody, "account").String()
	// if len(account) == 0 {
	// 	util.LogError("[alarm_controller.go UpdateAlarm] fail, account: ", account, "account can not be empty.")
	//
	// 	this.Data["json"] = util.JsonResult(0, nil, "account error. account can not be empty.")
	// 	this.ServeJSON()
	// 	return
	// }
	// alarm.Account = account

	alarm.BaseCur = gjson.GetBytes(reqBody, "basecur").String()
	alarm.TargetCur = gjson.GetBytes(reqBody, "targetcur").String()
	alarm.LbPrice = gjson.GetBytes(reqBody, "lbprice").Float()
	alarm.UbPrice = gjson.GetBytes(reqBody, "ubprice").Float()
	alarm.Enabled = gjson.GetBytes(reqBody, "enabled").Bool()
	alarm.UpdateTime = time.Now().Unix()

	err := redis.UpdateAlarm(alarm)
	if err != nil {
		util.LogError("[alarm_controller.go UpdateAlarm] fail, id: ", id, err)

		this.Data["json"] = util.JsonResult(0, nil, err.Error())
		this.ServeJSON()
		return
	}
	util.LogInfo("[alarm_controller.go UpdateAlarm] success, id: ", id)

	this.Data["json"] = util.JsonResult(1, nil, "update alarm success.")
	this.ServeJSON()
}

func (this *AlarmController) UpdateDevice() {
	alarm := model.Alarm{}
	reqBody := this.Ctx.Input.RequestBody

	// Account
	account := gjson.GetBytes(reqBody, "account").String()
	if len(account) == 0 {
		util.LogError("[alarm_controller.go UpdateDevice] fail, account: ", account, "account can not be empty.")

		this.Data["json"] = util.JsonResult(0, nil, "account error. Account can not be empty.")
		this.ServeJSON()
		return
	}

	alarm.Account = account
	alarm.DeviceId = gjson.GetBytes(reqBody, "devid").String()
	alarm.DeviceOS = gjson.GetBytes(reqBody, "devos").String()
	alarm.DeviceCountry = gjson.GetBytes(reqBody, "devcountry").String()
	alarm.DeviceLang = gjson.GetBytes(reqBody, "devlang").String()
	alarm.AppKey = gjson.GetBytes(reqBody, "appkey").String()
	alarm.UpdateTime = time.Now().Unix()

	err := redis.UpdateDevice(alarm)
	if err != nil {
		util.LogError("[alarm_controller.go UpdateDevice] fail, account: ", account, err)

		this.Data["json"] = util.JsonResult(0, nil, err.Error())
		this.ServeJSON()
		return
	}
	util.LogInfo("[alarm_controller.go UpdateDevice] success, account: ", account)

	this.Data["json"] = util.JsonResult(1, nil, "update device success.")
	this.ServeJSON()
}

func (this *AlarmController) ListAlarm() {
	account := this.Ctx.Input.Param(":account")

	alarms, err := redis.ListAlarm(account)
	if err != nil {
		util.LogError("[alarm_controller.go ListAlarm] fail, account: ", account, err)

		this.Data["json"] = util.JsonResult(0, nil, err.Error())
		this.ServeJSON()
		return
	}
	util.LogInfo("[alarm_controller.go ListAlarm] success, account: ", account)

	this.Data["json"] = util.JsonResult(1, alarms, nil)
	this.ServeJSON()
}

func (this *AlarmController) DelAlarm() {
	reqBody := this.Ctx.Input.RequestBody
	idArr := gjson.GetBytes(reqBody, "id")

	if !idArr.IsArray() {
		util.LogError("[alarm_controller.go DelAlarm] fail, id: ", idArr, "id is not a array.")

		this.Data["json"] = util.JsonResult(0, nil, "id is not a array.")
		this.ServeJSON()
		return
	}

	ids := []string{}
	for _, id := range idArr.Array() {
		ids = append(ids, id.String())
	}

	err := redis.DelAlarm(ids)
	if err != nil {
		util.LogError("[alarm_controller.go DelAlarm] fail, id: ", ids, err)

		this.Data["json"] = util.JsonResult(0, nil, err.Error())
		this.ServeJSON()
		return
	}
	util.LogInfo("[alarm_controller.go DelAlarm] success, id: ", ids)

	this.Data["json"] = util.JsonResult(1, nil, "delete alarm success.")
	this.ServeJSON()
}

func (this *AlarmController) AddPushMsg() {
	msg := model.PushMsg{}
	reqBody := this.Ctx.Input.RequestBody

	// Account
	account := gjson.GetBytes(reqBody, "account").String()
	if len(account) == 0 {
		errInfo := "[alarm_controller.go AddPushMsg] fail,  account can not be empty."
		util.LogError(errInfo)

		this.Data["json"] = util.JsonResult(0, nil, errInfo)
		this.ServeJSON()
		return
	}
	list, err := redis.ListAlarm(account)
	if err != nil {
		util.LogError(err)

		this.Data["json"] = util.JsonResult(0, nil, err)
		this.ServeJSON()
		return
	}
	if len(list) == 0 {
		errInfo := "[alarm_controller.go AddPushMsg] fail,  account not found."
		util.LogError(errInfo)

		this.Data["json"] = util.JsonResult(0, nil, errInfo)
		this.ServeJSON()
		return
	}

	alarm := list[0]
	msg.DeviceOS = alarm.DeviceOS
	msg.DeviceLang = alarm.DeviceLang
	msg.DeviceCountry = alarm.DeviceCountry

	msg.Account = account
	msg.DeviceId = gjson.GetBytes(reqBody, "token").String()
	msg.Title = gjson.GetBytes(reqBody, "title").String()
	msg.Body = gjson.GetBytes(reqBody, "body").String()
	msg.CreateTime = time.Now()

	id, err := redis.AddPushMsg(msg)
	if err != nil {
		util.LogError("[alarm_controller.go AddPushMsg] fail, account: ", account, err)

		this.Data["json"] = util.JsonResult(0, nil, err.Error())
		this.ServeJSON()
		return
	}
	util.LogInfo("[alarm_controller.go AddPushMsg] success, account: ", account)

	this.Data["json"] = util.JsonResult(1, id, "add push message success.")
	this.ServeJSON()
}

func (this *AlarmController) ListPushLog() {
	account := this.Ctx.Input.Param(":account")
	p := this.Ctx.Input.Param(":page")
	psize := this.Ctx.Input.Param(":pageSize")

	page := -1
	pageSize := -1
	if len(p) > 0 && len(psize) > 0 {
		p1, err1 := strconv.Atoi(p)
		p2, err2 := strconv.Atoi(psize)
		if err1 == nil && err2 == nil {
			page = p1
			pageSize = p2
		}
	}

	var j = jsontime.ConfigWithCustomTimeFormat
	pushLogs, err := pgclient.QueryPushMsg(account, page, pageSize)
	if err != nil {
		util.LogError("[alarm_controller.go ListPushLog] fail, account: ", account, err)
		byt, _ := j.Marshal(util.JsonResult(0, nil, err.Error()))
		this.Ctx.WriteString(string(byt))
		return
	}
	util.LogInfo("[alarm_controller.go ListPushLog] success, account: ", account)
	byt, _ := j.Marshal(util.JsonResult(1, pushLogs, nil))
	this.Ctx.WriteString(string(byt))
}
