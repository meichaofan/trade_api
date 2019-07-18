package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type TradeData1 struct {
	ID        int64   `json:"ID" bson:"ID"`
	Type      string  `json:"type" bson:"type"`
	Price     float64 `json:"price" bson:"price"`
	Amount    float64 `json:"amount" bson:"amount"`
	TradeTime int64   `json:"trade_time" bson:"trade_time"`
	PriceUsd  float64 `json:"price_usd" bson:"price_usd"`
	PriceCny  float64 `json:"price_cny" bson:"price_cny"`
	Is_buy    int     `json:"is_buy" bson:"is_buy"`
	Is_sell   int     `json:"is_sell" bson:"is_sell"`
}

type exchangetable1 struct {
	Platform    string  `json:"platform" bson:"platform"`             //name
	Pair_base   string  `json:"pair_base" bson:"pair_base"`           //name
	Pair_symbol string  `json:"pair_symbol" bson:"pair_symbol"`       //name
	Rank        int     `json:"rank" bson:"rank"`                     //network_address
	Total_usd   float64 `json:"volume_usd_24h" bson:"volume_usd_24h"` //network_address
	Total_cny   float64 `json:"volume_cny_24h" bson:"volume_cny_24h"` //network_address
	//Market_Cap_Usd float32 `json:"market_cap_usd"` //network_address
}

type coinbase struct {
	Rank   int    `json:"rank" bson:"rank"`     //network_address
	Symbol string `json:"symbol" bson:"symbol"` //network_address
}
type excharge struct {
	Platform string `json:"platform" bson:"platform"`
}
type ExchargeSlice []exchangetable1

func (c ExchargeSlice) Len() int {
	return len(c)
}
func (c ExchargeSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ExchargeSlice) Less(i, j int) bool {
	return c[i].Rank < c[j].Rank
}

var (
	mgouse string
)

func getDBMain(dbname string, mgoselect string) *mgo.Database {
	//session, err := mgo.Dial("mongodb://longman:lmNX2017@120.78.12.247:9017/reputation_test")
	//session, err := mgo.Dial("s6:27017")
	var session *mgo.Session
	var err error
	if mgoselect == "dev" {
		fmt.Println("mgomode:", "dev")
		session, err = mgo.Dial("mongodb://test:123456@dds-wz9dad308815f1a433270.mongodb.rds.aliyuncs.com:3717/ganlan_bc")
	} else if mgoselect == "test" {
		fmt.Println("mgomode:", "test")
		session, err = mgo.Dial("mongodb://test:123456@dds-wz94b98362f14d4433270.mongodb.rds.aliyuncs.com:3717/ganlan_bc")
	} else {
		fmt.Println("mgomode:", "local")
		session, err = mgo.Dial("s6:27017")
	}

	//session, err := mgo.Dial("mongodb://test:123456@dds-wz9dad308815f1a433270.mongodb.rds.aliyuncs.com:3717/ganlan_bc")
	//session, err := mgo.Dial("s6:27017")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbname)
	return db
}
func getDBTrade(dbname string, mgoselect string) *mgo.Database {
	//session, err := mgo.Dial("mongodb://longman:lmNX2017@120.78.12.247:9017/reputation_test")
	var session *mgo.Session
	var err error
	if mgoselect == "dev" {
		fmt.Println("mgomode:", "dev")
		session, err = mgo.Dial("mongodb://test:123456@dds-wz9b02f17154e93433290.mongodb.rds.aliyuncs.com:3717/ganlan_bc")
	} else if mgoselect == "test" {
		fmt.Println("mgomode:", "test")
		session, err = mgo.Dial("mongodb://test:123456@dds-wz9b02f17154e93433290.mongodb.rds.aliyuncs.com:3717/ganlan_bc")
	} else {
		fmt.Println("mgomode:", "local")
		session, err = mgo.Dial("s6:27017")
	}

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(1 * time.Hour)
	db := session.DB(dbname)
	return db
}

func update_excharge_pair(quote, base, platform string, db *mgo.Database, totalcny, totalusd float64, updatetime int64) error {
	record := make(map[string]interface{}, 1)
	record["volume_usd_24h"] = totalusd
	record["volume_cny_24h"] = totalcny
	record["pair_symbol"] = quote
	record["pair_base"] = base
	record["platform"] = platform
	_, err := db.C("exchange_pair").Upsert(bson.M{"pair_base": base, "pair_symbol": quote, "platform": platform}, bson.M{"$set": record})
	if err != nil {
		fmt.Println("update_err3:", quote, base)
		return err
	}
	return nil
}

