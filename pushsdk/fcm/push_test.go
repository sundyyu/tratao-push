package fcm

import (
	"testing"
	"xcurrency-push/config"
)

func TestFCM(t *testing.T) {

	config.LoadConfig("../../config/cfg.yaml")

	TPR := "f2IRZl96bC0:APA91bH7E8pize0YtewfCoCEIJzhK1asni44iR1kDL5h0XG9NFDxm-1SomAUy-nDwqc7NTA465Q19LF4YNcWzCdF07H-Q8c4BSyos2dy9SqK9mTB4vN5gzKxnnik7CnFiAKG5VtC77V1"
	pushService := &PushServiceImpl{}
	pushService.DoPush("极简汇率测试", "极简汇率第1次推送测试", TPR)

}

// ea1981d09b94d0a3e601fc369eed18c3ffa94ee492353ab8d05b9a8548b3a108
