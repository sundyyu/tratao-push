package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"strings"
	"sync"
	"xcurrency-push/util"
)

/**
* @des 根据go-config进行二次封装，生成GetStirng, GetInt, GetFloat64, GetBool四种简易接口
* @author 于朝鹏
* @date 2018年11月2日 18:10
 */
type TConfig struct {
	Config config.Config
}

var lock *sync.Mutex = new(sync.Mutex)
var tcfg *TConfig

// init config
func LoadConfig(path string) *TConfig {
	if tcfg == nil {
		lock.Lock()
		defer lock.Unlock()

		if tcfg == nil {
			tcfg = &TConfig{}
			tcfg.Config = config.NewConfig()

			// Load json config file
			tcfg.Config.Load(file.NewSource(
				file.WithPath(path),
			))

			util.LogInfo("config init")
		}
	}

	return tcfg
}

func GetConfig() *TConfig {
	return tcfg
}

// GetString("ip") or GetString("node.ip")
func (t TConfig) GetString(key string) string {
	if strings.Contains(key, ".") {
		return t.Config.Get(strings.Split(key, ".")...).String("")
	} else {
		return t.Config.Get("hosts", "database", "address").String("")
	}
}

// GetInt("port") or GetInt("node.port")
func (t TConfig) GetInt(key string) int {
	if strings.Contains(key, ".") {
		return t.Config.Get(strings.Split(key, ".")...).Int(-999)
	} else {
		return t.Config.Get("hosts", "database", "address").Int(-999)
	}
}

// GetFloat64("node.flt") or GetFloat64("node.flt")
func (t TConfig) GetFloat64(key string) float64 {
	if strings.Contains(key, ".") {
		return t.Config.Get(strings.Split(key, ".")...).Float64(-999)
	} else {
		return t.Config.Get("hosts", "database", "address").Float64(-999)
	}
}

// GetFloat64("lb") or GetFloat64("node.bl")
func (t TConfig) GetBool(key string) bool {
	if strings.Contains(key, ".") {
		return t.Config.Get(strings.Split(key, ".")...).Bool(false)
	} else {
		return t.Config.Get("hosts", "database", "address").Bool(false)
	}
}
