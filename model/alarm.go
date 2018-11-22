package model

import (
	"xcurrency-push/util"
)

type Alarm struct {
	Id            int64   `json:"id"`
	Account       string  `json:"account"`
	BaseCur       string  `json:"basecur"`
	TargetCur     string  `json:"targetcur"`
	Price         float64 `json:"price"`
	LbPrice       float64 `json:"lbprice"`
	UbPrice       float64 `json:"ubprice"`
	Enabled       bool    `json:"enabled"`
	DeviceId      string  `json:"devid"`
	DeviceOS      string  `json:"devos"`
	DeviceCountry string  `json:"devcountry"`
	DeviceLang    string  `json:"devlang"`
	AppKey        string  `json:"appkey"`
	Ltt           int64   `json:"ltt"` // last trigger time
	UpdateTime    int64   `json:"updatetime"`
	CreateTime    int64   `json:"createtime"`
}

func GetFields() []string {
	fields := []string{
		"Id",
		"Account",
		"BaseCur",
		"TargetCur",
		"Price",
		"LbPrice",
		"UbPrice",
		"Enabled",
		"DeviceId",
		"DeviceOS",
		"DeviceCountry",
		"DeviceLang",
		"AppKey",
		"Ltt",
		"UpdateTime",
		"CreateTime"}

	return fields
}

func ResultToAlarm(result []interface{}) Alarm {
	alarm := Alarm{}

	if alarm_id, err := util.RedisResult2Int64(result[0]); err == nil {
		alarm.Id = alarm_id
	}
	if account, ok := result[1].(string); ok {
		alarm.Account = account
	}
	if baseCur, ok := result[2].(string); ok {
		alarm.BaseCur = baseCur
	}
	if targetCur, ok := result[3].(string); ok {
		alarm.TargetCur = targetCur
	}
	if price, err := util.RedisResult2Float64(result[4]); err == nil {
		alarm.Price = price
	}
	if lbprice, err := util.RedisResult2Float64(result[5]); err == nil {
		alarm.LbPrice = lbprice
	}
	if ubprice, err := util.RedisResult2Float64(result[6]); err == nil {
		alarm.UbPrice = ubprice
	}
	if enable, err := util.RedisResult2Bool(result[7]); err == nil {
		alarm.Enabled = enable
	}
	if deviceId, ok := result[8].(string); ok {
		alarm.DeviceId = deviceId
	}
	if deviceOS, ok := result[9].(string); ok {
		alarm.DeviceOS = deviceOS
	}
	if deviceCountry, ok := result[10].(string); ok {
		alarm.DeviceCountry = deviceCountry
	}
	if deviceLang, ok := result[11].(string); ok {
		alarm.DeviceLang = deviceLang
	}
	if appKey, ok := result[12].(string); ok {
		alarm.AppKey = appKey
	}
	if ltt, err := util.RedisResult2Int64(result[13]); err == nil {
		alarm.Ltt = ltt
	}
	if updateTime, err := util.RedisResult2Int64(result[14]); err == nil {
		alarm.UpdateTime = updateTime
	}
	if createTime, err := util.RedisResult2Int64(result[15]); err == nil {
		alarm.CreateTime = createTime
	}
	return alarm
}
