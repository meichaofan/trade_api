package Cointiger

import (
	"github.com/tidwall/gjson"
	"strings"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://github.com/cointiger/api-docs/wiki
*/

const (
	//所有的symbol信息
	TickerAllHost = "https://www.cointiger.one/exchange/api/public/market/detail"
	//单个symbol
	TickerHost = "https://www.cointiger.one/exchange/trading/api/market/detail"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Cointiger struct {
}

func (c Cointiger) Name() string {
	return "cointiger"
}

//虚拟货币 -- 美元 汇率
func (c Cointiger) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := TickerHost + "?symbol=" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	code := ret.Get("code").Str
	if code == "0" {
		//btc - usdt
		//eth - usdt
		rate := ret.Get("data.trade_ticker_data.tick.close").Float()
		rateCoin.Store(symbol, rate)
		return rate
	}
	return 1
}

func (c Cointiger) PairHandler() []*data.ExchangeTicker {
	cnyUsdRate := common.CalRate("cny")
	var exchangeTickers []*data.ExchangeTicker
	url := TickerAllHost
	content := common.HttpGet(url)
	symbols := gjson.ParseBytes(content).Map()

	for symbol, value := range symbols {
		var quote string
		var base string
		s := symbol
		if strings.HasSuffix(s, "USDT") {
			quote = symbol[:len(symbol)-4]
			base = "USDT"
		} else if strings.HasSuffix(s, "BITCNY") {
			quote = symbol[:len(symbol)-6]
			base = "BITCNY"
		} else {
			//BTC - ETH -TRX
			quote = symbol[:len(symbol)-3]
			base = symbol[len(symbol)-3:]
		}

		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			AmountQuote:        value.Get("baseVolume").Float(),
			AmountBase:         value.Get("quoteVolume").Float(),
			Last:               value.Get("last").Float(),
			Time:               value.Get("id").Str, // id就是time
			PriceChangePercent: value.Get("percentChange").Float(),
		}

		//汇率
		if strings.ToUpper(base) == "USDT" {
			exchangeTicker.LastUsd = exchangeTicker.Last
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase
		} else if strings.ToUpper(base) == "BITCNY" { //人民币
			exchangeTicker.LastUsd = exchangeTicker.Last / cnyUsdRate
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase / cnyUsdRate
		} else {
			rate := c.GetRate(base, "USDT")
			exchangeTicker.LastUsd = exchangeTicker.Last * rate
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase * rate
		}
		exchangeTicker.AmountCny = exchangeTicker.AmountUsd * cnyUsdRate
		exchangeTicker.LastCny = exchangeTicker.LastUsd * cnyUsdRate
		exchangeTickers = append(exchangeTickers, exchangeTicker)
	}
	return exchangeTickers
}
