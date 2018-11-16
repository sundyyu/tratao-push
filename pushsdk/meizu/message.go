package meizu

type NotificationMsg struct {
  NoticeBarInfo    NoticeBarInfo    `json:"noticeBarInfo"`
  NoticeExpandInfo NoticeExpandInfo `json:"noticeExpandInfo"`
  ClickTypeInfo    ClickTypeInfo    `json:"clickTypeInfo"`
  PushTimeInfo     PushTimeInfo     `json:"pushTimeInfo"`
  AdvanceInfo      AdvanceInfo      `json:"AdvanceInfo"`
  Extra            Extra            `json:"Extra"`
}

type NoticeBarInfo struct {
  NoticeBarType int32  `json:"noticeBarType,omitempty"`                       // 通知栏样式(0, "标准")【int 非必填，值为0】
  Title         string `json:"title"`                                         // 推送标题, 【string 必填，字数限制1~32】
  Content       string `json:"content"`                                       // 推送内容, 【string 必填，字数限制1~100】
}

type NoticeExpandInfo struct {
  NoticeExpandType    int32  `json:"noticeExpandType,omitempty"`              // 展开方式 (0, "标准"),(1, "文本")【int 非必填，值为0、1】
  NoticeExpandContent string `json:"noticeExpandContent"`                     // 展开内容, 【string noticeExpandType为文本时，必填】
}

type ClickTypeInfo struct {
  ClickType       int32  `json:"clickType,omitempty"`                         // 展开方式 (0, "标准"),(1, "文本")【int 非必填，值为0、1】
  Url             string `json:"url"`                                         // URI页面地址, 【string clickType为打开URI页面时，必填, 长度限制1000】
  Parameters      string `json:"parameters,omitempty"`                        // 参数 【JSON格式】【非必填】
  Activity        string `json:"activity"`                                    // 应用页面地址 【string clickType为打开应用页面时，必填, 长度限制1000】
  CustomAttribute string `json:"customAttribute"`                             // 应用客户端自定义【 为应用客户端自定义时，必填， 输入长度为 字节以内】
}

type PushTimeInfo struct {
  OffLine   int8  `json:"offLine,omitempty"`                                  // 是否进离线消息(0 否 1 是[validTime]) 【int 非必填，默认值为1】
  ValidTime int32 `json:"validTime"`                                          // 有效时长 (1 72 小时内的正整数) 【int offLine值为1时，必填，默认24】
}

type AdvanceInfo struct {
  Suspend             int8             `json:"suspend,omitempty"`             // 是否通知栏悬浮窗显示 (1 显示  0 不显示) 【int 非必填，默认1】
  ClearNoticeBar      int8             `json:"clearNoticeBar,omitempty"`      // 是否可清除通知栏 (1 可以  0 不可以) 【int 非必填，默认1】
  FixDisplay          int8             `json:"fixDisplay,omitempty"`          // 是否定时展示 是 否 【 非必填，默认 】
  FixStartDisplayTime string           `json:"fixStartDisplayTime,omitempty"` // 定时展示开始时间 【 非必填】
  FixEndDisplayTime   string           `json:"fixEndDisplayTime,omitempty"`   // 定时展示结束时间 【 非必填】
  NotificationType    NotificationType `json:"notificationType"`
  NotifyKey           string           `json:"notifyKey"`                     // 非必填 默认空 分组合并推送的 凡是带有此 的通知栏消息 只会显示最后到达的一条。由数字大小写字母下划线和中划线组成长度不大于8个字符
}

type NotificationType struct {
  Vibrate int8 `json:"vibrate,omitempty"`                                     // 震动 (0关闭  1 开启), 【int 非必填，默认1】
  Lights  int8 `json:"lights,omitempty"`                                      // 闪光 (0关闭  1 开启), 【int 非必填，默认1】
  Sound   int8 `json:"sound,omitempty"`                                       // 声音 (0关闭  1 开启), 【int 非必填，默认1】
}

type Extra struct {
  Callback string `json:"callback"`                                           // （必填字段），第三方接收回执的HTTP接口，最大长度 字节
  CallbackParam string `json:"callback.param,omitempty"`                      // （可选字段）， 第三方自定义回执参数， 最大长度64字节
  CallbackType int8 `json:"callback.type,omitempty"`                          // （可选字段）， 回执类型（1-送达回执，2-点击回执，3-送达与点击回执，默认3
}

type Response struct {
  Code     string      `json:"code"`                                          // 必选，返回码
  Message  string      `json:"message"`                                       // 可选，返回消息，网页端接口出现错误时使用此消息展示给用户，手机端可 忽略此消息，甚至服务端不传输此消息
  Value    interface{} `json:"value"`                                         // 必选，返回结果
  Redirect string      `json:"redirect"`                                      // 可选，returnCode=300重定向时，使用此URL重新请求
  MsgId    string      `json:"msgId"`                                         // 可选，消息推送msgId
}