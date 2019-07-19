package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
	"trade_api/src/main/web/cli/exchange"
	"trade_api/src/main/web/cli/exchange/Bibox"
	"truxing/commons/log"
)

var (
//env string
)

func init() {
	//fmt.Println("Preparing ...")
	//flag.StringVar(&env, "e", "local", "local,test,dev")
	//flag.Parse()
}

/**
 */
func updatePair(exchange exchange.Exchange) {
	var s *mgo.Session
	var c *mgo.Collection
	var pairTickers []*data.ExchangeTicker
	//connect mongodb
	s, c = common.Connect("platform_pair", exchange.Name(), "local")
	defer s.Close()
	pairTickers = exchange.PairHandler()
	//插入或更新
	for _, pair := range pairTickers {
		_, err := c.Upsert(bson.M{"symbol": pair.Symbol}, bson.M{"$set": pair})
		if err != nil {
			log.Debugf("platform:%s symbol %s pair update failed", exchange.Name(), pair.Symbol)
		}
	}
}

func updateAmount(exchange exchange.Exchange) {
	var s *mgo.Session
	var c *mgo.Collection
	var trades []*data.TradeData
	s, c = common.Connect("platform_amount", exchange.Name(), "local")
	defer s.Close()
	trades = exchange.AmountHandler()
	fmt.Printf("the length is %d", len(trades))
	fmt.Println()
	for _, trade := range trades {
		_, err := c.Upsert(bson.M{"symbol": trade.Symbol}, bson.M{"$set": trade})
		if err != nil {
			log.Debugf("platform:%s symbol %s trade update failed", exchange.Name(), trade.Symbol)
		}
	}
}

func main() {
	ex := Bibox.Bibox{}
	updatePair(ex)
	updateAmount(ex)
}
