// Package bithumb huobi rest api package
package hitbtc

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	"time"
	"strconv"
)

const (
	API = "https://api.hitbtc.com/api/2/public/"
)

// Huobi API data
type Hitbtc struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Hitbtc {
	return &Hitbtc{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (hitbtc *Hitbtc) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "ticker/"+ strings.ToUpper(quote)+strings.ToUpper(base)

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

func (hitbtc *Hitbtc) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	recordsize := size
	if size > 1000 {
		recordsize = 1000
	}
	url := API + "candles/"+ strings.ToUpper(quote)+strings.ToUpper(base) + "?limit=" + strconv.Itoa(recordsize)
	if len(typ) != 0 {
		url += "&period=" + typ
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
	datajson := gjson.ParseBytes(body)
	if datajson.Exists() {
		re := datajson.Array()
		for _,v:=  range re {
			toBeCharge := v.Get("timestamp").String()
			timeLayout := "2006-01-02T15:04:05.000Z"                             //转化所需模板
			loc, _ := time.LoadLocation("Asia/Chongqing")                            //重要：获取时区
			theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
			sr := theTime.Unix()  +28800                                           //转化为时间戳 类型是int64
			record := model.Record{

				//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
				Open:  v.Get("open").Float(),
				High:  v.Get("max").Float(),
				Low:   v.Get("min").Float(),
				Close: v.Get("close").Float(),
				Vol:   v.Get("volume").Float(),
				Ktime:sr,
			}

			records = append(records, record)
		}
	}
/*
	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
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