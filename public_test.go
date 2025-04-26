package myokx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试：获取产品行情信息
func TestGetPublicInstruments(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	// p["instId"] = "BTC-USDT-SWAP"
	r, err := c.GetPublicInstruments(p)
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	println(jstr)
	WriteJSONFile("account-instruments.json", r)
}

// 测试：获取服务器时间
func TestGetServerTime(t *testing.T) {
	serverTime, err := c.GetServerTime()
	if err != nil {
		t.Error(err)
	}
	ts1 := serverTime.Data[0].Ts
	t.Log(TsToISOCST(ts1))
}

// 测试：获取永续合约的当前资金费率
func TestGetSwapFundingRate(t *testing.T) {
	r, err := c.GetFundingRate("DOGE-USDT-SWAP")
	assert.True(t, err == nil)
	t.Log(r.Data[0].FundingRate)
	t.Log(TsToISOCST(r.Data[0].FundingTime))
}

// 测试：获取永续合约历史资金费率
func TestGetSwapHistoryFundingRate(t *testing.T) {
	p := NewParams()
	p["limit"] = "10"
	r, err := c.GetFundingRateHistory("DOGE-USDT-SWAP", p)
	assert.True(t, err == nil)
	assert.True(t, len(r.Data) == 10)
	for _, v := range r.Data {
		fmt.Printf("%s %s %.4f%%\n", v.InstId, TsToISOCST(v.FundingTime), StringToFloat64(v.FundingRate)*100)
	}
}
