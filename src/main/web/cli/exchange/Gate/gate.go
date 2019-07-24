package Gate

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档: https://www.gateio.co/docs/futures/api/index.html
*/

const (
	ApiHost = "https://api.gateio.ws/api/v4"
)

type Gate struct {
}

func (c Gate) Name() string {
	return "gate"
}

//虚拟货币 -- 这里全部都是以usd计价
func (c Gate) GetRate(quote, base string) float64 {
	return 1
}

func (c Gate) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/futures/tickers"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	ret.ForEach(func(key, value gjson.Result) bool {
		symbol := value.Get("contract").Str
		quote := strings.Split(symbol, "_")[0]
		base := strings.Split(symbol, "_")[1]
		timeStr := strconv.FormatInt(time.Now().UnixNano(), 10)
		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			Volume:             value.Get("total_size").Float(), //交易量
			Last:               value.Get("last").Float(),
			Amount:             value.Get("volume_24h_usd").Float(),
			LastUsd:            value.Get("last").Float(),
			AmountUsd:          value.Get("volume_24h_usd").Float(),
			Time:               timeStr,
			PriceChangePercent: value.Get("change_percentage").Float(),
		}
		exchangeTickers = append(exchangeTickers, exchangeTicker)
		return true
	})

	return exchangeTickers
}

func (c Gate) AmountHandler() []*data.TradeData {
	var tradeDatas []*data.TradeData
	return tradeDatas
}
