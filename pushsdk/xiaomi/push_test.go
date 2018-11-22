package xiaomi

import (
	"testing"
	"xcurrency-push/config"
)

func TestXiaoMi(t *testing.T) {

	config.LoadConfig("../../config/cfg.yaml")

	// 小米
	// regId := "J5Ka5zXthn2pKd3nHmblcTJXR4bdYAs0xt1QVm8JXcrlwdzK5iH8jfkg/pMuPk7+"
	// NENUX
	// regId := "jlz+yEkw4BI+Lv8ZSWuGFbQPBg8CF4WFpBN3AsY8YZuLH0VON0dxJOq8cbB+w35s"
	// 三星
	TPR := "afjzj6tsrOhnw3jr5Fv6BL22wTBrldlDB/00p2l2P++YtZbP8nj3AP7Ub3TInch0"

	pushService := &PushServiceImpl{}
	pushService.DoPush("极简汇率测试", "极简汇率第1次推送测试", TPR)
}
