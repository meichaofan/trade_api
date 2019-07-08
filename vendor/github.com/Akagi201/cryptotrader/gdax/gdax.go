// Package bithumb huobi rest api package
package gdax

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
)

const (
	API   = "https://api.gdax.com/products/"
	KLINE = "https://api.gdax.com/products/"
	STATS = "https://api.gdax.com/products/"
)

// Huobi API data
type Gdax struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Gdax {
	return &Gdax{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (gdax *Gdax) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToUpper(quote) + "-" + strings.ToUpper(base) + "/ticker"

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
	vol := gjson.GetBytes(body, "volume").Float()
	url2 := API + strings.ToUpper(quote) + "-" + strings.ToUpper(base) + "/stats"

	log.Debugf("Request url: %v", url2)
	fmt.Println("Request url:", url2)
	resp2, err2 := http.Get(url2)
	if err2 != nil {
		return nil, err2
	}
	defer resp2.Body.Close()
	body2, err2 := ioutil.ReadAll(resp2.Body)
	if err2 != nil {
		return nil, err2
	}
	last := gjson.GetBytes(body2, "last").Float()
	low := gjson.GetBytes(body2, "low").Float()
	high := gjson.GetBytes(body2, "high").Float()
	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}

func (gdax *Gdax) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := KLINE + strings.ToUpper(quote) + "-" + strings.ToUpper(base) + "/candles/" + typ

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
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int(), 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
			Ktime: v.Array()[0].Int(),
		}

		recordscache = append(recordscache, record)
		return true
	})
	lenrecord := len(recordscache) - 1
	for i := lenrecord; i >= 0; i-- {
		records = append(records, recordscache[i])
	}
	return records, nil
}

func (gdax *Gdax) GetTrades(base string, quote string, typ string, since int, size int) ([]model.Trade, error) {
	url := KLINE + strings.ToUpper(quote) + "-" + strings.ToUpper(base) + "/trades"

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

	var records []model.Trade
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		timeLayout := "2006-01-02T15:04:05.000Z"                             //转化所需模板
		loc, _ := time.LoadLocation("Asia/Chongqing")                            //重要：获取时区
		theTime, _ := time.ParseInLocation(timeLayout, v.Get("time").String(), loc) //使用模板在对应时区转化为time.time类型
		sr := theTime.Unix()*1000
		record := model.Trade{
			ID:        v.Get("trade_id").Int(),
			Price:     v.Get("price").Float(),
			Amount:    v.Get("size").Float(),
			Type:      v.Get("side").String(),
			TradeTime: sr,
		}

		records = append(records, record)
		return true
	})
	return records, nil
}
