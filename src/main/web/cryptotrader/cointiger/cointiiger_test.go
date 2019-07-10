package cointiger_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/cointiger"
)

func TestBiki_GetTicker(t *testing.T) {
	ct := cointiger.CoinTiger{}
	ticker, err := ct.GetTicker("btc", "usdt")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestBiki_GetTades(t *testing.T) {
	ct := cointiger.CoinTiger{}
	trades, err := ct.GetTades("btc", "usdt", 100)
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range trades {
		fmt.Printf("%v", v)
	}
}

func TestCoinTiger_GetMarketDepth(t *testing.T) {
	ct := cointiger.CoinTiger{}
	orderBook, err := ct.GetMarketDepth("btc", "usdt", "step0")
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", orderBook.Time)
}

func TestCoinTiger_GetMarkets(t *testing.T) {
	ct := cointiger.CoinTiger{}
	pair, err := ct.GetMarkets()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", pair)
	fmt.Printf("the length is %d", len(pair))
}
