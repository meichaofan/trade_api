package Coinall

import (
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"sync"
	"time"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
)

/**
文档：https://www.coinall.live/docs/zh/
*/

const (
	ApiHost = "https://www.coinall.live"
)

var (
	//虚拟货币 -- 美元 汇率
	rateCoin sync.Map
)

type Coinall struct {
}

func (c Coinall) Name() string {
	return "coinall"
}

//虚拟货币 -- 美元 汇率
func (c Coinall) GetRate(quote, base string) float64 {
	symbol := strings.ToUpper(quote + "-" + base)
	if rate, ok := rateCoin.Load(symbol); ok {
		r := rate.(float64)
		return r
	}
	url := ApiHost + "/api/spot/v3/instruments/" + symbol + "/ticker"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	if ret.Exists() {
		//btc - usdt
		//eth - usdt
		rate := ret.Get("last").Float()
		rateCoin.Store(symbol, rate)
		return rate
	}
	return 1;
}

func (c Coinall) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "/api/spot/v3/instruments/ticker"
	content := common.HttpGet(url)
	ret := gjson.ParseBytes(content)
	ret.ForEach(func(key, value gjson.Result) bool {
		var quote string
		var base string
		var symbol string = value.Get("instrument_id").Str
		quote = strings.Split(symbol, "-")[0]
		base = strings.Split(symbol, "-")[1]

		//时间
		timeLayout := "2006-01-02T15:04:05Z" //转化所需模板
		loc, _ := time.LoadLocation("Local") //获取时区
		tmp, _ := time.ParseInLocation(timeLayout, value.Get("timestamp").Str, loc)
		timeStr := strconv.FormatInt(tmp.Unix(), 10) + "000" //转化为13位时间戳

		//涨幅
		openPrice := value.Get("open_24h").Float()
		last := value.Get("last").Float()
		pcg := (last - openPrice) / openPrice

		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(symbol),
			Quote:              strings.ToUpper(quote),
			Base:               strings.ToUpper(base),
			Volume:             value.Get("base_volume_24h").Float(),
			Amount:             value.Get("quote_volume_24h").Float(),
			Last:               last,
			Time:               timeStr,
			PriceChangePercent: pcg,
		}

		//汇率
		if strings.ToUpper(base) == "USDT" {
			exchangeTicker.LastUsd = exchangeTicker.Last
			exchangeTicker.AmountUsd = exchangeTicker.Amount
		} else if strings.ToUpper(base) == "USDK" {  //怎么出来这么一个货币 -)'
			rate := c.GetRate("USDT", "USDK")
			exchangeTicker.LastUsd = exchangeTicker.Last / rate
			exchangeTicker.AmountUsd = exchangeTicker.Amount / rate
		} else {
			rate := c.GetRate(base, "USDT")
			exchangeTicker.LastUsd = exchangeTicker.Last * rate
			exchangeTicker.AmountUsd = exchangeTicker.Amount * rate
		}
		exchangeTickers = append(exchangeTickers, exchangeTicker)
		return true
	})

	return exchangeTickers
}

func (c Coinall) AmountHandler() []*data.TradeData {
	var tradeDatas []*data.TradeData
	return tradeDatas
}
