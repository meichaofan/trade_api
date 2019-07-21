package Biki_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"trade_api/src/main/web/cli/exchange/Biki"
)

func TestGetSymbol(t *testing.T) {
	paris := Biki.GetSymbol()
	for _, v := range paris {
		fmt.Println("symbol: ", v)
	}
}

func TestSign(t *testing.T) {
	params := map[string]string{
		"api_key": Biki.ApiKey,
		"time":    strconv.Itoa(int(time.Now().Unix())),
		"symbol":  "btcusdt",
	}
	r := Biki.Sign(params)
	fmt.Println(r)
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
		fmt.Printf("quote %s base %s last %f last_usd %f", v.Quote, v.Base, v.Last, v.LastUsd)
		fmt.Println()
	}
}

func TestBiki_AmountHandler(t *testing.T) {
	biki := Biki.Biki{}
	tds := biki.AmountHandler()
	i := 0
	for _, v := range tds {
		fmt.Printf("time %s : symbol %s price %f price_usd %f amount %f amount_usd %f",v.TradeTime, v.Symbol, v.Price, v.PriceUsd, v.Amount, v.AmountUsd)
		fmt.Println()
		i++
	}
	fmt.Printf("the length is %d", i)
}
