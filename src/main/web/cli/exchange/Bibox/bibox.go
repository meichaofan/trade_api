package Bibox

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

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
	coinRate = make(map[string]float64)
)

type Bibox struct {
}

func (bibox Bibox) Name() string {
	return "bibox"
}

func (bibox Bibox) GetRate(quote, base string) float64 {
	var exchangeRate float64
	if rate, ok := coinRate[base]; ok == true {
		exchangeRate = rate
	} else if strings.ToUpper(base) == "USDT" {
		url := ApiHost + "?cmd=ticker&pair=" + strings.ToUpper(quote) + "_" + strings.ToUpper(base)
		content := common.HttpGet(url)
		last := gjson.ParseBytes(content).Get("result.last").Float()
		coinRate[quote] = last
		exchangeRate = 1
	} else {
		nUrl := ApiHost + "?cmd=ticker&pair=" + strings.ToUpper(base) + "_USDT"
		nContent := common.HttpGet(nUrl)
		nLast := gjson.ParseBytes(nContent).Get("result.last").Float()
		coinRate[base] = nLast
		exchangeRate = nLast
	}
	return exchangeRate
}

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
			Vol:                value.Get("vol24H").Float(),
			Amount:             value.Get("amount").Float(),
			Last:               value.Get("last").Float(),
			LastUSD:            value.Get("last_usd").Float(),
			PriceChangePercent: priceChangePercent,
			Time:               time.Now().Unix(),
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
	pairs := bibox.GetAllPair()
	for _, pair := range pairs {
		url := ApiHost + "?cmd=deals&pair=" + pair.Quote + "_" + pair.Base
		content := common.HttpGet(url)
		gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
			trade := &data.TradeData{
				Symbol:    pair.Quote + "_" + pair.Base,
				PairQuote: pair.Quote,
				PairBase:  pair.Base,
				Amount:    value.Get("amount").Float(),
				Price:     value.Get("price").Float(),
				Type:      sideMap[value.Get("side").String()],
				TradeTime: value.Get("time").Int() / 1000,
			}
			if strings.ToUpper(pair.Base) == "USDT" {
				trade.PriceUsd = trade.Price
				trade.AmountUsd = trade.Amount
			} else {
				rate := bibox.GetRate(pair.Quote, pair.Base)
				trade.PriceUsd = rate * trade.Price
				trade.AmountUsd = rate * trade.Amount
			}
			tradeDatas = append(tradeDatas, trade)
			return true
		})
	}
	return tradeDatas
}

/**
获取平台交易对，及其价格
*/
func (bibox Bibox) GetAllPair() []*data.MarketPair {
	var pairs []*data.MarketPair
	url := ApiHost + "?cmd=pairList"
	content := common.HttpGet(url)
	gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
		quote, base := strings.Split(value.Get("pair").String(), "_")[0], strings.Split(value.Get("pair").String(), "_")[1]
		pair := &data.MarketPair{
			Base:  base,
			Quote: quote,
		}
		pairs = append(pairs, pair)
		return true
	})
	return pairs
}
