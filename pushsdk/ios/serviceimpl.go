package ios

import (
	"errors"
	"tratao-push/config"
	"tratao-push/util"
)

type PushServiceImpl struct {
}

func (service *PushServiceImpl) DoPush(title string, body string, TPR string) error {
	// cfg := config.NewConfig("../../config/cfg.yaml")
	cfg := config.GetConfig()

	client := NewClient(cfg.GetString("ios.cerPath"), cfg.GetString("ios.cerPass"))
	n := &Notification{title, body}
	resp, err := client.Push(n, TPR)

	util.LogInfo("ios push response: ", util.ToJson(resp))
	if err != nil {
		return err
	}

	if len(resp.Message) <= 0 {
		return errors.New("device [" + TPR + "] for IOS push failed.")
	}

	return nil
}
