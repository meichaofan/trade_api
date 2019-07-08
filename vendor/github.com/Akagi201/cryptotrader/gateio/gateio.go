// Package gateio gate.io rest api package
package gateio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"sort"
	"strconv"
	"time"

	"github.com/Akagi201/cryptotrader/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	RestHost = "data.gate.io"
	ApiVer   = "api2/1"
	API      = "https://data.gateio.io/api2/1/"
)

// Client OkEx client
type Client struct {
	URL        url.URL
	HTTPClient *http.Client
	AccessKey  string
	SecretKey  string
}

type DeepDataResp struct {
	Ts   int      `json:"ts" bson:"ts"`
	Tick DeepTick `json:"data" bson:"data"`
}
type DeepTick struct {
	Id   string     `json:"id" bson:"id"`
	Ts   int        `json:"ts" bson:"ts"`
	Bids [][]string `json:"bids" bson:"bids"`
	Asks [][]string `json:"asks" bson:"asks"`
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

// New creates a new OkEx Client
func New(accessKey string, secretKey string) *Client {
	u := url.URL{
		Scheme: "http",
		Host:   RestHost,
		Path:   "/",
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

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

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

// GetPairs 返回所有系统支持的交易对, for http://data.gate.io/api2/1/pairs
func (c *Client) GetPairs(ctx context.Context) ([]string, error) {
	req, err := c.newRequest(ctx, "GET", "pairs", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))
	var pairs []string
	for _, v := range gjson.ParseBytes(body).Array() {
		pairs = append(pairs, v.String())
	}
	return pairs, nil
}

// GetMarketInfo 交易市场订单参数, 返回所有系统支持的交易市场的参数信息，包括交易费，最小下单量，价格精度等。for http://data.gate.io/api2/1/marketinfo
func (c *Client) GetMarketInfo(ctx context.Context) ([]model.MarketInfo, error) {
	req, err := c.newRequest(ctx, "GET", "marketinfo", nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	var marketInfos []model.MarketInfo
	var marketInfo model.MarketInfo

	for _, v := range gjson.GetBytes(body, "pairs").Array() {
		v.ForEach(func(key, value gjson.Result) bool {
			marketInfo.Symbol = key.String()
			marketInfo.DecimalPlaces = value.Get("decimal_places").Int()
			marketInfo.MinAmount = value.Get("min_amount").Float()
			marketInfo.Fee = value.Get("fee").Float()
			marketInfos = append(marketInfos, marketInfo)
			return true // keep iterating
		})
	}

	return marketInfos, nil
}

// GetTicker 返回最新, 最高, 最低 交易行情和交易量, 每 10 秒钟更新, for http://data.gate.io/api2/1/ticker/[quote]_[base]
func (c *Client) GetTicker(ctx context.Context, quote string, base string) (*model.Ticker, error) {
	req, err := c.newRequest(ctx, "GET", "ticker/"+strings.ToUpper(quote)+"_"+strings.ToUpper(base), nil, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.getResponse(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body))

	buy := gjson.GetBytes(body, "highestBid").Float()
	sell := gjson.GetBytes(body, "lowestAsk").Float()
	last := gjson.GetBytes(body, "last").Float()
	low := gjson.GetBytes(body, "low24hr").Float()
	high := gjson.GetBytes(body, "high24hr").Float()
	vol := gjson.GetBytes(body, "baseVolume").Float()
	percentChange := gjson.GetBytes(body, "percentChange").Float()
	percentChange = percentChange / 100.00
	return &model.Ticker{
		Buy:                buy,
		Sell:               sell,
		Last:               last,
		Low:                low,
		High:               high,
		Vol:                vol,
		Raw:                string(body),
		PriceChangePercent: percentChange,
	}, nil
}

func (c *Client) GetDeep(quote string, base string) (model.OrderBook, error) {
	records := model.OrderBook{}
	url := API + "orderBook/" + strings.ToLower(quote) + "_" + strings.ToLower(base)
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

	transferdata := DeepTick{}

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
	for _, v := range transferdata.Asks {
		v0, _ := strconv.ParseFloat(v[0], 64)
		v1, _ := strconv.ParseFloat(v[1], 64)
		marketorder := model.MarketOrder{
			Price:  v0,
			Amount: v1,
		}
		asks = append(asks, marketorder)
		//fmt.Println(v[0])
	}
	for _, v := range transferdata.Bids {
		v0, _ := strconv.ParseFloat(v[0], 64)
		v1, _ := strconv.ParseFloat(v[1], 64)
		marketorder := model.MarketOrder{
			Price:  v0,
			Amount: v1,
		}
		bids = append(bids, marketorder)
		//fmt.Println(v[0])
	}
	sort.Sort(asks)
	sort.Sort(bids)
	records.Asks = asks
	records.Bids = bids

	//fmt.Println(records)
	return records, nil
}

func (c *Client) GetTrades(quote string, base string) ([]model.Trade, error) {
	var records []model.Trade
	url := API + "tradeHistory/" + strings.ToLower(quote) + "_" + strings.ToLower(base)
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
				ID:        v2.Get("tradeID").Int(),
				Price:     v2.Get("rate").Float(),
				Amount:    v2.Get("amount").Float(),
				Type:      v2.Get("type").String(),
				TradeTime: v2.Get("timestamp").Int() * 1000,
			}
			//fmt.Println(v.Get("direction").String())
			records = append(records, trade)

		}
	}
	//fmt.Println(records)
	return records, nil
}

func (c *Client) GetMarkets() ([]model.MarketPairInfo, error) {

	var records []model.MarketPairInfo
	url := "https://data.gateio.io/api2/1/pairs"
	resp, err := http.Get(url)
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}

	//log.Debugf("Response body: %v", string(body))
	fmt.Println(string(body))

	//去掉【】, "
	result1 := strings.Replace(string(body), "[", "", -1)
	result2 := strings.Replace(result1, "]", "", -1)
	result3 := strings.Replace(result2, "\"", "", -1)
	ret := strings.Split(result3, ",")
	for _, k := range ret {
		marketjson := strings.Split(k, "_")
		records = append(records, model.MarketPairInfo{
			Quote: marketjson[0],
			Base:  marketjson[1],
		})

	}
	return records, nil
}
