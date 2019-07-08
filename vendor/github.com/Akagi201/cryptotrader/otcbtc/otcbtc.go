// Package huobi huobi rest api package
package otcbtc

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"strconv"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/spf13/cast"
	//"time"
)

const (
	API = "https://bb.otcbtc.com/"
)

// Huobi API data
type Otcbtc struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Otcbtc {
	return &Otcbtc{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (hb *Otcbtc) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + "/api/v2/tickers/" + strings.ToLower(quote) + strings.ToLower(base)

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
	buy := gjson.GetBytes(body, "ticker.buy").Array()[0].Float()
	sell := gjson.GetBytes(body, "ticker.sell").Array()[0].Float()
	last := gjson.GetBytes(body, "ticker.last").Float()
	low := gjson.GetBytes(body, "ticker.low").Float()
	high := gjson.GetBytes(body, "ticker.high").Float()
	vol := gjson.GetBytes(body, "ticker.vol").Float()
	/*records,err := hb.GetRecords(base,quote,"1440",0,1)
	var precent24 float64
	if err == nil {
		open := records[0].Open
		clo := records[0].Close
		precent24 = (clo - open) / open
	}*/

	return &model.Ticker{
		Buy:  buy,
		Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
		//PriceChangePercent:precent24 * 100,
	}, nil
}

func (z *Otcbtc) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := API + "/api/v2/klines?market=" + strings.ToLower(quote) + strings.ToLower(base)

	if len(typ) != 0 {
		url += "&period=" + typ
	}


	if size != 0 {
		url += "&limit=" + strconv.Itoa(size)
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
	var record model.Record
	gjson.ParseBytes(body).ForEach(func(key, value gjson.Result) bool {
		record.Open = cast.ToFloat64(value.Array()[1].String())
		record.High = cast.ToFloat64(value.Array()[2].String())
		record.Low = cast.ToFloat64(value.Array()[3].String())
		record.Close = cast.ToFloat64(value.Array()[4].String())
		record.Vol = cast.ToFloat64(value.Array()[5].String())
		record.Time = cast.ToTime(value.Array()[0].Int() )
		record.Ktime = value.Array()[0].Int()
		records = append(records, record)
		return true // keep iterating
	})
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

