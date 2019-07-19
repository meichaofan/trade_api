package coinall

import (
	"testing"
)

func TestCoinAll_GetMarkets(t *testing.T) {
	ca := &CoinAll{}
	ca.GetMarkets()
}
func TestCoinAll_GetDeep(t *testing.T) {
	ca := &CoinAll{}
	ca.GetDepth("usdt","btc")
}