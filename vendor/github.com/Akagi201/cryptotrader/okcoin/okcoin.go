// Package okcoin okcoin rest api package
package okcoin

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"fmt"
	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
	"time"
	"encoding/json"
)

const (
	API = "https://www.okcoin.com/api/v1/"
)

// OkCoin API data
type OkCoin struct {
	AccessKey string
	SecretKey string
}

// New create new OkCoin API data
func New(accessKey string, secretKey string) *OkCoin {
	return &OkCoin{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (oc *OkCoin) GetTicker(base string, quote string) (*model.Ticker, error) {
	/*url := API + "ticker.do" + "?symbol=" + quote + "_" + base

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

	buyRes := gjson.GetBytes(body, "ticker.buy").String()
	buy, err := strconv.ParseFloat(buyRes, 64)
	if err != nil {
		return nil, err
	}

	sellRes := gjson.GetBytes(body, "ticker.sell").String()
	sell, err := strconv.ParseFloat(sellRes, 64)
	if err != nil {
		return nil, err
	}

	lastRes := gjson.GetBytes(body, "ticker.last").String()
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

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil*/
	url := "https://www.okcoin.com/v2/spot/markets/tickers"
	symbol := strings.ToLower(quote) + "_" + strings.ToLower(base)
	log.Debugf("Request url: %v", url)
	fmt.Println("url:", url)
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

	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		re := datajson.Array()
		for _, v := range re {
			if v.Get("symbol").String() == symbol {
				changeperstring := v.Get("changePercentage").String()
				changeperstring = strings.Replace(changeperstring, "%", "", 1)
				changeper, _ := strconv.ParseFloat(changeperstring, 64)
				changeper = changeper / 100
				return &model.Ticker{
					Buy:                v.Get("buy").Float(),
					Sell:               v.Get("sell").Float(),
					Last:               v.Get("last").Float(),
					Low:                v.Get("low").Float(),
					High:               v.Get("high").Float(),
					Vol:                v.Get("volume").Float(),
					PriceChangePercent: changeper * 100,
					//Raw:  string(body),
				}, nil
			}
		}
	}
	return &model.Ticker{}, nil
}

func (oc *OkCoin) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	since = since * 1000
	url := API + "kline.do" + "?symbol=" + quote + "_" + base + "&type=" + typ + "&size=" + strconv.Itoa(size) + "&since=" + strconv.Itoa(since)

	fmt.Println("url:", url)

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

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		record := model.Record{
			Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[2].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[4].Float(),
			Vol:   v.Array()[5].Float(),
			Ktime: v.Array()[0].Int() / 1000,
		}

		records = append(records, record)
		return true
	})

	return records, nil
}

func (oc *OkCoin) GetTrade(base string, quote string, typ string, since int, size int) ([]model.Trade, error) {
	since = since * 1000
	url := API + "trades.do" + "?symbol=" + quote + "_" + base

	fmt.Println("url:", url)

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
	//fmt.Println("Response body: %v", string(body))
	var records []model.Trade

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		trade := model.Trade{
			ID:        v.Get("tid").Int(),
			Price:     v.Get("price").Float(),
			Amount:    v.Get("amount").Float(),
			Type:      v.Get("type").String(),
			TradeTime: v.Get("date_ms").Int(),
		}

		records = append(records, trade)
		return true
	})

	return records, nil
}

type DeepDataResp struct {
	Ts int `json:"ts" bson:"ts"`
	Tick DeepTick `json:"tick" bson:"tick"`
	Data DeepTick
	Bids [][]float64 `json:"bids" bson:"bids"`
	Asks [][]float64 `json:"asks" bson:"asks"`

}
type DeepTick struct {
	Id string `json:"id" bson:"id"`
	Ts int `json:"ts" bson:"ts"`
	Bids [][]float64 `json:"bids" bson:"bids"`
	Asks [][]float64 `json:"asks" bson:"asks"`

}

func (oc *OkCoin) GetDeep(base string, quote string) (model.OrderBook, error) {
	//var records model.OrderBook
	records := model.OrderBook{}
	//url := "https://www.okcoin.com/v2/spot/markets/deep-deal?symbol=" + + strings.ToLower(quote) +"_" + strings.ToLower(base)
	url := API + "depth.do?symbol=" + strings.ToLower(quote) +"_" + strings.ToLower(base)
	resp, err := http.Get(url)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}
	fmt.Println("get url",url)
	//log.Debugf("Response body: %v", string(body))
	//fmt.Println("Response body: %v", string(body))

	//datajson := gjson.GetBytes(body, "data.tick.bids")
	transferdata:= DeepDataResp{}

	//json.Unmarshal(body, &records)
	err = json.Unmarshal(body, &transferdata)
	if err != nil {
		fmt.Println(string(body))
		fmt.Println("json unmarshal err:", err.Error())
		//return transferdata.Data, err
	}
	//fmt.Println("transferdata",transferdata)
	var askss []model.MarketOrder
	var asks []model.MarketOrder
	var bids []model.MarketOrder
	timestamp := time.Now().Unix()
	records.OrderTime = int64(timestamp)
	for _,v := range transferdata.Asks {
		marketorder := model.MarketOrder{
			Price:v[0],
			Amount:v[1],
		}
		askss = append(askss, marketorder)
		//fmt.Println(v[0])
	}
	askslen := len(askss) - 1
	for k,_ :=range askss {

		asks = append(asks,askss[askslen-k])
	}
	for _,v := range transferdata.Bids {
		marketorder := model.MarketOrder{
			Price:v[0],
			Amount:v[1],
		}
		bids = append(bids, marketorder)
		//fmt.Println(v[0])
	}
	records.Asks = asks
	records.Bids = bids
	fmt.Println(records)
	/*if datajson.Exists() {
		re := datajson.Array()
		for _, v := range re {

			datajson2 := gjson.Get(v.String(), "data")
			//	fmt.Println("datajson2:",datajson2)
			if datajson2.Exists() {

				re2 := datajson2.Array()
				for _, v2 := range re2 {
					//fmt.Println("datajson2:",v2.String())
					trade := model.Trade{
						ID :v2.Get("id").Int(),
						Price:v2.Get("price").Float(),
						Amount:v2.Get("amount").Float(),
						Type:v2.Get("direction").String(),
						TradeTime:v2.Get("ts").Int(),
					}
					//fmt.Println(v.Get("direction").String())
					records = append(records, trade)
				}
			}

		}
	}*/


	return records, nil
}