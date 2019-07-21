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

/**
type Exchange interface {
	//交易所名称
	Name() string
	//获取交易对的汇率
	GetRate(quote, base string) float64
	//交易所所有交易对
	PairHandler() []*data.ExchangeTicker
	//交易所交易额
	AmountHandler() []*data.TradeData
}
*/

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
			rate := usd / last
			exchangeTicker := &data.ExchangeTicker{
				Symbol:             strings.ToUpper(symbol),
				Quote:              strings.ToUpper(quote),
				Base:               strings.ToUpper(base),
				Volume:             detail.Get("volume").Float(),
				Amount:             detail.Get("quoteVolume").Float(),
				AmountUsd:          detail.Get("quoteVolume").Float() * rate,
				Last:               last,
				LastUsd:            usd,
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

/**
交易量 ：https://apidoc.bitz.com/cn/market-quotation-data/Get-the-tickerall.html
 */
func (b Bitz) AmountHandler() []*data.TradeData {
	var trades []*data.TradeData
	url := ApiHost + "/Market/tickerall"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Int()
	if status == 200 {
		symbols := ret.Get("data").Map()
		for symbol, detail := range symbols {
			quote := strings.Split(symbol, "_")[0]
			base := strings.Split(symbol, "_")[1]
			last := detail.Get("now").Float()
			usd := detail.Get("usd").Float()
			rate := usd / last
			tradeData := &data.TradeData{
				Symbol:    strings.ToUpper(symbol),
				Quote:     strings.ToUpper(quote),
				Base:      strings.ToUpper(base),
				Volume:    detail.Get("volume").Float(),
				Amount:    detail.Get("quoteVolume").Float(),
				AmountUsd: detail.Get("quoteVolume").Float() * rate,
				Price:     last,
				PriceUsd:  usd,
			}
			trades = append(trades, tradeData)
		}
	} else {
		common.MessageHandler(b.Name(), ret.Get("msg").Str)
	}
	return trades
}
