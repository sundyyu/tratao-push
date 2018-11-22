package xiaomi

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

	cli := NewClient(cfg.GetString("xiaomi.appSecret"), []string{cfg.GetString("xiaomi.appPkgName")})
	msg := NewMessage(title, body).
		SetPayload(body).
		SetPassThrough(0).
		StartLauncherActivity()
	resp, err := cli.MultiSend(msg, []string{TPR})

	util.LogInfo("xiaomi push response: ", util.ToJson(resp))
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return errors.New("device [" + TPR + "] for xiaomi push failed.")
	}
	return nil
}
