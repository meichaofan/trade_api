// Package coincheck coincheck rest api
package coincheck

import (
	coincheck "github.com/Akagi201/coincheckgo"
	"github.com/Akagi201/cryptotrader/model"
	"github.com/tidwall/gjson"
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	//"math/rand"
	"net/http"
	log "github.com/sirupsen/logrus"
)

// Coincheck API data
type Coincheck struct {
	coincheck.CoinCheck
}

const (
	KLINE = "https://coincheck.com/api/charts/candle_rates?limit="
	TRADEURL = "https://coincheck.com/api/trades?pair="
	//https://cex.io/api/ohlcv2/hd/20180318/BTC/USD
)
// New create new Allcoin API data
func New(accessKey string, secretKey string) *Coincheck {
	client := new(coincheck.CoinCheck).NewClient(accessKey, secretKey)

	return &Coincheck{
		client,
	}
}

// GetTicker 行情
func (cc *Coincheck) GetTicker(base string, quote string) (*model.Ticker, error) {

	resp := cc.Ticker.All()
	fmt.Println("resp:",resp)
	buy := gjson.Get(resp, "bid").Float()
	sell := gjson.Get(resp, "ask").Float()
	last := gjson.Get(resp, "last").Float()
	low := gjson.Get(resp, "low").Float()
	high := gjson.Get(resp, "high").Float()
	vol := gjson.Get(resp, "volume").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}

func (cc *Coincheck) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := KLINE  + typ +"&market=coincheck&pair="+ strings.ToLower(quote) + "_" + strings.ToLower(base) +"&unit=" + typ +"&v2=true"

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
			Ktime:v.Array()[0].Int(),
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


func (cc *Coincheck) GetTrade(base string, quote string, typ string, since int, size int) ([]model.Trade, error) {
	url := TRADEURL+ strings.ToLower(quote) + "_" + strings.ToLower(base)

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

	//log.Debugf("Response body: %v", string(body))

	var records []model.Trade
	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		re := datajson.Array()
		for _, v := range re {
			/*open := v.Get("open").Float()
			close := v.Get("close").Float()
			precent24 := (close - open) / open*/
			timeLayout := "2006-01-02T15:04:05.000Z"                             //转化所需模板
			loc, _ := time.LoadLocation("Asia/Chongqing")                            //重要：获取时区
			theTime, _ := time.ParseInLocation(timeLayout, v.Get("created_at").String(), loc) //使用模板在对应时区转化为time.time类型
			sr := theTime.Unix()*1000
			trade := model.Trade{
				ID :v.Get("id").Int(),
				Price:v.Get("rate").Float(),
				Amount:v.Get("amount").Float(),
				Type:v.Get("order_type").String(),
				TradeTime:sr,
			}

			records = append(records, trade)
		}
	}


	return records, nil
}