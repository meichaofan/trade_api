package Bitfinex

import (
	"github.com/tidwall/gjson"
	"strings"
	"sync"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://docs.bitfinex.com/docs
*/

const (
	ApiHost = "https://api.bitfinex.com/v1"
)

var (
	rateCoin sync.Map
	once     sync.Once
)

type Bitfinex struct {
}

func (b Bitfinex) Name() string {
	return "bitfinex"
}

/**
获取该交易所支持的交易对
https://docs.bitfinex.com/reference#rest-public-symbols
*/
func GetSymbol() []*data.MarketPair {
	var pairs []*data.MarketPair
	url := ApiHost + "/symbols"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	symbols := ret.Array()
	for _, s := range symbols {
		symbol := s.Str
		quote := symbol[:len(symbol)-3]
		base := symbol[len(symbol)-3:]
		pair := &data.MarketPair{
			Symbol: strings.ToUpper(symbol),
			Quote:  strings.ToUpper(quote),
			Base:   strings.ToUpper(base),
		}
		pairs = append(pairs, pair)
	}
	return pairs[:len(pairs)-2]
}

func initCoinRate() {
	var pairs []*data.MarketPair
	pairs = GetSymbol()
	for _, pair := range pairs {
		if pair.Symbol == "BTCUSD" ||
			pair.Symbol == "DAIUSD" ||
			pair.Symbol == "EOSUSD" ||
			pair.Symbol == "ETHUSD" ||
			pair.Symbol == "USTUSD" ||
			pair.Symbol == "XLMUSD" {
			url := ApiHost + "/pubticker/" + pair.Symbol
			content := common.HttpGet(url)
			ret := gjson.ParseBytes(content)
			rate := ret.Get("last_price").Float()
			rateCoin.Store(pair.Symbol, rate)
		}
	}
}

func (b Bitfinex) GetRate(quote, base string) float64 {
	symbol := strings.ToUpper(quote + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	return 0
}

func getCountCoinRate(rate chan<- map[string]float64) {
	countRate := make(map[string]float64)
	countCoin := []string{
		"EUR",
		"GBP",
		"JPY",
		"CNY",
	}
	for _, v := range countCoin {
		countRate[v] = common.CalRate(v)
	}
	rate <- countRate
}

/**
此接口请注意，Ratelimit: 20 req/min
因此放弃goroutine
*/
func (b Bitfinex) PairHandler() []*data.ExchangeTicker {
	//初始化汇率
	once.Do(initCoinRate)
	//获取 cny、eur、jpy、gbp - usd 的汇率
	r := make(chan map[string]float64, 1)
	getCountCoinRate(r)
	countRate := <-r
	var exchangeTickers []*data.ExchangeTicker
	var pairs []*data.MarketPair
	pairs = GetSymbol()

	for _, pair := range pairs {
		url := ApiHost + "/pubticker/" + pair.Symbol
		content := common.HttpGet(url)
		ret := gjson.ParseBytes(content)
		amountQuote := ret.Get("volume").Float()
		highPriceRate := ret.Get("high").Float()
		lowPriceRate := ret.Get("low").Float()
		amountBase := amountQuote * (highPriceRate + lowPriceRate) / 2
		exchangeTicker := &data.ExchangeTicker{
			Symbol:      strings.ToUpper(pair.Symbol),
			Quote:       pair.Quote,
			Base:        pair.Base,
			AmountQuote: amountQuote,
			AmountBase:  amountBase,
			Last:        ret.Get("last_price").Float(),
			Time:        ret.Get("timestamp").Str,
		}
		//汇率
		if strings.ToUpper(pair.Base) == "USD" {
			exchangeTicker.LastUsd = exchangeTicker.Last
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase
		} else if strings.ToUpper(pair.Base) == "EUR" ||
			strings.ToUpper(pair.Base) == "GBP" ||
			strings.ToUpper(pair.Base) == "JPY" {
			if rate, ok := countRate[pair.Base]; ok {
				exchangeTicker.LastUsd = exchangeTicker.Last / rate
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase / rate
			}
		} else {
			rate := b.GetRate(pair.Base, "usd")
			exchangeTicker.LastUsd = exchangeTicker.Last * rate
			exchangeTicker.AmountUsd = exchangeTicker.AmountBase * rate
		}
		cnyRate := countRate["CNY"]
		exchangeTicker.AmountCny = exchangeTicker.AmountUsd * cnyRate
		exchangeTicker.LastCny = exchangeTicker.LastUsd * cnyRate
		exchangeTickers = append(exchangeTickers, exchangeTicker)
		time.Sleep(4 * time.Second)
	}

	return exchangeTickers
}
