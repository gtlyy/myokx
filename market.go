package myokx

// 获取指定类型的行情信息
func (client *Client) GetTickerType(instType string) (r TickerResult, err error) {
	uri := MARKET_TICKERS + "?instType=" + instType
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return TickerResult{}, err
	}
	return r, nil
}

// 获取单个产品行情信息
func (client *Client) GetTicker(instId string) (r TickerResult, err error) {
	uri := MARKET_TICKER + "?instId=" + instId
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return TickerResult{}, err
	}
	return r, nil
}
