package data

//-----------------------------------------------
//某个时间，交易单
type TradeData struct {
	ID        int64   //交易ID
	PairQuote string  //交易货币
	PairBase  string  // 计价货币
	Symbol    string  //交易对
	Type      string  //交易类型 [买-卖]
	Price     float64 //交易价格
	Volume    float64 //交易量
	Amount    float64 //交易额
	AmountUsd float64 //交易额(美元)
	PriceUsd  float64 //交易兑换成美元
	PriceCny  float64 //交易兑换成人民币
	TradeTime int64   // 交易时间
}

//平台交易额
type ExchangeAmount struct {
	Platform string  //交易平台
	TotalUsd float64 //交易额(美元)
	TotalCny float64 //交易额(人民币)
}

//----------------------------------------------

//-----------------交易所所支持的交易对最新价格--------------------
//平台交易所
type ExchangeTicker struct {
	Symbol             string  //交易对
	Quote              string  //交易货币
	Base               string  //计价货币
	Vol                float64 // 24成交量
	Amount             float64 // 24h成交额
	Last               float64 // 最新价格
	LastUSD            float64 // 最新价格折换成美元
	PriceChangePercent float64 //涨跌幅
	Time               int64
}

type MarketPair struct {
	Quote string
	Base  string
}
