package bibox

import (
	"github.com/Akagi201/cryptotrader/model"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
	"truxing/commons/log"
)

/**
文档
https://github.com/Biboxcom/API_Docs/wiki
*/

const (
	RestHost = "https://api.bibox365.com/v1/mdata"
	PlatForm = "bibox"
	DbPair   = "platform_pair"
	DbAmount = "platform_amount"
)

//定义taker的成交方向
var sideMap = map[string]string{
	"1": "buy",
	"2": "sell",
}

// BitBox client
type BitBox struct {
	ApiKey string
	Secret string
}

func New(apiKey string, secret string) *BitBox {
	return &BitBox{
		ApiKey: apiKey,
		Secret: secret,
	}
}

func getSymbol(base, quote string) string {
	return strings.ToUpper(strings.TrimSpace(quote)) + "_" + strings.ToUpper(strings.TrimSpace(base))
}

//pair ETH_BTC
func splitSymbol(pair string) (base, quote string) {
	r := strings.Split(pair, "_")
	return r[0], r[1]
}

/**
获取最新价
*/
func (bb *BitBox) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "?cmd=ticker&pair=" + getSymbol(base, quote)
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
	percent := result.Get("result.percent").Str
	priceChangePercent, _ := strconv.ParseFloat(strings.TrimRight(percent, "%"), 64)
	ticker = &model.Ticker{
		Buy:                result.Get("result.buy").Float(),
		Sell:               result.Get("result.sell").Float(),
		Last:               result.Get("result.last").Float(),
		Low:                result.Get("result.low").Float(),
		High:               result.Get("result.high").Float(),
		Vol:                result.Get("result.vol").Float(),
		Time:               cast.ToTime(result.Get("result.timestamp").Int() / 1000),
		PriceChangePercent: priceChangePercent,
		Raw:                result.Get("result").Str,
	}
	return ticker, nil
}

/**
获取成交记录
https://github.com/Biboxcom/API_Docs/wiki/REST_API_Reference#%E6%9F%A5%E8%AF%A2%E6%88%90%E4%BA%A4%E8%AE%B0%E5%BD%95
*/
func (bb *BitBox) GetTades(base, quote string, size int) ([]model.Trade, error) {
	url := RestHost + "?cmd=deals&pair=" + getSymbol(base, quote)
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

	result.Get("result").ForEach(func(key, value gjson.Result) bool {
		trade := model.Trade{
			Amount:    value.Get("amount").Float(),
			Price:     value.Get("price").Float(),
			Type:      sideMap[value.Get("side").String()],
			Raw:       value.String(),
			TradeTime: value.Get("time").Int() / 1000,
			Time:      cast.ToTime(value.Get("time").Int() / 1000),
		}
		trades = append(trades, trade)
		return true
	})

	return trades, nil
}

/**
查询买卖盘深度
https://github.com/Biboxcom/API_Docs/wiki/REST_API_Reference#%E6%9F%A5%E8%AF%A2%E5%B8%82%E5%9C%BA%E6%B7%B1%E5%BA%A6
*/
func (bb *BitBox) GetMarketDepth(base, quote string, size int) (model.OrderBook, error) {
	orderBook := model.OrderBook{}
	url := RestHost + "?cmd=depth&pair=" + getSymbol(base, quote)
	if size != 0 {
		url += "&size" + strconv.Itoa(size)
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

	gjson.GetBytes(body, "result.bids").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Get("price").Float()
		order.Amount = value.Get("volume").Float()
		orderBook.Bids = append(orderBook.Bids, order)
		return true // keep iterating
	})

	gjson.GetBytes(body, "result.asks").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Get("price").Float()
		order.Amount = value.Get("volume").Float()
		orderBook.Asks = append(orderBook.Asks, order)
		return true // keep iterating
	})

	orderBook.OrderTime = gjson.GetBytes(body, "result.update_time").Int() / 1000
	orderBook.Time = cast.ToTime(orderBook.OrderTime)
	orderBook.Raw = gjson.ParseBytes(body).Get("result").Raw

	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
查询系统支持的所有交易对及精度
*/
func (bb *BitBox) GetMarkets() ([]model.MarketPairInfo, error) {
	var tradePairs []model.MarketPairInfo
	url := RestHost + "?cmd=pairList"
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
	gjson.ParseBytes(body).Get("result").ForEach(func(key, value gjson.Result) bool {
		base, quote := splitSymbol(value.Get("pair").String())
		tradePair := model.MarketPairInfo{
			Base:  base,
			Quote: quote,
		}
		tradePairs = append(tradePairs, tradePair)
		return true
	})
	return tradePairs, nil
}

