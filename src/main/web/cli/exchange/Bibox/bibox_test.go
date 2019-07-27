package Bibox_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/exchange/Bibox"
)

func TestBibox_GetRate(t *testing.T) {
	bibox := Bibox.Bibox{}
	rate := bibox.GetRate("btc", "usdt")
	fmt.Printf("rate: %f", rate)
}

func TestBibox_PairHandler(t *testing.T) {
	bibox := Bibox.Bibox{}
	exchangeTickers := bibox.PairHandler()
	for _, v := range exchangeTickers {
		fmt.Printf("symbol:%s\tlast:%f\tlast_usd:%f\tlast_cny:%f\tamount_%s:%f\tamount_%s:%f\tamount_usd:%f\tamount_cny:%f\tpcg:%f\ttime:%s\n",
			v.Symbol, v.Last, v.LastUsd, v.LastCny, v.Quote, v.AmountQuote, v.Base, v.AmountBase, v.AmountUsd, v.AmountCny, v.PriceChangePercent, v.Time)
	}
}

func Test_AmountUsd(t *testing.T) {
	bibox := Bibox.Bibox{}
	exchangeTickers := bibox.PairHandler()
	var amountUsd float64 = 0
	for _, v := range exchangeTickers {
		amountUsd += v.AmountUsd
	}
	fmt.Printf("amount_usd:%f\n", amountUsd)
}
