// 功能：测试数据库操作函数。 ps:这里可以直接写入数据库的，yeah.
package myokx

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gtlyy/mytime"

	// "mytushare"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试：增加，通过 id, bar, start, end
func TestInsertUseIdAndBar1(t *testing.T) {
	id := "DOGE-USDT-SWAP"
	bar := "15m"
	start := "-1" // -1:auto update, or Default: "2021-01-01T00:00:00Z"
	// start := "2024-02-25T00:00:00Z"
	end := "2021-02-01T00:15:05Z" // -1: auto update to now
	// end := "-1"
	table := strings.Replace(id+bar, "-", "", -1)
	if !maria.CheckTableExists(table) {
		maria.CreateTable(table)
	}
	err := maria.InsertUseIdAndBar(c, id, bar, start, end, table, true) // true 表示不更新到变化中的当前kline
	assert.True(t, err == nil)
}

// 测试：增加，通过 id, bar, start, end
func TestInsertUseIdAndBar2(t *testing.T) {
	ids := []string{"DOGE", "ETC", "BTC", "ETH", "SOL", "LTC", "XRP", "BCH", "MATIC", "OKB"}
	// ids := []string{"DOGE", "ETC", "BTC"}
	bars := []string{"1H", "1D", "15m"}
	for _, id := range ids {
		id = id + "-USDT"
		for _, bar := range bars {
			start := "-1" // -1:auto update, or Default: "2021-01-01T00:00:00Z"
			// start := "2024-01-01T00:00:00Z"
			end := "-1" // -1: auto update to now
			// end := "2023-01-01T00:00:00Z"
			table := strings.Replace(id+bar, "-", "", -1)
			err := maria.InsertUseIdAndBar(c, id, bar, start, end, table, true) // true 表示不更新到变化中的当前kline
			assert.True(t, err == nil)
		}
	}
}

// 测试：删除
func TestDeleteNum(t *testing.T) {
	id := "DOGE-USDT-SWAP"
	bar := "1H"
	table := strings.Replace(id+bar, "-", "", -1)
	orderby := "ts"
	n_del := 1
	err := maria.DeleteNum(table, n_del, -1, orderby)
	assert.True(t, err == nil)
}

// 测试：查询1
func TestQueryNum(t *testing.T) {
	id := "KAITO-USDT-SWAP"
	bar := "15m"
	table := strings.Replace(id+bar, "-", "", -1)
	var r []KlineDataS
	n := 1000
	err := maria.QueryNum(&r, table, -1, n)
	assert.True(t, err == nil)
	assert.True(t, len(r) == n)
	for _, v := range r {
		fmt.Printf("%s,%s,%s,%s,%s,%s\n", v.O, v.H, v.L, v.C, v.Vol, v.Tstocst)
	}
}

// 测试：查询2
func TestQueryStartEnd(t *testing.T) {
	end := "2022-01-01T05:00:00.000Z"
	start := "2022-01-01T00:00:00.000Z"
	id := "DOGE-USDT-SWAP"
	bar := "1H"
	table := strings.Replace(id+bar, "-", "", -1)
	var r []KlineDataS
	err := maria.QueryStartEnd(&r, table, start, end)
	t.Log(TsToISOCST(r[0].Ts), TsToISOCST(r[len(r)-1].Ts))
	assert.True(t, start == TsToISOCST(r[0].Ts))
	assert.True(t, end == TsToISOCST(r[len(r)-1].Ts))
	assert.True(t, err == nil)
}

// 测试： 获取close价格，以用于计算macd等。
func TestQueryClose(t *testing.T) {
	// end := "2022-01-01T05:00:00.000Z"
	end := mytime.ISONowCST()
	start := "2021-01-01T00:00:00.000Z"
	id := "DOGE-USDT-SWAP"
	bar := "1H"
	table := IdAndBarToTable(id, bar)
	cs, err := maria.QueryClose(table, start, end)
	assert.True(t, err == nil)
	t.Log(len(cs))
}

// 测试：Query() 随机n条 ok
func TestQuerySql(t *testing.T) {
	var r1 []KlineDataS
	var r2 []KlineDataS
	n := 10
	query1 := "SELECT ts FROM DOGEUSDTSWAP1H ORDER BY RAND() LIMIT 1"
	maria.Query(&r1, query1)
	t.Log(r1[0].Ts)

	query2 := "SELECT * FROM DOGEUSDTSWAP1H WHERE ts >= " + r1[0].Ts + " LIMIT " + IntToString(n)
	maria.Query(&r2, query2)
	t.Log(r2)
}

// 测试：QueryRand() ok
func TestQueryRand(t *testing.T) {
	var r []KlineDataS
	id := "DOGE-USDT-SWAP"
	bar := "1H"
	table := IdAndBarToTable(id, bar)
	n := 10
	// query1 := "SELECT ts FROM DOGEUSDTSWAP1H ORDER BY RAND() LIMIT 1"
	// maria.Query(&r, query1)
	// t.Log(r[0].Ts)

	// query2 := "SELECT * FROM DOGEUSDTSWAP1H WHERE ts >= " + r[0].Ts + " LIMIT " + IntToString(n)
	// maria.Query(&r, query2)
	// t.Log(r)

	err := maria.QueryRand(&r, table, n)
	assert.True(t, err == nil)
	assert.True(t, len(r) == n)
}

// 测试：获取A股的代码，上海。混合大A和加密币。 这个增加：返回股票名称
func TestCreateTradeGameData3(t *testing.T) {
	// 使用随机数生成索引来选择一个随机的 bar
	bars := []string{"15m", "1H", "1D"}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(bars))
	bar := bars[randomIndex]

	for i := 0; i < 10; i++ {
		r, stock, name := maria.CreateTradeGameData3(true, true, bar)
		t.Log(i, r[0].C, stock, name)
	}
}

// 测试：InsertUseIdAndBarSync 增加，通过 id, bar, start, end
// 功能：可以通过这个测试函数，直接更新数据库。
func TestInsertUseIdAndBarSyncCh(t *testing.T) {
	id := "BTC-USDT-SWAP"
	bar := "1H"   //15m, 1H, 1D
	start := "-1" // -1:auto update, or Default: "2021-01-01T00:00:00Z"
	maria.InsertUseIdAndBarSyncCh(c, id, bar, start)
}

func TestGetRowCount(t *testing.T) {
	table := "000001SZ1D"
	n, err := maria.GetRowCount(table)
	assert.True(t, err == nil)
	t.Log(n)
	assert.True(t, n == 5738)
}
