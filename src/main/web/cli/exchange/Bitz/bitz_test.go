package Bitz_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Bitz"
)

func TestBitz_PairHandlert(t *testing.T) {
	b := Bitz.Bitz{}
	ets := b.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f percent:%f time:%s", v.Symbol, v.Amount, v.AmountUsd, v.Last, v.LastUsd, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}

func TestBitz_AmountHandler(t *testing.T) {
	b := Bitz.Bitz{}
	ets := b.AmountHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f", v.Symbol, v.Amount, v.AmountUsd, v.Price, v.PriceUsd)
		fmt.Println()
	}
}
