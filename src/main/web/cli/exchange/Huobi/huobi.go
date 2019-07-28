package Huobi

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档: https://huobiapi.github.io/docs/spot/v1/cn/#5ea2e0cde2
*/

const (
	ApiHost = "https://api.huobi.pro"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Huobi struct {
}

func (c Huobi) Name() string {
	return "huobi"
}

//虚拟货币 -- 美元 汇率
func (c Huobi) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/market/detail/merged?symbol=" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Str
	if status == "ok" {
		//btc - usdt
		//eth - usdt
		rate := ret.Get("tick.close").Float()
		rateCoin.Store(symbol, rate)
		return rate
	}
	return 1;
}

func (c Huobi) PairHandler() []*data.ExchangeTicker {
	cnyUsdRate := common.CalRate("cny")
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/market/tickers"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Str
	timeStr := strconv.FormatInt(ret.Get("ts").Int(), 10)
	if status == "ok" {
		ret.Get("data").ForEach(func(key, value gjson.Result) bool {
			var quote string
			var base string
			symbol := value.Get("symbol").Str
			if strings.HasSuffix(symbol, "usdt") {
				quote = symbol[:len(symbol)-4]
				base = "USDT"
			} else if strings.HasSuffix(symbol, "ht") {
				quote = symbol[:len(symbol)-2]
				base = "HT"
			} else {
				//BTC - ETH -TRX
				quote = symbol[:len(symbol)-3]
				base = symbol[len(symbol)-3:]
			}
			//涨幅
			openPrice := value.Get("open").Float()
			last := value.Get("close").Float()
			pcg := (last - openPrice) / openPrice
			exchangeTicker := &data.ExchangeTicker{
				Symbol:             strings.ToUpper(symbol),
				Quote:              strings.ToUpper(quote),
				Base:               strings.ToUpper(base),
				AmountQuote:        value.Get("amount").Float(), //这里取值注意一下
				AmountBase:         value.Get("vol").Float(),
				Last:               last,
				Time:               timeStr,
				PriceChangePercent: pcg,
			}
			//汇率
			if strings.ToUpper(base) == "USDT" {
				exchangeTicker.LastUsd = exchangeTicker.Last
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase
			} else if strings.ToUpper(base) == "HT" {
				rate := c.GetRate("HT", "USDT")
				exchangeTicker.LastUsd = exchangeTicker.Last * rate
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase * rate
			} else {
				rate := c.GetRate(base, "USDT")
				exchangeTicker.LastUsd = exchangeTicker.Last * rate
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase * rate
			}
			exchangeTicker.AmountCny = exchangeTicker.AmountUsd * cnyUsdRate
			exchangeTicker.LastCny = exchangeTicker.LastUsd * cnyUsdRate
			exchangeTickers = append(exchangeTickers, exchangeTicker)
			return true
		})
	}
	return exchangeTickers
}
