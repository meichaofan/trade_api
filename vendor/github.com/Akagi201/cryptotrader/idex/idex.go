// Package bithumb huobi rest api package
package idex

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	"bytes"
	"errors"
)

const (
	API = "https://api.idex.market/returnTicker"
	TRADEAPI = "https://api.idex.market/returnTradeHistory"
)

// Huobi API data
type Idex struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Idex {
	return &Idex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (idex *Idex) GetTicker(base string, quote string) (*model.Ticker, error) {
	//url := API + strings.ToLower(quote)+"_"+strings.ToLower(base)
	/*url := API
	post := "{\"market\":\"" + strings.ToUpper(base) + "_" +strings.ToUpper(quote) + "\"}"
	log.Debugf("Request url: %v", url)
	fmt.Println("Request post:", post)
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
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
	//buy := gjson.GetBytes(body, "data.buy").Float()
	//sell := gjson.GetBytes(body, "data.sell").Float()
	last := gjson.GetBytes(body, "last").Float()
	low := gjson.GetBytes(body, "low").Float()
	high := gjson.GetBytes(body, "high").Float()
	vol := gjson.GetBytes(body, "quoteVolume").Float()
	priceChangePercent := gjson.GetBytes(body, "percentChange").Float() /100.00

	return &model.Ticker{
		//Buy:  buy,
		//Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		PriceChangePercent:priceChangePercent,
	}, nil
*/
	url := API

	log.Debugf("Request url: %v", url)
	fmt.Println("url:",url)
	var jsonStr2 = []byte("")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr2))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
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
	key := strings.ToUpper(base) + "_" + strings.ToUpper(quote)
	if !tickers[key].Exists() {
		return nil, errors.New("The ticker not exists")
	}

	v := tickers[key]
	if v.Get("last").String() == "" {
		return &model.Ticker{} ,nil
	}
	return &model.Ticker{
		//Buy:  v.Get("highestBid").Float(),
		//Sell: v.Get("lowestAsk").Float(),
		Last: v.Get("last").Float(),
		Low:  v.Get("low").Float(),
		High: v.Get("high").Float(),
		Vol:  v.Get("quoteVolume").Float(),
		PriceChangePercent:v.Get("percentChange").Float()/100,
	}, nil
}
