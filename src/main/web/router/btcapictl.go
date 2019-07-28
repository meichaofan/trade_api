package web

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"trade_api/src/main/web/cli/common"
	"trade_api/src/main/web/cli/data"
	"trade_api/src/main/web/cryptotrader/bibox"
	"trade_api/src/main/web/cryptotrader/biki"
	"trade_api/src/main/web/cryptotrader/bitfinex"
	"trade_api/src/main/web/cryptotrader/bitz"
	coinall2 "trade_api/src/main/web/cryptotrader/coinall"
	"trade_api/src/main/web/cryptotrader/coinbene"
	"trade_api/src/main/web/cryptotrader/cointiger"
	"truxing/commons/log"

	"github.com/Akagi201/cryptotrader/binance"

	"github.com/Akagi201/cryptotrader/huobi"

	"fmt"
	"strconv"
	"time"
	dconf "trade_api/src/main/conf"
	"trade_api/src/main/lib/cache"
	bconf "truxing/commons/conf"

	"github.com/Akagi201/cryptotrader/coincheck"
	"github.com/Akagi201/cryptotrader/fcoin"
	"github.com/Akagi201/cryptotrader/gateio"
	"github.com/Akagi201/cryptotrader/gdax"
	"github.com/Akagi201/cryptotrader/mxc"
	"github.com/Akagi201/cryptotrader/okcoin"
	"github.com/Akagi201/cryptotrader/okex"
	"github.com/Akagi201/cryptotrader/zb"
	"github.com/gin-gonic/gin"

	"github.com/Akagi201/cryptotrader/model"
)

var (
	incache  cache.Cache
	resource *bconf.ResourceConfig
)

func init() {

	log.Info("btcapi start")
	dconf.InitConfig()
	resource = dconf.Conf().Resource
	redisCfg := resource.Redis["base"]
	incache = cache.NewRedisCache(fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port))
	log.Info("[btcapi] init redis configuration end...")
}
func set_redis_cache(key, value string) error {
	reply := incache.SetExpire(key, value, 30*time.Second)
	if reply != nil {
		log.Error("setredis to redis error,%s", reply)
		return reply
	}

	return nil
}

func set_exchange_cache(key, value string) error {
	reply := incache.SetExpire(key, value, 14400*time.Second)
	if reply != nil {
		log.Error("setredis to redis error,%s", reply)
		return reply
	}

	return nil
}

func get_redis_cache(key string) (value string, err error) {
	reply, err := incache.Get(key)
	//log.Errorln("value:%s", reply.Val())
	if err != nil {
		log.Error("get redis error:", err)
		return "", err
	}
	if len(reply) == 0 {
		log.Error("get redis null")
		return "", nil
	}
	return reply, nil
}

func set_redis_cache_for_exptime(key, value string, exptime string) error {
	exp, _ := time.ParseDuration(exptime)
	reply := incache.SetExpire(key, value, exp)
	if reply != nil {
		log.Error("setredis to redis error,%s", reply)
		return reply
	}

	return nil
}

func get_excharge_cache(key string) (value string, err error) {
	reply, err := incache.Get(key)
	//log.Errorln("value:%s", reply.Val())
	if err != nil {
		log.Error("get redis error:", err)
		return "", err
	}
	if len(reply) == 0 {
		log.Error("get redis null")
		return "", nil
	}
	return reply, nil
}

