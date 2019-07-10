package bitfinex

import (
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
https://docs.bitfinex.com/docs
*/

const (
	RestHost = "https://api.bitfinex.com/v1"
)

//定义taker的成交方向
var sideMap = map[string]string{
	"1": "buy",
	"2": "sell",
}

// Bitfinex client
type Bitfinex struct {
	ApiKey string
	Secret string
}

func New(apiKey string, secret string) *Bitfinex {
	return &Bitfinex{
		ApiKey: apiKey,
		Secret: secret,
	}
}

func getSymbol(base, quote string) string {
	return strings.ToLower(strings.TrimSpace(base)) + strings.ToLower(strings.TrimSpace(quote))
}

/**
获取最新价
https://docs.bitfinex.com/reference#rest-public-ticker
*/
func (b *Bitfinex) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "/pubticker/" + getSymbol(base, quote)
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
	ticker = &model.Ticker{
		Buy:  result.Get("bid").Float(),
		Sell: result.Get("ask").Float(),
		Last: result.Get("last_price").Float(),
		Low:  result.Get("low").Float(),
		High: result.Get("high").Float(),
		Vol:  result.Get("volume").Float(),
		Time: cast.ToTime(strings.Split(result.Get("timestamp").String(), ".")[0]),
		Raw:  result.Raw,
	}
	return ticker, nil
}

/**
获取成交记录
https://docs.bitfinex.com/reference#rest-public-trades
*/
func (b *Bitfinex) GetTades(base, quote string, size int) ([]model.Trade, error) {
	url := RestHost + "/trades/" + getSymbol(base, quote)
	if size != 0 {
		url += "?limit_trades=" + strconv.Itoa(size)
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

	result.Get("result").ForEach(func(key, value gjson.Result) bool {
		trade := model.Trade{
			Amount:    value.Get("amount").Float(),
			Price:     value.Get("price").Float(),
			Type:      sideMap[value.Get("type").String()],
			ID:        value.Get("tid").Int(),
			TradeTime: value.Get("timestamp").Int(),
			Time:      cast.ToTime(value.Get("timestamp").Int()),
			Raw:       value.String(),
		}
		trades = append(trades, trade)
		return true
	})

	return trades, nil
}

/**
查询买卖盘深度
https://docs.bitfinex.com/reference#rest-public-orderbook
*/
func (b *Bitfinex) GetMarketDepth(base, quote string) (model.OrderBook, error) {
	orderBook := model.OrderBook{}
	url := RestHost + "/book/" + getSymbol(base, quote)
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

	gjson.GetBytes(body, "bids").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Get("price").Float()
		order.Amount = value.Get("amount").Float()
		orderBook.Bids = append(orderBook.Bids, order)
		return true // keep iterating
	})

	gjson.GetBytes(body, "asks").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Get("price").Float()
		order.Amount = value.Get("amount").Float()
		orderBook.Asks = append(orderBook.Asks, order)
		return true // keep iterating
	})

	orderBook.OrderTime = time.Now().Unix()
	orderBook.Time = time.Now()
	orderBook.Raw = gjson.ParseBytes(body).Raw

	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
查询系统支持的所有交易对及精度
https://docs.bitfinex.com/reference#rest-public-symbols
*/
func (b *Bitfinex) GetMarkets() (interface{}, error) {
	url := RestHost + "/symbols"
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

	r := gjson.ParseBytes(body).Value()
	return r, nil
}
