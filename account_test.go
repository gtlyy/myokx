package myokx

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/gtlyy/mytime"
	"github.com/stretchr/testify/assert"
)

// 测试：获取当前账户可交易产品的信息
func TestGetAccountInstruments(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	// p["instId"] = "BTC-USDT-SWAP"
	r, err := c.GetAccountInstruments(p)
	assert.True(t, err == nil)
	WriteJSONFile("account-instruments.json", r)
}

// 测试：获取交易价格精度
func TestGetTickSz(t *testing.T) {
	ticksz, err := c.GetTickSz("SWAP", "BTC-USDT-SWAP")
	assert.True(t, err == nil)
	assert.True(t, ticksz == "0.1")

	ticksz, err = c.GetTickSz("SWAP", "ETH-USDT-SWAP")
	assert.True(t, err == nil)
	assert.True(t, ticksz == "0.01")

	ticksz, err = c.GetTickSz("SWAP", "DOGE-USDT-SWAP")
	assert.True(t, err == nil)
	assert.True(t, ticksz == "0.00001")
}

// 测试：从本地json文件，获取交易产品的价格精度
func TestGetTickSzFromJson(t *testing.T) {
	ticksz := ""
	instId := "DOGE-USDT-SWAP"
	filename := "account-instruments.json"
	ticksz = GetTickSzFromJson(filename, instId)
	assert.True(t, ticksz == "0.00001")
}

// 测试：查看账户余额：单个币种
func TestGetAccountBalance(t *testing.T) {
	ccy := "TRX"
	balance := "0"
	r, err := c.GetAccountBalance(ccy)
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	println(jstr)
	if len(r.Data[0].Details) > 0 {
		balance = r.Data[0].Details[0].Eq
	}
	fmt.Printf("Balance=%s USD\n", r.Data[0].TotalEq)
	fmt.Printf("%s's Balance=%s %s\n", ccy, balance, ccy)
}

// 测试：查看账户余额：多个币种
func TestGetAccountBalance2(t *testing.T) {
	s := []string{"BTC", "ETH", "TRX", "USDT"}
	r, err := c.GetAccountBalance(strings.Join(s, ","))
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	println(jstr)
}

