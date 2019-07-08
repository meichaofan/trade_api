// Package poloniex poloniex rest api package
package poloniex

import (
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	//"time"
	"fmt"
)

const (
	API = "https://poloniex.com/public"
)

// Poloniex API data
type Poloniex struct {
	AccessKey string
	SecretKey string
}

// New create new Poloniex API data
func New(accessKey string, secretKey string) *Poloniex {
	return &Poloniex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (pl *Poloniex) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "?command=returnTicker"

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

	tickers := gjson.ParseBytes(body).Map()
	key := strings.ToUpper(base) + "_" + strings.ToUpper(quote)
	if !tickers[key].Exists() {
		return nil, errors.New("The ticker not exists")
	}

	v := tickers[key]
	return &model.Ticker{
		Buy:  v.Get("highestBid").Float(),
		Sell: v.Get("lowestAsk").Float(),
		Last: v.Get("last").Float(),
		Low:  v.Get("low24hr").Float(),
		High: v.Get("high24hr").Float(),
		Vol:  v.Get("baseVolume").Float(),
		PriceChangePercent:v.Get("percentChange").Float(),
	}, nil
}

func (pl *Poloniex) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + "?command=returnChartData&currencyPair=" + strings.ToUpper(base) + "_" + strings.ToUpper(quote) + "&start=" + strconv.Itoa(since)

	if len(typ) != 0 {
		url += "&period=" + typ
	}

	fmt.Println("url:", url)
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
		fmt.Println(len(v.Array()))
		record := model.Record{
			//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
		}

		records = append(records, record)
		return true
	})

	return records, nil
}
