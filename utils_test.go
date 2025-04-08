package myokx

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

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

// 数据转换 ========================================================================= Start
// 测试 int ---> string
func TestInt2tring(t *testing.T) {
	var a int = 123456789
	s_right := "123456789"
	ss := Int2String(a)
	assert.True(t, s_right == ss)
}

// 测试 int64 ---> string
func TestInt642tring(t *testing.T) {
	var a int64 = 123456789
	s_right := "123456789"
	ss := Int642String(a)
	assert.True(t, s_right == ss)
}

// 测试 JsonBytes ---> Struct
func TestJsonBytes2Struct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	jsonBytes := []byte(`{"name": "Alice", "age": 30}`)
	var result Person
	err := JsonBytes2Struct(jsonBytes, &result)
	assert.NoError(t, err, "expected no error during unmarshalling")
	assert.Equal(t, Person{Name: "Alice", Age: 30}, result, "expected result does not match")
}

// 测试 Struct ---> JsonString
func TestStruct2JsonString(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	person := Person{Name: "Alice", Age: 30}
	expectedJson := `{"name":"Alice","age":30}`
	jsonString, err := Struct2JsonString(person)
	assert.NoError(t, err, "expected no error during marshalling")
	assert.JSONEq(t, expectedJson, jsonString, "expected JSON string does not match")
}

// 数据转换 ========================================================================= End

// 时间函数 ========================================================================= Start
// 测试EpochTime的返回格式是否正确
func TestEpochTime(t *testing.T) {
	epochTime := EpochTime()
	assert.Contains(t, epochTime, ".", "epochTime should contain a dot")
	parts := strings.Split(epochTime, ".")
	assert.Len(t, parts, 2, "epochTime should have two parts")
	_, err := strconv.Atoi(parts[0])
	assert.NoError(t, err, "the first part of epochTime should be a valid integer")
	_, err = strconv.Atoi(parts[1])
	assert.NoError(t, err, "the second part of epochTime should be a valid integer")
}

// 测试：1540365300000 -> 2018-10-24 15:15:00 +0800 CST
func TestLongTimeToUTC8(t *testing.T) {
	longTime := int64(1540365300000) // Corresponds to 2018-10-24 15:15:00 UTC+8
	expectedTime := time.Date(2018, 10, 24, 15, 15, 0, 0, time.FixedZone("CST", 8*3600))
	result := LongTimeToUTC8(longTime)
	assert.Equal(t, expectedTime, result, "the converted time should match the expected time")
}

// 测试："2018-11-18T16:51:55.933Z" -> 2018-11-18 16:51:55.000000933 +0000 UTC
func TestIsoToTime(t *testing.T) {
	iso := "2018-11-18T16:51:55.933Z"
	expectedTime := time.Date(2018, 11, 18, 16, 51, 55, 933000000, time.UTC)
	result, err := IsoToTime(iso)
	assert.NoError(t, err, "expected no error while parsing ISO string")
	assert.Equal(t, expectedTime, result, "the parsed time does not match the expected time")
}

// 测试 IsoTime() 获取当前时间 iso
func TestIsoTime(t *testing.T) {
	iso1 := IsoTime()
	b := strings.ContainsAny(iso1, "TZ")
	assert.True(t, b)
}

// 时间函数 ========================================================================= End

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

// 测试：构造请求参数，版本1
// go test -run ^TestBuildParams$ okex -v
func TestBuildParams(t *testing.T) {
	params := NewParams()
	params["depth"] = "200"
	params["conflated"] = "0"
	url := BuildParams("/api/futures/v3/products/BTC-USD-0310/book", params)
	str_right := "/api/futures/v3/products/BTC-USD-0310/book?conflated=0&depth=200"
	assert.True(t, url == str_right)
}

// 测试：构造请求参数，版本2
func TestBuildParams2(t *testing.T) {
	params := NewParams()
	params["depth"] = "200"
	params["conflated"] = "0"
	url := BuildParams2("uri?instId=BTC-USDT&", params)
	str_right := "uri?instId=BTC-USDT&conflated=0&depth=200"
	assert.True(t, url == str_right)
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

// 测试：实现三目运算：a == b ? c : d
func TestT3O(t *testing.T) {
	type testStruct struct{ field int }

	tests := []struct {
		name      string
		condition bool
		trueVal   interface{}
		falseVal  interface{}
		want      interface{}
	}{
		{"true returns int", true, 42, 0, 42},
		{"false returns string", false, "apple", "banana", "banana"},
		{"true with nil", true, nil, "not-nil", nil},
		{"false with nil", false, "not-nil", nil, nil},
		{"different types", true, 3.14, "pi", 3.14},
		{"struct comparison", true, testStruct{5}, testStruct{10}, testStruct{5}},
		{"pointer comparison",
			true,
			&testStruct{7},
			&testStruct{9},
			&testStruct{7}},
		{"typed nil pointer",
			false,
			(*testStruct)(nil),
			(*testStruct)(nil),
			(*testStruct)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := T3O(tt.condition, tt.trueVal, tt.falseVal)

			// 特殊处理nil比较
			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v (%T)", got, got)
				}
				return
			}

			// 使用反射处理复杂类型比较
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("T3O() = %v (%T), want %v (%T)",
					got, got, tt.want, tt.want)
			}
		})
	}
}
