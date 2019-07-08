package coinbene

import (
	"errors"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"truxing/commons/log"
)

const (
	RestHost = "http://api.coinbene.com/v1/market/"
)

// CoinBene client
type CoinBene struct {
	AppId  string
	Secret string
}

func New(appId string, secret string) *CoinBene {
	return &CoinBene{
		AppId:  appId,
		Secret: secret,
	}
}

func getSymbol(base, quote string) string {
	return strings.ToUpper(strings.TrimSpace(base)) + strings.ToUpper(strings.TrimSpace(quote))
}

/**
获取最新价
https://github.com/Coinbene/API-Documents-CHN/wiki/1.1.0-%E8%8E%B7%E5%8F%96%E6%9C%80%E6%96%B0%E4%BB%B7(Ticker)-%E3%80%90%E8%A1%8C%E6%83%85%E3%80%91
*/
func (cb *CoinBene) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "ticker?symbol=" + getSymbol(base, quote)
	resp, err := http.Get(url)
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
	status := result.Get("status").String()
	if status == "ok" {
		buy := result.Get("ticker.0.bid").Float()
		sell := result.Get("ticker.0.ask").Float()
		last := result.Get("ticker.0.last").Float()
		low := result.Get("ticker.0.24hrLow").Float()
		high := result.Get("ticker.0.24hrHigh").Float()
		vol := result.Get("ticker.0.24hrVol").Float()
		time := cast.ToTime(result.Get("timestamp").Int() / 1000)
		raw := string(body)
		//pricechangepercent
		ticker = &model.Ticker{
			Buy:  buy,
			Sell: sell,
			Last: last,
			Low:  low,
			High: high,
			Vol:  vol,
			Time: time,
			Raw:  raw,
		}
	} else {
		return nil, errors.New(result.Get("description").String())
	}
	return ticker, nil
}

/**
获取成交记录
https://github.com/Coinbene/API-Documents-CHN/wiki/1.1.2-%E8%8E%B7%E5%8F%96%E6%88%90%E4%BA%A4%E8%AE%B0%E5%BD%95%E3%80%90%E8%A1%8C%E6%83%85%E3%80%91
*/
func (cb *CoinBene) GetTades(base, quote string, size int) ([]model.Trade, error) {
	url := RestHost + "trades?symbol=" + getSymbol(base, quote)
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

	status := result.Get("status").String()

	if status == "ok" {
		result.Get("trades").ForEach(func(key, value gjson.Result) bool {
			trade := model.Trade{
				Amount: value.Get("quantity").Float(),
				Price:  value.Get("price").Float(),
				//此处ID溢出了
				ID:   value.Get("tradeId").Int(),
				Type: value.Get("take").String(),
				Time: cast.ToTime(value.Get("time").Int() / 1000),
				Raw:  value.String(),
			}
			trades = append(trades, trade)
			return true
		})
	} else {
		return nil, errors.New(result.Get("description").Str)
	}
	return trades, nil
}

/**
获取挂单 (委托单)
https://github.com/Coinbene/API-Documents-CHN/wiki/1.1.1-%E8%8E%B7%E5%8F%96%E6%8C%82%E5%8D%95%E3%80%90%E8%A1%8C%E6%83%85%E3%80%91
*/
func (cb *CoinBene) GetOrderBook(base, quote string, depth int) (model.OrderBook, error) {
	orderBook := model.OrderBook{}

	url := RestHost + "orderbook?symbol=" + getSymbol(base, quote)
	if depth != 0 {
		url += "&depth" + strconv.Itoa(depth)
	}
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

	if gjson.GetBytes(body, "status").Str == "ok" {

		gjson.GetBytes(body, "orderbook.bids").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Get("price").Float()
			order.Amount = value.Get("quantity").Float()
			orderBook.Bids = append(orderBook.Bids, order)
			return true // keep iterating
		})

		gjson.GetBytes(body, "orderbook.asks").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Get("price").Float()
			order.Amount = value.Get("quantity").Float()
			orderBook.Asks = append(orderBook.Asks, order)
			return true // keep iterating
		})
		orderTime := gjson.GetBytes(body, "timestamp").Int() / 1000
		orderBook.OrderTime = orderTime
		orderBook.Time = cast.ToTime(orderTime)
		orderBook.Raw = gjson.ParseBytes(body).Raw
	} else {
		return orderBook, errors.New(gjson.GetBytes(body, "description").Str)
	}
	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
交易对信息 【行情】
https://github.com/Coinbene/API-Documents-CHN/wiki/1.1.3-%E4%BA%A4%E6%98%93%E5%AF%B9%E4%BF%A1%E6%81%AF-%E3%80%90%E8%A1%8C%E6%83%85%E3%80%91
*/
func (cb *CoinBene) GetMarkets() ([]model.MarketPairInfo, error) {
	var tradePairs []model.MarketPairInfo
	url := RestHost + "symbol"
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
	if gjson.GetBytes(body, "status").Str == "ok" {
		gjson.ParseBytes(body).Get("symbol").ForEach(func(key, value gjson.Result) bool {
			tradePair := model.MarketPairInfo{
				Quote: value.Get("quoteAsset").String(),
				Base:  value.Get("baseAsset").String(),
			}
			tradePairs = append(tradePairs, tradePair)
			return true
		})
	} else {
		return nil, errors.New(gjson.GetBytes(body, "description").Str)
	}
	return tradePairs, nil
}
