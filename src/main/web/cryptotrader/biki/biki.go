package biki

import (
	"errors"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
	"truxing/commons/log"
)

/**
文档
https://bikicoin.oss-cn-hangzhou.aliyuncs.com/web_doc/openapi.pdf
*/

const (
	RestHost = "https://api.biki.com/"
	PlatForm = "biki"
	DbPair   = "platform_pair"
)

var (
	// 定义虚拟货币汇率映射表 key: 表示虚拟货币, value是美元
	exchangeRate sync.Map
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
	return strings.ToLower(strings.TrimSpace(quote)) + strings.ToLower(strings.TrimSpace(base))
}

/**
获取最新价
*/
func (biki *Biki) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "open/api/get_ticker?symbol=" + getSymbol(base, quote)
	log.Debugf("url: %s",url)
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
				Base:  value.Get("count_coin").String(),
				Quote: value.Get("base_coin").String(),
			}
			tradePairs = append(tradePairs, tradePair)
			return true
		})
	} else {
		return nil, errors.New(gjson.GetBytes(body, "msg").Str)
	}
	return tradePairs, nil
}

/**
获取k线
 */
func (bk *Biki) GetRecords(base, quote, period string) ([]model.Record, error) {
	url := RestHost + "open/api/get_records?symbol=" + getSymbol(base, quote) + "&period=" + period
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
	timeLayout := "2006-01-02T15:04:05.000Z" //转化所需模板
	loc, _ := time.LoadLocation("Local")     //重要：获取时区
	if gjson.GetBytes(body, "code").Int() == 0 {
		gjson.ParseBytes(body).Get("data").ForEach(func(key, value gjson.Result) bool {
			record.Open = value.Array()[1].Float()
			record.High = value.Array()[2].Float()
			record.Low = value.Array()[3].Float()
			record.Close = value.Array()[4].Float()
			record.Vol = value.Array()[5].Float()
			theTime, _ := time.ParseInLocation(timeLayout, value.Array()[0].Str, loc) //使用模板在对应时区转化为time.time类型
			sr := theTime.Unix()                                                      //转化为时间戳 类型是int64
			record.Ktime = sr
			records = append(records, record)
			return true // keep iterating
		})
	} else {
		return records, errors.New(gjson.GetBytes(body, "msg").Str)

	}

	return records, nil
}


/**
获取当前平台所有交易额(从数据库获取)
*/
func (bb *Biki) GetExchangeAmountFormDb() (*data.ExchangeAmount, error) {
	var tradeDatas []*data.ExchangeTicker
	var amountUsd float64 = 0
	var amountCny float64 = 0
	s, c := common.Connect(DbPair, "biki", "local")
	defer s.Close()
	err := c.Find(bson.M{}).All(&tradeDatas)
	if err != nil {
		log.Debug("error:%s", err)
		return nil, err
	}
	for _, tradeData := range tradeDatas {
		amountUsd += tradeData.AmountUsd
		amountCny += tradeData.AmountCny
	}
	exchangeAmount := &data.ExchangeAmount{
		Platform: PlatForm,
		TotalUsd: amountUsd,
		TotalCny: amountCny,
	}
	return exchangeAmount, nil
}

/**
交易所所有交易对及其价格（从数据库获取）
完全没问题
*/
func (bb *Biki) GetExchangeTickersFromDb() ([]*data.ExchangeTicker, error) {
	var exchangeTicker []*data.ExchangeTicker
	s, c := common.Connect(DbPair, "biki", "local")
	defer s.Close()
	err := c.Find(bson.M{}).All(&exchangeTicker)
	if err != nil {
		log.Debug("error:%s", err)
		return nil, err
	}
	return exchangeTicker, nil
}
