package fcm

import (
	"errors"
	"strings"
	"xcurrency-push/config"
	"xcurrency-push/util"
)

type PushServiceImpl struct {
}

func (service *PushServiceImpl) DoPush(title string, body string, TPR string) error {
	// cfg := config.LoadConfig("../../config/cfg.yaml")
	cfg := config.GetConfig()

	client := NewClient(cfg.GetString("fcm.cfgPath"))
	n := &Notification{title, body}
	resp, err := client.Push(n, TPR)

	util.LogInfo("fcm push response: ", util.ToJson(resp))
	if err != nil {
		return err
	}

	util.LogInfo(resp.Message)

	if !strings.Contains(resp.Message, "messages/0:") {
		return errors.New("device [" + TPR + "] for FCM push failed.")
	}

	return nil
}
