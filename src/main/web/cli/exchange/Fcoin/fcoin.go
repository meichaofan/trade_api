package Fcoin

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://developer.fcoin.com/zh.html
需翻墙
*/

const (
	ApiHost = "https://api.fcoin.com/v2"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Fcoin struct {
}

func (c Fcoin) Name() string {
	return "fcoin"
}

//虚拟货币 -- 美元 汇率
func (c Fcoin) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/market/ticker/" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Int()
	if status == 0 {
		rate := ret.Get("data.ticker").Array()[0].Float()
		defer rateCoin.Store(symbol, rate)
		return rate
	}
	return 0
}

func (c Fcoin) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/market/all-tickers"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	ret.Get("tickers").ForEach(func(key, value gjson.Result) bool {
		var quote string
		var base string
		symbol := value.Get("symbol").Str

		if strings.HasSuffix(symbol, "usdt") ||
			strings.HasSuffix(symbol, "tusd") ||
			strings.HasSuffix(symbol, "usdc") {
			quote = symbol[:len(symbol)-4]
			base = symbol[len(symbol)-4:]
		} else {
			//BTC - ETH -TRX
			quote = symbol[:len(symbol)-3]
			base = symbol[len(symbol)-3:]
		}

		timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

		//涨幅
		openPrice := value.Get("ticker").Array()[6].Float() //24小时前成交价
		last := value.Get("ticker").Array()[0].Float()
		pcg := (last - openPrice) / openPrice

		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			AmountQuote:        value.Get("ticker").Array()[9].Float(),
			AmountBase:         value.Get("ticker").Array()[10].Float(),
			Last:               last,
			Time:               timeStr,
			PriceChangePercent: pcg,
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
		return true
	})

	return exchangeTickers
}