// 测试：查看持仓信息：指定 id
func TestGetPositionsId(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	p["instId"] = "DOGE-USDT-SWAP"
	r, err := c.GetPositions(p)
	assert.True(t, err == nil)
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

// 测试：查看持仓盈利情况
func Test2GetPositions2(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	imrs := 0.0
	upls := 0.0
	for i := 0; i < 10; i++ {
		imrs = 0.0
		upls = 0.0
		r, err := c.GetPositions(p)
		assert.True(t, err == nil)
		for _, v := range r.Data {
			imrs += StringToFloat64(v.Imr)
			upls += StringToFloat64(v.Upl)
		}
		log.Printf("Input=%.2f, Upl=%.2f, Ratio=%.2f%%\n", imrs, upls, 100*upls/imrs)
		time.Sleep(2 * time.Second)
	}
}

// 测试：查看账单流水查询 近7天
func TestGetBills(t *testing.T) {
	p := NewParams()
	p["type"] = "2" // 2：交易，8：资金费
	// p["subType"] = "174"	//  173:资金费支出     174:资金费收入
	r, err := c.GetBills(p)
	assert.True(t, err == nil)
	sum := 0.0
	for _, x := range r.Data {
		fmt.Println(x.InstId, x.Fee, mytime.TsToISOCST(x.Ts))
		sum += StringToFloat64(x.Fee)
	}
	sum_pnl := 0.0
	for _, x := range r.Data {
		fmt.Println(x.InstId, x.Pnl, mytime.TsToISOCST(x.Ts))
		sum_pnl += StringToFloat64(x.Pnl)
	}
	fmt.Println("nums of trade:", len(r.Data))
	fmt.Println("sum of fee =", sum)
	fmt.Println("sum of pnl =", sum_pnl)
	fmt.Println("pnl - fee =", sum_pnl+sum)
}

// 测试：查看账单流水查询 近3个月
func TestGet3MonthsBills(t *testing.T) {
	p := NewParams()
	p["type"] = "2"
	r, err := c.Get3MonthsBills(p)
	assert.True(t, err == nil)
	sum := 0.0
	for _, x := range r.Data {
		fmt.Println(x.InstId, x.Fee, mytime.TsToISOCST(x.Ts))
		sum += StringToFloat64(x.Fee)
	}
	sum_pnl := 0.0
	for _, x := range r.Data {
		fmt.Println(x.InstId, x.Pnl, mytime.TsToISOCST(x.Ts))
		sum_pnl += StringToFloat64(x.Pnl)
	}
	fmt.Println("nums of trade:", len(r.Data))
	fmt.Println("sum of fee =", sum)
	fmt.Println("sum of pnl =", sum_pnl)
	fmt.Println("pnl - fee =", sum_pnl+sum)
}

// 测试：GetBillsPlus
func TestGetBillsPlus(t *testing.T) {
	p := NewParams()
	p["type"] = "2"
	p["end"] = ISOCSTToTs("2024-11-02T22:55:00.001Z")
	p["start"] = ISOCSTToTs("2024-09-25T00:00:00.001Z")
	r, err := c.GetBillsPlus(p)
	assert.True(t, err == nil)

	sum, sum_pnl := 0.0, 0.0
	for _, x := range r {
		fmt.Println(x.InstId, x.Pnl, mytime.TsToISOCST(x.Ts))
		sum += StringToFloat64(x.Fee)
		sum_pnl += StringToFloat64(x.Pnl)
	}
	fmt.Println("nums of trade:", len(r))
	fmt.Println("sum of fee =", sum)
	fmt.Println("sum of pnl =", sum_pnl)
	fmt.Println("pnl - fee =", sum_pnl+sum)
}

// 测试：GetBillsPlus ： 统计资金费收入
func TestGetBillsPlus2(t *testing.T) {
	p := NewParams()
	p["type"] = "8"
	p["end"] = mytime.ISOCSTToTs("2025-04-25T22:55:00.001Z")
	p["start"] = mytime.ISOCSTToTs("2024-09-25T00:00:00.001Z")
	r, err := c.GetBillsPlus(p)
	assert.True(t, err == nil)

	sum := 0.0
	for _, x := range r {
		fmt.Println(x.InstId, x.Pnl, mytime.TsToISOCST(x.Ts))
		sum += StringToFloat64(x.Pnl)
	}
	fmt.Println("nums of trade:", len(r))
	fmt.Println("sum of 资金费 =", sum)
}

// 测试：查看账户配置
func TestGetAccountConfig(t *testing.T) {
	r, err := c.GetAccountConfig()
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	fmt.Println(jstr)
}

// 测试：获取最大可下单数量
func TestGetMaxSize(t *testing.T) {
	p := NewParams()
	p["instId"] = "ETH-USDT-SWAP"
	p["tdMode"] = "cross" // cross: 全仓；isolated: 逐仓
	r, err := c.GetMaxSize(p)
	assert.True(t, err == nil)
	fmt.Println(r.Data[0].MaxBuy, r.Data[0].MaxSell)
}

// 测试：获取最大可用余额/保证金
func TestGetMaxAvailSize(t *testing.T) {
	p := NewParams()
	p["instId"] = "DOGE-USDT-SWAP"
	p["tdMode"] = "cross" // cross: 全仓 ； isolated: 逐仓
	r, err := c.GetMaxAvailSize(p)
	assert.True(t, err == nil)
	fmt.Println(r.Data[0].AvailBuy, r.Data[0].AvailSell)
}

// 测试：获取当前账户交易手续费率
func TestGetTradeFee(t *testing.T) {
	p := NewParams()
	p["instType"] = "SWAP"
	// p["instId"] = "DOGE-USDT"
	r, err := c.GetTradeFee(p)
	assert.True(t, err == nil)
	fmt.Println(r.Data[0].Level, r.Data[0].MakerU, r.Data[0].TakerU)
}

// 测试：获取杠杆倍数
func TestGetLeverInfo(t *testing.T) {
	p := NewParams()
	p["instId"] = "DOGE-USDT-SWAP"
	p["mgnMode"] = "cross"
	rd, err := c.GetLeverInfo(p)
	assert.True(t, err == nil)
	for _, r := range rd.Data {
		fmt.Printf("id=%s,保证金模式=%s,方向=%s,杠杆=%s\n", r.InstId, r.MgnMode, r.PosSide, r.Lever)
	}
}

// 测试：设置杠杆倍数
func TestSetLever(t *testing.T) {
	p := NewParams()
	p["instId"] = "XRP-USDT-SWAP"
	p["mgnMode"] = "cross"
	p["lever"] = "50"
	r, err := c.SetLever(p)
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	fmt.Println(jstr)
}
