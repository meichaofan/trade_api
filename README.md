btc_api
===

 btc_api是一个可以取币的日常情况和K线等数据



HTTP服务访问
---

```shell
1、行情接口
$ curl -XPOST 127.0.0.1:25433/btc/ticker -d "quote=eth&base=btc&platform=huobipro"

$ {
    "data": {
        "Buy": 0,
        "Sell": 0.076379,
        "Last": 0.076273,
        "Low": 0.0743,
        "High": 0.077463,
        "Vol": 1478.5294659496,
        "Time": "0001-01-01T00:00:00Z",
        "Raw": ""
    },
    "status": {
        "status": "success",
        "status_code": 200
    }
}
```

#### 请求参数
quote base platform 在exchange_pair表里

**quote** [string] 币种 比特币 btc

**base** [string] 对应现实的货币，人民币 cny

**platform** [string] 选择不同的交易市场


火币  hb
zb  zb  失效
allcoin  allcoin   base  ck.usd
bigone   bigone   USDT
binance  binance
bitfinex bitfinex  USD 美元
bitflyer  bitflyer
bittrex   bittrex  usdt
cex     cex  usd
coincheck  coincheck  usd
coinegg  coinegg   base etc
gateio   gateio  usdt  https://gate.io/api2
lbank    lbank  btc  https://www.lbank.info/api/api-ticker
liqui    liqui  btc https://liqui.io/api
okcoin   okcoin  usd  https://support.okcoin.com/hc/zh-cn/articles/360000697832-REST-API%E5%8F%82%E8%80%83
okex     okex  usdt
poloniex  poloniex  usdt
Bithumb   bithumb  https://www.bithumb.com/u1/US127    https://crix-api-endpoint.upbit.com/v1/crix/candles/lines?code=CRIX.UPBIT.KRW-ETH
Upbit   upbit  https://crix-api-endpoint.upbit.com/v1/crix/trades/ticks?code=CRIX.UPBIT.KRW-BTC&count=1
HitBTC https://api.hitbtc.com/#symbols
bitz https://www.bit-z.com/api.html#ticker
Bitstamp  bitstamp
Bitbank   bitbank
Zaif      zaif
require 	Bithumb  	GDAX  	Upbit  	HitBTC  Bit-Z  	Bitbank  Bitstamp  CoinsBank（not）  	Zaif(not)   IDEX
not add CoinsBank

#### 返回格式

**buy** : 买入量

**sell** : 卖出量

**last** : 最后交易价格

**low** : 24h最低价格

**high**: 24h最高价格

**vol**: 24h交易量

**status**: 状态 success 成功，error失败



```shell
2、K线接口
$ curl -XPOST 127.0.0.1:25433/btc/records -d "quote=eth&base=btc&since=1520906400&type=1min"

$ {
    "data":  [
        {
            "Open": 0.07561841,
            "High": 0.07561841,
            "Low": 0.07555555,
            "Close": 0.07555555,
            "Vol": 29.80180156,
            "Time": "2018-03-12T00:59:00Z",
            "Raw": "[1520816340000,\"0.07561841\",\"0.07561841\",\"0.07555555\",\"0.07555555\",\"29.80180156\"]"
        }],
    "status": {
        "status": "success",
        "status_code": 200
    }
}
```

#### 请求参数
quote base platform 在exchange_pair表里

**quote** [string] 币种 比特币 btc

**base** [string] 对应现实的货币，人民币 cny

**platform** [string] 选择不同的交易市场

**type** [string] 时间段
  - 1h 小时
  - 1d 天
  - 1w 周
  
**since** 从这个时间戳之后的  秒

**size** 返回数量 最大2000，最少60

#### 返回格式

**open** : 开盘价

**high** : 最高价格

**low** : 最低价格

**close**: 收盘价格

**vol**: 交易量

**status**: 状态 success 成功，error失败

```shell
3、获取连接接口
$ curl -XPOST 127.0.0.1:25433/btc/geturl -d "url=http://www.baidu.com&valid=1s"


```

#### 请求参数

**url** [string] 请求URL

**valid** [string] 1m 1分钟 1h 1小时






