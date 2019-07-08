// Package bitfinex bitfinex rest api
package bitfinex

import (
	"strconv"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v1"
	"net/http"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
	"fmt"
)

// Bitfinex API data
type Bitfinex struct {
	bitfinex.Client
}
const (

	Kline = "https://api.bitfinex.com/v2/candles/"
	API = "https://api.bitfinex.com/v2/tickers?symbols=t"
)
// New create new Allcoin API data
func New(accessKey string, secretKey string) *Bitfinex {
	var client *bitfinex.Client
	if accessKey != "" && secretKey != "" {
		client = bitfinex.NewClient().Auth(accessKey, secretKey)
	} else {
		client = bitfinex.NewClient()
	}

	return &Bitfinex{
		*client,
	}
}

// GetTicker 行情
func (bf *Bitfinex) GetTicker(base string, quote string) (*model.Ticker, error) {
	tick, err := bf.Ticker.Get(strings.ToUpper(quote) + strings.ToUpper(base))
	fmt.Println("resq:",tick)
	buy, err := strconv.ParseFloat(tick.Bid, 64)
	if err != nil {
		return nil, err
	}

	sell, err := strconv.ParseFloat(tick.Ask, 64)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseFloat(tick.LastPrice, 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(tick.Low, 64)
	if err != nil {
		return nil, err
	}

	high, err := strconv.ParseFloat(tick.High, 64)
	if err != nil {
		return nil, err
	}

	vol, err := strconv.ParseFloat(tick.Volume, 64)
	if err != nil {
		return nil, err
	}
	url := API + strings.ToUpper(quote) + strings.ToUpper(base)

	log.Debugf("Request url: %v", url)
	fmt.Println("Request url:", url)
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
	fmt.Println("Response body:", string(body))
	var pcg float64
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
			pcg= v.Array()[6].Float()
			//Raw:  string(body),
		return true
	})

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		PriceChangePercent:pcg,
	}, nil
}

func (bf *Bitfinex) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := Kline + "trade:" + typ + ":t"+ strings.ToUpper(quote) + strings.ToUpper(base) + "/hist"


	fmt.Println("url:",url)
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

	var records []model.Record

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[3].Float(),
			Low:   v.Array()[4].Float(),
			Close: v.Array()[2].Float(),
			Vol:   v.Array()[5].Float(),
			Ktime:v.Array()[0].Int()/1000,
		}

		records = append(records, record)
		return true
	})

	return records, nil
}