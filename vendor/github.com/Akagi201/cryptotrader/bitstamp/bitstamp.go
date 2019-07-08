// Package bithumb huobi rest api package
package bitstamp

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
	API = "https://www.bitstamp.net/api/v2/ticker/"
)

// Huobi API data
type Bitstamp struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Bitstamp {
	return &Bitstamp{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (bitstamp *Bitstamp) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToLower(quote)+strings.ToLower(base)

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
	buy := gjson.GetBytes(body, "bid").Float()
	sell := gjson.GetBytes(body, "ask").Float()
	last := gjson.GetBytes(body, "last").Float()
	low := gjson.GetBytes(body, "low").Float()
	high := gjson.GetBytes(body, "high").Float()
	vol := gjson.GetBytes(body, "volume").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
