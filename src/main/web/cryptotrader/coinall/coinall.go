package coinall

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

const (
	RestHost = "https://www.CoinAll.live"
	//ApiVer   = "v3"
)

// CoinAll client
type CoinAll struct {
	AccessKey string
	SecretKey string
}

func New(accessKey string, secretKey string) *CoinAll {
	return &CoinAll{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func getInstrumentId(base, quote string) string {
	return strings.ToUpper(base) + "_" + strings.ToUpper(quote)
}

/**
https://www.coinall.live/docs/zh/#spot-some
GetTicker 获取某个ticker信息  /api/spot/v3/instruments/<instrument-id>/ticker
*/
func (ca *CoinAll) GetTicker(base, quote string) (*model.Ticker, error) {
	url := RestHost + "/api/spot/v3/products/" + getInstrumentId(base, quote) + "/ticker"
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
	buy := gjson.GetBytes(body, "best_ask").Float()
	sell := gjson.GetBytes(body, "best_bid").Float()
	last := gjson.GetBytes(body, "last").Float()
	low := gjson.GetBytes(body, "low_24h").Float()
	high := gjson.GetBytes(body, "high_24h").Float()
	vol := gjson.GetBytes(body, "base_volume_24h").Float()
	open_24h := gjson.GetBytes(body, "open_24h").Float()
	pricechangepercent := (last - open_24h) / open_24h
	return &model.Ticker{
		Buy:                buy,
		Sell:               sell,
		Last:               last,
		Low:                low,
		High:               high,
		Vol:                vol,
		PriceChangePercent: pricechangepercent * 100,
	}, nil
}

/**
https://www.coinall.live/docs/zh/#spot-line
GetRecords 获取币币K线数据 /api/spot/v3/instruments/<instrument_id>/candles
*/
func (ca *CoinAll) GetRecords(base, quote, typ string, since int) ([]model.Record, error) {
	url := RestHost + "/api/spot/v3/instruments/" + getInstrumentId(base, quote) + "/candles"
	start := time.Unix(int64(since), 0).UTC().Format(time.RFC3339)
	end := time.Now().UTC().Format(time.RFC3339)
	url += "?granularity=" + typ + "&start=" + start + "&end=" + end
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
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		record.Open = cast.ToFloat64(value.Array()[1].String())
		record.High = cast.ToFloat64(value.Array()[2].String())
		record.Low = cast.ToFloat64(value.Array()[3].String())
		record.Close = cast.ToFloat64(value.Array()[4].String())
		record.Vol = cast.ToFloat64(value.Array()[5].String())
		record.Time = cast.ToTime(value.Array()[0].Int() / 1000)
		record.Ktime = value.Array()[0].Int() / 1000
		record.Raw = value.String()
		records = append(records, record)
		return true // keep iterating
	})
	return records, nil
}

/**
https://www.coinall.live/docs/zh/#spot-deal_information
GetTades 获取大单动向（成交数据）/api/spot/v3/instruments/<instrument_id>/trades
*/
func (ca *CoinAll) GetTades(base, quote string, since int) ([]model.Trade, error) {
	url := RestHost + "/api/spot/v3/instruments/" + getInstrumentId(base, quote) + "/trades"
	if since != 0 {
		url += "&after=" + strconv.Itoa(since)
	}
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
	var trades []model.Trade
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		trade := model.Trade{
			Amount: value.Get("size").Float(),
			Price:  value.Get("price").Float(),
			ID:     value.Get("trade_id").Int(),
			Type:   value.Get("side").String(),
			Time:   value.Get("timestamp").Time(),
			Raw:    value.String(),
		}
		trades = append(trades, trade)
		return true
	})
	return trades, nil
}
/**
  https://www.coinall.live/docs/zh/#spot-currency
  GetMarkets 获取所有交易对 /api/spot/v3/instruments
*/

func (ca *CoinAll) GetMarkets() ([]model.MarketPairInfo, error) {
	var tradePairs []model.MarketPairInfo
	url := RestHost + "/api/spot/v3/instruments"
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
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		tradePair := model.MarketPairInfo{
			Quote: value.Get("quote_currency").String(),
			Base:  value.Get("base_currency").String(),
		}
		tradePairs = append(tradePairs, tradePair)
		return true
	})
	return tradePairs, nil
}

/**
https://www.coinall.live/docs/zh/#spot-data
GetDepth获取深度数据 /api/spot/v3/instruments/BTC-USDT/book?size=5&depth=0.2
*/
func (ca *CoinAll) GetDepth(base, quote string) (model.OrderBook, error) {
	orderBook := model.OrderBook{}
	url := RestHost + "/api/spot/v3/instruments/" + getInstrumentId(base, quote) + "/book"
	log.Debugf("url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return orderBook, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return orderBook, err
	}
	log.Debugf("Response body: %v", string(body))

	var order model.MarketOrder

	gjson.GetBytes(body, "bids").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Array()[0].Float()
		order.Amount = value.Array()[1].Float()
		orderBook.Bids = append(orderBook.Bids, order)
		return true // keep iterating
	})

	gjson.GetBytes(body, "asks").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Array()[0].Float()
		order.Amount = value.Array()[1].Float()
		orderBook.Asks = append(orderBook.Asks, order)
		return true // keep iterating
	})

	log.Debugf("%v", orderBook)
	return orderBook, nil
}

/**
https://www.coinall.live/docs/zh/#spot-all
 */