package Bitfinex

import (
	"github.com/tidwall/gjson"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

const (
	ApiHost = "https://api.bitfinex.com/v1"
)

var (
	//计价货币
	CountRate = []string{
		"usd",
		"btc",
		"eur",
		"jpy",
		"eth",
		"gbp",
		"dai",
		"eos",
		"xch",
	}
)

type Bitfinex struct {
}

func (b Bitfinex) Name() string {
	return "bitfinex"
}

/**
获取该交易所支持的交易对
https://docs.bitfinex.com/reference#rest-public-symbols
*/
func GetSymbol() []*data.MarketPair {
	var pairs []*data.MarketPair
	url := ApiHost + "/symbols"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	symbols := ret.Array()

	ret.Get("data").ForEach(func(key, value gjson.Result) bool {
		//symbol := value.Get("symbol").Str
		pair := &data.MarketPair{
			Symbol: value.Get("symbol").Str,
			Quote:  value.Get("base_coin").Str,
			Base:   value.Get("count_coin").Str,
		}
		pairs = append(pairs, pair)
		return true
	})
	return pairs
}