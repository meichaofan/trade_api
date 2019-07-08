// Package bithumb huobi rest api package
package bithumb

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
)

const (
	API = "https://api.bithumb.com/public/ticker/"
)

// Huobi API data
type Bithumb struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Bithumb {
	return &Bithumb{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (bithumb *Bithumb) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToLower(quote)

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
	buy := gjson.GetBytes(body, "data.buy_prise").Float()
	sell := gjson.GetBytes(body, "data.sell_price").Float()
	last := gjson.GetBytes(body, "data.closing_price").Float()
	low := gjson.GetBytes(body, "data.min_price").Float()
	high := gjson.GetBytes(body, "data.max_price").Float()
	vol := gjson.GetBytes(body, "data.volume_1day").Float()
	open := gjson.GetBytes(body, "data.opening_price").Float()
	close := gjson.GetBytes(body, "data.closing_price").Float()
	pcg := (open-close) / open

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
