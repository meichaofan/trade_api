package main

import (
	"fmt"
	"trade_api/src/main/conf"
	handler "trade_api/src/main/web/router"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println(conf.GetMode())
	fmt.Println(conf.GetPort())

	router := gin.Default()
	//交易单
	router.POST("/btc/trade", handler.GetTradeHandler)
	//委托单
	router.POST("/btc/order", handler.GetOrderHandler)
	//交易对
	router.POST("/btc/marketpairinfo", handler.GetMarketPairInfoHandler)
	//平台交易额
	router.POST("/btc/exchangeamount", handler.GetExchangeAmount)
	//平台交易对、最新价
	router.POST("/btc/exchangeticker", handler.GetExchangeTicker)
	//交易所列表
	router.POST("/btc/platformlist", handler.GetPlatformList)

	router.Run(conf.GetPort())
}
