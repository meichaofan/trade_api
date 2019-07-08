// Package lbank lbank rest api package
package lbank

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
	"fmt"
)

const (
	API = "https://api.lbank.info/v1"
)

// LBank API data
type LBank struct {
	AccessKey string
	SecretKey string
}

// New create new LHang API data
func New(accessKey string, secretKey string) *LBank {
	return &LBank{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (*LBank) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "/ticker.do" + "?symbol=" + quote + "_" + base

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

	// TODO: use market depth API.
	buyRes := gjson.GetBytes(body, "ticker.latest").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "ticker.latest").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "ticker.latest").String()
	last, err := strconv.ParseFloat(lastRes, 64)
	if err != nil {
		return nil, err
	}

	lowRes := gjson.GetBytes(body, "ticker.low").String()
	low, err := strconv.ParseFloat(lowRes, 64)
	if err != nil {
		return nil, err
	}

	highRes := gjson.GetBytes(body, "ticker.high").String()
	high, err := strconv.ParseFloat(highRes, 64)
	if err != nil {
		return nil, err
	}

	volRes := gjson.GetBytes(body, "ticker.vol").String()
	vol, err := strconv.ParseFloat(volRes, 64)
	if err != nil {
		return nil, err
	}
	changesRes := gjson.GetBytes(body, "ticker.change").String()
	changes, err := strconv.ParseFloat(changesRes, 64)
	if err != nil {
		return nil, err
	}
	changes = changes / 100.00
	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		PriceChangePercent:changes,
	}, nil
}

func (*LBank) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + "/kline.do" + "?symbol=" + quote + "_" + base

	if len(typ) != 0 {
		url += "&type=" + typ
	}

	if since != 0 {
		url += "&time=" + strconv.Itoa(since)
	}

	if size != 0 {
		url += "&size=" + strconv.Itoa(size)
	}
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
	var recordscache []model.Record
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
			Ktime:v.Array()[0].Int()/1000,
		}

		recordscache = append(recordscache, record)
		return true
	})
	lenrecord := len(recordscache) -1
	for i:= lenrecord;i>=0 ;i-- {
		records = append(records,recordscache[i])
	}
	return records, nil
}