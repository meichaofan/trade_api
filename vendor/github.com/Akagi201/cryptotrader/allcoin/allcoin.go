// Package allcoin allcoin rest api package
package allcoin

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"fmt"
	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	//"time"
	//"bytes"
	"encoding/json"
	"net/url"
)

const (
	API          = "https://api.allcoin.com/api/v1/"
	Kline        = "https://api.allcoin.com/api/v1/kline"
	PRECENTAPI   = "https://www.allcoin.com/marketoverviews/MarketAreas/"
	NEWKLINE     = "https://www.allcoin.com/MarketPeriods/Kpro/?needTickers=1&id="
	ALLCOINPRICE = "https://www.allcoin.ca/Api_Market/getPriceList"
	KLINEV2 = "https://www.allcoin.com/market/nkline"
)

// Allcoin API data
type Allcoin struct {
	AccessKey string
	SecretKey string
}

// New create new Allcoin API data
func New(accessKey string, secretKey string) *Allcoin {
	return &Allcoin{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (ac *Allcoin) GetTicker(base string, quote string) (*model.Ticker, error) {

	var changeprice float64
	var buy float64
	var sell float64
	var last float64
	var low float64
	var high float64
	var vol float64
	url2 := ALLCOINPRICE

	log.Debugf("Request url: %v", url2)
	fmt.Println("url:", url2)
	resp, err := http.Get(url2)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	base = strings.Replace(base,".","",-1)
	fmt.Println(base)
	datajson := gjson.GetBytes(body, strings.ToLower(base))
	fmt.Println(datajson.String())
	if datajson.Exists() {
		re := datajson.Array()
		fmt.Println(len(re))
		for _, v := range re {
			fmt.Println(v.Get("coin_from").String())
			if v.Get("coin_from").String() == strings.ToLower(quote) {

				changeprice = v.Get("change_24h").Float()
				buy = v.Get("buy").Float()
				sell = v.Get("sell").Float()
				low = v.Get("min").Float()
				high = v.Get("max").Float()
				last = v.Get("current").Float()
				vol = v.Get("current").Float()
				break

			}

		}
	}
	return &model.Ticker{
		Buy:                buy,
		Sell:               sell,
		Last:               last,
		Low:                low,
		High:               high,
		Vol:                vol,
		PriceChangePercent: changeprice,
	}, nil
}

func (ac *Allcoin) GetRecords(base string, quote string, typ string, since int, size int) ([]model.Record, error) {
	quoteid := ""
	url2 := KLINEV2
	base = strings.Replace(base,".","",-1)
	log.Debugf("Request url: %v", url2)
	fmt.Println("url:", url2)
	needsymbol := strings.ToLower(quote) +"2"+ strings.ToLower(base)
	postValue := url.Values{}
	postValue.Set("needTickers", "1")
	postValue.Set("symbol", needsymbol)
	postValue.Set("type", typ)
	postValue.Set("size", strconv.Itoa(size))
	fmt.Println("symbol",needsymbol)
	fmt.Println("type", typ)
	fmt.Println("size", strconv.Itoa(size))
	fmt.Println(strings.NewReader(postValue.Encode()))
	resp, err := http.Post(url2, "application/json;charset=utf-8", strings.NewReader(postValue.Encode()))
	if err != nil {
		fmt.Println("post err:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("body:",string(body))
	datajson := gjson.GetBytes(body, "MarketAreas")
	if datajson.Exists() {
		re := datajson.Array()
		for _, v := range re {
			if v.Get("Secondary").String() == strings.ToUpper(base) {
				ree := v.Get("Markets").Array()
				for _, vv := range ree {
					if vv.Get("Primary").String() == strings.ToUpper(quote) {
						quoteid = vv.Get("Id").String()
					}
				}
			}

		}
	}
	url := NEWKLINE + quoteid + "&type=" + typ + "&size=" + strconv.Itoa(size)
	fmt.Println("url:", url)
	resp2, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return nil, err
	}

	log.Debugf("Response body: %v", string(body2))
	var records []model.Record
	datajson2 := gjson.GetBytes(body2, "datas.data")
	//fmt.Println(datajson2.String())

	if datajson2.Exists() {
		ret := string(datajson2.String())

		var res []interface{}
		err = json.Unmarshal([]byte(ret), &res)
		if err != nil {
			return nil, err
		}
		//var sentence = (res[0]).([]interface{})
		//fmt.Println(sentence[0].(float64)/1000)
		//var translated = make([]string,0)
		fmt.Println(len(res))
		for i := 0; i < len(res); i++ {
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
				Ktime: ktime / 1000,
			}

			records = append(records, record)
			//translated = append(translated,n)
		}


	}
	return records, nil

}
