package Zb

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
文档：https://web.zb.cn/i/developer/restApi
 */

const (
	ApiHost = "http://api.zb.cn"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Zb struct {
}

func (c Zb) Name() string {
	return "zb"
}

//虚拟货币 -- 美元 汇率
func (c Zb) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + "_" + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/data/v1/ticker?market=" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	if ret.Exists() {
		rate := ret.Get("ticker.last").Float()
		rateCoin.Store(symbol, rate)
		return rate
	}
	return 1
}

func (c Zb) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/data/v1/allTicker"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content).Map()

	for symbol, value := range ret {
		s := symbol
		var quote string
		var base string
		if strings.HasSuffix(s, "usdt") {
			quote = symbol[:len(symbol)-4]
			base = "usdt"
		} else if strings.HasSuffix(symbol, "qc") ||
			strings.HasSuffix(symbol, "zb") {
			quote = symbol[:len(symbol)-2]
			base = symbol[len(symbol)-2:]
		} else {
			//BTC - ETH -TRX
			quote = symbol[:len(symbol)-3]
			base = symbol[len(symbol)-3:]
		}

		//计算交易额
		vol := value.Get("vol").Float()
		highPrice := value.Get("high").Float()
		lowPrice := value.Get("low").Float()
		amount := (highPrice + lowPrice) * vol / 2

		timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)

		exchangeTicker := &data.ExchangeTicker{
			Symbol: strings.ToUpper(symbol),
			Quote:  strings.ToUpper(quote),
			Base:   strings.ToUpper(base),
			Volume: value.Get("vol").Float(),
			Amount: amount,
			Last:   value.Get("last").Float(),
			Time:   timeStr,
		}

		//汇率
		if strings.ToUpper(base) == "USDT" {
			exchangeTicker.LastUsd = exchangeTicker.Last
			exchangeTicker.AmountUsd = exchangeTicker.Amount
		} else if strings.ToUpper(base) == "QC" || strings.ToUpper(base) == "ZB" {
			rate := c.GetRate("USDT", base) //这里注意  usdtqc usdtzb
			exchangeTicker.LastUsd = exchangeTicker.Last / rate
			exchangeTicker.AmountUsd = exchangeTicker.Amount / rate
		} else {
			rate := c.GetRate(base, "USDT")
			exchangeTicker.LastUsd = exchangeTicker.Last * rate
			exchangeTicker.AmountUsd = exchangeTicker.Amount * rate
		}
		exchangeTickers = append(exchangeTickers, exchangeTicker)
	}
	return exchangeTickers
}

func (c Zb) AmountHandler() []*data.TradeData {
	var tradeDatas []*data.TradeData
	return tradeDatas
}
