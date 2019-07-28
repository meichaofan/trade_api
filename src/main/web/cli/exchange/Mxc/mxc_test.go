package Mxc_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Mxc"
)


func Test(t *testing.T) {
	c := Mxc.Mxc{}
	rate := c.GetRate("btc", "usdt")
	fmt.Printf("rate: %f", rate)
}

func TestMxc_PairHandler(t *testing.T) {
	c := Mxc.Mxc{}
	ets := c.PairHandler()
	i := 0
	for _, v := range ets {
		fmt.Printf("symbol:%s\tlast:%f\tlast_usd:%f\tlast_cny:%f\tamount_%s:%f\tamount_%s:%f\tamount_usd:%f\tamount_cny:%f\tpcg:%f\ttime:%s\n",
			v.Symbol, v.Last, v.LastUsd, v.LastCny, v.Quote, v.AmountQuote, v.Base, v.AmountBase, v.AmountUsd, v.AmountCny, v.PriceChangePercent, v.Time)
		i++
	}
	fmt.Printf("i := %d\n", i)
}
