package huawei

import (
	"fmt"
	"sync"
	"testing"
	"tratao-push/config"
)

func Token(t *testing.T) {
	wg := new(sync.WaitGroup)
	var (
		t1 *token
		t2 *token
	)

	wg.Add(1)
	go func() {
		t1, err := accessToken.get("100384977", "976b53ce5707fff89181f6412418102d")
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("token:", t1)
		}
		wg.Done()
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		accessToken.expireImmediately()
		t2, err := accessToken.get("100384977", "976b53ce5707fff89181f6412418102d")
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("token:", t2)
		}
		wg.Done()
	}()

	wg.Wait()

	if t1 != t2 {
		t.Error()
	}
}

func TestHuaWei(t *testing.T) {

	config.LoadConfig("../../config/cfg.yaml")
	//0866321030619723300002279300CN01    0860983035024104300002279300CN01
	TPR := "0866321030619723300002279300CN01"
	pushService := &PushServiceImpl{}
	pushService.DoPush("极简汇率测试", "极简汇率第1次推送测试", TPR)

}
