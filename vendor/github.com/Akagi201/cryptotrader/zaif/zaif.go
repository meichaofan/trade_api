// Package zaif huobi rest api package
package zaif

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

const (
	API   = "https://api.zaif.jp/api/1/ticker/"
	KLINE = "https://zaif.jp/zaif_chart_api/v1/history?symbol="
)

// Huobi API data
type Zaif struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Zaif {
	return &Zaif{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (zaif *Zaif) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToLower(quote) + "_" + strings.ToLower(base)

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

func (zaif *Zaif) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {


	nowtime := time.Now().Unix()
	//step,_ := strconv.Atoi(typ)
	starttime := int(nowtime) - 86400
	url := KLINE + strings.ToUpper(quote) + "_" + strings.ToUpper(base) + "&resolution="+typ+"&from=" + strconv.Itoa(int(starttime)) + "&to=" + strconv.Itoa(int(nowtime))

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
	//fmt.Println(string(body))
	stringbody := strings.Replace(string(body),"\\","",-1)
	var records []model.Record
	//fmt.Println(gjson.Get(string(body), "ohlc_data"))
	gjson.Get(stringbody, "ohlc_data").ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Get("time").Int()/1000, 0),
			Open:  v.Get("open").Float(),
			High:  v.Get("high").Float(),
			Low:   v.Get("low").Float(),
			Close: v.Get("close").Float(),
			Vol:   v.Get("volume").Float(),
			Ktime:v.Get("time").Int()/1000,
		}

		records = append(records, record)
		return true
	})

	return records, nil
}
