// Package okex okex rest api package
package okex

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/Akagi201/utilgo/signs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"encoding/json"
)

const (
	RestHost = "www.okex.com"
	ApiVer   = "v1"
)

// Client OkEx client
type Client struct {
	URL        url.URL
	HTTPClient *http.Client
	AccessKey  string
	SecretKey  string
}

// New creates a new OkEx Client
func New(accessKey string, secretKey string) *Client {
	u := url.URL{
		Scheme: "https",
		Host:   RestHost,
		Path:   "/api/",
	}

	c := Client{
		URL:        u,
		HTTPClient: &http.Client{},
		AccessKey:  accessKey,
		SecretKey:  secretKey,
	}

	return &c
}

func (c *Client) newRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader) (*http.Request, error) {
	u := c.URL
	u.Path = path.Join(c.URL.Path, ApiVer, spath)
	u.RawQuery = values.Encode()
	log.Debugf("Request URL: %#v", u.String())
	fmt.Println("url:",u.String())
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) SignParams(values *url.Values) error {
	values.Set("api_key", c.AccessKey)
	payload := values.Encode()
	payload = payload + "&secret_key=" + c.SecretKey
	//payload, _ = url.QueryUnescape(payload)

	sign, err := signs.GetMd5Sign(payload)
	if err != nil {
		return err
	}

	values.Add("sign", strings.ToUpper(sign))
	return nil
}

func (c *Client) newPrivateRequest(ctx context.Context, method string, spath string, values url.Values, body io.Reader) (*http.Request, error) {
	if values == nil {
		values = url.Values{}
	}

	c.SignParams(&values)
	//req.URL.RawQuery = values.Encode()

	log.Debugf("Encode values: %v", values.Encode())

	req, err := c.newRequest(ctx, method, spath, values, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) getResponse(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		log.Errorf("body: %v", string(body))
		return nil, errors.New(fmt.Sprintf("status code: %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetTicker 获取OKEx最新币币行情数据, for Get /api/v1/ticker.do
func (c *Client) GetTicker(ctx context.Context, quote string, base string) (*model.Ticker, error) {
	/*v := url.Values{}
	v.Set("symbol", strings.ToLower(quote)+"_"+strings.ToLower(base))

	req, err := c.newRequest(ctx, "GET", "ticker.do", v, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
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
		Raw:  string(body),
	}, nil*/
	url := "https://www.okex.com/v2/spot/markets/tickers"
	symbol := strings.ToLower(quote)+"_"+strings.ToLower(base)
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

	datajson := gjson.GetBytes(body, "data")
	if datajson.Exists() {
		re := datajson.Array()
		for _,v:=  range re {
			if v.Get("symbol").String() == symbol {
				changeperstring := v.Get("changePercentage").String()
				changeperstring = strings.Replace(changeperstring,"%","",1)
				changeper, _ := strconv.ParseFloat(changeperstring, 64)
				changeper = changeper /100
				return &model.Ticker{
					Buy:  v.Get("buy").Float(),
					Sell: v.Get("sell").Float(),
					Last: v.Get("last").Float(),
					Low:  v.Get("low").Float(),
					High: v.Get("high").Float(),
					Vol:  v.Get("volume").Float(),
					PriceChangePercent:changeper * 100,
					//Raw:  string(body),
				}, nil
			}
		}
	}
	return &model.Ticker{},nil
}

// GetDepth 获取币币市场深度, for Get /api/v1/depth
func (c *Client) GetDepth(ctx context.Context, quote string, base string) (*model.OrderBook, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToLower(quote)+"_"+strings.ToLower(base))

	req, err := c.newRequest(ctx, "GET", "depth.do", v, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var order model.MarketOrder
	var orderBook model.OrderBook
	gjson.GetBytes(body, "bids").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Array()[0].Float()
		order.Amount = value.Array()[1].Float()
		orderBook.Bids = append(orderBook.Bids, order)
		return true // keep iterating
	})

	gjson.GetBytes(body, "asks").ForEach(func(key, value gjson.Result) bool {
		order.Price = value.Array()[0].Float()
		order.Amount = value.Array()[1].Float()
		orderBook.Asks = append(orderBook.Asks, order)
		return true // keep iterating
	})

	return &orderBook, nil
}

// GetTrades 获取币币交易信息, for GET https://www.okex.com/api/v1/trades.do
func (c *Client) GetTrades(ctx context.Context, quote string, base string,since int) ([]model.Trade, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToLower(quote)+"_"+strings.ToLower(base))
	//v.Set("since", strconv.Itoa(since))
	req, err := c.newRequest(ctx, "GET", "trades.do", v, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var trade model.Trade
	var trades []model.Trade
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		trade.ID = value.Get("tid").Int()
		trade.Price = value.Get("price").Float()
		trade.Amount = value.Get("amount").Float()
		trade.Type = value.Get("type").String()
		trade.Time = cast.ToTime(cast.ToInt64(value.Get("date").String()))
		trade.TradeTime = value.Get("date_ms").Int()
		trades = append(trades, trade)
		return true // keep iterating
	})

	return trades, nil
}

