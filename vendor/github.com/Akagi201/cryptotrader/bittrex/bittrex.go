// Package bittrex bittrex rest api package
package bittrex

import (
	"io/ioutil"
	"net/http"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	"strings"
	"time"
)

const (
	API = "https://bittrex.com/api/v1.1"
	KLINE = "https://bittrex.com/Api/v2.0/pub/market/GetTicks?marketName="
)

// Bittrex API data
type Bittrex struct {
	AccessKey string
	SecretKey string
}

// New create new Yunbi API data
func New(accessKey string, secretKey string) *Bittrex {
	return &Bittrex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (yb *Bittrex) GetTicker(base string, quote string) (ticker *model.Ticker, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			ticker = nil
			rerr = err.(error)
		}
	}()

	url := API + "/public/getmarketsummary?market=" + base + "-" + quote

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

	buy := gjson.GetBytes(body, "result.#.Ask").Array()[0].Float()
	sell := gjson.GetBytes(body, "result.#.Bid").Array()[0].Float()
	last := gjson.GetBytes(body, "result.#.Last").Array()[0].Float()
	low := gjson.GetBytes(body, "result.#.Low").Array()[0].Float()
	high := gjson.GetBytes(body, "result.#.High").Array()[0].Float()
	vol := gjson.GetBytes(body, "result.#.BaseVolume").Array()[0].Float()
	prevday := gjson.GetBytes(body, "result.#.PrevDay").Array()[0].Float()
	pricechangepercent := (last-prevday)/prevday

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		PriceChangePercent:pricechangepercent * 100,
	}, nil
}

func (yb *Bittrex) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := KLINE  + strings.ToUpper(base) + "-" + strings.ToUpper(quote) + "&tickInterval=" + typ

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
	datajson := gjson.GetBytes(body, "result")
	if datajson.Exists() {
		re := datajson.Array()
		for _,v:=  range re {
			//fmt.Println(v.Get("T").String())
			 toBeCharge := v.Get("T").String()
			timeLayout := "2006-01-02T15:04:05"                             //转化所需模板
			loc, _ := time.LoadLocation("Asia/Chongqing")                            //重要：获取时区
			theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
			sr := theTime.Unix()+28800                                         //转化为时间戳 类型是int64
			//fmt.Println(theTime)                                            //打印输出theTime 2015-01-01 15:15:00 +0800 CST
			//fmt.Println(sr)
			record := model.Record{

				Time:  time.Unix(v.Get("T").Int(), 0),
				//Time:v.Get("T").Int(),
				//Time:  time.Time(v.Get("T").String()),
				Open:  v.Get("O").Float(),
				High:  v.Get("H").Float(),
				Low:   v.Get("L").Float(),
				Close: v.Get("C").Float(),
				Vol:   v.Get("V").Float(),
				Ktime:sr,
			}

			records = append(records, record)
		}
	}
	/*
		gjson.GetBytes(body, "candles").ForEach(func(k, v gjson.Result) bool {
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
