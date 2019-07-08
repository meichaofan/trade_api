// Package bithumb huobi rest api package
package coinsbank

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
	API = "https://coinsbank.com/sapi/trade/ohlcv?pairCode="
	CHARGEAPI = "https://coinsbank.com/sapi/head"
)

// Huobi API data
type Coinsbank struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Coinsbank {
	return &Coinsbank{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (coinsbank *Coinsbank) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToUpper(quote)+strings.ToUpper(base) + "&interval=60"
	chargeapi :=CHARGEAPI
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
	//buy := gjson.GetBytes(body, ".#.Ask").Array()[0].Float()
	//sell := gjson.GetBytes(body, ".#.Ask").Array()[0].Float()
	last := gjson.GetBytes(body, "#.c").Array()[0].Float()
	low := gjson.GetBytes(body, "#.l").Array()[0].Float()
	high := gjson.GetBytes(body, "#.h").Array()[0].Float()
	vol := gjson.GetBytes(body, "#.v").Array()[0].Float()
	respc, err := http.Get(chargeapi)
	if err != nil {
		return nil, err
	}
	defer respc.Body.Close()
	bodyc, err := ioutil.ReadAll(respc.Body)
	if err != nil {
		return nil, err
	}
	bodycpath := "ratesChanges." + strings.ToUpper(quote)+strings.ToUpper(base)
	ptc := gjson.GetBytes(bodyc, bodycpath).Float()/100
	return &model.Ticker{
		//Buy:  buy,
		//Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		PriceChangePercent:ptc,
	}, nil
}

func (coinsbank *Coinsbank) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + strings.ToUpper(quote)+strings.ToUpper(base) + "&interval=" + typ

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
	var recordscache []model.Record
	gjson.Parse(string(body)).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Get("date").Int()/1000, 0),
			Open:  v.Get("o").Float(),
			High:  v.Get("h").Float(),
			Low:   v.Get("l").Float(),
			Close: v.Get("c").Float(),
			Vol:   v.Get("v").Float(),
			Ktime:v.Get("date").Int()/1000,
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
