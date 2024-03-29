package Mxc

import (
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://github.com/mxcdevelop/APIDoc
这个交易所脚本goland接口请求异常，但是浏览器没有问题.....
*/

const (
	ApiHost = " https://www.mxc.com/open/api/v1"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Mxc struct {
}

func (c Mxc) Name() string {
	return "mxc"
}

//虚拟货币 -- 美元 汇率
func (c Mxc) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + "_" + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/data/ticker?market=" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	code := ret.Get("code").Int()
	if code == 200 {
		rate := ret.Get("data.last").Float()
		rateCoin.Store(symbol, rate)
		return rate
	} else {
		common.ErrorHandler(errors.New(ret.Get("msg").Str))
	}
	return 1;
}

func (c Mxc) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/data/ticker"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	code := ret.Get("code").Int()
	if code == 200 {
		symbols := ret.Get("data").Map()
		for symbol, value := range symbols {
			s := symbol
			quote := strings.Split(s, "_")[0]
			base := strings.Split(s, "_")[1]
			timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

			amountQuote := value.Get("volume").Float()
			highPrice := value.Get("high").Float()
			lowPrice := value.Get("low").Float()
			amountBase := amountQuote * (highPrice + lowPrice) / 2

			exchangeTicker := &data.ExchangeTicker{
				Symbol:             strings.ToUpper(symbol),
				Quote:              strings.ToUpper(quote),
				Base:               strings.ToUpper(base),
				AmountQuote:        amountQuote,
				AmountBase:         amountBase,
				Last:               value.Get("last").Float(),
				Time:               timeStr,
				PriceChangePercent: value.Get("percentChange").Float(),
			}

			//汇率
			if strings.ToUpper(base) == "USDT" {
				exchangeTicker.LastUsd = exchangeTicker.Last
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase
			} else {
				rate := c.GetRate(base, "USDT")
				exchangeTicker.LastUsd = exchangeTicker.Last * rate
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase * rate
			}
			exchangeTicker.AmountCny = exchangeTicker.AmountUsd * common.CnyUsdRate
			exchangeTicker.LastCny = exchangeTicker.LastUsd * common.CnyUsdRate
			exchangeTickers = append(exchangeTickers, exchangeTicker)
		}
	} else {
		common.ErrorHandler(errors.New(ret.Get("msg").Str))
	}
	return exchangeTickers
}
