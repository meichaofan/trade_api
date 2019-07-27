package exchange

import "trade_api/src/main/web/cli/data"

/**
交易所
*/

type Exchange interface {
	//交易所名称
	Name() string
	//获取交易对的汇率
	GetRate(quote, base string) float64
	//交易所所有交易对
	PairHandler() []*data.ExchangeTicker
}
