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
	i := 0
	for _, v := range ets {
		fmt.Printf("symbol:%s\tlast:%f\tlast_usd:%f\tlast_cny:%f\tamount_%s:%f\tamount_%s:%f\tamount_usd:%f\tamount_cny:%f\tpcg:%f\ttime:%s\n",
			v.Symbol, v.Last, v.LastUsd, v.LastCny, v.Quote, v.AmountQuote, v.Base, v.AmountBase, v.AmountUsd, v.AmountCny, v.PriceChangePercent, v.Time)
		i++
	}
	fmt.Printf("i := %d\n", i)
}

func Test_AmountUsd(t *testing.T) {
	c := Zb.Zb{}
	ets := c.PairHandler()
	var amountUsd float64 = 0
	for _, v := range ets {
		amountUsd += v.AmountUsd
		//fmt.Printf("symbol:%s , amount_usd:%f\n", v.Symbol, v.AmountUsd)
	}
	fmt.Printf("amount_usd:%f\n", amountUsd)
}
