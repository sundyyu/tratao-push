package push

import (
	"github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"strconv"
	"strings"
	"xcurrency-push/model"
	"xcurrency-push/module/pgsql"
	"xcurrency-push/pushsdk"
	"xcurrency-push/pushsdk/fcm"
	"xcurrency-push/pushsdk/huawei"
	"xcurrency-push/pushsdk/ios"
	"xcurrency-push/pushsdk/meizu"
	"xcurrency-push/pushsdk/xiaomi"
	"xcurrency-push/util"
)

type AlarmReceive struct {
	CallChan chan int
}

func (rs AlarmReceive) Call(msg amqp.Delivery) {
	// 占用通道
	rs.CallChan <- 1
	// 线程中执行推送
	go DoPush(msg, rs.CallChan)
}

func DoPush(msg amqp.Delivery, callChan chan int) {
	defer ack(msg)
	defer util.ReadChan(callChan)

	// 获取数据
	data := msg.Body
	// 解析信息
	alarm := model.Alarm{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &alarm); err != nil {
		util.LogError(err)
		return
	}

	baseCur := alarm.BaseCur
	targetCur := alarm.TargetCur
	price := alarm.Price

	deviceId := alarm.DeviceId
	deviceOS := strings.ToLower(alarm.DeviceOS)
	// deviceLang := alarm.DeviceLang
	deviceCountry := strings.ToLower(alarm.DeviceCountry)

	title := "极简汇率"
	p := strconv.FormatFloat(price, 'E', -1, 64)
	body := "您关注的汇率[" + baseCur + "/" + targetCur + "] 当前值为：" + p + "," + " 已达到你设置的预警阈值。"
	util.LogInfo("push alarm: ", alarm)

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
	if err := pushSerice.DoPush(title, body, deviceId); err == nil {
		pushmsg := model.PushMsg{}
		pushmsg.Title = title
		pushmsg.Body = body
		pushmsg.Account = alarm.Account
		pushmsg.DeviceId = alarm.DeviceId
		pushmsg.DeviceOS = alarm.DeviceOS
		pushmsg.DeviceCountry = alarm.DeviceCountry
		pushmsg.DeviceLang = alarm.DeviceLang
		pgclient.InsertPushMsg(pushmsg)
	} else {
		util.LogError(err)
	}
}

// 应答消息
func ack(msg amqp.Delivery) {
	msg.Ack(false)
}
