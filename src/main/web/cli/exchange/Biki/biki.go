package Biki

import (
	"crypto/md5"
	"fmt"
	"github.com/tidwall/gjson"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档 ：https://bikicoin.oss-cn-hangzhou.aliyuncs.com/web_doc/openapi.pdf
*/

const (
	ApiHost   = "https://api.biki.com"
	ApiKey    = "e001cacf1d6e0ce79521ca024b70f7f6"
	SecretKey = "b1bf7a9c25f409fa52d5f8d01414e173"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Biki struct {
}

func (biki Biki) Name() string {
	return "biki"
}

/**
计算签名
*/
func Sign(params map[string]string) string {
	var keys []string
	var mystring string
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		mystring += k + params[k]
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(mystring+SecretKey)))
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
func initCoinRate(done chan<- struct{}) {
	timeStr := strconv.Itoa(int(time.Now().Unix()))
	params := map[string]string{
		"api_key": ApiKey,
		"time":    timeStr,
	}
	sign := Sign(params)
	url := ApiHost + "/open/api/market?api_key=" + ApiKey + "&time=" + timeStr + "&sign=" + sign
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	code := ret.Get("code").Int()
	if code == 0 {
		symbolsRate := ret.Get("data").Map()
		for k, v := range symbolsRate {
			rateCoin.Store(strings.ToUpper(k), v.Float())
		}
	} else {
		common.MessageHandler("biki getRate", ret.Get("msg").Str)
	}
	done <- struct {
	}{}
}

func (biki Biki) GetRate(quote, base string) float64 {
	symbol := quote + base
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	return 0
}

func (biki Biki) PairHandler() []*data.ExchangeTicker {
	//初始化汇率
	done := make(chan struct{}, 1)
	initCoinRate(done)
	<-done
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

			highPrice := ret.Get("data.high").Float()
			lowPrice := ret.Get("data.low").Float()
			vol := ret.Get("data.vol").Float()
			amount := (highPrice + lowPrice) * vol
			exchangeTicker := &data.ExchangeTicker{
				Symbol:             strings.ToUpper(pair.Symbol),
				Quote:              pair.Quote,
				Base:               pair.Base,
				Volume:             vol,
				Amount:             amount,
				Last:               ret.Get("data.last").Float(),
				PriceChangePercent: ret.Get("data.rose").Float(),
				Time:               strconv.FormatInt(ret.Get("data.time").Int(), 10),
				//汇率
			}
			if strings.ToUpper(pair.Base) == "USDT" {
				exchangeTicker.LastUsd = exchangeTicker.Last
				exchangeTicker.AmountUsd = exchangeTicker.Amount
			} else {
				rate := biki.GetRate(pair.Base, "USDT")
				//fmt.Printf("rate: %f",rate)
				//fmt.Println()
				exchangeTicker.LastUsd = exchangeTicker.Last * rate
				exchangeTicker.AmountUsd = exchangeTicker.Amount * rate
			}
			exchangeTickers = append(exchangeTickers, exchangeTicker)
		}(pair)
	}
	wg.Wait()
	return exchangeTickers
}

func (biki Biki) AmountHandler() []*data.TradeData {
	//初始化汇率
	done := make(chan struct{}, 1)
	initCoinRate(done)
	<-done
	var tradeDatas []*data.TradeData
	var pairs []*data.MarketPair
	pairs = GetSymbol()
	var wg sync.WaitGroup
	wg.Add(len(pairs))
	for _, pair := range pairs {
		go func(pair *data.MarketPair) {
			defer wg.Done()
			url := ApiHost + "/open/api/get_trades?symbol=" + pair.Symbol
			content := common.HttpGet(url)
			ret := gjson.ParseBytes(content)
			code := ret.Get("code").Int()
			if code == 0 {
				ret.Get("data").ForEach(func(key, value gjson.Result) bool {
					volume := value.Get("amount").Float()
					price := value.Get("price").Float()
					amount := volume * price
					tradeData := &data.TradeData{
						ID:        value.Get("id").String(),
						Symbol:    strings.ToUpper(pair.Symbol),
						Quote:     pair.Quote,
						Base:      pair.Base,
						Type:      value.Get("type").String(),
						Volume:    volume,
						Price:     price,
						Amount:    amount,
						TradeTime: strconv.FormatInt(value.Get("ctime").Int(), 10),
						TradeTs:   value.Get("ctime").Int() / 1000,
					}
					if pair.Base == "USDT" {
						tradeData.PriceUsd = tradeData.Price
						tradeData.AmountUsd = tradeData.Amount
					} else {
						rate := biki.GetRate(pair.Base, "USDT")
						tradeData.PriceUsd = tradeData.Price * rate
						tradeData.AmountUsd = tradeData.Amount * rate
					}
					tradeDatas = append(tradeDatas, tradeData)
					return true
				})
			}
		}(pair)
	}
	wg.Wait()
	return tradeDatas
}
