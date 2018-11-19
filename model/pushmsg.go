package model

import (
	"tratao-push/util"
)

type PushMsg struct {
	Id         int64  `json:"id"`
	Account    string `json:"account"`
	DeviceId   string `json:"devid"`
	Message    string `json:"message"`
	CreateTime int64  `json:"createtime"`
}

func GetPushMsgFields() []string {
	fields := []string{
		"Id",
		"Account",
		"DeviceId",
		"Message",
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
	if message, ok := result[3].(string); ok {
		msg.Message = message
	}

	if createTime, err := util.RedisResult2Int64(result[4]); err == nil {
		msg.CreateTime = createTime
	}
	return msg
}
