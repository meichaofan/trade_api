package Mxc_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Mxc"
)
/*
func Test(t *testing.T) {
	c := Mxc.Mxc{}
	rate := c.GetRate("btc", "usd")
	fmt.Printf("rate: %f", rate)
}*/

func TestMxc_PairHandler(t *testing.T) {
	c := Mxc.Mxc{}
	ets := c.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f percent:%f time:%s", v.Symbol, v.Amount, v.AmountUsd, v.Last, v.LastUsd, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}