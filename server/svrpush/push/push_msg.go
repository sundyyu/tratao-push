package push

import (
	"github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"strings"
	"xcurrency-push/model"
	"xcurrency-push/pushsdk"
	"xcurrency-push/pushsdk/fcm"
	"xcurrency-push/pushsdk/huawei"
	"xcurrency-push/pushsdk/ios"
	"xcurrency-push/pushsdk/meizu"
	"xcurrency-push/pushsdk/xiaomi"
	"xcurrency-push/util"
)

type PushMsgReceive struct {
	CallChan chan int
}

func (rs PushMsgReceive) Call(msg amqp.Delivery) {
	// 占用通道
	rs.CallChan <- 1
	// 线程中执行推送
	go DoPushMsg(msg, rs.CallChan)
}

func DoPushMsg(msg amqp.Delivery, callChan chan int) {
	defer ackMsg(msg)
	defer util.ReadChan(callChan)

	// 获取数据
	data := msg.Body
	// 解析信息
	pushMsg := model.PushMsg{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &pushMsg); err != nil {
		util.LogError(err)
		return
	}

	deviceId := pushMsg.DeviceId
	deviceOS := strings.ToLower(pushMsg.DeviceOS)
	deviceCountry := strings.ToLower(pushMsg.DeviceCountry)
	// deviceLang := pushMsg.DeviceLang

	title := pushMsg.Title
	body := pushMsg.Body
	util.LogInfo("push message: ", pushMsg)

	// 调用推送SDK进行消息推送
	var pushSerice pushsdk.PushService
	if strings.Contains(deviceOS, "huawei") { // 华为推送
		pushSerice = &huawei.PushServiceImpl{}
	} else if strings.Contains(deviceOS, "meizu") { // 魅族推送
		pushSerice = &meizu.PushServiceImpl{}
	} else if strings.Contains(deviceOS, "xiaomi") { // 小米推送
		pushSerice = &xiaomi.PushServiceImpl{}
	} else if strings.Contains(deviceOS, "ios") {
		pushSerice = &ios.PushServiceImpl{}
	} else if strings.Contains(deviceCountry, "cn") { // 中国的其他手机小米推送
		pushSerice = &xiaomi.PushServiceImpl{}
	} else { // 国外的手机FCM推送
		pushSerice = &fcm.PushServiceImpl{}
	}

	// 测试使用
	// deviceId = ""

	if len(deviceId) <= 0 {
		return
	}
	if err := pushSerice.DoPush(title, body, deviceId); err != nil {
		util.LogError(err)
	}

}

// 应答消息
func ackMsg(msg amqp.Delivery) {
	msg.Ack(false)
}
