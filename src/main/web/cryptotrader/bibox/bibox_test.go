package bibox_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cryptotrader/bibox"
)

func TestBiki_GetTicker(t *testing.T) {
	bitbox := bibox.BitBox{}
	ticker, err := bitbox.GetTicker("ETH", "RED") //RED_ETH
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%v", ticker)
}

func TestBiki_GetTades(t *testing.T) {
	bitbox := bibox.BitBox{}
	trades, err := bitbox.GetTades("usdt", "btc", 100)
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range trades {
		fmt.Printf("%v", v)
	}
}

func TestBitBox_GetMarketDepth(t *testing.T) {
	bitbox := bibox.BitBox{}
	orderBook, err := bitbox.GetMarketDepth("usdt", "btc", 100)
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
	for _, v := range pair {
		fmt.Printf("base: %s , quote: %s", v.Base, v.Quote)
		fmt.Println()
	}
	fmt.Printf("the length is %d", len(pair))
}

func TestBitBox_GetExchangeTickers(t *testing.T) {
	bitbox := bibox.BitBox{}
	exchangeTickers, err := bitbox.GetExchangeTickers()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	for _, v := range exchangeTickers {
		fmt.Printf("quote: %s,base: %s , last: %f ,last_usd: %f", v.MarketPair.Quote, v.MarketPair.Base, v.Last, v.LastUSD)
		fmt.Println()
	}
}

func TestBitBox_GetExchangeAmount(t *testing.T) {
	bitbox := bibox.BitBox{}
	exchangeAmount, err := bitbox.GetExchangeAmount()
	if err != nil {
		fmt.Printf("error:%s", err)
	}
	fmt.Printf("platform:%s", exchangeAmount.PlatForm)
	fmt.Println()
	fmt.Printf("amount:%f", exchangeAmount.AmountUSD)
}
