package model

import (
	"time"
	"xcurrency-push/util"
)

type PushMsg struct {
	Id            int64     `json:"id"`
	Account       string    `json:"account"`
	DeviceId      string    `json:"token"`
	DeviceOS      string    `json:"os"`
	DeviceCountry string    `json:"country"`
	DeviceLang    string    `json:"lang"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	CreateTime    time.Time `json:"createtime" time_format:"2006-01-02 15:04:05" time_location:"UTC"`
}

func GetPushMsgFields() []string {
	fields := []string{
		"Id",
		"Account",
		"DeviceId",
		"DeviceOS",
		"DeviceCountry",
		"DeviceLang",
		"Title",
		"Body",
		"CreateTime"}

	return fields
}

func ResultToPushMsg(result []interface{}) PushMsg {
	msg := PushMsg{}

	if mid, err := util.RedisResult2Int64(result[0]); err == nil {
		msg.Id = mid
	}
	if account, ok := result[1].(string); ok {
		msg.Account = account
	}
	if deviceId, ok := result[2].(string); ok {
		msg.DeviceId = deviceId
	}
	if deviceOS, ok := result[3].(string); ok {
		msg.DeviceOS = deviceOS
	}
	if deviceCountry, ok := result[4].(string); ok {
		msg.DeviceCountry = deviceCountry
	}
	if deviceLang, ok := result[5].(string); ok {
		msg.DeviceLang = deviceLang
	}
	if title, ok := result[6].(string); ok {
		msg.Title = title
	}
	if body, ok := result[7].(string); ok {
		msg.Body = body
	}

	// if createTime, err := util.RedisResult2Int64(result[5]); err == nil {
	// 	msg.CreateTime = createTime
	// }
	return msg
}
