package bitfinex_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/bitfinex"
)

func TestBitfinex_GetTicker(t *testing.T) {
	b := bitfinex.Bitfinex{}
	ticker, err := b.GetTicker( "usd","btc")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestBitfinex_GetTades(t *testing.T) {
	b := bitfinex.Bitfinex{}
	trades, err := b.GetTades("usd", "btc", 100)
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range trades {
		fmt.Printf("%v", v)
	}
}

func TestBitfinex_GetMarkets(t *testing.T) {
	b := bitfinex.Bitfinex{}
	pair, err := b.GetMarkets()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", pair)
}
