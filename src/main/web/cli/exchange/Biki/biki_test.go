package Biki_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Biki"
)

func TestGetSymbol(t *testing.T) {
	paris := Biki.GetSymbol()
	for _, v := range paris {
		fmt.Println("symbol: ", v)
	}
}

/*
func TestGetRate(t *testing.T) {
	f := Biki.GetRate("eos", "btc")
	fmt.Printf("%f", f)
}
*/
func TestBiki_PairHandler(t *testing.T) {
	biki := Biki.Biki{}
	ets := biki.PairHandler()
	for _, v := range ets {
		fmt.Printf("symbol:%s\tlast:%f\tlast_usd:%f\tlast_cny:%f\tamount_%s:%f\tamount_%s:%f\tamount_usd:%f\tamount_cny:%f\tpcg:%f\ttime:%s\n",
			v.Symbol, v.Last, v.LastUsd, v.LastCny, v.Quote, v.AmountQuote, v.Base, v.AmountBase, v.AmountUsd, v.AmountCny, v.PriceChangePercent, v.Time)
	}
}

func Test_AmountUsd(t *testing.T) {
	biki := Biki.Biki{}
	exchangeTickers := biki.PairHandler()
	var amountUsd float64 = 0
	for _, v := range exchangeTickers {
		amountUsd += v.AmountUsd
	}
	fmt.Printf("amount_usd:%f\n", amountUsd)
}
