package meizu

import (
	"errors"
	"xcurrency-push/config"
	"xcurrency-push/util"
)

type PushServiceImpl struct {
}

func (service *PushServiceImpl) DoPush(title string, body string, TPR string) error {
	// cfg := config.NewConfig("../../config/cfg.yaml")
	cfg := config.GetConfig()

	cli := NewClient(cfg.GetString("meizu.appId"), cfg.GetString("meizu.appSecret"))
	msg := NotificationMsg{}
	msg.NoticeBarInfo.Title = title
	msg.NoticeBarInfo.Content = body
	msg.PushTimeInfo.OffLine = 0
	msg.PushTimeInfo.ValidTime = 1
	resp, err := cli.Push(&msg, []string{TPR})

	util.LogInfo("meizu push response: ", util.ToJson(resp))
	if err != nil {
		return err
	}

	util.LogInfo(resp)

	if resp.Code != "200" {
		return errors.New("device [" + TPR + "] for meizu push failed.")
	}
	return nil
}
