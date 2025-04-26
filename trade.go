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
func (client *Client) GetPendingOrders(param map[string]string) (r SpendingOrderResult, err error) {
	uri := BuildUri(TRADE_PENDING_ORDER, "", param)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return SpendingOrderResult{}, err
	}
	return r, nil
}

// 撤销订单
func (client *Client) CancelOrders(param interface{}) (r CancelOrderResult, err error) {
	if _, err = client.Request(POST, TRADE_CANCEL_ORDER, param, &r); err != nil {
		return CancelOrderResult{}, err
	}
	return r, nil
}
