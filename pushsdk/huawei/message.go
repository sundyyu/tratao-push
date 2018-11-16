package huawei

import (
	"encoding/json"
	"net/url"
	"time"
)

type Notification struct {
	deviceTokenList []string
	expireTime      *time.Time
	payload         payload
}

type payload struct {
	Hps hps `json:"hps"`
}

type hps struct {
	Msg message   `json:"msg"`           // Push消息定义。包括：消息类型, 消息内容, 消息动作
	Ext extension `json:"ext,omitempty"` // 扩展信息，含BI消息统计，特定展示风格，消息折叠。
}

type message struct {
	Type   int8          `json:"type"`             // 取值含义和说明：1 透传异步消息; 3 系统通知栏异步消息. 注意：2和4以后为保留后续扩展使用
	Body   messageBody   `json:"body"`             // 消息内容。注意：对于透传类的消息可以是字符串，不必是JSON Object。
	Action messageAction `json:"action,omitempty"` // 消息点击动作
}

type messageBody struct {
	Content string `json:"content"` // 消息内容体
	Title   string `json:"title"`   // 消息标题
}

type messageAction struct {
	Type  int8               `json:"type,omitempty"`  // 1 自定义行为：行为由参数intent定义; 2 打开URL：URL地址由参数url定义; 3 打开APP：默认值，打开App的首页. 注意：富媒体消息开放API不支持。
	Param messageActionParam `json:"param,omitempty"` // 关于消息点击动作的参数
}

type messageActionParam struct {
	Intent     string `json:"intent,omitempty"`     // Action的type为1的时候表示自定义行为。开发者可以自定义Intent，用户收到通知栏消息后点击通知栏消息打开应用定义的这个Intent页面，该过程通过context.startActivity(intent)实现，具体代码如
	Url        string `json:"url,omitempty"`        // Action的type为2的时候表示打开URL地址
	AppPkgName string `json:"appPkgName,omitempty"` // 需要拉起的应用包名，必须和注册推送的包名一致。
}

type extension struct {
	BiTag     string                 `json:"biTag,omitempty"`     // 设置消息标签，如果带了这个标签，会在回执中推送给CP用于检测某种类型消息的到达率和状态。注意：BigTag不能携带下面几个保留字符：逗号‘，’，竖线‘|’，长度不能超过100个字符。
	Customize map[string]interface{} `json:"customize,omitempty"` // 扩展样例：[{"season":"Spring"},{"weather":"raining"}]. 说明：这个字段类型必须是JSON Array，里面是key-value的一组扩展信息。
}

type Response struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
}

func NewNotification(title string, content string) *Notification {
	n := Notification{}
	n.deviceTokenList = make([]string, 0)
	n.payload.Hps.Msg.Type = 3
	n.payload.Hps.Msg.Body.Title = title
	n.payload.Hps.Msg.Body.Content = content
	return &n
}

func (n *Notification) AddDeviceToken(token string) *Notification {
	n.deviceTokenList = append(n.deviceTokenList, token)
	return n
}

func (n *Notification) ExpireIn(t time.Duration) *Notification {
	expTime := time.Now().Add(t)
	n.expireTime = &expTime
	return n
}

func (n *Notification) ExpireAt(t time.Time) *Notification {
	n.expireTime = &t
	return n
}

func (n *Notification) StartIntent(intent string) *Notification {
	n.payload.Hps.Msg.Action.Type = 1
	n.payload.Hps.Msg.Action.Param.Intent = intent
	return n
}

func (n *Notification) OpenUrl(url string) *Notification {
	n.payload.Hps.Msg.Action.Type = 2
	n.payload.Hps.Msg.Action.Param.Url = url
	return n
}

func (n *Notification) StartApp(appPkgName string) *Notification {
	n.payload.Hps.Msg.Action.Type = 3
	n.payload.Hps.Msg.Action.Param.AppPkgName = appPkgName
	return n
}

func (n *Notification) Form() url.Values {
	m := url.Values{}

	jsonStr, _ := json.Marshal(n.deviceTokenList)
	m.Add("device_token_list", string(jsonStr))

	if n.expireTime != nil {
		timeStr := n.expireTime.Format("2006-01-02T15:04")
		m.Add("expire_time", timeStr)
	}

	jsonStr, _ = json.Marshal(n.payload)
	m.Add("payload", string(jsonStr))

	return m
}
