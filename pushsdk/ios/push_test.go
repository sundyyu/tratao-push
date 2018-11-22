package ios

import (
	"testing"
	"xcurrency-push/util"
)

func TestIOS(t *testing.T) {

	TPR := "1b75335f6d40eab6ebff39ba5ebfedb8e291f784f025e2617758a29e24a0d763"
	pushService := &PushServiceImpl{}
	err := pushService.DoPush("极简汇率测试3", "IOS第3次推送测试", TPR)
	util.LogError(err)
}

// ea1981d09b94d0a3e601fc369eed18c3ffa94ee492353ab8d05b9a8548b3a108
