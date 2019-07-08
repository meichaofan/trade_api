// Package huobi huobi rest api package
package huobi

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"strconv"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"encoding/json"
	//"time"
)

const (
	API = "https://api.huobi.pro/"
)

// Huobi API data
type Huobi struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Huobi {
	return &Huobi{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

type DeepDataResp struct {
	Ts int `json:"ts" bson:"ts"`
	Tick DeepTick `json:"tick" bson:"tick"`

}
type DeepTick struct {
	Id string `json:"id" bson:"id"`
	Ts int `json:"ts" bson:"ts"`
	Bids [][]float64 `json:"bids" bson:"bids"`
	Asks [][]float64 `json:"asks" bson:"asks"`

}


/*type DeepBids struct {
	Data [2]float64
}*/
// GetTicker 行情
func (hb *Huobi) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "market/detail/merged?symbol=" + strings.ToLower(quote) + strings.ToLower(base)

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
	buy := gjson.GetBytes(body, "tick.bid").Array()[0].Float()
	sell := gjson.GetBytes(body, "tick.ask").Array()[0].Float()
	last := gjson.GetBytes(body, "tick.close").Float()
	low := gjson.GetBytes(body, "tick.low").Float()
	high := gjson.GetBytes(body, "tick.high").Float()
	vol := gjson.GetBytes(body, "tick.amount").Float()
	open := gjson.GetBytes(body, "tick.open").Float()
	precent24 := (last - open) / open
	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		PriceChangePercent:precent24 * 100,
	}, nil
}

func (z *Huobi) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + "market/history/kline?symbol=" + strings.ToLower(quote) + strings.ToLower(base)

	if len(typ) != 0 {
		url += "&period=" + typ
	}

	if since != 0 {
		url += "&since=" + strconv.Itoa(since)
	}

	if size != 0 {
		if size > 1000 {
			size = 1000
		}
		url += "&size=" + strconv.Itoa(size)
	}
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
	var recordscache []model.Record
	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		re := datajson.Array()
		for _, v := range re {
			/*open := v.Get("open").Float()
			close := v.Get("close").Float()
			precent24 := (close - open) / open*/
			record := model.Record{

				//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
				Open:  v.Get("open").Float(),
				High:  v.Get("high").Float(),
				Low:   v.Get("low").Float(),
				Close: v.Get("close").Float(),
				Vol:   v.Get("vol").Float(),
				Ktime: v.Get("id").Int(),
				//PriceChangePercent:precent24 * 100,
			}

			recordscache = append(recordscache, record)
		}
	}
	lenrecord := len(recordscache) - 1
	for i := lenrecord; i >= 0; i-- {
		records = append(records, recordscache[i])
	}
	/*gjson.ParseBytes([]byte(datajson.String())).ForEach(func(k, v gjson.Result) bool {
		fmt.Println(len(v.Array()))
		record := model.Record{

			//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Array()[1].Float(),
			High:  v.Array()[4].Float(),
			Low:   v.Array()[3].Float(),
			Close: v.Array()[2].Float(),
			Vol:   v.Array()[5].Float(),
		}

		records = append(records, record)
		return true
	})*/

	return records, nil
}

func (z *Huobi) GetTrades(base string, quote string) ([]model.Trade, error) {
	url := API + "market/history/trade?symbol=" + strings.ToLower(quote) + strings.ToLower(base) +"&size=300"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("get url",url)
	//log.Debugf("Response body: %v", string(body))
	//fmt.Println("Response body: %v", string(body))
	var records []model.Trade
	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		re := datajson.Array()
		for _, v := range re {
			/*open := v.Get("open").Float()
			close := v.Get("close").Float()
			precent24 := (close - open) / open*/
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
	}
	if len(records) == 0 {
		fmt.Println("Response body: %v", string(body))
	}

	return records, nil
}

func (z *Huobi) GetDeep(base string, quote string) (model.OrderBook, error) {
	//var records model.OrderBook
	records := model.OrderBook{}
	url := API + "market/depth?symbol=" + strings.ToLower(quote) + strings.ToLower(base) +"&type=step0"
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
	var asks []model.MarketOrder
	var bids []model.MarketOrder
	records.OrderTime = int64(transferdata.Ts)
	for _,v := range transferdata.Tick.Asks {
		marketorder := model.MarketOrder{
			Price:v[0],
			Amount:v[1],
		}
		asks = append(asks, marketorder)
		//fmt.Println(v[0])
	}
	for _,v := range transferdata.Tick.Bids {
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