package coinbene_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/coinbene"
)

/*func TestGetSymbol(t *testing.T) {
	r := coinbene.GetSymbol("btc ", " usdt")
	if r == "BTCUSDT" {
		fmt.Println("success")
	} else {
		fmt.Println("error")
	}
}*/

func TestCoinBene_GetTicker(t *testing.T) {
	cb := coinbene.CoinBene{}
	ticker, err := cb.GetTicker( "usdt","btc")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestCoinBene_GetTades(t *testing.T) {
	cb := coinbene.CoinBene{}
	trades, err := cb.GetTades( "usdt", "btc",200)
	if err!=nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", trades[0])
}

func TestCoinBene_GetOrderBook(t *testing.T) {
	cb := coinbene.CoinBene{}
	orderBook, err := cb.GetOrderBook( "usdt", "btc",200)
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", orderBook)
}

func TestCoinBene_GetMarkets(t *testing.T) {
	cb := coinbene.CoinBene{}
	pair, err := cb.GetMarkets()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", pair)
	fmt.Printf("the length is %d", len(pair))
}