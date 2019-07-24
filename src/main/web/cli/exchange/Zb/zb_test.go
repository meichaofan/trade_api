package Zb_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Zb"
)

func TestZb_GetRate(t *testing.T) {
	c := Zb.Zb{}
	rate := c.GetRate("usdt", "qc")
	fmt.Printf("rate: %f", rate)
}


func TestZb_PairHandler(t *testing.T) {
	c := Zb.Zb{}
	ets := c.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f percent:%f time:%s", v.Symbol, v.Amount, v.AmountUsd, v.Last, v.LastUsd, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}

