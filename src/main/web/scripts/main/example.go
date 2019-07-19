package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

/**
设计思路：交易平台交易额
1.新建数据库 platform_amount , 然后再创建 bibox 、coinall等集合
2.集合下实时插入数据......

交易平台交易对行情
1.新建数据库 platform_pair , 然后在创建 bibox、coinall等集合
2.集合中插入交易对
*/

var (
	mgouse string
	//对接平台
	platforms = []string{"bibox"}
)


func init() {
	fmt.Println("Preparing ...")
	flag.StringVar(&mgouse, "mgomode", "local", "local,test,dev")
}

//连接mongodb
func connect(dbName, cName string) (*mgo.Session, *mgo.Collection) {
	var session *mgo.Session
	var err error
	if mgouse == "dev" {
		fmt.Println("mgomode:", "dev")
		session, err = mgo.Dial("")
	} else if mgouse == "test" {
		fmt.Println("mgomode:", "test")
		session, err = mgo.Dial("")
	} else {
		fmt.Println("mgomode:", "local")
		session, err = mgo.Dial("192.168.244.128:27017")
	}
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(1 * time.Hour)
	collection := session.DB(dbName).C(cName)
	return session, collection
}

func httpGet(url string) []byte {
	log.Printf("url: %s", url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	return content
}

/**
写一个函数，实时查询汇率，最多支持两级
*/
func GetExchangeRate(quote, base, platform string) float64 {
	var exchangeRate float64
	switch platform {
	case "bibox":
		//假设是bibox平台
		url := "https://api.bibox365.com/v1/mdata?cmd=ticker&pair=" + strings.ToUpper(quote) + "_" + strings.ToUpper(base)
		content := httpGet(url)
		if strings.ToUpper(base) != "USDT" {
			nUrl := "https://api.bibox365.com/v1/mdata?cmd=ticker&pair=" + strings.ToUpper(base) + "_USDT"
			nContent := httpGet(nUrl)
			exchangeRate = gjson.ParseBytes(content).Get("result.last").Float() * gjson.ParseBytes(nContent).Get("result.last").Float()
		}
		exchangeRate = gjson.ParseBytes(content).Get("result.last").Float()
	default:
		exchangeRate = 0
	}
	return exchangeRate
}

/*func insertTradeDataToCollection(data TradeData, platform string) {
	//s, c := connect("platforms", platform)
	defer s.Close()
	//c.Upsert()
}
*/
func InsertPlatformPair() {

}