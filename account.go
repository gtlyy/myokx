package myokx

import (
	"github.com/gtlyy/myfun"
)

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

// 将获取的可交易产品信息并写入json文件，如写入 account-instruments.json
func (client *Client) GetAndWriteAccountInstruments(instType, filename string) {
	p := NewParams()
	p["instType"] = instType
	r, err := client.GetAccountInstruments(p)
	myfun.IfError("Error:Get,Write account-instruments.", err)
	myfun.WriteJSONFile(filename, r)
}

// 查看持仓信息
func (client *Client) GetPositions(optionalParams map[string]string) (r PositionsResult, err error) {
	uri := BuildUri(ACCOUNT_POSITIONS, "", optionalParams)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return PositionsResult{}, err
	}
	return r, nil
}