// GetRecords 获取币币K线数据, for GET https://www.okex.com/api/v1/kline.do
func (c *Client) GetRecords(ctx context.Context, quote string, base string, interval string, since int64, limit int64) ([]model.Record, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToLower(quote)+"_"+strings.ToLower(base))
	v.Set("type", interval)

	if since != 0 {
		v.Set("startTime", strconv.FormatInt(since, 10))
	}

	if limit != 0 {
		v.Set("limit", strconv.FormatInt(limit, 10))
	}

	req, err := c.newRequest(ctx, "GET", "kline.do", v, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	//log.Debugf("Response body: %v", string(body))

	var record model.Record
	var records []model.Record

	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		//fmt.Println(len(value.Array()))
		record.Open = value.Array()[1].Float()
		record.High = value.Array()[2].Float()
		record.Low = value.Array()[3].Float()
		record.Close = value.Array()[4].Float()
		record.Vol = value.Array()[5].Float()
		record.Time = cast.ToTime(value.Array()[0].Int() / 1000)
		record.Raw = value.String()
		record.Ktime = value.Array()[0].Int() / 1000
		records = append(records, record)
		return true // keep iterating
	})
	return records, nil
}

// GetAccount 获取用户信息, for POST https://www.okex.com/api/v1/userinfo.do
func (c *Client) GetAccount(ctx context.Context) ([]model.Balance, error) {
	req, err := c.newPrivateRequest(ctx, "POST", "userinfo.do", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Infof("Response body: %v", string(body))

	//var balance model.Balance
	//var balances []model.Balance
	//gjson.GetBytes(body, "data").ForEach(func(key, value gjson.Result) bool {
	//	balance.Currency = value.Get("account_type").String()
	//	balance.Free = cast.ToFloat64(value.Get("active_balance").String())
	//	balance.Frozen = cast.ToFloat64(value.Get("frozen_balance").String())
	//
	//	balances = append(balances, balance)
	//	return true // keep iterating
	//})

	//return balances, nil
	return nil, nil
}

type DeepDataResp struct {
	Ts int `json:"ts" bson:"ts"`
	Tick DeepTick `json:"tick" bson:"tick"`
	Bids [][]float64 `json:"bids" bson:"bids"`
	Asks [][]float64 `json:"asks" bson:"asks"`

}
type DeepTick struct {
	Id string `json:"id" bson:"id"`
	Ts int `json:"ts" bson:"ts"`
	Bids [][]float64 `json:"bids" bson:"bids"`
	Asks [][]float64 `json:"asks" bson:"asks"`

}

func (c *Client) GetDeep(base string, quote string) (model.OrderBook, error) {
	//var records model.OrderBook
	apiurl := "https://www.okex.com/api/v1/depth.do?symbol="
	records := model.OrderBook{}
	url := apiurl  + strings.ToLower(quote) +"_" + strings.ToLower(base)
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
	records.OrderTime = int64(transferdata.Ts)
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