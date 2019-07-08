// Package bithumb huobi rest api package
package bitbank

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	"time"
	//"strconv"
)

const (
	API = "https://public.bitbank.cc/"
)

// Huobi API data
type Bitbank struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Bitbank {
	return &Bitbank{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (bitbank *Bitbank) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToLower(quote)+"_"+strings.ToLower(base) + "/ticker"

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
	buy := gjson.GetBytes(body, "data.buy").Float()
	sell := gjson.GetBytes(body, "data.sell").Float()
	last := gjson.GetBytes(body, "data.last").Float()
	low := gjson.GetBytes(body, "data.low").Float()
	high := gjson.GetBytes(body, "data.high").Float()
	vol := gjson.GetBytes(body, "data.vol").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}


func (bitbank *Bitbank) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	table_time := ""
	if since == 0 {
		since = int(time.Now().Unix())
	}
	if typ == "1hour" || typ == "1min" {
		table_time = time.Unix(int64(since),0).Format("20060102")
	} else {
		table_time = time.Unix(int64(since),0).Format("2006")
	}

	url := API + strings.ToLower(quote)+"_"+strings.ToLower(base) + "/candlestick" + "/" + typ + "/" + table_time
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
	datajson := gjson.GetBytes(body, "data.candlestick.#.ohlcv")
	fmt.Println(datajson)
	gjson.Parse(datajson.Array()[0].String()).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[5].Int()/1000, 0),
			Open:  v.Array()[0].Float(),
			High:  v.Array()[1].Float(),
			Low:   v.Array()[2].Float(),
			Close: v.Array()[3].Float(),
			Vol:   v.Array()[4].Float(),
			Ktime:v.Array()[5].Int()/1000-28800,
		}

		records = append(records, record)
		return true
	})
/*
	gjson.GetBytes(body, "data").ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
		}

		records = append(records, record)
		return true
	})
*/
	return records, nil
}