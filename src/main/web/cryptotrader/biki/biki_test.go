package biki_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/biki"
)

func TestBiki_GetTicker(t *testing.T) {
	biki := biki.Biki{}
	ticker, err := biki.GetTicker("usdt", "btc")
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("last %f ", ticker.Last)
}

func TestBiki_GetTades(t *testing.T) {
	biki := biki.Biki{}
	trades, err := biki.GetTades("usdt", "btc")
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range trades {
		fmt.Printf("%v", v)
	}
}

func TestBiki_GetMarketDepth(t *testing.T) {
	biki := biki.Biki{}
	orderBook, err := biki.GetMarketDepth("usdt", "btc", "")
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
	for _,v := range pair{
		fmt.Printf("%s %s", v.Quote,v.Base)
		fmt.Println()
	}
}
/*
func TestBiki_GetExchangeTickers(t *testing.T) {
	biki := biki.Biki{}
	tickers, err := biki.GetExchangeTickers()
	if err!=nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range tickers {
		fmt.Printf("quote:%s base:%s last:%f last_usd:%f", v.MarketPair.Quote, v.MarketPair.Base, v.Last, v.LastUSD)
		fmt.Println()
	}
}
*/

func TestBitBox_GetExchangeTickersFromDb(t *testing.T) {
	c := biki.Biki{}
	exchangeTickers, err := c.GetExchangeTickersFromDb()
	if err != nil {
		panic(err)
	}
	for _, v := range exchangeTickers {
		fmt.Printf("quote: %s,base: %s , last: %f ,last_usd: %f", v.Quote, v.Base, v.Last, v.LastUsd)
		fmt.Println()
	}
}

func TestBitBox_GetExchangeAmountFormDb(t *testing.T) {
	c := biki.Biki{}
	exchangeAmount, err := c.GetExchangeAmountFormDb()
	if err != nil {
		panic(err)
	}
	fmt.Printf("platform:%s total_usd:%f", exchangeAmount.Platform, exchangeAmount.TotalUsd)
	fmt.Println()
}
