package Bibox

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
文档:
*/
const (
	ApiHost = "https://api.bibox365.com/v1/mdata"
)

var (
	//定义taker的成交方向
	sideMap = map[string]string{
		"1": "buy",
		"2": "sell",
	}
	//虚拟货币兑换美元的汇率
	coinRate sync.Map
)

type Bibox struct {
}

func (bibox Bibox) Name() string {
	return "bibox"
}

/**
获取汇率
*/
func (bibox Bibox) GetRate(quote, base string) float64 {
	var exchangeRate float64
	if rate, ok := coinRate.Load(base); ok == true {
		exchangeRate = rate.(float64)
	} else if strings.ToUpper(base) == "USDT" || strings.ToUpper(base) == "GUSD" || strings.ToUpper(base) == "DAI" { //稳定币
		/*		url := ApiHost + "?cmd=ticker&pair=" + strings.ToUpper(quote) + "_" + strings.ToUpper(base)
				content := common.HttpGet(url)
				last := gjson.ParseBytes(content).Get("result.last").Float()
				coinRate.Store(quote, last)*/
		exchangeRate = 1
	} else {
		nUrl := ApiHost + "?cmd=ticker&pair=" + strings.ToUpper(base) + "_USDT"
		//log.Debugf("---quote:%s----base:%s",quote,base)
		nContent := common.HttpGet(nUrl)
		nLast := gjson.ParseBytes(nContent).Get("result.last").Float()
		coinRate.Store(base, nLast)
		exchangeRate = nLast
	}
	return exchangeRate
}

/**
获取平台交易对及其价格
*/
func (bibox Bibox) PairHandler() []*data.ExchangeTicker {
	var exchangeTickers []*data.ExchangeTicker
	url := ApiHost + "?cmd=marketAll"
	content := common.HttpGet(url)
	gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
		percent := value.Get("percent").Str
		priceChangePercent, _ := strconv.ParseFloat(strings.TrimRight(percent, "%"), 64)
		exchangeTicker := &data.ExchangeTicker{
			Symbol:             strings.ToUpper(value.Get("coin_symbol").Str + "_" + value.Get("currency_symbol").Str),
			Quote:              value.Get("coin_symbol").Str,
			Base:               value.Get("currency_symbol").Str,
			Volume:             value.Get("vol24H").Float(),
			Amount:             value.Get("amount").Float(),
			Last:               value.Get("last").Float(),
			LastUSD:            value.Get("last_usd").Float(),
			PriceChangePercent: priceChangePercent,
			Time:               strconv.FormatInt(time.Now().Unix(), 10),
		}
		exchangeTickers = append(exchangeTickers, exchangeTicker)
		return true
	})
	return exchangeTickers
}

/**
获取平台交易额
*/
func (bibox Bibox) AmountHandler() []*data.TradeData {
	//首先获取所有的交易对
	var tradeDatas []*data.TradeData
	pairs := bibox.GetAllPair()

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(pairs))

	for _, pair := range pairs {
		go func(quote, base string) {
			defer waitGroup.Done()
			url := ApiHost + "?cmd=deals&pair=" + quote + "_" + base
			content := common.HttpGet(url)
			gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
				trade := &data.TradeData{
					Symbol:    quote + "_" + base,
					Quote:     quote,
					Base:      base,
					Amount:    value.Get("amount").Float(),
					Price:     value.Get("price").Float(),
					Type:      sideMap[value.Get("side").String()],
					TradeTime: strconv.FormatInt(value.Get("time").Int(), 10),
					TradeTs:   value.Get("time").Int() / 1000,
				}
				if strings.ToUpper(base) == "USDT" {
					trade.PriceUsd = trade.Price
					trade.AmountUsd = trade.Amount
				} else {
					rate := bibox.GetRate(quote, base)
					trade.PriceUsd = rate * trade.Price
					trade.AmountUsd = rate * trade.Amount
				}
				tradeDatas = append(tradeDatas, trade)
				return true
			})
		}(pair.Quote, pair.Base)
	}
	waitGroup.Wait()

	/*coinRate.Range(func(key, value interface{}) bool {
		fmt.Printf("key %s", key)
		fmt.Println()
		fmt.Printf("value %f", value)
		fmt.Println()
		return true
	})*/

	return tradeDatas
}

/**
获取平台交易对，及其价格
*/
func (bibox Bibox) GetAllPair() []*data.MarketPair {
	var pairs []*data.MarketPair
	url := ApiHost + "?cmd=pairList"
	content := common.HttpGet(url)
	gjson.ParseBytes(content).Get("result").ForEach(func(key, value gjson.Result) bool {
		quote, base := strings.Split(value.Get("pair").String(), "_")[0], strings.Split(value.Get("pair").String(), "_")[1]
		pair := &data.MarketPair{
			Base:  base,
			Quote: quote,
		}
		pairs = append(pairs, pair)
		return true
	})
	return pairs
}
