package bitz

import (
	"errors"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"truxing/commons/log"
)

const (
	RestHost = "https://api.bitzapi.com/"
)

// CoinBene client
type BitZ struct {
	AppId  string
	Secret string
}

func New(appId string, secret string) *BitZ {
	return &BitZ{
		AppId:  appId,
		Secret: secret,
	}
}

func getSymbol(base, quote string) string {
	return strings.ToLower(strings.TrimSpace(base)) + "_" + strings.ToLower(strings.TrimSpace(quote))
}

/**
获取最新价
https://apidoc.bitz.com/cn/market-quotation-data/Get-ticker-data.html
*/
func (bz *BitZ) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "Market/ticker?symbol=" + getSymbol(base, quote)
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
	status := result.Get("status").Int()
	if status == 200 {
		buy := result.Get("data.bidPrice").Float()
		sell := result.Get("data.askPrice").Float()
		last := result.Get("data.now").Float()
		low := result.Get("data.low").Float()
		high := result.Get("data.high").Float()
		vol := result.Get("data.quoteVolume").Float()
		time := cast.ToTime(result.Get("time").Int())
		priceChangePercent := result.Get("data.priceChange24h").Float()

		raw := string(body)
		//pricechangepercent
		ticker = &model.Ticker{
			Buy:                buy,
			Sell:               sell,
			Last:               last,
			Low:                low,
			High:               high,
			Vol:                vol,
			Time:               time,
			Raw:                raw,
			PriceChangePercent: priceChangePercent,
		}
	} else {
		return nil, errors.New(result.Get("msg").String())
	}
	return ticker, nil
}

/**
获取成交记录
https://apidoc.bitz.com/cn/market-quotation-data/Get-orders-data.html
*/
func (bz *BitZ) GetTades(base, quote string) ([]model.Trade, error) {
	url := RestHost + "Market/order?symbol=" + getSymbol(base, quote)
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

	status := result.Get("status").Int()

	if status == 200 {
		result.Get("data").ForEach(func(key, value gjson.Result) bool {
			trade := model.Trade{
				Amount: value.Get("n").Float(),
				Price:  value.Get("p").Float(),
				ID:   value.Get("id").Int(),
				Type: value.Get("s").String(),
				Time: cast.ToTime(value.Get("T").Int()),
				Raw:  value.String(),
			}
			trades = append(trades, trade)
			return true
		})
	} else {
		return nil, errors.New(result.Get("msg").Str)
	}
	return trades, nil
}

/**
获取挂单 (委托单)
https://apidoc.bitz.com/cn/market-quotation-data/Get-depth-data.html
*/
func (bz *BitZ) GetOrderBook(base, quote string) (model.OrderBook, error) {
	orderBook := model.OrderBook{}

	url := RestHost + "Market/depth?symbol=" + getSymbol(base, quote)

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

	if gjson.GetBytes(body, "status").Int() == 200 {

		gjson.GetBytes(body, "data.bids").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Array()[0].Float()
			order.Amount = value.Array()[1].Float()
			orderBook.Bids = append(orderBook.Bids, order)
			return true // keep iterating
		})

		gjson.GetBytes(body, "data.asks").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Array()[0].Float()
			order.Amount = value.Array()[1].Float()
			orderBook.Asks = append(orderBook.Asks, order)
			return true // keep iterating
		})
		orderTime := gjson.GetBytes(body, "time").Int()
		orderBook.OrderTime = orderTime
		orderBook.Time = cast.ToTime(orderTime)
		orderBook.Raw = gjson.ParseBytes(body).Raw
	} else {
		return orderBook, errors.New(gjson.GetBytes(body, "msg").Str)
	}
	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
交易对信息 【行情】
https://apidoc.bitz.com/cn/market-quotation-data/Get-symbolList.html
*/
func (bz *BitZ) GetMarkets() ([]model.MarketPairInfo, error) {
	var tradePairs []model.MarketPairInfo
	url := RestHost + "Market/symbolList"
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
	if gjson.GetBytes(body, "status").Int() == 200 {
		symbols := gjson.ParseBytes(body).Get("data").Map()
		for _, v := range symbols {
			tradePair := model.MarketPairInfo{
				Base:  gjson.Get(v.String(), "coinFrom").String(),
				Quote: gjson.Get(v.String(), "coinTo").String(),
			}
			tradePairs = append(tradePairs, tradePair)
		}
	} else {
		return nil, errors.New(gjson.GetBytes(body, "msg").Str)
	}
	return tradePairs, nil
}
