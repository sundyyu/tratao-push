package check

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"xcurrency-push/config"
	"xcurrency-push/util"
)

const (
	APIURL        = "https://xcr.tratao.com/api/ver2/exchange/yahoo/latest"
	BASECUR       = "USD"
	CUR_SEPARATOR = "/"
)

type ExrateYahoo struct {
	Prices map[string]float64
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

	// 本币对本币汇率为1
	if strings.ToUpper(targetCur) == strings.ToUpper(baseCur) {
		return 1, nil
	}

	// 目标货币是美元
	if strings.ToUpper(targetCur) == BASECUR {
		key = targetCur + CUR_SEPARATOR + baseCur
	}

	// 查找以美元为基础的货币
	if strings.ToUpper(targetCur) == BASECUR || strings.ToUpper(baseCur) == BASECUR {
		if price, ok := prices[key]; ok {

			// 目标货币是美元
			if strings.ToUpper(targetCur) == BASECUR {
				return calculatePrice(price, 1), nil
			}
			return price, nil
		} else {
			return -1, errors.New("No exchange rate for currency [" + key + "] was found.")
		}
	}
	return calculate(prices, baseCur, targetCur)
}

func calculate(prices map[string]float64, base string, target string) (float64, error) {
	baseKey := BASECUR + CUR_SEPARATOR + base
	targetKey := BASECUR + CUR_SEPARATOR + target

	var basePrice float64
	var targetPrice float64
	var ok bool

	if basePrice, ok = prices[baseKey]; !ok {
		return -1, errors.New("No exchange rate for base currency [" + baseKey + "] was found.")
	}
	if targetPrice, ok = prices[targetKey]; !ok {
		return -1, errors.New("No exchange rate for target currency [" + targetKey + "] was found.")
	}
	return calculatePrice(basePrice, targetPrice), nil
}

// 利用高精度类decimal 进行汇率计算
func calculatePrice(base float64, target float64) float64 {
	baseDecimal := decimal.NewFromFloat(base)
	targetDecimal := decimal.NewFromFloat(target)
	dec := targetDecimal.DivRound(baseDecimal, 10)

	price, _ := dec.Float64()
	return price
}

func (yahoo *ExrateYahoo) Update() {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "recover ExrateYahoo Update error.")
		}
	}()

	res, err := http.Get(APIURL)
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
}

func (yahoo *ExrateYahoo) Loop() {
	util.LogInfo("yahoo exrate loop start.")

	cfg := config.GetConfig()
	t := cfg.GetInt("check.exrateTickTime") // 定时执行时间（秒）
	TickLoop(yahoo, t)
}