func ZbTradeapi(base, quote, typ, acckey, secret string, since, size int, platform string) ([]model.Trade, error) {
	var trades []model.Trade
	var err error
	switch platform {
	case "zb":
		api := zb.New(acckey, secret)
		trades, err = api.GetTrades(base, quote, int(since))
	case "huobipro":
		api := huobi.New(acckey, secret)
		trades, err = api.GetTrades(base, quote)
	case "binance":
		api := binance.New(acckey, secret)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		since = since * 1000
		trades, err = api.GetTrades(ctx, quote, base, 0, int64(since), 0, int64(size))
	case "coincheck":

		api := coincheck.New(acckey, secret)
		trades, err = api.GetTrade(base, quote, typ, int(since), int(size))

	case "okcoin-intl":

		api := okcoin.New(acckey, secret)
		trades, err = api.GetTrade(base, quote, typ, int(since), int(size))
	case "okex":

		api := okex.New(acckey, secret)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		trades, err = api.GetTrades(ctx, quote, base, int(since))
	case "gdax":

		api := gdax.New(acckey, secret)
		trades, err = api.GetTrades(base, quote, typ, int(since), int(size))
	case "gate-io":
		api := gateio.New(acckey, secret)
		trades, err = api.GetTrades(quote, base)
	case "fcoin":
		api := fcoin.New(acckey, secret)
		trades, err = api.GetTrades(quote, base)
	case "coinall":
		api := coinall2.New(acckey, secret)
		trades, err = api.GetTades(base, quote, int(since))
	case "bibox":
		api := bibox.New("", "")
		trades, err = api.GetTades(base, quote, size)
	case "biki":
		api := biki.New("", "")
		trades, err = api.GetTades(base, quote)
	case "bitz":
		api := bitz.New("", "")
		trades, err = api.GetTades(base, quote)
	case "coinbene":
		api := coinbene.New("", "")
		trades, err = api.GetTades(base, quote, size)
	case "cointiger":
		api := cointiger.New("", "")
		trades, err = api.GetTades(base, quote, size)
	case "bitfinex":
		api := bitfinex.New("", "")
		trades, err = api.GetTades(base, quote, size)
	default:
		trades = []model.Trade{}
		err = nil
	}
	fmt.Println(len(trades))

	return trades, err
}

func ZbOrderBook(base, quote, typ, acckey, secret string, since, size int, platform string) (model.OrderBook, error) {
	var trades model.OrderBook
	var err error
	switch platform {
	case "huobipro":
		api := huobi.New(acckey, secret)
		trades, err = api.GetDeep(base, quote)
	case "binance":
		api := binance.New(acckey, secret)
		trades, err = api.GetDeep(base, quote)
	/*case "coincheck":

	api := coincheck.New(acckey, secret)
	trades, err = api.GetTrade(base, quote, typ, int(since), int(size))
	*/
	case "okcoin-intl":

		api := okcoin.New(acckey, secret)
		trades, err = api.GetDeep(base, quote)
	case "okex":

		api := okex.New(acckey, secret)
		trades, err = api.GetDeep(base, quote)
	/*case "gdax":

	api := gdax.New(acckey, secret)
	trades, err = api.GetTrades(base, quote, typ, int(since), int(size))*/
	case "gate-io":

		api := gateio.New(acckey, secret)
		trades, err = api.GetDeep(quote, base)
	case "zb":
		api := zb.New(acckey, secret)
		trades, err = api.GetOrderBook(base, quote)
	case "fcoin":
		api := fcoin.New(acckey, secret)
		trades, err = api.GetDeep(quote, base)
	case "coinall":
		api := coinall2.New(acckey, secret)
		trades, err = api.GetDepth(base, quote)
	case "bibox":
		api := bibox.New("", "")
		trades, err = api.GetMarketDepth(base, quote, size)
	case "biki":
		api := biki.New("", "")
		trades, err = api.GetMarketDepth(base, quote, "step0")
	case "bitz":
		api := bitz.New("", "")
		trades, err = api.GetOrderBook(base, quote)
	case "coinbene":
		api := coinbene.New("", "")
		trades, err = api.GetOrderBook(base, quote, size)
	case "cointiger":
		api := cointiger.New("", "")
		trades, err = api.GetMarketDepth(base, quote, "step0")
	case "bitfinex":
		api := bitfinex.New("", "")
		trades, err = api.GetMarketDepth(base, quote)
	default:
		trades = model.OrderBook{}
		err = nil
	}

	return trades, err
}

func GetMarketPairInfo(platform string) ([]model.MarketPairInfo, error) {
	var marketpairinfo []model.MarketPairInfo
	var err error
	switch platform {
	case "zb":
		api := zb.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "fcoin":
		api := fcoin.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "gate-io":
		api := gateio.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "mxc":
		api := mxc.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "coinall":
		api := coinall2.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "bibox":
		api := bibox.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "biki":
		api := biki.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "bitz":
		api := bitz.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "coinbene":
		api := coinbene.New("", "")
		marketpairinfo, err = api.GetMarkets()
	case "cointiger":
		api := cointiger.New("", "")
		marketpairinfo, err = api.GetMarkets()
	default:

		marketpairinfo = []model.MarketPairInfo{}
		err = nil
	}

	return marketpairinfo, err
}

