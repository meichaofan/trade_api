package main

import (
	"fmt"
	conf "trade_api/src/main/conf"
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
	router.Run(conf.GetPort())
}
