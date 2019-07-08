// Package coinegg coinegg rest api package
package coinegg

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	//"time"
	"errors"
)

const (
	API = "https://api.coinegg.com"
	KLINE = "https://k.coinegg.com/period/data?type=kline&symbol=coinegg_"
	API2 = "https://www.coinegg.com/coin/"
	BTCPRICEAPI = "https://www.coinegg.com/index/pricebtc"
)

// Coinegg API data
type Coinegg struct {
	AccessKey string
	SecretKey string
}

// New create new Coinegg API data
func New(accessKey string, secretKey string) *Coinegg {
	return &Coinegg{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (*Coinegg) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "/api/v1/ticker" + "?coin=" + strings.ToLower(quote)

	log.Debugf("Request url: %v", url)
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
	fmt.Println("response:",string(body))
	buyRes := gjson.GetBytes(body, "buy").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "sell").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	/*lastRes := gjson.GetBytes(body, "last").String()
	last, err := strconv.ParseFloat(lastRes, 64)
	if err != nil {
		return nil, err
	}

	lowRes := gjson.GetBytes(body, "low").String()
	low, err := strconv.ParseFloat(lowRes, 64)
	if err != nil {
		return nil, err
	}

	highRes := gjson.GetBytes(body, "high").String()
	high, err := strconv.ParseFloat(highRes, 64)
	if err != nil {
		return nil, err
	}

	vol := gjson.GetBytes(body, "vol").Float()*/

	url2 := API2 + strings.ToLower(base) + "/allcoin"

	log.Debugf("Request url: %v", url2)
	fmt.Println("url:",url2)
	resp, err = http.Get(url2)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	tickers := gjson.ParseBytes(body).Map()
	key :=  strings.ToLower(quote)
	if !tickers[key].Exists() {
		return nil, errors.New("The ticker not exists")
	}

	v := tickers[key]
	return &model.Ticker{
		//Buy:  v.Array()[].Float(),
		//Sell: v.Get("lowestAsk").Float(),
		Buy:  buy,
		Sell: sell,
		Last: v.Array()[1].Float(),
		Low:  v.Array()[5].Float(),
		High: v.Array()[4].Float(),
		Vol:  v.Array()[6].Float(),
		PriceChangePercent:v.Array()[8].Float()/100,
	}, nil
}

func (*Coinegg) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := KLINE  + strings.ToLower(quote) + "_" + strings.ToLower(base) +"&step=" + typ +"&size=" + strconv.Itoa(size)

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

	var records []model.Record

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			//Time:  time.Unix(v.Array()[0].Int(), 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
			Ktime:v.Array()[0].Int()/1000,
		}

		records = append(records, record)
		return true
	})

	return records, nil
}

func (*Coinegg) GetBTCPRICE() (*model.Ticker, error) {
	url := BTCPRICEAPI

	log.Debugf("Request url: %v", url)
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
	fmt.Println("response:",string(body))
	buyRes := gjson.GetBytes(body, "data.usd").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}
	fmt.Println("btcprice:",buy)
	return &model.Ticker{
		//Buy:  v.Array()[].Float(),
		//Sell: v.Get("lowestAsk").Float(),
		Buy:  buy,
	}, nil
}

func (*Coinegg) GetETHPRICE() (*model.Ticker, error) {

	url2 := API2 + strings.ToLower("btc") + "/allcoin"

	log.Debugf("Request url: %v", url2)
	fmt.Println("url:",url2)
	resp, err := http.Get(url2)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	tickers := gjson.ParseBytes(body).Map()
	key :=  strings.ToLower("eth")
	if !tickers[key].Exists() {
		fmt.Println("get eth error")
		return nil, errors.New("The ticker not exists")
	}

	v := tickers[key]
	fmt.Println("ethprice:",v.Array()[1].Float())
	return &model.Ticker{
		//Buy:  v.Array()[].Float(),
		//Sell: v.Get("lowestAsk").Float(),
		Buy:  v.Array()[1].Float(),
	}, nil
}