func update_excharge(platform string, db *mgo.Database, totalcny, totalusd float64, updatetime int64) error {
	record := make(map[string]interface{}, 1)
	record["volume_usd_24h"] = totalusd
	record["volume_cny_24h"] = totalcny
	record["platform"] = platform
	_, err := db.C("exchange").Upsert(bson.M{"platform": platform}, bson.M{"$set": record})
	if err != nil {
		fmt.Println("update exchange err:", err)
		return err
	}
	return nil
}
func sum_excharge_pair(platform string, db *mgo.Database) (float64, float64) {
	var sumcny float64
	var sumusd float64
	result := make([]exchangetable1, 1000)
	err := db.C("exchange_pair").Find(bson.M{"platform": platform}).All(&result)
	//err := db3.C("coin_trade").Find(bson.M{"quote": v.Pair_symbol, "base": v.Pair_base, "platform": v.Platform, "trade_time": bson.M{"$gte": usetime3}}).All(&result)
	if err != nil {
		fmt.Println(err)
		return sumcny, sumusd
	}
	for _, v := range result {
		sumcny = sumcny + v.Total_cny
		sumusd = sumusd + v.Total_usd
	}
	return sumcny, sumusd
}

func getrank(symbol string, db *mgo.Database) int {
	result := coinbase{}
	err := db.C("coin_base2").Find(bson.M{"symbol": symbol}).One(&result)
	if err != nil {
		fmt.Println(err)
		return 50000
	}
	return result.Rank
}

func sumbykline(quote, base, platform string) (float64, float64) {
	sendurl := "http://120.79.181.36:35433" + "/btc/vol" + "?quote=" + quote + "&base=" + base + "&platform=" + platform
	resp, err := http.Get(sendurl)
	if err != nil {
		return 0, 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0
	}
	//fmt.Println(body[100000000])

	volcny := gjson.GetBytes(body, "data.VolCNY").Float()
	volusd := gjson.GetBytes(body, "data.VolUSD").Float()
	//fmt.Println("platform:",typ)

	//log.Infof("Get exchangeticker: %+v", exchangeticker.Sell/exchangeticker.Buy)

	return volcny, volusd

}

func init() {
	fmt.Println("Preparing ...")
	flag.StringVar(&mgouse, "mgomode", "local", "local,test,dev")
	flag.Parse()
}

func main() {
	db := getDBMain("ganlan_bc", mgouse)
	db2 := getDBTrade("ganlan_bc", mgouse)
	nowtime := time.Now().Unix()

	yestime := (nowtime - 86400) * 1000
	result4 := make([]exchangetable1, 1000)
	var result3 ExchargeSlice

	err := db.C("exchange_pair_trade").Find(bson.M{"show_flag": 1, "$or": []bson.M{bson.M{"platform": "huobipro"}, bson.M{"platform": "binance"}, bson.M{"platform": "fcoin"}, bson.M{"platform": "zb"}, bson.M{"platform": "okex"}, bson.M{"platform": "gate-io"}}}).All(&result4)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range result4 {
		rank := getrank(v.Pair_symbol, db)
		result3 = append(result3, exchangetable1{
			Platform:    v.Platform,
			Pair_base:   v.Pair_base,
			Pair_symbol: v.Pair_symbol,
			Rank:        rank,
		})
	}
	sort.Sort(result3)

	for _, v := range result3 {
		result := make([]TradeData1, 100000)

		err := db2.C("coin_trade").Find(bson.M{"quote": v.Pair_symbol, "base": v.Pair_base, "platform": v.Platform, "trade_time": bson.M{"$gte": yestime}}).All(&result)

		if err != nil {
			fmt.Println(err)
		}
		var totalcny float64
		var totalusd float64

		totalcny = 0
		totalusd = 0
		if v.Pair_base == "USDT" && v.Platform == "zb" && v.Platform == "fcoin" {
			fmt.Println("platform", v.Platform)
			totalcny, totalusd = sumbykline(v.Pair_symbol, v.Pair_base, v.Platform)
			fmt.Println("totalcny", totalcny)
			fmt.Println("totalusd", totalusd)

		} else {
			for _, v := range result {
				single_cny := v.PriceCny * v.Amount
				single_usd := v.PriceUsd * v.Amount
				totalcny = totalcny + single_cny
				totalusd = totalusd + single_usd

			}
		}

		if totalusd > 0 {

			update_excharge_pair(v.Pair_symbol, v.Pair_base, v.Platform, db, totalcny, totalusd, nowtime)
		}

	}
	var result_excharge []exchangetable1
	err = db.C("exchange").Find(bson.M{"show_flag": 1}).All(&result_excharge)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range result_excharge {
		fmt.Println("start sum platform:", v.Platform)
		cny, usd := sum_excharge_pair(v.Platform, db)
		update_excharge(v.Platform, db, cny, usd, nowtime)
	}
}
