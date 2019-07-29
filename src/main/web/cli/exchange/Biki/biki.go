package Biki

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档 ：https://bikicoin.oss-cn-hangzhou.aliyuncs.com/web_doc/openapi.pdf
*/

const (
	ApiHost = "https://api.biki.com"
	//ApiKey    = "e001cacf1d6e0ce79521ca024b70f7f6"
	//SecretKey = "b1bf7a9c25f409fa52d5f8d01414e173"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin   sync.Map
	cnyUsdRate float64
)

type Biki struct {
}

func (biki Biki) Name() string {
	return "biki"
}

/**
获取该交易所支持的交易对
https://api.biki.com/open/api/common/symbols
*/
func GetSymbol() []*data.MarketPair {
	var pairs []*data.MarketPair
	url := ApiHost + "/open/api/common/symbols"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	code := ret.Get("code").Int()
	if code == 0 {
		ret.Get("data").ForEach(func(key, value gjson.Result) bool {
			//symbol := value.Get("symbol").Str
			pair := &data.MarketPair{
				Symbol: value.Get("symbol").Str,
				Quote:  value.Get("base_coin").Str,
				Base:   value.Get("count_coin").Str,
			}
			pairs = append(pairs, pair)
			return true
		})
	} else {
		common.MessageHandler("biki", ret.Get("msg").Str)
	}
	return pairs
}

/**
计算汇率
BTC - USDT
ETH - USDT
TRX - USDT
*/
func initCnyUsdRate(rate chan<- float64) {
	rate <- common.CnyUsdRate
}

func (biki Biki) GetRate(quote, base string) float64 {
	symbol := quote + base
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/open/api/get_ticker?symbol=" + symbol
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	code := ret.Get("code").Int()
	if code == 0 {
		//btc - usdt
		//eth - usdt
		rate := ret.Get("data.last").Float()
		rateCoin.Store(symbol, rate)
		return rate
	}
	return 0
}

func (biki Biki) PairHandler() []*data.ExchangeTicker {
	//初始化cny-usd汇率
	rate := make(chan float64, 1)
	initCnyUsdRate(rate)
	cnyUsdRate := <-rate
	var exchangeTickers []*data.ExchangeTicker
	var pairs []*data.MarketPair
	pairs = GetSymbol()
	var wg sync.WaitGroup
	wg.Add(len(pairs))
	for _, pair := range pairs {
		go func(pair *data.MarketPair) {
			defer wg.Done()
			url := ApiHost + "/open/api/get_ticker?symbol=" + pair.Symbol
			content := common.HttpGet(url)
			ret := gjson.ParseBytes(content)
			amountQuote := ret.Get("data.vol").Float()
			highPriceRate := ret.Get("data.high").Float()
			lowPriceRate := ret.Get("data.low").Float()
			amountBase := amountQuote * (highPriceRate + lowPriceRate) / 2

			exchangeTicker := &data.ExchangeTicker{
				Symbol:             strings.ToUpper(pair.Symbol),
				Quote:              pair.Quote,
				Base:               pair.Base,
				AmountQuote:        amountQuote,
				AmountBase:         amountBase,
				Last:               ret.Get("data.last").Float(),
				PriceChangePercent: ret.Get("data.rose").Float(),
				Time:               strconv.FormatInt(ret.Get("data.time").Int(), 10),
			}

			//汇率
			if strings.ToUpper(pair.Base) == "USDT" {
				exchangeTicker.LastUsd = exchangeTicker.Last
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase
			} else {
				baseUsdRate := biki.GetRate(pair.Base, "USDT")
				exchangeTicker.LastUsd = exchangeTicker.Last * baseUsdRate
				exchangeTicker.AmountUsd = exchangeTicker.AmountBase * baseUsdRate
			}
			exchangeTicker.AmountCny = exchangeTicker.AmountUsd * cnyUsdRate
			exchangeTicker.LastCny = exchangeTicker.LastUsd * cnyUsdRate
			exchangeTickers = append(exchangeTickers, exchangeTicker)
		}(pair)
	}
	wg.Wait()
	return exchangeTickers
}
