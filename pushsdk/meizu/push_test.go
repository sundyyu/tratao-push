package meizu

import (
	"testing"
	"xcurrency-push/config"
)

func TestMeiZu(t *testing.T) {
	config.LoadConfig("../../config/cfg.yaml")

	TPR := "S5Q4b726e754068797c5e455a040007585c4271657a42"
	pushService := &PushServiceImpl{}
	pushService.DoPush("极简汇率测试", "极简汇率第1次推送测试", TPR)

}
