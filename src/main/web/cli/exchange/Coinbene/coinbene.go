package Coinbene

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://github.com/Coinbene/API-Documents-CHN/wiki
*/

const (
	ApiHost = "http://api.coinbene.com/v1/market"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Coinbene struct {
}

func (c Coinbene) Name() string {
	return "coinbene"
}

//虚拟货币 -- 美元 汇率
func (c Coinbene) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/ticker?symbol=" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Str
	if status == "ok" {
		//btc - usdt
		//eth - usdt
		rate := ret.Get("ticker").Array()[0].Get("last").Float()
		rateCoin.Store(symbol, rate)
		return rate
	}
	return 1
}

func (c Coinbene) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/ticker?symbol=all"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Str
	if status == "ok" {
		ret.Get("ticker").ForEach(func(key, value gjson.Result) bool {
			var quote string
			var base string
			symbol := value.Get("symbol").Str
			if strings.HasSuffix(symbol, "USDT") {
				quote = symbol[:len(symbol)-4]
				base = "USDT"
			} else { //BTC - ETH
				quote = symbol[:len(symbol)-3]
				base = symbol[len(symbol)-3:]
			}
			exchangeTicker := &data.ExchangeTicker{
				Symbol: strings.ToUpper(symbol),
				Quote:  strings.ToUpper(quote),
				Base:   strings.ToUpper(base),
				Volume: value.Get("24hrVol").Float(),
				Amount: value.Get("24hrAmt").Float(),
				Last:   value.Get("last").Float(),
				Time:   strconv.FormatInt(ret.Get("timestamp").Int(), 10),
			}
			//汇率
			if strings.ToUpper(base) == "USDT" {
				exchangeTicker.LastUsd = exchangeTicker.Last
				exchangeTicker.AmountUsd = exchangeTicker.Amount
			} else {
				rate := c.GetRate(base, "USDT")
				exchangeTicker.LastUsd = exchangeTicker.Last * rate
				exchangeTicker.AmountUsd = exchangeTicker.Amount * rate
			}
			exchangeTickers = append(exchangeTickers, exchangeTicker)
			return true
		})
	} else {
		common.MessageHandler(c.Name(), ret.Get("description").Str)
	}
	return exchangeTickers
}

func (c Coinbene) AmountHandler() []*data.TradeData {
	var tradeDatas []*data.TradeData
	url := ApiHost + "/ticker?symbol=all"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status").Str
	if status == "ok" {
		ret.Get("ticker").ForEach(func(key, value gjson.Result) bool {
			var quote string
			var base string
			symbol := value.Get("symbol").Str
			if strings.HasSuffix(symbol, "USDT") {
				quote = symbol[:len(symbol)-4]
				base = "USDT"
			} else { //BTC - ETH
				quote = symbol[:len(symbol)-3]
				base = symbol[len(symbol)-3:]
			}
			tradeData := &data.TradeData{
				Symbol: strings.ToUpper(symbol),
				Quote:  strings.ToUpper(quote),
				Base:   strings.ToUpper(base),
				Volume: value.Get("24hrVol").Float(),
				Amount: value.Get("24hrAmt").Float(),
				Price:  value.Get("last").Float(),
			}
			if strings.ToUpper(base) == "USDT" {
				tradeData.PriceUsd = tradeData.Price
				tradeData.AmountUsd = tradeData.Amount
			} else {
				rate := c.GetRate(base, "USDT")
				tradeData.PriceUsd = tradeData.Price * rate
				tradeData.AmountUsd = tradeData.Amount * rate
			}
			tradeDatas = append(tradeDatas, tradeData)
			return true
		})
	} else {
		common.MessageHandler(c.Name(), ret.Get("description").Str)
	}
	return tradeDatas
}