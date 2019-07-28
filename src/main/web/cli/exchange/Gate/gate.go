package Gate

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
文档: https://www.gateio.co/api2
*/

const (
	ApiHost = "https://data.gateio.co/api2/1"
)

var (
	rateCoin sync.Map
)

type Gate struct {
}

func (c Gate) Name() string {
	return "gate"
}

//虚拟货币 -- 这里全部都是以usd计价
func (c Gate) GetRate(quote, base string) float64 {
	symbol := strings.ToLower(quote + "_" + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/ticker/" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	rate := ret.Get("last").Float()
	defer rateCoin.Store(symbol, rate)
	return rate
}

func (c Gate) PairHandler() []*data.ExchangeTicker {
	cnyUsdRate := common.CalRate("cny")
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/tickers"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)

	symbols := ret.Map()
	for symbol, value := range symbols {
		quote := strings.Split(symbol, "_")[0]
		base := strings.Split(symbol, "_")[1]
		timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)
		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			AmountQuote:        value.Get("quoteVolume").Float(), //交易量
			AmountBase:         value.Get("baseVolume").Float(),
			Last:               value.Get("last").Float(),
			Time:               timeStr,
			PriceChangePercent: value.Get("percentChange").Float(),
		}
		//汇率
		if strings.ToUpper(base) == "USDT" {
			exchangeTicker.LastUsd = exchangeTicker.Last
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase
		} else if strings.ToUpper(base) == "CNYX" {
			usdtCnyxRate := c.GetRate("USDT", "CNYX")
			exchangeTicker.LastUsd = exchangeTicker.Last / usdtCnyxRate
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase / usdtCnyxRate
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
