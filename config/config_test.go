package config

import (
	"fmt"
	// "path/filepath"
	"testing"
	"xcurrency-push/model"
	// "tratao-push/redis"
)

func TestConfig2(t *testing.T) {

	conf := LoadConfig("path")

	fmt.Println(conf.GetString("node.addr"))
	fmt.Println(conf.GetInt("node.rpcTimeout"))
	fmt.Println(conf.GetFloat64("test.percent"))
	fmt.Println(conf.GetBool("test.alived"))
	fmt.Println(conf.GetInt("http.listenAddr"))

	a := make([]model.Alarm, 0, 100)
	fmt.Println(a)

	m := map[string]interface{}{"Ltt": t}
	fmt.Println(m)

	var str string
	fmt.Println(str, len(str))

}
