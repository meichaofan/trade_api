package model

import "time"

// Ticker 行情数据
type Ticker struct {
	Buy                  float64 // 买一价
	Sell                 float64 // 卖一价
	Last                 float64 // 最新成交价
	Low                  float64 // 最低价
	High                 float64 // 最高价
	Vol                  float64 // 成交量(最近 24 小时)
	Time                 time.Time
	Raw                  string // exchange original info
	PriceChangePercent   float64
	Exchange             float64
	PriceChangePercent1h float64
	ExchargeBit          float64 `json:"ExchangeBit"`
}

type ExchangeTickers []*ExchangeTicker;

// 交易所不同交易对行情
type ExchangeTicker struct {
	MarketPair         MarketPairInfo //交易对
	Vol                float64        // 24成交量
	Amount             float64        // 24h成交额
	Last               float64        // 最新价格
	LastUSD            float64        // 最新价格折换成美元
	PriceChangePercent float64        //涨跌幅
	Time               time.Time
}

// 平台交易额
type ExchangeAmount struct {
	PlatForm  string
	AmountUSD float64
}

type SimpleTicker struct {
	Price  float64
	Symbol string
}

type BookTicker struct {
	Symbol    string
	BidPrice  float64
	BidAmount float64
	AskPrice  float64
	AskAmount float64
}

// Trades 多个历史成交
type Trades []*Trade

// 历史成交
type Trade struct {
	ID        int64 // trade id
	Type      string
	Price     float64
	Amount    float64
	Time      time.Time
	Raw       string // exchange original info
	TradeTime int64  `json:"trade_time"`
}

// Kline OHLC struct
type Record struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
	Vol   float64
	Time  time.Time
	Raw   string
	Ktime int64 `json:"Ktime"`
}

type MarketOrder struct {
	Price  float64
	Amount float64
}

type OrderBook struct {
	Asks      []MarketOrder // 卖方深度
	Bids      []MarketOrder // 买方深度
	Time      time.Time
	OrderTime int64  `json:"Otime"`
	Raw       string // exchange original info
}

type Order struct {
	ID         int64
	Amount     float64
	DealAmount float64
	Price      float64
	Status     string
	Type       string
	Side       string
	Raw        string
}

type Balance struct {
	Currency string
	Free     float64
	Frozen   float64
}

type MarketInfo struct {
	Symbol        string
	DecimalPlaces int64
	MinAmount     float64
	Fee           float64
}

type Excharge struct {
	Exchargebtcusd  float64 // 买一价
	Exchargeethusd  float64 // 买一价
	Exchargedkktusd float64 // 买一价

}

type MarketPairInfo struct {
	Quote string `json:"quote" bson:"quote"`
	Base  string `json:"base" bson:"base"`
}
