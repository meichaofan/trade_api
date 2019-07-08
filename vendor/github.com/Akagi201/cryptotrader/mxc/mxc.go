// Package huobi huobi rest api package
package mxc

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"fmt"
	"strconv"

	"encoding/json"

	"sort"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	//"time"
)

const (
	API = "https://www.mxc.com/"
)

// Huobi API data
type Mxc struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Mxc {
	return &Mxc{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

type DeepDataResp struct {
	//Ts   int      `json:"ts" bson:"ts"`
	Tick DeepTick `json:"data" bson:"data"`
}
type DeepTick struct {
	Id   string    `json:"id" bson:"id"`
	Ts   int       `json:"ts" bson:"ts"`
	Bids []float64 `json:"bids" bson:"bids"`
	Asks []float64 `json:"asks" bson:"asks"`
}
type AsksSlice []model.MarketOrder
type BidsSlice []model.MarketOrder

func (c AsksSlice) Len() int {
	return len(c)
}
func (c AsksSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c AsksSlice) Less(i, j int) bool {
	return c[i].Price < c[j].Price
}

func (c BidsSlice) Len() int {
	return len(c)
}
func (c BidsSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c BidsSlice) Less(i, j int) bool {
	return c[i].Price > c[j].Price
}

// GetTicker 行情
func (hb *Mxc) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "market/ticker/" + strings.ToLower(quote) + strings.ToLower(base)

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
	fmt.Println("Response body:", string(body))
	buy := gjson.GetBytes(body, "data.ticker").Array()[2].Float()
	sell := gjson.GetBytes(body, "data.ticker").Array()[4].Float()
	last := gjson.GetBytes(body, "data.ticker").Array()[0].Float()
	low := gjson.GetBytes(body, "data.ticker").Array()[8].Float()
	high := gjson.GetBytes(body, "data.ticker").Array()[7].Float()
	vol := gjson.GetBytes(body, "data.ticker").Array()[1].Float()
	open := gjson.GetBytes(body, "data.ticker").Array()[6].Float()
	precent24 := (last - open) / open
	//fmt.Println("ticker open:", open)
	return &model.Ticker{
		Buy:                buy,
		Sell:               sell,
		Last:               last,
		Low:                low,
		High:               high,
		Vol:                vol,
		PriceChangePercent: precent24 * 100,
	}, nil
}

func (z *Mxc) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + "market/candles"

	if len(typ) != 0 {
		url += "/" + typ + "/" + strings.ToLower(quote) + strings.ToLower(base)
	}

	if size != 0 {
		if size > 1000 {
			size = 1000
		}
		url += "?limit=" + strconv.Itoa(size)
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

	//log.Debugf("Response body: %v", string(body))

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
				Vol:   v.Get("count").Float(),
				Ktime: v.Get("id").Int(),
				Time:  cast.ToTime(v.Get("id").Int()),
				//PriceChangePercent:precent24 * 100,
			}

			recordscache = append(recordscache, record)
		}
	}
	lenrecord := len(recordscache) - 1
	for i := lenrecord; i >= 0; i-- {
		records = append(records, recordscache[i])
	}

	return records, nil
}

func (c *Mxc) GetTrades(quote string, base string) ([]model.Trade, error) {
	var records []model.Trade
	url := API + "market/trades/" + strings.ToLower(quote) + strings.ToLower(base) + "?limit=500"
	fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}
	//fmt.Println("url", url)
	//json.Valid(body)
	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		re := datajson.Array()
		for _, v2 := range re {
			/*open := v.Get("open").Float()
			close := v.Get("close").Float()
			precent24 := (close - open) / open*/

			//fmt.Println("datajson2:",v2.String())
			trade := model.Trade{
				ID:        v2.Get("id").Int(),
				Price:     v2.Get("price").Float(),
				Amount:    v2.Get("amount").Float(),
				Type:      v2.Get("side").String(),
				TradeTime: v2.Get("ts").Int(),
			}
			//fmt.Println(v.Get("direction").String())
			records = append(records, trade)

		}
	}
	//fmt.Println(records)
	return records, nil
}

func (c *Mxc) GetDeep(quote string, base string) (model.OrderBook, error) {
	records := model.OrderBook{}
	url := API + "market/depth/L20/" + strings.ToLower(quote) + strings.ToLower(base)
	fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}
	//fmt.Println("url", url)
	//json.Valid(body)
	log.Debugf("Response body: %v", string(body))

	transferdata := DeepDataResp{}

	//json.Unmarshal(body, &records)
	err = json.Unmarshal(body, &transferdata)
	if err != nil {
		fmt.Println(string(body))
		fmt.Println("json unmarshal err:", err.Error())
		//return transferdata.Data, err
	}
	//fmt.Println("transferdata",transferdata)
	var asks AsksSlice
	var bids BidsSlice
	records.OrderTime = time.Now().Unix()
	for k, v := range transferdata.Tick.Asks {
		if k%2 == 0 {
			marketorder := model.MarketOrder{
				Price:  v,
				Amount: transferdata.Tick.Asks[k+1],
			}
			asks = append(asks, marketorder)
		}

		//fmt.Println(v[0])
	}
	for k, v := range transferdata.Tick.Bids {

		if k%2 == 0 {
			marketorder := model.MarketOrder{
				Price:  v,
				Amount: transferdata.Tick.Bids[k+1],
			}
			bids = append(bids, marketorder)
		}

		//fmt.Println(v[0])
	}
	sort.Sort(asks)
	sort.Sort(bids)
	records.Asks = asks
	records.Bids = bids

	//fmt.Println(records)
	return records, nil
}

func (c *Mxc) GetMarkets() ([]model.MarketPairInfo, error) {
	var records []model.MarketPairInfo
	url := "https://www.mxc.com/open/api/v1/data/markets"
	fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}
	//fmt.Println("url", url)
	//json.Valid(body)
	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		datajson := gjson.GetBytes(body, "data")
		fmt.Println("data", datajson.String())
		data := strings.Replace(datajson.String(), "[", "", -1)
		data = strings.Replace(data, "]", "", -1)
		data = strings.Replace(data, "\"", "", -1)
		fmt.Println("datarel", data)
		dataarr := strings.Split(data, ",")
		for _, v2 := range dataarr {
			quote_base := strings.Split(v2, "_")
			fmt.Println("quote_base", quote_base)
			trade := model.MarketPairInfo{
				Quote: strings.ToUpper(quote_base[0]),
				Base:  strings.ToUpper(quote_base[1]),
			}
			records = append(records, trade)

		}
	}
	return records, nil
}
