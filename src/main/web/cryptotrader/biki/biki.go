package biki

import (
	"errors"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"truxing/commons/log"
)

/**
文档
https://bikicoin.oss-cn-hangzhou.aliyuncs.com/web_doc/openapi.pdf
*/

const (
	RestHost = "https://api.biki.com/"
)

// Biki client
type Biki struct {
	ApiKey string
	Sign   string
}

func New(apiKey string, sign string) *Biki {
	return &Biki{
		ApiKey: apiKey,
		Sign:   sign,
	}
}

func getSymbol(base, quote string) string {
	return strings.ToLower(strings.TrimSpace(base)) + strings.ToLower(strings.TrimSpace(quote))
}

/**
获取最新价
*/
func (biki *Biki) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "open/api/get_ticker?symbol=" + getSymbol(base, quote)
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
	code := result.Get("code").Int()
	if code == 0 {
		buy := result.Get("data.buy").Float()
		sell := result.Get("data.sell").Float()
		last := result.Get("data.last").Float()
		low := result.Get("data.low").Float()
		high := result.Get("data.high").Float()
		vol := result.Get("data.vol").Float()
		time := cast.ToTime(result.Get("data.time").Int() / 1000)
		raw := result.Get("data").String()
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
		return nil, errors.New(result.Get("msg").String())
	}
	return ticker, nil
}

/**
获取成交记录
*/
func (biki *Biki) GetTades(base, quote string) ([]model.Trade, error) {
	url := RestHost + "open/api/get_trades?symbol=" + getSymbol(base, quote)
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

	code := result.Get("code").Int()

	if code == 0 {
		result.Get("data").ForEach(func(key, value gjson.Result) bool {
			trade := model.Trade{
				Amount: value.Get("amount").Float(),
				Price:  value.Get("price").Float(),
				//此处ID溢出了
				ID:        value.Get("id").Int(),
				Type:      value.Get("type").String(),
				Raw:       value.String(),
				TradeTime: value.Get("ctime").Int() / 1000,
				Time:      cast.ToTime(value.Get("ctime").Int() / 1000),
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
查询买卖盘深度
@type 深度类型，step0，step1，step2（合并深度0-2）;step0时，精度最高
*/
func (biki *Biki) GetMarketDepth(base, quote string, typ string) (model.OrderBook, error) {
	orderBook := model.OrderBook{}
	url := RestHost + "open/api/market_dept?symbol=" + getSymbol(base, quote)
	if typ == "" {
		typ = "step0"
	}
	url += "&type=" + typ
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

	if gjson.GetBytes(body, "code").Int() == 0 {

		gjson.GetBytes(body, "data.tick.bids").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Array()[1].Float()
			order.Amount = value.Array()[0].Float()
			orderBook.Bids = append(orderBook.Bids, order)
			return true // keep iterating
		})

		gjson.GetBytes(body, "data.tick.asks").ForEach(func(key, value gjson.Result) bool {
			order.Price = value.Array()[1].Float()
			order.Amount = value.Array()[0].Float()
			orderBook.Asks = append(orderBook.Asks, order)
			return true // keep iterating
		})

		orderBook.OrderTime = time.Now().Unix()
		local, _ := time.LoadLocation("Local")
		orderBook.Time = time.Now().In(local)
		orderBook.Raw = gjson.ParseBytes(body).Get("data").Raw
	} else {
		return orderBook, errors.New(gjson.GetBytes(body, "msg").Str)
	}
	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
查询系统支持的所有交易对及精度
*/
func (biki *Biki) GetMarkets() ([]model.MarketPairInfo, error) {
	var tradePairs []model.MarketPairInfo
	url := RestHost + "open/api/common/symbols"
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
	if gjson.GetBytes(body, "code").Int() == 0 {
		gjson.ParseBytes(body).Get("data").ForEach(func(key, value gjson.Result) bool {
			tradePair := model.MarketPairInfo{
				Base:  value.Get("base_coin").String(),
				Quote: value.Get("count_coin").String(),
			}
			tradePairs = append(tradePairs, tradePair)
			return true
		})
	} else {
		return nil, errors.New(gjson.GetBytes(body, "msg").Str)
	}
	return tradePairs, nil
}
