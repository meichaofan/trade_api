package Coinbene_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Coinbene"
)

func TestCoinbene_GetRate(t *testing.T) {
	c := Coinbene.Coinbene{}
	rate := c.GetRate("eth", "usdt")
	fmt.Printf("rate: %f", rate)
}

func TestCoinbene_PairHandler(t *testing.T) {
	c := Coinbene.Coinbene{}
	ets := c.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f percent:%f time:%s", v.Symbol, v.Amount, v.AmountUsd, v.Last, v.LastUsd, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}

func TestBitz_AmountHandler(t *testing.T) {
	c := Coinbene.Coinbene{}
	ets := c.AmountHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f", v.Symbol, v.Amount, v.AmountUsd, v.Price, v.PriceUsd)
		fmt.Println()
	}
}
