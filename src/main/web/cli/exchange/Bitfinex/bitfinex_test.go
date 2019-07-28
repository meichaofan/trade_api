package Bitfinex_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Bitfinex"
)

func TestGetSymbol(t *testing.T) {
	paris := Bitfinex.GetSymbol()
	for _, v := range paris {
		fmt.Printf("%v\n", v)
	}
}

/*
func TestGetCountCoinRate(t *testing.T) {
	s := Bitfinex.GetCountCoinRate()
	for k, v := range s {
		fmt.Printf("coin: %s\t rate :%f\n", k, v)
	}
}
*/

func TestBitfinex_PairHandler(t *testing.T) {
	b := Bitfinex.Bitfinex{}
	ets := b.PairHandler()
	i := 0
	for _, v := range ets {
		fmt.Printf("symbol:%s\tlast:%f\tlast_usd:%f\tlast_cny:%f\tamount_%s:%f\tamount_%s:%f\tamount_usd:%f\tamount_cny:%f\tpcg:%f\ttime:%s\n",
			v.Symbol, v.Last, v.LastUsd, v.LastCny, v.Quote, v.AmountQuote, v.Base, v.AmountBase, v.AmountUsd, v.AmountCny, v.PriceChangePercent, v.Time)
		i++
	}
	fmt.Printf("i := %d\n", i)
}
