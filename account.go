package myokx

import (
	"log"
	"time"
)

// 获取当前账户可交易产品的信息
func (client *Client) GetAccountInstruments(params map[string]string) (r AccountInstrumentsResult, err error) {
	uri := BuildUri(ACCOUNT_INSTRUMENTS, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return AccountInstrumentsResult{}, err
	}
	return r, err
}

// 获取下单价格精度
func (client *Client) GetTickSz(instType, instId string) (ticksz string, err error) {
	p := NewParams()
	p["instType"] = instType
	p["instId"] = instId
	r, err := client.GetAccountInstruments(p)
	if err != nil {
		log.Println("Error: In GetTickSz().", err)
		return "", err
	}
	ticksz = r.Data[0].TickSz
	return ticksz, err
}

// 快速获取交易价格精度
func GetTickSzQuick(id1 string) string {
	if id1 == "DOGE-USDT-SWAP" {
		return "0.00001"
	} else if id1 == "ETH-USDT-SWAP" {
		return "0.01"
	} else if id1 == "BTC-USDT-SWAP" {
		return "0.1"
	} else if id1 == "TRUMP-USDT-SWAP" {
		return "0.001"
	} else {
		return GetTickSzFromJson("account-instruments.json", id1)
	}
}

// 获取下单价格精度，从本地文件
func GetTickSzFromJson(filename, instId string) string {
	var a AccountInstrumentsResult
	ReadJSONFile(filename, &a)
	ticksz := ""
	for i := range len(a.Data) {
		if instId == a.Data[i].InstId {
			ticksz = a.Data[i].TickSz
		}
	}
	return ticksz
}

// 获取合约面值，从本地文件
func GetctValFromJson(filename, instId string) string {
	var a AccountInstrumentsResult
	ReadJSONFile(filename, &a)
	ctVal := ""
	for i := range len(a.Data) {
		if instId == a.Data[i].InstId {
			ctVal = a.Data[i].CtVal
		}
	}
	return ctVal
}

// 将获取的可交易产品信息并写入json文件，如写入 account-instruments.json
func (client *Client) GetAndWriteAccountInstruments(instType, filename string) {
	p := NewParams()
	p["instType"] = instType
	r, err := client.GetAccountInstruments(p)
	IfError("Error:Get,Write account-instruments.", err)
	WriteJSONFile(filename, r)
}

// 获取交易账户中资金余额信息
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

// 获取持仓信息
func (client *Client) GetPositions(optionalParams map[string]string) (r PositionsResult, err error) {
	uri := BuildUri(ACCOUNT_POSITIONS, "", optionalParams)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return PositionsResult{}, err
	}
	return r, nil
}

// 获取近7天的账户流水
func (client *Client) GetBills(optionalParams map[string]string) (r BillsResult, err error) {
	uri := BuildUri(ACCOUNT_BILLS, "", optionalParams)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return BillsResult{}, err
	}
	return r, nil
}

// 获取近3个月的账户流水
func (client *Client) Get3MonthsBills(optionalParams map[string]string) (r BillsResult, err error) {
	uri := BuildUri(ACCOUNT_BILLS_ARCHIVE, "", optionalParams)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return BillsResult{}, err
	}
	return r, nil
}

// 获取账户流水 Plus版 ： 任意时间段
func (client *Client) GetBillsPlus(optionalParams map[string]string) ([]Bills, error) {
	start := optionalParams["start"]
	var err error
	bs := make([]Bills, 0, 1440) // 1440 no meanings.

	i := 0
	for {
		var r BillsResult
		uri := BuildUri(ACCOUNT_BILLS_ARCHIVE, "", optionalParams)
		_, err = client.Request(GET, uri, nil, &r)
		bs = append(bs, r.Data...)
		if len(r.Data) < 100 || start > r.Data[len(r.Data)-1].Ts {
			break
		}
		optionalParams["after"] = r.Data[len(r.Data)-1].BillId
		delete(optionalParams, "end")
		if (i+1)%5 == 0 {
			time.Sleep(time.Second * 2)
		}
		i += 1
	}

	// 保留 start 至 end
	for i = 0; i < len(bs); i++ {
		if start > bs[i].Ts {
			break
		}
	}
	return bs[0:i], err
}

// 获取账户配置
func (client *Client) GetAccountConfig() (r ConfigResult, err error) {
	uri := ACCOUNT_CONFIG
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return ConfigResult{}, err
	}
	return r, nil
}

// 获取最大可下单数量：对应下单时的 sz 字段
func (client *Client) GetMaxSize(params map[string]string) (r MaxSizeResult, err error) {
	uri := BuildUri(Max_SIZE, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return MaxSizeResult{}, err
	}
	return r, nil
}

// 获取最大可用余额/保证金
func (client *Client) GetMaxAvailSize(params map[string]string) (r MaxAvailSizeResult, err error) {
	uri := BuildUri(Max_Avail_SIZE, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return MaxAvailSizeResult{}, err
	}
	return r, nil
}

// 获取当前账户交易手续费率
func (client *Client) GetTradeFee(params map[string]string) (r TradeFeeResult, err error) {
	uri := BuildUri(TRADE_FEE, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return TradeFeeResult{}, err
	}
	return r, nil
}

// 获取杠杆倍数
func (client *Client) GetLeverInfo(optionalParams map[string]string) (r LeverInfoResult, err error) {
	uri := BuildUri(ACCOUNT_LEVERINFO, "", optionalParams)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return LeverInfoResult{}, err
	}
	return r, nil
}

// 设置杠杠倍数
func (client *Client) SetLever(leverParam interface{}) (r LeverResult, err error) {
	if _, err = client.Request(POST, SET_LEVER, leverParam, &r); err != nil {
		return LeverResult{}, err
	}
	return r, nil
}
