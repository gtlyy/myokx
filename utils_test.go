package myokx

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 请求验证 ========================================================================= Start
// 测试：生成签名用的message
func TestPreHashString(t *testing.T) {
	timestamp := "2022-03-08T10:59:25.789Z"
	method := "POST"
	request_path := "/orders?before=2&limit=30"
	body := `{"product_id":"BTC-USDT-0309","order_id":"377454671037444"}`
	result := `2022-03-08T10:59:25.789ZPOST/orders?before=2&limit=30{"product_id":"BTC-USDT-0309","order_id":"377454671037444"}`
	r1 := PreHashString(timestamp, method, request_path, body)
	assert.True(t, r1 == result)
}

// 测试：签名
func TestHmacSha256Base64Signer(t *testing.T) {
	message := "2022-02-12T09:15:43.729ZGET/api/v5/public/time"
	secretKey := "1E3B0C3C2770150248D7331ABD83B32D"
	signRight := "VwF3ol6s0NmEuWi5h7TGNHaqNcf1FJErmPSEPEAtW2U="
	sign, err := HmacSha256Base64Signer(message, secretKey)
	assert.Nil(t, err)
	assert.True(t, sign == signRight)
}

// 请求验证 ========================================================================= End

// TestParseRequestParams tests the ParseRequestParams function.
func TestParseRequestParams(t *testing.T) {
	// Define a sample struct for testing
	type Sample struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	params := Sample{Name: "Alice", Age: 30}
	jsonBody, binBody, err := ParseRequestParams(params)

	assert.NoError(t, err, "expected no error while parsing request parameters")

	expectedJson := `{"name":"Alice","age":30}`
	assert.Equal(t, expectedJson, jsonBody, "the JSON body does not match the expected value")

	binData, err := json.Marshal(params) // binData: []byte
	assert.NoError(t, err, "expected no error while marshaling params")
	actualBinData := make([]byte, len(binData))
	t.Log(string(actualBinData))
	_, err = binBody.Read(actualBinData) // 从 binBody 中读取数据到 actualBinData
	t.Log(string(actualBinData))
	assert.NoError(t, err, "expected no error while reading from binBody")
	assert.Equal(t, binData, actualBinData, "the byte array content does not match the expected value")

	// Test with nil parameter
	_, _, err = ParseRequestParams(nil)
	assert.Error(t, err, "expected an error for nil parameter")
	assert.Equal(t, "illegal parameter", err.Error(), "expected specific error message")
}

// 测试：BuildUri函数
func TestBuildUri(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		instId   string
		params   map[string]string
		expected string
	}{
		// 测试用例组
		{
			name:     "only instId",
			uri:      "/api/v1",
			instId:   "BTC-USD",
			params:   nil,
			expected: "/api/v1?instId=BTC-USD",
		},
		{
			name:     "only params",
			uri:      "/api/v1",
			instId:   "",
			params:   map[string]string{"type": "spot", "interval": "1h"},
			expected: "/api/v1?interval=1h&type=spot",
		},
		{
			name:     "both instId and params",
			uri:      "/api/v1",
			instId:   "ETH-USD",
			params:   map[string]string{"type": "swap"},
			expected: "/api/v1?instId=ETH-USD&type=swap",
		},
		{
			name:     "no params",
			uri:      "/api/v1",
			instId:   "",
			params:   nil,
			expected: "/api/v1",
		},
		{
			name:     "special characters",
			uri:      "/api/v1",
			instId:   "BTC-USD",
			params:   map[string]string{"filter": "price>100"},
			expected: "/api/v1?instId=BTC-USD&filter=price%3E100",
		},
		{
			name:     "multiple params",
			uri:      "/api/v1",
			instId:   "XRP-USD",
			params:   map[string]string{"from": "2023-01-01", "to": "2023-01-31"},
			expected: "/api/v1?instId=XRP-USD&from=2023-01-01&to=2023-01-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildUri(tt.uri, tt.instId, tt.params)

			// 测试参数顺序无关（因为参数顺序不影响功能）
			expected := strings.Split(tt.expected, "?")
			if len(expected) > 1 {
				params := strings.Split(expected[1], "&")
				gotParams := strings.Split(strings.Split(got, "?")[1], "&")

				if !unorderedSliceEqual(params, gotParams) {
					t.Errorf("params mismatch. expected: %s, got: %s", tt.expected, got)
				}
			} else if got != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, got)
			}
		})
	}
}

// 辅助函数：判断两个字符串切片是否相等（顺序无关）
func unorderedSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	items := make(map[string]int)
	for _, v := range a {
		items[v]++
	}

	for _, v := range b {
		if items[v] == 0 {
			return false
		}
		items[v]--
	}

	return true
}

// 测试：从本地json文件，获取交易产品的信息，比如交易价格精度
func TestGetTickSzFromJson(t *testing.T) {
	ticksz := ""
	instId := "DOGE-USDT-SWAP"
	filename := "account-instruments.json"
	ticksz = GetTickSzFromJson(filename, instId)
	assert.True(t, ticksz == "0.00001")
}
