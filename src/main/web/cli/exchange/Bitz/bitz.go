package Bitz

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://apidoc.bitz.com/cn/
*/

const (
	ApiHost = "https://api.bitzapi.com"
	//APiHost = "https://apiv2.bitz.com" 都可用
)

type Bitz struct {
}

func (b Bitz) Name() string {
	return "bitz"
}

/**
注：行情接口里有汇率 赞 ^_^
*/
func (b Bitz) GetRate(quote, base string) float64 {
	return 0
}

/**
交易对：https://apidoc.bitz.com/cn/market-quotation-data/Get-the-tickerall.html
*/
func (b Bitz) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/Market/tickerall"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Int()
	if status == 200 {
		symbols := ret.Get("data").Map()
		time := strconv.FormatInt(ret.Get("time").Int(), 10)
		for symbol, detail := range symbols {
			quote := strings.Split(symbol, "_")[0]
			base := strings.Split(symbol, "_")[1]
			last := detail.Get("now").Float()
			usd := detail.Get("usd").Float()
			cny := detail.Get("cny").Float()
			rate := usd / last
			cnyRate := cny / last
			exchangeTicker := &data.ExchangeTicker{
				Symbol:             strings.ToUpper(symbol),
				Quote:              strings.ToUpper(quote),
				Base:               strings.ToUpper(base),
				AmountQuote:        detail.Get("volume").Float(),
				AmountBase:         detail.Get("quoteVolume").Float(),
				AmountUsd:          detail.Get("quoteVolume").Float() * rate,
				AmountCny:          detail.Get("quoteVolume").Float() * cnyRate,
				Last:               last,
				LastUsd:            usd,
				LastCny:            cny,
				PriceChangePercent: detail.Get("priceChange").Float(),
				Time:               time,
			}
			exchangeTickers = append(exchangeTickers, exchangeTicker)
		}
	} else {
		common.MessageHandler(b.Name(), ret.Get("msg").Str)
	}
	return exchangeTickers
}
