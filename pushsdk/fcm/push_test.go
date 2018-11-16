package fcm

import (
	"testing"
)

func TestFCM(t *testing.T) {

	TPR := "drH1k3PVhgk:APA91bERRAfrZx8M0NneyS1kYVkmDSCcJdfUvHsv1urEzTD88QVeu8SfMX7dQH16D1EOcmlfwCoJuJueMV4EhNEMAPXMzir73nP-9EVx_E9m5OkhfGa6yrtdxii62bYB26tTRsEZZlmc"
	pushService := &PushServiceImpl{}
	pushService.DoPush("极简汇率测试", "极简汇率第1次推送测试", TPR)

}

// ea1981d09b94d0a3e601fc369eed18c3ffa94ee492353ab8d05b9a8548b3a108
