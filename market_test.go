package myokx

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"testing"

	"github.com/gtlyy/myfun"
	"github.com/stretchr/testify/assert"
)

// 测试：获取指定类型的行情信息
func TestGetTickerType(t *testing.T) {
	r, err := c.GetTickerType("SWAP")
	assert.True(t, err == nil)
	sum := 0
	IDS := []string{"DOGE-USDT-SWAP", "ETH-USDT-SWAP"}
	for _, v := range r.Data {
		if myfun.In(v.InstId, IDS) {
			sum += 1
		}
	}
	fmt.Println(sum)
}

// 测试：获取单个产品行情信息
func TestGetTicker(t *testing.T) {
	r, err := c.GetTicker("ETH-USDT-SWAP")
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	println(jstr)
}

// 测试：获取产品深度列表
func TestGetBooks(t *testing.T) {
	p := NewParams()
	p["instId"] = "ETH-USDT-SWAP"
	p["sz"] = "3" // 最大 400
	r, err := c.GetBooks(p)
	assert.True(t, err == nil)
	jstr, _ := Struct2JsonString(r)
	println(jstr)
}

// Klines ================================================================================= Start:
// 测试：获取K线数据 (end, start)  0->1  new->old
func TestGetKlines1(t *testing.T) {
	p := NewParams()
	// p["limit"] = "5"		// 单次最多300
	p["after"] = ISOCSTToTs("2024-03-06T00:00:00.000Z")  // end
	p["before"] = ISOCSTToTs("2024-03-01T00:00:00.000Z") // start
	p["bar"] = "1D"
	r, err := c.GetKlines("DOGE-USDT-SWAP", p)
	assert.True(t, err == nil)
	fmt.Println("len =", len(r.Data))
	for i := 0; i < len(r.Data); i++ {
		k := r.Data[i]
		fmt.Printf("ts=%s, ts_iso=%s, confirm=%s\n", k[0], TsToISOCST(k[0]), k[8])
	}
}

// 测试：获取最近n条k线数据，其中最近那条是未完成的，即 confirm=0
func TestGetKlines2(t *testing.T) {
	p := NewParams()
	p["limit"] = "5" // 单次最多300
	// p["after"] = ISOCSTToTs("2024-03-06T00:00:00.000Z")  // end
	// p["before"] = ISOCSTToTs("2024-03-01T00:00:00.000Z") // start
	p["bar"] = "1m"
	r, err := c.GetKlines("DOGE-USDT-SWAP", p)
	assert.True(t, err == nil)
	fmt.Println("len =", len(r.Data))
	for i := 0; i < len(r.Data); i++ {
		k := r.Data[i]
		fmt.Printf("ts=%s, ts_iso=%s, confirm=%s\n", k[0], TsToISOCST(k[0]), k[8])
	}
}

// 测试：获取K线历史数据  (end, start) 0->1  new->old
func TestGetKlinesHistory1(t *testing.T) {
	p := NewParams()
	// p["limit"] = "5"                                     // 单次最多100
	p["after"] = ISOCSTToTs("2024-01-06T00:00:00.000Z")  // end
	p["before"] = ISOCSTToTs("2024-01-01T00:00:00.000Z") // start
	p["bar"] = "1D"
	r, err := c.GetKlinesHistory("DOGE-USDT-SWAP", p)
	assert.True(t, err == nil)
	fmt.Println("len =", len(r.Data))
	for i := 0; i < len(r.Data); i++ {
		k := r.Data[i]
		fmt.Printf("ts=%s, ts_iso=%s, confirm=%s\n", k[0], TsToISOCST(k[0]), k[8])
	}
}

// 测试：获取最近n条历史k线数据，历史数据都是 confirm=1
func TestGetKlinesHistory2(t *testing.T) {
	p := NewParams()
	p["limit"] = "5" // 单次最多100
	// p["after"] = ISOCSTToTs("2024-01-06T00:00:00.000Z")  // end
	// p["before"] = ISOCSTToTs("2024-01-01T00:00:00.000Z") // start
	p["bar"] = "1m"
	r, err := c.GetKlinesHistory("DOGE-USDT-SWAP", p)
	assert.True(t, err == nil)
	fmt.Println("len =", len(r.Data))
	for i := 0; i < len(r.Data); i++ {
		k := r.Data[i]
		fmt.Printf("ts=%s, ts_iso=%s, confirm=%s\n", k[0], TsToISOCST(k[0]), k[8])
	}
}

// 测试TestTimeSplit1函数
func TestTimeSplit1(t *testing.T) {
	end := "2021-01-02T00:00:00.000Z"
	start := "2021-01-01T00:00:00.000Z"
	bar := "5m"
	r := TimeSplit1(ISOCSTToTs(end), ISOCSTToTs(start), bar)
	fmt.Println(len(r))
	for i := 0; i < len(r); i++ {
		fmt.Println(TsToISOCST(r[i][0]), TsToISOCST(r[i][1]))
	}
}

// 测试TestTimeSplit2函数
func TestTimeSplit2(t *testing.T) {
	end := "2021-01-02T00:00:00.000Z"
	start := "2021-01-01T00:00:00.000Z"
	bar := "5m"
	r := TimeSplit2(ISOCSTToTs(end), ISOCSTToTs(start), bar)
	fmt.Println(len(r))
	for i := 0; i < len(r); i++ {
		fmt.Println(TsToISOCST(r[i][0]), TsToISOCST(r[i][1]))
	}
}

// 测试：获取Kline历史，加强版1  (end, start]
func TestGetKlinesHistoryPlus1(t *testing.T) {
	p := NewParams()
	p["after"] = ISOCSTToTs("2023-12-31T00:00:00.000Z")  // end
	p["before"] = ISOCSTToTs("2023-01-01T00:00:00.000Z") // start
	p["bar"] = "1D"
	r, err := c.GetKlinesHistoryPlus1("DOGE-USDT-SWAP", p)
	assert.True(t, err == nil)
	t.Log("获取的k线数量 =", len(r))
	for _, v := range r {
		fmt.Printf("%s %s\n", v[0], TsToISOCST(v[0]))
	}
}

// 测试：获取Kline历史，加强版2。 flag=0, (end, start) ; flag=1, (end, start] ; flag=2, [end, start]
func TestGetKlinesHistoryPlus2(t *testing.T) {
	p := NewParams()
	p["after"] = ISOCSTToTs("2024-01-10T00:00:00.000Z")  // end
	p["before"] = ISOCSTToTs("2024-01-01T00:00:00.000Z") // start
	p["bar"] = "1D"
	flag := 1
	r, err := c.GetKlinesHistoryPlus2("DOGE-USDT", p, flag)
	assert.True(t, err == nil)
	fmt.Println("len of klines =", len(r))
	for _, v := range r {
		fmt.Printf("%s %s\n", v[0], TsToISOCST(v[0]))
	}
}

// 测试 GetKlinesSync 函数
func TestGetKlinesSync(t *testing.T) {
	SetMyLog()
	p := NewParams()
	bar := "1m"
	var num int64 = 10
	p["after"] = TsNow()
	p["before"] = Int64ToString(StringToInt64(p["after"]) - TimeToMs(bar)*num)
	p["bar"] = bar

	ch := make(chan []KlineData, 10)
	go func() {
		for klineData := range ch {
			for _, kline := range klineData {
				log.Println(kline[0], TsToISOCST(kline[0]))
			}
		}
	}()

	go c.GetKlinesSync("DOGE-USDT", p, ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	fmt.Println(" Received interrupt signal, exiting...")
}

// Klines ================================================================================= End.
