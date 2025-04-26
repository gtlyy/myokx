package myokx

// 获取公开产品的信息
func (client *Client) GetPublicInstruments(params map[string]string) (r AccountInstrumentsResult, err error) {
	uri := BuildUri(PUBLIC_INSTRUMENTS, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return AccountInstrumentsResult{}, err
	}
	return r, err
}

// 将获取公开产品信息并写入json文件，如写入 account-instruments.json
func (client *Client) GetAndWritePublicInstruments(instType, filename string) {
	p := NewParams()
	p["instType"] = instType
	r, err := client.GetPublicInstruments(p)
	IfError("Error:Get,Write account-instruments.", err)
	WriteJSONFile(filename, r)
}

// 获取服务器时间
func (client *Client) GetServerTime() (serverTime ServerTime, err error) {
	_, err = client.Request(GET, OKEX_TIME_URI, nil, &serverTime)
	return serverTime, err
}

// 获取永续合约当前资金费率
func (client *Client) GetFundingRate(instId string) (r FundingRate, err error) {
	uri := SWAP_FUNDING_RATE + "?instId=" + instId
	_, err = client.Request(GET, uri, nil, &r)
	return
}

// 获取永续合约历史资金费率
func (client *Client) GetFundingRateHistory(instId string, optionalParams map[string]string) (r FundingRate, err error) {
	uri := BuildUri(SWAP_HISTORY_FUNDING_RATE, instId, optionalParams)
	_, err = client.Request(GET, uri, nil, &r)
	return
}
