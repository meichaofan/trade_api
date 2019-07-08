// Package cex cex rest api package
package cex

import (
	"io/ioutil"
	//"math/rand"
	"net/http"
	//"strconv"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	"strings"
	"time"
	"encoding/json"
)

const (
	API = "https://cex.io/api/ticker/"
	KLINE = "https://cex.io/api/ohlcv2/d/"
	KILNEOTHER = "https://cex.io/api/ohlcv/hd/"
	//https://cex.io/api/ohlcv2/hd/20180318/BTC/USD
)

// Cex API data
type Cex struct {
	AccessKey string
	SecretKey string
}

// New create new Cex API data
func New(accessKey string, secretKey string) *Cex {
	return &Cex{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
/*
func (bt *Cex) GetTicker(base string, quote string) (ticker *model.Ticker, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			ticker = nil
			rerr = err.(error)
		}
	}()

	var url string
	if quote == "pay" {
		url = API + "/trade_43.js?v=" + strconv.FormatFloat(rand.Float64(), 'g', 1, 64)
	} else if quote == "omg" {
		url = API + "/trade_41.js?v=" + strconv.FormatFloat(rand.Float64(), 'g', 1, 64)
	}
	fmt.Println("url:",url)
	log.Debugf("Request url: %v", url)

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
	fmt.Println("reponse:",string(body))
	buy := gjson.GetBytes(body, "depth.1.#.price").Array()[0].Float()
	sell := gjson.GetBytes(body, "depth.2.#.price").Array()[0].Float()
	last := gjson.GetBytes(body, "cmark.new_price").Float()
	low := gjson.GetBytes(body, "cmark.min_price").Float()
	high := gjson.GetBytes(body, "cmark.max_price").Float()
	vol := gjson.GetBytes(body, "cmark.H24_done_num").Float()

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}
*/

func (bt *Cex) GetTicker(base string, quote string) (ticker *model.Ticker, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			ticker = nil
			rerr = err.(error)
		}
	}()

	url := API  + strings.ToUpper(quote) + "/" + strings.ToUpper(base)

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
	fmt.Println("reponse:", string(body))
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

func (bt *Cex) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	//url := KLINE  + typ +"/"+ strings.ToUpper(quote) + "/" + strings.ToUpper(base)
	//since = since-86400
	if since == 0 {
		if typ == "data1m" {
			since = int(time.Now().Unix()) - 2000
		} else {
			since = int(time.Now().Unix()) - 86400
		}

	}
	table_time := time.Unix(int64(since),0).Format("20060102")
	url := KILNEOTHER + table_time + "/"+strings.ToUpper(quote) + "/" + strings.ToUpper(base)
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
	datajson := gjson.GetBytes(body, typ)
	//fmt.Println(datajson)
	if datajson.Exists() {
		ret := string(datajson.String())

		var res []interface{}
		err = json.Unmarshal([]byte(ret),&res)
		if err != nil {
			return nil, err
		}
		//var sentence = (res[0]).([]interface{})
		//fmt.Println(sentence[0].(float64)/1000)
		//var translated = make([]string,0)
		fmt.Println(len(res))
		for i:=0;i<len(res);i++ {
			n := (res[i]).([]interface{})
			//fmt.Println(n[0].(float64)/1000)
			ktime := int64(n[0].(float64))
			//ktimeint, _ := strconv.ParseInt(ktime, 10, 0)

			record := model.Record{

				//Time:  time.Unix(v.Array()[0].Int(), 0),
				//Time:v.Get("T").Int(),
				//Time:  time.Time(v.Get("T").String()),
				Open:  n[1].(float64),
				High:  n[2].(float64),
				Low:   n[3].(float64),
				Close: n[4].(float64),
				Vol:   n[5].(float64),
				Ktime: ktime ,
			}

			records = append(records, record)
		}
		/*
		re := datajson.Array()
		for _,v:=  range re {
			//fmt.Println(v.Get("T").String())
			fmt.Println(v.Array())
			record := model.Record{

				Time:  time.Unix(v.Array()[0].Int(), 0),
				//Time:v.Get("T").Int(),
				//Time:  time.Time(v.Get("T").String()),
				Open:  v.Array()[1].Float(),
				High:  v.Array()[2].Float(),
				Low:   v.Array()[3].Float(),
				Close: v.Array()[4].Float(),
				Vol:   v.Array()[5].Float(),
				Ktime:v.Array()[0].Int(),
			}

			records = append(records, record)
		}*/
	}
/*
	gjson.GetBytes(body, "ohlcv").ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int(), 0),
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