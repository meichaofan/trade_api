package data

//-----------------------------------------------
//某个时间，交易单
type TradeData struct {
	ID        string  `bson:"id"`         //交易ID
	Quote     string  `bson:"quote"`      //交易货币
	Base      string  `bson:"base"`       // 计价货币
	Symbol    string  `bson:"symbol"`     //交易对
	Type      string  `bson:"type"`       //交易类型 [买-卖]
	Volume    float64 `bson:"volume"`     //交易量
	Price     float64 `bson:"price"`      //交易价格
	PriceUsd  float64 `bson:"price_usd"`  //交易兑换成美元
	Amount    float64 `bson:"amount"`     //交易额
	AmountUsd float64 `bson:"amount_usd"` //交易额(美元)
	TradeTime string  `bson:"trade_time"` // 交易时间 string类型的时间 "1563522757819" 为了mongodb存储
	TradeTs   int64   `bson:"trade_ts"`   // int类型，用于时间比较，取出1天的数据
}

//平台交易额
type ExchangeAmount struct {
	Platform string  //交易平台
	TotalUsd float64 //交易额(美元)
	//TotalCny float64 //交易额(人民币)
}

//----------------------------------------------

//-----------------交易所所支持的交易对最新价格--------------------
//平台交易所
type ExchangeTicker struct {
	Symbol             string  `bson:"symbol"`               //交易对
	Quote              string  `bson:"quote"`                //交易货币
	Base               string  `bson:"base"`                 //计价货币
	AmountQuote        float64 `bson:"amount_quote"`         // 24h成交额(交易货币)
	AmountBase         float64 `bson:"amount_base"`          // 24h成交额 (计价货币)
	AmountUsd          float64 `bson:"amount_usd"`           // 24h成交额 (美元)
	AmountCny          float64 `bson:"amount_cny"`           // 24h成交额 (人民币)
	Last               float64 `bson:"last"`                 // 最新价格 (计价货币)
	LastUsd            float64 `bson:"last_usd"`             // 最新价格(美元)
	LastCny            float64 `bson:"last_cny"`             // 最新价格(人民币)
	PriceChangePercent float64 `bson:"price_change_percent"` //涨跌幅
	Time               string  `bson:"time"`                 //string类型的时间 "1563522757819" 为了mongodb存储
}

type MarketPair struct {
	Symbol string
	Quote  string
	Base   string
}
