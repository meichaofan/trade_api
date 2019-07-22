package Cointiger_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Cointiger"
)

func TestCointiger_GetRate(t *testing.T) {
	c := Cointiger.Cointiger{}
	rate := c.GetRate("eth", "usdt")
	fmt.Printf("rate: %f", rate)
}

func TestCointiger_PairHandler(t *testing.T) {
	c := Cointiger.Cointiger{}
	ets := c.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f percent:%f time:%s", v.Symbol, v.Amount, v.AmountUsd, v.Last, v.LastUsd, v.PriceChangePercent, v.Time)
		fmt.Println()
	}
}

func TestCointiger_AmountHandler(t *testing.T) {
	c := Cointiger.Cointiger{}
	ets := c.AmountHandler()
	for _, v := range ets {
		fmt.Printf("symbol: %s,amount:%f amount_usd:%f last: %f last_usd:%f", v.Symbol, v.Amount, v.AmountUsd, v.Price, v.PriceUsd)
		fmt.Println()
	}
}
