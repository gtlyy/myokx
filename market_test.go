package myokx

import (
	"fmt"
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
