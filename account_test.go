package myokx

import (
	"testing"

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
