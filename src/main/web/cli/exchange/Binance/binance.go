package Binance

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
需要翻墙
*/

const (
	ApiHost = "https://api.binance.com"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
	one      sync.Once
)

type Binance struct {
}

func (c Binance) Name() string {
	return "binance"
}

func initCoinRate() {
	url := ApiHost + "/api/v3/ticker/price"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	ret.ForEach(func(key, value gjson.Result) bool {
		symbol := value.Get("symbol").Str
		if strings.HasSuffix(symbol, "USDT") {
			rate := value.Get("price").Float()
			rateCoin.Store(symbol, rate)
		}
		return true
	})
}

//虚拟货币 -- 美元 汇率
func (c Binance) GetRate(quote, base string) float64 {
	//initCoinRate只执行一次
	one.Do(initCoinRate)
	symbol := strings.ToUpper(quote + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	return 1
}

func (c Binance) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/api/v1/ticker/24hr"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	ret.ForEach(func(key, value gjson.Result) bool {
		var quote string
		var base string
		var symbol = value.Get("symbol").Str

		// USDT TUSD USDC USDS
		if strings.HasSuffix(symbol, "USDT") ||
			strings.HasSuffix(symbol, "TUSD") ||
			strings.HasSuffix(symbol, "USDC") ||
			strings.HasSuffix(symbol, "USDS") {
			quote = symbol[:len(symbol)-4]
			base = symbol[len(symbol)-4:]
		} else {
			//BTC - ETH -TRX
			quote = symbol[:len(symbol)-3]
			base = symbol[len(symbol)-3:]
		}

		//quote
		amountQuote := value.Get("volume").Float()
		//base
		amountBase := value.Get("quoteVolume").Float()
		//last
		last := value.Get("lastPrice").Float()
		//pcg
		pcg := value.Get("priceChangePercent").Float()
		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			AmountQuote:        amountQuote,
			AmountBase:         amountBase,
			Last:               last,
			Time:               strconv.FormatInt(value.Get("closeTime").Int(), 10),
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
