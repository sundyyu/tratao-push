package check

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"tratao-push/config"
	"tratao-push/util"
)

const (
	YAHOO_APIURL  = "https://xcr.tratao.com/api/ver2/exchange/yahoo/latest"
	YAHOO_BASECUR = "USD"
	CUR_SEPARATOR = "/"
)

type ExrateYahoo struct {
	StopChan chan int
	Prices   map[string]float64
}

var lock sync.Mutex

func (yahoo *ExrateYahoo) GetPrice(baseCur string, targetCur string) (float64, error) {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "recover ExrateYahoo GetPrice error.")
		}
	}()

	prices := yahoo.Prices
	key := baseCur + CUR_SEPARATOR + targetCur

	if price, ok := prices[key]; ok {
		return price, nil
	} else {

		baseKey := YAHOO_BASECUR + CUR_SEPARATOR + baseCur
		targetKey := YAHOO_BASECUR + CUR_SEPARATOR + targetCur

		var basePrice float64
		var targetPrice float64

		if basePrice, ok = prices[baseKey]; !ok {
			return -1, errors.New("No exchange rate for base currency was found.")
		}
		if targetPrice, ok = prices[targetKey]; !ok {
			return -1, errors.New("No exchange rate for target currency was found.")
		}

		// 利用高精度类decimal 进行汇率计算
		baseDecimal := decimal.NewFromFloat(basePrice)
		targetDecimal := decimal.NewFromFloat(targetPrice)
		dec := targetDecimal.DivRound(baseDecimal, 10)
		if price, ok := dec.Float64(); ok {
			return price, nil
		}
		return -1, nil
	}
}

func (yahoo *ExrateYahoo) Update() error {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "recover ExrateYahoo Update error.")
		}
	}()

	res, err := http.Get(YAHOO_APIURL)
	util.LogError(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	util.LogError(err)

	if yahoo.Prices == nil {
		yahoo.Prices = make(map[string]float64, 300)
	}

	result := gjson.GetBytes(body, "resources")
	result.ForEach(func(key, resource gjson.Result) bool {
		name := resource.Get("resource.fields.name").String()
		price := resource.Get("resource.fields.price").Float()

		yahoo.Prices[name] = price
		return true
	})
	util.LogInfoF("update exrate at %v\n", time.Now())

	return nil
}

func (yahoo *ExrateYahoo) Loop() {
	util.LogInfo("exrate loop start.")

	// cfg := config.NewConfig("../../config/cfg.yaml")
	cfg := config.GetConfig()
	t := cfg.GetInt("check.exrateTickTime") // 定时执行时间（秒）

	yahoo.StopChan = make(chan int, 1)
	ticker := time.NewTicker(time.Second * time.Duration(t))

	go func() {
		for {
			select {
			case <-yahoo.StopChan:
				util.LogInfo("exrate loop stop.")
				return
			case <-ticker.C:
				yahoo.Update()
			}
		}
	}()
}

func (yahoo *ExrateYahoo) StopLoop() {
	yahoo.StopChan <- 1
}
