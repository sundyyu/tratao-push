package xiaomi

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	MaxTimeToSend = time.Hour * 24 * 7
	MaxTimeToLive = time.Hour * 24 * 7 * 2
)

type Message struct {
	RestrictedPkgName string            `json:"restricted_package_name,omitempty"` // 设置app的多包名packageNames（多包名发送广播消息）。p
	Payload           string            `json:"payload,omitempty"`                 // 消息内容payload
	Title             string            `json:"title,omitempty"`                   // 通知栏展示的通知的标题
	Desc              string            `json:"description,omitempty"`             // 通知栏展示的通知的描述
	PassThrough       int32             `json:"pass_through"`                      // 是否通过透传的方式送给app，1表示透传消息，0表示通知栏消息。
	NotifyType        int32             `json:"notify_type,omitempty"`             // DEFAULT_ALL = -1; DEFAULT_SOUND  = 1;   // 使用默认提示音提示 DEFAULT_VIBRATE = 2;   // 使用默认震动提示 DEFAULT_LIGHTS = 4;    // 使用默认led灯光提示
	TTL               int64             `json:"time_to_live,omitempty"`            // 可选项。如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
	TimeToSend        int64             `json:"time_to_send,omitempty"`            // 可选项。定时发送消息。timeToSend是以毫秒为单位的时间戳。注：仅支持七天内的定时消息。
	NotifyId          int64             `json:"notify_id"`                         // 可选项。默认情况下，通知栏只显示一条推送消息。如果通知栏要显示多条推送消息，需要针对不同的消息设置不同的notify_id（相同notify_id的通知栏消息会覆盖之前的）。
	Extra             map[string]string `json:"extra,omitempty"`                   // 可选项，对app提供一些扩展的功能，请参考2.2。除了这些扩展功能，开发者还可以定义一些key和value来控制客户端的行为。注：key和value的字符数不能超过1024，至多可以设置10个key-value键值对。
}

type BaseResponse struct {
	Result      string `json:"result"`
	MessageID   string `json:"trace_id"`
	Code        int64  `json:"code"`
	Description string `json:"description,omitempty"`
	Info        string `json:"info,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

type Response struct {
	BaseResponse
	Data struct {
		ID string `json:"id,omitempty"`
	} `json:"data,omitempty"`
}

func (m *Message) SetRestrictedPackageName(pkgNames []string) *Message {
	m.RestrictedPkgName = strings.Join(pkgNames, ",")
	return m
}

func (m *Message) SetPassThrough(passThrough int32) *Message {
	m.PassThrough = passThrough
	return m
}

func (m *Message) SetNotifyType(notifyType int32) *Message {
	m.NotifyType = notifyType
	return m
}

func (m *Message) SetTimeToSend(tts int64) *Message {
	if time.Since(time.Unix(0, tts*int64(time.Millisecond))) > MaxTimeToSend {
		m.TimeToSend = time.Now().Add(MaxTimeToSend).UnixNano() / 1e6
	} else {
		m.TimeToSend = tts
	}
	return m
}

func (m *Message) SetTimeToLive(ttl int64) *Message {
	if time.Since(time.Unix(0, ttl*int64(time.Millisecond))) > MaxTimeToLive {
		m.TTL = time.Now().Add(MaxTimeToLive).UnixNano() / 1e6
	} else {
		m.TTL = ttl
	}
	return m
}

func (m *Message) SetNotifyId(notifyId int64) *Message {
	m.NotifyId = notifyId
	return m
}

func (m *Message) EnableFlowControl() *Message {
	m.Extra["flow_control"] = "1"
	return m
}

func (m *Message) DisableFlowControl() *Message {
	delete(m.Extra, "flow_control")
	return m
}

// 开发者在发送消息时可以设置消息的组ID（JobKey），带有相同的组ID的消息会被聚合为一个消息组。
// 系统支持按照消息组展示消息详情以及计划推送／送达数量／送达曲线等统计信息。
// 另外，相同JobKey的消息在客户端会进行去重，只展示其中的第一条。
// 这样如果发送时同JobKey中不慎有重复的设备也不用担心用户会收到重复的通知。
func (m *Message) SetJobKey(jobKey string) *Message {
	m.Extra["jobkey"] = jobKey
	return m
}

// 小米推送服务器每隔1s将已送达或已点击的消息ID和对应设备的regid或alias通过调用第三方
// http接口传给开发者。
func (m *Message) SetCallback(callbackURL string) *Message {
	m.Extra["callback"] = callbackURL
	m.Extra["callback.type"] = "3" // 1:送达回执, 2:点击回执, 3:送达和点击回执
	return m
}

// 添加自定义字段, 客户端使用
func (m *Message) AddExtra(key, value string) *Message {
	m.Extra[key] = value
	return m
}

func (m *Message) JSON() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (m *Message) Form() url.Values {
	f := url.Values{}
	if m.TTL > 0 {
		f.Add("time_to_live", strconv.FormatInt(m.TTL, 10))
	}
	if len(m.Payload) > 0 {
		f.Add("payload", m.Payload)
	}
	if len(m.Title) > 0 {
		f.Add("title", m.Title)
	}
	if len(m.Desc) > 0 {
		f.Add("description", m.Desc)
	}
	f.Add("notify_type", strconv.FormatInt(int64(m.NotifyType), 10))
	f.Add("pass_through", strconv.FormatInt(int64(m.PassThrough), 10))
	if m.NotifyId != 0 {
		f.Add("notify_id", strconv.FormatInt(int64(m.NotifyId), 10))
	}
	if m.TimeToSend > 0 {
		f.Add("time_to_send", strconv.FormatInt(int64(m.TimeToSend), 10))
	}
	if len(m.Extra) > 0 {
		for k, v := range m.Extra {
			f.Add("extra."+k, v)
		}
	}
	return f
}

// 发送给Android设备的Message对象
func NewMessage(title string, desc string) *Message {
	return &Message{
		Payload:     "",
		Title:       title,
		Desc:        desc,
		PassThrough: 0,
		NotifyType:  -1, // default notify type
		TTL:         0,
		TimeToSend:  0,
		NotifyId:    0,
		Extra:       make(map[string]string),
	}
}

// 打开当前app对应的Launcher Activity。
func (m *Message) StartLauncherActivity() *Message {
	m.Extra["notify_effect"] = "1"
	return m
}

// 打开当前app内的任意一个Activity。
func (m *Message) StartActivity(value string) *Message {
	m.Extra["notify_effect"] = "2"
	m.Extra["intent_uri"] = value
	return m
}

// 打开网页
func (m *Message) OpenUrl(value string) *Message {
	m.Extra["notify_effect"] = "3"
	m.Extra["web_uri"] = value
	return m
}

func (m *Message) SetPayload(payload string) *Message {
	m.Payload = payload
	return m
}
