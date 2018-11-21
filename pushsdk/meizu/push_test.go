package meizu

import (
	"testing"
	"tratao-push/config"
)

func TestMeiZu(t *testing.T) {
	config.LoadConfig("../../config/cfg.yaml")

	TPR := "UCI4e0f4d7a027f4948017d70416a627f48470c467106"
	pushService := &PushServiceImpl{}
	pushService.DoPush("极简汇率测试", "极简汇率第1次推送测试", TPR)

}
