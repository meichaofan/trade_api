package Bibox

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
文档:
*/
const (
	ApiHost = "https://api.bibox365.com/v1/mdata"
)

var (
	//定义taker的成交方向
	sideMap = map[string]string{
		"1": "buy",
		"2": "sell",
	}
	//虚拟货币兑换美元的汇率
	coinRate sync.Map
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
		percent := value.Get("percent").Str
		priceChangePercent, _ := strconv.ParseFloat(strings.TrimRight(percent, "%"), 64)
		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(value.Get("coin_symbol").Str + "_" + value.Get("currency_symbol").Str),
			Quote:              value.Get("coin_symbol").Str,
			Base:               value.Get("currency_symbol").Str,
			Volume:             value.Get("vol24H").Float(),
			Amount:             value.Get("amount").Float(),
			Last:               value.Get("last").Float(),
			LastUSD:            value.Get("last_usd").Float(),
			PriceChangePercent: priceChangePercent,
			Time:               strconv.FormatInt(time.Now().Unix(), 10),
		}
		exchangeTickers = append(exchangeTickers, exchangeTicker)
		return true
	})
	return exchangeTickers
}

/**
获取平台交易额
*/
func (bibox Bibox) AmountHandler() []*data.TradeData {
	//首先获取所有的交易对
	var tradeDatas []*data.TradeData
	url := ApiHost + "?cmd=marketAll"
	content := common.HttpGet(url)
	gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
		quote := value.Get("coin_symbol").Str
		base := value.Get("currency_symbol").Str
		symbol := strings.ToUpper(quote + "_" + base)
		last := value.Get("last").Float()
		lastUsd := value.Get("last_usd").Float()
		rate := lastUsd / last
		tradeData := &data.TradeData{
			Symbol:    symbol,
			Quote:     quote,
			Base:      base,
			Volume:    value.Get("vol24H").Float(),
			Amount:    value.Get("amount").Float(),
			AmountUsd: value.Get("amount").Float() * rate,
		}
		tradeDatas = append(tradeDatas, tradeData)
		return true
	})
	return tradeDatas
}
