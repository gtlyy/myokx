package myokx

import (
	"fmt"
	"testing"

	"github.com/gtlyy/myfun"
	"github.com/stretchr/testify/assert"
)

// 测试：查看账户余额：单个币种
func TestGetAccountBalance(t *testing.T) {
	r, err := c.GetAccountBalance("USDT")
	assert.True(t, err == nil)
	// fmt.Printf("%+v, %+v\n", r, err) // %+v 打印结构体及字段名，没有+的话只是打印结构体。
	jstr, _ := Struct2JsonString(r)
	println(jstr)
	println("余额=", r.Data[0].TotalEq)
}

// 测试：获取当前账户可交易产品的信息列表
func TestGetAccountInstruments(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	// p["instId"] = "TRUMP-USDT-SWAP"
	r, err := c.GetAccountInstruments(p)
	assert.True(t, err == nil)
	// jstr, _ := Struct2JsonString(r)
	// fmt.Println(jstr)
	myfun.WriteJSONFile("account-instruments.json", r)
}

// 测试：查看持仓信息：指定 id
func TestGetPositionsId(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	p["instId"] = "DOGE-USDT-SWAP"
	r, err := c.GetPositions(p)
	assert.True(t, err == nil)
	// fmt.Printf("%+v, %+v\n", r, err)
	jstr, _ := Struct2JsonString(r)
	fmt.Println(jstr)
}

// 测试：查看持仓信息：只指定类型
func TestGetPositionsType(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	r, err := c.GetPositions(p)
	assert.True(t, err == nil)
	fmt.Println("Id             持仓量        开仓均价    最新价     浮盈（亏）      盈亏率")
	for _, v := range r.Data {
		fmt.Println(v.InstId, v.AvailPos, v.AvgPx, v.Last, v.Upl, v.UplRatio)
	}
}
