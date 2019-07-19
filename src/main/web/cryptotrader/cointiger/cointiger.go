package cointiger

import (
	"errors"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"truxing/commons/log"
)

/**
文档
*/

const (
	TickerUrl = "https://www.cointiger.one/exchange/api/public/market/detail"
	TradeUrl  = "https://www.cointiger.one/exchange/trading/api/market/history/trade"
	DepthUrl  = "https://www.cointiger.one/exchange/trading/api/market/depth"
	MarketUrl = "https://www.cointiger.one/exchange/trading/api/v2/currencys/v2"
	KLineURL  = "https://www.cointiger.one/exchange/trading/api/market/history/kline"
)

//定义taker的成交方向
var sideMap = map[string]string{
	"1": "buy",
	"2": "sell",
}

// BitBox client
type CoinTiger struct {
	ApiKey string
	Secret string
}

func New(apiKey string, secret string) *CoinTiger {
	return &CoinTiger{
		ApiKey: apiKey,
		Secret: secret,
	}
}

func getSymbol(base, quote string) string {
	return strings.ToLower(strings.TrimSpace(quote)) + strings.ToLower(strings.TrimSpace(base))
}

/**
获取最新价
https://github.com/cointiger/api-docs/wiki/%E5%89%8D24%E5%B0%8F%E6%97%B6%E8%A1%8C%E6%83%85-(%E9%80%82%E7%94%A8%E4%BA%8E%E8%A1%8C%E6%83%85%E5%B1%95%E7%A4%BA%E5%B9%B3%E5%8F%B0%E4%BD%BF%E7%94%A8)
*/
func (ct *CoinTiger) GetTicker(base, quote string) (*model.Ticker, error) {
	url := TickerUrl
	resp, err := http.Get(url)
	log.Debugf("%s", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ticker *model.Ticker
	result := gjson.ParseBytes(body)
	symbol := strings.ToUpper(base) + strings.ToUpper(quote)
	if result.Get(symbol).Exists() {
		t := result.Get(symbol)
		ticker = &model.Ticker{
			Buy:                t.Get("highestBid").Float(),
			Sell:               t.Get("lowestAsk").Float(),
			Last:               t.Get("last").Float(),
			Low:                t.Get("low24hr").Float(),
			High:               t.Get("high24hr").Float(),
			Vol:                t.Get("baseVolume").Float(),
			Exchange:           t.Get("quoteVolume").Float(),
			Time:               cast.ToTime(time.Now()),
			PriceChangePercent: t.Get("percentChange").Float(),
			Raw:                t.Str,
		}
	} else {
		return nil, errors.New(result.Get("msg").String())
	}

	return ticker, nil
}

/**
获取成交记录
https://github.com/cointiger/api-docs/wiki/REST-%E6%88%90%E4%BA%A4%E5%8E%86%E5%8F%B2%E6%95%B0%E6%8D%AE
*/
func (ct *CoinTiger) GetTades(base, quote string, size int) ([]model.Trade, error) {
	url := TradeUrl + "?symbol=" + getSymbol(base, quote)
	if size != 0 {
		url += "&size=" + strconv.Itoa(size)
	}
	log.Debugf("url:%s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var trades []model.Trade

	result := gjson.ParseBytes(body)

	log.Debugf("%s", result)

	if result.Get("code").Str == "0" {
		result.Get("data.trade_data").ForEach(func(key, value gjson.Result) bool {
			trade := model.Trade{
				ID:        value.Get("id").Int(),
				Amount:    value.Get("amount").Float(),
				Price:     value.Get("price").Float(),
				Type:      value.Get("side").String(),
				TradeTime: value.Get("ts").Int() / 1000,
				Time:      cast.ToTime(value.Get("ts").Int() / 1000),
				Raw:       value.String(),
			}
			trades = append(trades, trade)
			return true
		})
	} else {
		return nil, errors.New(result.Get("msg").String())
	}

	return trades, nil
}

/**
查询买卖盘深度
https://github.com/cointiger/api-docs/wiki/REST-%E6%B7%B1%E5%BA%A6%E7%9B%98%E5%8F%A3
*/
func (ct *CoinTiger) GetMarketDepth(base, quote string, typ string) (model.OrderBook, error) {
	orderBook := model.OrderBook{}
	url := DepthUrl + "?symbol=" + getSymbol(base, quote)
	if typ == "" {
		typ = "step0"
	}
	url += "&type=" + typ
	log.Debugf("url:%s", url)
	resp, err := http.Get(url)
	if err != nil {
		return orderBook, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return orderBook, err
	}

	log.Debugf("%s", string(body))

	var order model.MarketOrder
	result := gjson.ParseBytes(body)
	if result.Get("code").Str == "0" {
		gjson.GetBytes(body, "data.depth_data.buys").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Array()[0].Float()
			order.Amount = value.Array()[1].Float()
			orderBook.Bids = append(orderBook.Bids, order)
			return true // keep iterating
		})

		gjson.GetBytes(body, "data.depth_data.asks").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Array()[0].Float()
			order.Amount = value.Array()[1].Float()
			orderBook.Asks = append(orderBook.Asks, order)
			return true // keep iterating
		})

		orderBook.OrderTime = time.Now().Unix()
		orderBook.Time = time.Now()
		orderBook.Raw = result.Get("data").Raw
	} else {
		return orderBook, errors.New(result.Get("msg").String())
	}
	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
