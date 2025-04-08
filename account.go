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
