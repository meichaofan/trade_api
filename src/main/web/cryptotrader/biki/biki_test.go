package biki_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/biki"
)

func TestBiki_GetTicker(t *testing.T) {
	biki := biki.Biki{}
	ticker, err := biki.GetTicker("btc", "usdt")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestBiki_GetTades(t *testing.T) {
	biki := biki.Biki{}
	trades, err := biki.GetTades("btc", "usdt")
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range trades {
		fmt.Printf("%v", v)
	}
}

func TestBiki_GetMarketDepth(t *testing.T) {
	biki := biki.Biki{}
	orderBook, err := biki.GetMarketDepth("btc", "usdt","")
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", orderBook.Time)
}

func TestBiki_GetMarkets(t *testing.T) {
	biki := biki.Biki{}
	pair, err := biki.GetMarkets()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", pair)
	fmt.Printf("the length is %d", len(pair))
}
