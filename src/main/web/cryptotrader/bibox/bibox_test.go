package bibox_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/bibox"
)

func TestBiki_GetTicker(t *testing.T) {
	bitbox := bibox.BitBox{}
	ticker, err := bitbox.GetTicker("btc", "usdt")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestBiki_GetTades(t *testing.T) {
	bitbox := bibox.BitBox{}
	trades, err := bitbox.GetTades("btc", "usdt", 100)
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range trades {
		fmt.Printf("%v", v)
	}
}

func TestBitBox_GetMarketDepth(t *testing.T) {
	bitbox := bibox.BitBox{}
	orderBook, err := bitbox.GetMarketDepth("btc", "usdt",100)
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", orderBook.Time)
}

func TestBiki_GetMarkets(t *testing.T) {
	bitbox := bibox.BitBox{}
	pair, err := bitbox.GetMarkets()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("%v", pair)
	fmt.Printf("the length is %d", len(pair))
}
