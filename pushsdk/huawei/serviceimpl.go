package huawei

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

	cli := NewClient(cfg.GetString("huawei.appId"), cfg.GetString("huawei.appSecret"))
	n := NewNotification(title, body).
		AddDeviceToken(TPR).
		StartApp(cfg.GetString("huawei.appPkgName"))

	resp, err := cli.Push(n)

	util.LogInfo("huawei push response: ", resp)
	if err != nil {
		return err
	}

	if resp.Code != "200" {
		return errors.New("device [" + TPR + "] for huawei push failed.")
	}
	return nil
}
