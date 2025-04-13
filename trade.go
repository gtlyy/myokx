package myokx

// 下单函数
func (client *Client) TradeOrder(orderParam interface{}) (r OrderResult, err error) {
	if _, err = client.Request(POST, TRADE_ORDER, orderParam, &r); err != nil {
		return OrderResult{}, err
	}
	return r, nil
}

// 获取订单信息
func (client *Client) GetTradeOrderInfo(param map[string]string) (r TradeOrderInfoResult, err error) {
	uri := BuildUri(TRADE_ORDER, "", param)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return TradeOrderInfoResult{}, err
	}
	return r, nil
}

// 获取未成交订单列表
/*
instId string 否
ordType String 	否  订单类型。
			market：市价单。
			limit：限价单。
			post_only：只做maker单。
			fok：全部成交或立即取消。
			ioc：立即成交并取消剩余。
			optimal_limit_ioc：市价委托立即成交并取消剩余（仅适用交割、永续）
instType String 否  产品类型。SPOT币币  MARGIN杠杆  SWAP永续合约  FUTURES交割合约  OPTION期权
*/
func (client *Client) GetPendingOrders(param map[string]string) (r SpendingOrderResult, err error) {
	uri := BuildUri(TRADE_PENDING_ORDER, "", param)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return SpendingOrderResult{}, err
	}
	return r, nil
}

/*
撤单
instId
ordId
*/
func (client *Client) CancelOrders(param interface{}) (r CancelOrderResult, err error) {
	if _, err = client.Request(POST, TRADE_CANCEL_ORDER, param, &r); err != nil {
		return CancelOrderResult{}, err
	}
	return r, nil
}
