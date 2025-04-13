package myokx

// Swap 下单参数
type SwapOrderParam struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode"`  // isolated, cross
	Side    string `json:"side"`    // buy, sell
	PosSide string `json:"posSide"` // long, short
	OrdType string `json:"OrdType"` //market, limit, post_only, fok, ioc
	Sz      string `json:"sz"`      // 数量
	Px      string `json:"px"`      // 价格
}

// 撤单参数
type CancelOrderParam struct {
	InstId string `json:"instId"`
	OrdId  string `json:"ordId"`
}

// 下单返回的数据结构，主要的
type BaseOrderResult struct {
	OrderId string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	Tag     string `json:"tag"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

// 下单返回的数据结构，完整的
type OrderResult struct {
	ApiCodeMsg
	Data []BaseOrderResult `json:"data"`
}

// 查询订单时，返回的结构体
type BaseTradeOrderInfoResult struct {
	OrdId     string `json:"ordId"`
	State     string `json:"state"`
	AccFillSz string `json:"accFillSz"`
}

type TradeOrderInfoResult struct {
	ApiCodeMsg
	Data []BaseTradeOrderInfoResult `json:"data"`
}

// 查询未成交订单时，返回的结构体
type BaseSpendingOrderResult struct {
	InstId  string `json:"instId"`
	OrdId   string `json:"ordId"`
	Px      string `json:"px"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	PosSide string `json:"posside"`
	CTime   string `json:"cTime"`
}

type SpendingOrderResult struct {
	ApiCodeMsg
	Data []BaseSpendingOrderResult `json:"data"`
}

// 撤单返回的结构体
type BaseCancelOrderResult struct {
	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
}

type CancelOrderResult struct {
	ApiCodeMsg
	Data []BaseCancelOrderResult `json:"data"`
}