func GetTradeHandler(c *gin.Context) {
	base := c.PostForm("base")
	quote := c.PostForm("quote")
	acckey := c.PostForm("accessKey")
	secret := c.PostForm("secretKey")
	sinceform := c.PostForm("since")
	sizeform := c.PostForm("size")
	typ := c.PostForm("type")
	platform := c.PostForm("platform")
	since, err := strconv.Atoi(sinceform)
	if err != nil {
		log.Errorf("Get since failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "error",
			},
			"data": "",
		})
		return
	}
	size, err := strconv.Atoi(sizeform)
	if err != nil {
		log.Errorf("Get size failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"data": "",
		})
		return
	}
	result, err := ZbTradeapi(base, quote, typ, acckey, secret, since, size, platform)

	if err != nil || len(result) == 0 {
		fmt.Println(base)
		fmt.Println(quote)
		log.Errorf("Get result failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": 500,
				"status":      "error",
			},
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"status_code": http.StatusOK,
			"status":      "success",
		},
		"data": result,
	})
}

func GetOrderHandler(c *gin.Context) {
	base := c.PostForm("base")
	quote := c.PostForm("quote")
	acckey := c.PostForm("accessKey")
	secret := c.PostForm("secretKey")
	sinceform := c.PostForm("since")
	sizeform := c.PostForm("size")
	typ := c.PostForm("type")
	platform := c.PostForm("platform")
	since, err := strconv.Atoi(sinceform)
	if err != nil {
		log.Errorf("Get since failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "error",
			},
			"data": "",
		})
		return
	}
	size, err := strconv.Atoi(sizeform)
	if err != nil {
		log.Errorf("Get size failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"data": "",
		})
		return
	}
	result, err := ZbOrderBook(base, quote, typ, acckey, secret, since, size, platform)

	if err != nil {
		fmt.Println(base)
		fmt.Println(quote)
		log.Errorf("Get result failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": 500,
				"status":      "error",
			},
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"status_code": http.StatusOK,
			"status":      "success",
		},
		"data": result,
	})
}

func GetMarketPairInfoHandler(c *gin.Context) {

	platform := c.PostForm("platform")

	result, err := GetMarketPairInfo(platform)

	if err != nil {
		log.Errorf("Get result failed, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": 500,
				"status":      "error",
			},
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"status_code": http.StatusOK,
			"status":      "success",
		},
		"data": result,
	})
}

/**
平台交易额
*/
func GetExchangeAmount(c *gin.Context) {
	platform := c.PostForm("platform")
	fmt.Printf("P:%s\n", platform)
	var exchangeTickers []*data.ExchangeTicker
	var amountUsd float64 = 0
	var amountCny float64 = 0
	s, col := common.Connect("platform_pair", platform, "local")
	defer s.Close()
	err := col.Find(bson.M{}).All(&exchangeTickers)
	if err != nil {
		log.Debug("error:%s", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": 500,
				"status":      "error",
			},
			"data": "",
		})
		return
	}
	for _, tradeData := range exchangeTickers {
		amountUsd += tradeData.AmountUsd
		amountCny += tradeData.AmountCny
	}
	exchangeAmount := &data.ExchangeAmount{
		Platform: platform,
		TotalUsd: amountUsd,
		TotalCny: amountCny,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"status_code": http.StatusOK,
			"status":      "success",
		},
		"data": exchangeAmount,
	})
}

/**
平台交易对
*/
func GetExchangeTicker(c *gin.Context) {
	platform := c.PostForm("platform")
	fmt.Printf("P:%s\n", platform)
	var exchangeTicker []*data.ExchangeTicker
	s, col := common.Connect("platform_pair", platform, "local")
	defer s.Close()
	err := col.Find(bson.M{}).All(&exchangeTicker)
	if err != nil {
		log.Debug("error:%s", err)
		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": 500,
				"status":      "error",
			},
			"data": "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": gin.H{
			"status_code": http.StatusOK,
			"status":      "success",
		},
		"data": exchangeTicker,
	})
}
