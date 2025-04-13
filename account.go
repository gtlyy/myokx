package myokx

// 查看交易账户中资金余额信息
func (client *Client) GetAccountBalance(ccy string) (r AccountBalanceResult, err error) {
	uri := ACCOUNT_BALANCE
	if ccy != "" {
		uri = uri + "?ccy=" + ccy
	}
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return AccountBalanceResult{}, err
	}
	return r, err
}

// 获取当前账户可交易产品的信息列表
func (client *Client) GetAccountInstruments(params map[string]string) (r AccountInstrumentsResult, err error) {
	uri := BuildUri(ACCOUNT_INSTRUMENTS, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return AccountInstrumentsResult{}, err
	}
	return r, err
}

// 查看持仓信息
func (client *Client) GetPositions(optionalParams map[string]string) (r PositionsResult, err error) {
	uri := BuildUri(ACCOUNT_POSITIONS, "", optionalParams)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return PositionsResult{}, err
	}
	return r, nil
}
