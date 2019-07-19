package Bibox_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Bibox"
)

func TestBibox_GetRate(t *testing.T) {
	bibox := Bibox.Bibox{}
	rate := bibox.GetRate("btc", "usdt")
	fmt.Printf("rate: %f", rate)
}

func TestBibox_PairHandler(t *testing.T) {
	bibox := Bibox.Bibox{}
	exchangeTickers := bibox.PairHandler()
	for _, v := range exchangeTickers {
		fmt.Printf("symbol:%s , quote:%s , base:%s,vol:%f,amount:%f,last:%f,last_usd:%f,pcg:%f,time:%d", v.Symbol, v.Quote, v.Base, v.Vol, v.Amount, v.Last, v.LastUSD, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}

func TestBibox_AmountHandler(t *testing.T) {
	bibox := Bibox.Bibox{}
	tradeDatas := bibox.AmountHandler()
	for _, v := range tradeDatas {
		fmt.Printf("symbol:%s,price:%f,price_usd:%f,amount:%f,amount_usd:%f", v.Symbol, v.Price, v.PriceUsd, v.Amount, v.AmountUsd)
		fmt.Println()
	}
}

func TestBibox_GetAllPair(t *testing.T) {
	bibox := Bibox.Bibox{}
	pairs := bibox.GetAllPair()
	for _, v := range pairs {
		fmt.Printf("quote:%s - base:%s", v.Quote, v.Base)
		fmt.Println()
	}
}
