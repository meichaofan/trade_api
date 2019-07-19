package bitz_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/bitz"
)

func TestBitZ_GetTicker(t *testing.T) {
	bz := bitz.BitZ{}
	ticker, err := bz.GetTicker("usdt","btc")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestBitZ_GetTades(t *testing.T) {
	bz := bitz.BitZ{}
	trades, err := bz.GetTades( "usdt","btc")
	if err!=nil {
		fmt.Printf("error:%s", err)
	}
	for _,v := range trades{
		fmt.Printf("%v", v)
	}
}

func TestBitZ_GetOrderBook(t *testing.T) {
	bz := bitz.BitZ{}
	orderBook, err := bz.GetOrderBook( "usdt","btc")
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", orderBook.Time)
}

func TestBitZ_GetMarkets(t *testing.T) {
	bz := bitz.BitZ{}
	pair, err := bz.GetMarkets()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", pair)
	fmt.Printf("the length is %d", len(pair))
}