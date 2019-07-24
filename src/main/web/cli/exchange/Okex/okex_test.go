package Okex_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Okex"
)


func TestOkex_GetRate(t *testing.T) {
	c := Okex.Okex{}
	rate := c.GetRate("eth", "usdt")
	fmt.Printf("rate: %f", rate)
}

func TestOkex_PairHandler(t *testing.T) {
	c := Okex.Okex{}
	ets := c.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f percent:%f time:%s", v.Symbol, v.Amount, v.AmountUsd, v.Last, v.LastUsd, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}


