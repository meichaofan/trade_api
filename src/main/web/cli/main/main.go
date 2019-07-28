package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
	"trade_api/src/main/web/cli/exchange"
	"trade_api/src/main/web/cli/exchange/Bibox"
	"trade_api/src/main/web/cli/exchange/Biki"
	"trade_api/src/main/web/cli/exchange/Binance"
	"trade_api/src/main/web/cli/exchange/Bitz"
	"trade_api/src/main/web/cli/exchange/Coinall"
	"trade_api/src/main/web/cli/exchange/Coinbene"
	"trade_api/src/main/web/cli/exchange/Cointiger"
	"trade_api/src/main/web/cli/exchange/Fcoin"
	"trade_api/src/main/web/cli/exchange/Gate"
	"trade_api/src/main/web/cli/exchange/Huobi"
	"trade_api/src/main/web/cli/exchange/Okex"
	"trade_api/src/main/web/cli/exchange/Zb"
	"truxing/commons/log"
)

var (
	//env string
	platforms = []exchange.Exchange{
		Bibox.Bibox{}, //bibox
		Biki.Biki{},           //biki
		Bitz.Bitz{}, //bitz
		//Bitfinex.Bitfinex{},
		Coinbene.Coinbene{},   //coinbene
		Cointiger.Cointiger{}, //cointiger
		Coinall.Coinall{},
		Huobi.Huobi{},
		Binance.Binance{},
		Zb.Zb{},
		Okex.Okex{},
		Fcoin.Fcoin{},
		Gate.Gate{},
		//Mxc.Mxc{},
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

func main() {
	var wg sync.WaitGroup
	wg.Add(len(platforms))
	for _, v := range platforms {
		go func(e exchange.Exchange) {
			defer wg.Done()
			updatePair(e)
		}(v)
	}
	wg.Wait()
}
