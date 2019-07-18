package main

import (
	"fmt"
	"testing"
)

func TestExchangeRate(t *testing.T) {
	r := ExchangeRate("omg", "usdt")	// 1.402088
	fmt.Printf("rate %f", r)
}