查询系统支持的所有交易对及精度
https://github.com/CoinTiger/api-docs/wiki/GET-%7BTrading_Macro_v2%7D-currencys-%E6%9F%A5%E8%AF%A2cointiger%E7%AB%99%E6%94%AF%E6%8C%81%E7%9A%84%E6%89%80%E6%9C%89%E5%B8%81%E7%A7%8D-(V2%E7%89%88%E6%9C%AC)
*/
func (ct *CoinTiger) GetMarkets() ([]model.MarketPairInfo, error) {
	var tradePairs []model.MarketPairInfo
	url := MarketUrl
	log.Debugf("url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("Response body: %v", string(body))
	result := gjson.ParseBytes(body)
	if result.Get("code").Str == "0" {
		result.Get("data").ForEach(func(key, value gjson.Result) bool {
			value.ForEach(func(key, value gjson.Result) bool {
				tradePair := model.MarketPairInfo{
					Base:  value.Get("baseCurrency").Str,
					Quote: value.Get("quoteCurrency").Str,
				}
				tradePairs = append(tradePairs, tradePair)
				return true
			})
			return true
		})
	} else {
		return tradePairs, errors.New(result.Get("msg").String())
	}
	return tradePairs, nil
}

/**
K线
https://github.com/cointiger/api-docs/wiki/REST-K%E7%BA%BF%E5%8E%86%E5%8F%B2%E6%95%B0%E6%8D%AE
*/
func (ct *CoinTiger) GetRecords(base, quote, period string, size int) ([]model.Record, error) {
	url := KLineURL + "?symbol=" + getSymbol(base, quote) + "&period=" + period
	if size != 0 {
		url += "&size=" + strconv.Itoa(size)
	}
	log.Debugf("Request url:%v", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("Response body: %v", string(body))

	var record model.Record
	var records []model.Record
	//2019-07-03T04:00:00.000Z
	//timeLayout := "2006-01-02T15:04:05.000Z" //转化所需模板
	//loc, _ := time.LoadLocation("Local")     //重要：获取时区
	result := gjson.ParseBytes(body)
	if result.Get("code").Str == "0" {
		result.Get("data.kline_data").ForEach(func(key, value gjson.Result) bool {
			record.Open = value.Get("open").Float()
			record.High = value.Get("high").Float()
			record.Low = value.Get("low").Float()
			record.Close = value.Get("close").Float()
			record.Vol = value.Get("vol").Float()
			record.Ktime = time.Now().Unix()
			records = append(records, record)
			return true // keep iterating
		})
	} else {
		return records, errors.New(result.Get("msg").String())
	}
	return records, nil
}
