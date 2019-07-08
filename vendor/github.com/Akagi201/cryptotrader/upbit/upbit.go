// Package bithumb huobi rest api package
package upbit

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
	API = "https://crix-api-endpoint.upbit.com/v1/crix/candles/lines?code=CRIX.UPBIT.KRW-"
	KILNE = "https://crix-api-endpoint.upbit.com/v1/crix/candles/"
)

// Huobi API data
type Upbit struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Upbit {
	return &Upbit{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (upbit *Upbit) GetTicker(base string, quote string) (*model.Ticker, error) {
	url := API + strings.ToUpper(quote)

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
	//buy := gjson.GetBytes(body, "askBiD").Float()
	//sell := gjson.GetBytes(body, "ask").Float()
	//buy := gjson.GetBytes(body, "result.#.Ask").Array()[0].Float()
	last := gjson.GetBytes(body, "candles.#.tradePrice").Array()[0].Float()
	low := gjson.GetBytes(body, "candles.#.lowPrice").Array()[0].Float()
	high := gjson.GetBytes(body, "candles.#.highPrice").Array()[0].Float()
	vol := gjson.GetBytes(body, "candles.#.candleAccTradeVolume").Array()[0].Float()

	return &model.Ticker{
		//Buy:  buy,
		//Sell: sell,
		Last: last,
		Low:  low,
		High: high,
		Vol:  vol,
	}, nil
}


func (upbit *Upbit) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	url := KILNE + typ + "?code=CRIX.UPBIT."+strings.ToUpper(quote) + "-" + strings.ToUpper(base) + "&count=" + strconv.Itoa(size) +"&ciqrandom=" + strconv.Itoa(since)

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
	/*datajson := gjson.GetBytes(body, "candles")
	if datajson.Exists() {
		re := datajson.Array()
		for _,v:=  range re {
			record := model.Record{

				//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
				Open:  v.Get("openingPrice").Float(),
				High:  v.Get("highPrice").Float(),
				Low:   v.Get("lowPrice").Float(),
				Close: v.Get("tradePrice").Float(),
				Vol:   v.Get("candleAccTradeVolume").Float(),
			}

			records = append(records, record)
		}
	}*/

	gjson.ParseBytes(body).ForEach(func(k, v gjson.Result) bool {
		toBeCharge := v.Get("candleDateTime").String()
		timeLayout := "2006-01-02T15:04:05+00:00"                             //转化所需模板
		loc, _ := time.LoadLocation("Asia/Chongqing")                            //重要：获取时区
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
		sr := theTime.Unix()  +28800                                           //转化为时间戳 类型是int64
		record := model.Record{
			//Time:  time.Unix(v.Array()[0].Int()/1000, 0),
			Open:  v.Get("openingPrice").Float(),
			High:  v.Get("highPrice").Float(),
			Low:   v.Get("lowPrice").Float(),
			Close: v.Get("tradePrice").Float(),
			Vol:   v.Get("candleAccTradeVolume").Float(),
			Ktime:sr,
		}

		records = append(records, record)
		return true
	})

	return records, nil
}