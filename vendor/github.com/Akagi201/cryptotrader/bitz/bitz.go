// Package bithumb huobi rest api package
package bitz

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	//"time"
)

const (
	API = "https://www.bit-z.com/api_v1/"
)

// Huobi API data
type Bitz struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Bitz {
	return &Bitz{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (bitz *Bitz) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "ticker?coin="+strings.ToLower(quote)+"_"+strings.ToLower(base)

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

	url2 := "https://www.bit-z.com/index/infofresh"
	symbol := strings.ToLower(quote)+"_"+strings.ToLower(base)
	log.Debugf("Request url: %v", url)
	fmt.Println("url:",url2)
	resp, err = http.Get(url2)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))
	pathgetbyte := "data." + symbol
	datajson := gjson.GetBytes(body, pathgetbyte)
	if datajson.Exists() {
		re := datajson.String()
		re = strings.Replace(re,"\\","",-1)
		fmt.Println(re)

		changeper := gjson.Get(re, "float_24").Float()
				//changeperstring = strings.Replace(changeperstring,"%","",1)
				//changeper, _ := strconv.ParseFloat(changeperstring, 64)
				changeper = changeper /100
				return &model.Ticker{
					//Buy:  v.Get("buy").Float(),
					//Sell: v.Get("sell").Float(),
					Buy:  buy,
					Sell: sell,
					Last: last,
					Low:  low,
					High: high,
					Vol:  vol,
					PriceChangePercent:changeper * 100,
					//Raw:  string(body),
				}, nil


	}
	return &model.Ticker{},nil
}


func (bitz *Bitz) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + "kline?coin="+strings.ToLower(quote)+"_"+strings.ToLower(base)

	if len(typ) != 0 {
		url += "&type=" + typ
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
	datajson := gjson.GetBytes(body, "data.datas.data")
	//fmt.Println(datajson)
	gjson.Parse(datajson.String()).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
			Ktime:v.Array()[0].Int()/1000,
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

func (bitz *Bitz) GetExchange(base string, quote string) (*model.Excharge, error) {
	//url := API + "ticker?coin="+strings.ToLower(quote)+"_"+strings.ToLower(base)
	url := "https://www.bit-z.com/ajax/exchangeRate"
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
	exchargedkkt := gjson.GetBytes(body, "data.dkkt_usd").Float()
	exchargeeth := gjson.GetBytes(body, "data.eth_usd").Float()
	exchargebtc := gjson.GetBytes(body, "data.btc_usd").Float()
		return &model.Excharge{
			//Buy:  v.Get("buy").Float(),
			//Sell: v.Get("sell").Float(),
			Exchargebtcusd:exchargebtc,
			Exchargedkktusd:exchargedkkt,
			Exchargeethusd:exchargeeth,
			//Raw:  string(body),
		}, nil
}