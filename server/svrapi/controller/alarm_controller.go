package controller

import (
	"github.com/astaxie/beego"
	"github.com/tidwall/gjson"
	"time"
	"tratao-push/model"
	"tratao-push/redis"
	"tratao-push/util"
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
	account := gjson.GetBytes(reqBody, "account").String()
	if len(account) == 0 {
		util.LogError("[alarm_controller.go UpdateAlarm] fail, account: ", account, "account can not be empty.")

		this.Data["json"] = util.JsonResult(0, nil, "account error. account can not be empty.")
		this.ServeJSON()
		return
	}
	alarm.Account = account
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
