// Package huobi huobi rest api package
package exchange

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Akagi201/cryptotrader/model"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"fmt"
	//"strconv"
	//"time"
)

const (
	//API = "https://api.coinmarketcap.com/v1/ticker/bitcoin/?convert="
	//API = "http://api.fixer.io/latest?base=USD"
	API = "http://www.apilayer.net/api/live?access_key=05d21d06397253f53268622c38c63a40"
)

// Huobi API data
type Exchange struct {
	AccessKey string
	SecretKey string
}

// New create new Huobi API data
func New(accessKey string, secretKey string) *Exchange {
	return &Exchange{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// GetTicker 行情
func (hb *Exchange) GetTicker(base string, quote string,cachestring string) (*model.Ticker, error, string) {
	var body []byte
	if cachestring == ""{
		url := API

		log.Debugf("Request url: %v", url)
		fmt.Println("Request url:", url)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err,""
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err,""
		}
		fmt.Println("for url")
	} else {
		body = []byte(cachestring)
	}


	log.Debugf("Response body: %v", string(body))
	fmt.Println("Response body:", string(body))
	/*old
	sellpath := "rates." + strings.ToUpper(quote)
	sell := gjson.GetBytes(body, sellpath).Array()[0].Float()*/
	sellpath := "quotes.USD" + strings.ToUpper(quote)
	sell := gjson.GetBytes(body, sellpath).Array()[0].Float()

	return &model.Ticker{
		Buy:  1,
		Sell: sell,
	}, nil,string(body)
}