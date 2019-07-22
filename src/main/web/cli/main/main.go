package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
	"trade_api/src/main/web/cli/exchange"
	"trade_api/src/main/web/cli/exchange/Bibox"
	"trade_api/src/main/web/cli/exchange/Bitz"
	"trade_api/src/main/web/cli/exchange/Coinbene"
	"trade_api/src/main/web/cli/exchange/Cointiger"
	"truxing/commons/log"
)

var (
	//env string
	platforms = []exchange.Exchange{
		Bibox.Bibox{}, //bibox
		//Biki.Biki{},           //biki
		Bitz.Bitz{},           //bitz
		Coinbene.Coinbene{},   //coinbene
		Cointiger.Cointiger{}, //cointiger
	}
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
	var err error
	s, c = common.Connect("platform_amount", exchange.Name(), "local")
	defer s.Close()
	trades = exchange.AmountHandler()
	fmt.Printf("the length is %d", len(trades))
	fmt.Println()
	for _, trade := range trades {
		//不同交易所的uniq field不一样
		switch exchange.Name() {
		case "bibox":
		case "bitz":
		case "coinbene":
		case "cointiger":
			_, err = c.Upsert(bson.M{"symbol": trade.Symbol}, bson.M{"$set": trade})
		case "biki":
			_, err = c.Upsert(bson.M{"id": trade.ID}, bson.M{"$set": trade})
		default:
			panic("err")
		}
		if err != nil {
			log.Debugf("platform:%s symbol %s trade update failed", exchange.Name(), trade.Symbol)
		}
	}
}

func updatePairAndAmount(e exchange.Exchange) {
	updatePair(e)
	updateAmount(e)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(len(platforms))
	for _, v := range platforms {
		go func(e exchange.Exchange) {
			defer wg.Done()
			updatePairAndAmount(e)
		}(v)
	}
	wg.Wait()
}