/**
交易所所有交易对及其价格
*/
func (bb *BitBox) GetExchangeTickers() (model.ExchangeTickers, error) {
	var exchangeTickers model.ExchangeTickers
	url := RestHost + "?cmd=marketAll"
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
	gjson.ParseBytes(body).Get("result").ForEach(func(key, value gjson.Result) bool {

		marketPair := model.MarketPairInfo{
			Base:  value.Get("currency_symbol").String(),
			Quote: value.Get("coin_symbol").String(),
		}

		exchangeTicker := &model.ExchangeTicker{
			MarketPair:         marketPair,
			Vol:                value.Get("vol24H").Float(),
			Amount:             value.Get("amount").Float(),
			Last:               value.Get("last").Float(),
			LastUSD:            value.Get("last_usd").Float(),
			PriceChangePercent: value.Get("percent").Float(),
			Time:               time.Now(),
		}

		exchangeTickers = append(exchangeTickers, exchangeTicker)
		return true
	})
	return exchangeTickers, nil
}

/**
交易所所有交易对及其价格（从数据库获取）
完全没问题
*/

func (bb *BitBox) GetExchangeTickersFromDb() ([]*data.ExchangeTicker, error) {
	var exchangeTicker []*data.ExchangeTicker
	s, c := common.Connect(DbPair, "bibox", "local")
	defer s.Close()
	err := c.Find(bson.M{}).All(&exchangeTicker)
	if err != nil {
		log.Debug("error:%s", err)
		return nil, err
	}
	return exchangeTicker, nil
}

/**
获取当前平台所有交易额
https://github.com/Biboxcom/API_Docs/wiki/REST_API_Reference#%E6%9F%A5%E8%AF%A2%E6%88%90%E4%BA%A4%E8%AE%B0%E5%BD%95
这个👍
*/
func (bb *BitBox) GetExchangeAmount() (model.ExchangeAmount, error) {
	var exchangeAmount model.ExchangeAmount
	url := RestHost + "?cmd=marketAll"
	log.Debugf("url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return exchangeAmount, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return exchangeAmount, err
	}
	log.Debugf("Response body: %v", string(body))
	gjson.ParseBytes(body).Get("result").ForEach(func(key, value gjson.Result) bool {
		last := value.Get("last").Float()
		lastUsd := value.Get("last_usd").Float()
		rate := lastUsd / last
		exchangeAmount.AmountUSD += value.Get("amount").Float() * rate
		return true
	})
	exchangeAmount.PlatForm = PlatForm
	return exchangeAmount, nil
}

/**
获取当前平台所有交易额(从数据库获取)
*/
func (bb *BitBox) GetExchangeAmountFormDb() (*data.ExchangeAmount, error) {
	var tradeDatas []*data.TradeData
	var amount float64 = 0
	s, c := common.Connect(DbAmount, "bibox", "local")
	defer s.Close()
	err := c.Find(bson.M{}).All(&tradeDatas)
	if err != nil {
		log.Debug("error:%s", err)
		return nil, err
	}
	for _, tradeData := range tradeDatas {
		amount += tradeData.AmountUsd
	}
	exchangeAmount := &data.ExchangeAmount{
		Platform: PlatForm,
		TotalUsd: amount,
	}
	return exchangeAmount, nil
}

/**
获取k线 https://github.com/Biboxcom/API_Docs/wiki/REST_API_Reference#%E6%9F%A5%E8%AF%A2k%E7%BA%BF
*/
func (bb *BitBox) GetRecords(base, quote, period string, size int) ([]model.Record, error) {
	url := RestHost + "?cmd=kline&pair=" + getSymbol(base, quote) + "&period=" + period
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
	gjson.ParseBytes(body).Get("result").ForEach(func(key, value gjson.Result) bool {
		record.Open = value.Get("open").Float()
		record.High = value.Get("high").Float()
		record.Low = value.Get("low").Float()
		record.Close = value.Get("close").Float()
		record.Vol = value.Get("vol").Float()
		record.Ktime = value.Get("time").Int() / 1000
		records = append(records, record)
		return true // keep iterating
	})

	return records, nil
}
