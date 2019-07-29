package Bibox

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://github.com/Biboxcom/API_Docs/wiki
*/
const (
	ApiHost = "https://api.bibox365.com/v1/mdata"
)

type Bibox struct {
}

func (bibox Bibox) Name() string {
	return "bibox"
}

/**
获取汇率
*/
func (bibox Bibox) GetRate(quote, base string) float64 {
	return 0
}

/**
获取平台交易对及其价格
*/
func (bibox Bibox) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "?cmd=marketAll"
	content := common.HttpGet(url)
	gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
		//涨幅
		percent := value.Get("percent").Str
		priceChangePercent, _ := strconv.ParseFloat(strings.TrimRight(percent, "%"), 64)
		//symbol
		quote := value.Get("coin_symbol").Str
		base := value.Get("currency_symbol").Str
		Symbol := quote + "_" + base
		//计价货币和美元之间汇率 rate
		last := value.Get("last").Float()
		lastUsd := value.Get("last_usd").Float()
		baseUsdRate := lastUsd / last
		lastCny := lastUsd * common.CnyUsdRate
		//amount 看文档 对比 mytoken App 可知
		amountQuote := value.Get("vol24H").Float()
		amountBase := value.Get("amount").Float()
		amountUsd := amountBase * baseUsdRate
		amountCny := amountUsd * common.CnyUsdRate
		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(Symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			AmountQuote:        amountQuote,
			AmountBase:         amountBase,
			AmountUsd:          amountUsd,
			AmountCny:          amountCny,
			Last:               last,
			LastUsd:            lastUsd,
			LastCny:            lastCny,
			PriceChangePercent: priceChangePercent,
			Time:               strconv.FormatInt(time.Now().Unix(), 10),
		}
		exchangeTickers = append(exchangeTickers, exchangeTicker)
		return true
	})
	return exchangeTickers
}
