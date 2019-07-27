package Binance_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Binance"
)

func TestBinance_PairHandler(t *testing.T) {
	b := Binance.Binance{}
	ets := b.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol:%s\tlast:%f\tlast_usd:%f\tlast_cny:%f\tamount_%s:%f\tamount_%s:%f\tamount_usd:%f\tamount_cny:%f\tpcg:%f\ttime:%s\n",
			v.Symbol, v.Last, v.LastUsd, v.LastCny, v.Quote, v.AmountQuote, v.Base, v.AmountBase, v.AmountUsd, v.AmountCny, v.PriceChangePercent, v.Time)
	}
}

func Test_AmountUsd(t *testing.T) {
	b := Binance.Binance{}
	ets := b.PairHandler()
	var amountUsd float64 = 0
	for _, v := range ets {
		amountUsd += v.AmountUsd
	}
	fmt.Printf("amount_usd:%f\n", amountUsd)
}

func TestLast(t *testing.T) {
	Binance.Last()
}