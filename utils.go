// 功能：提供请求验证、数据转换、时间处理等相关函数。
package myokx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 请求验证 ========================================================================= Start
// 生成签名用的 message
func PreHashString(timestamp string, method string, requestPath string, body string) string {
	return timestamp + strings.ToUpper(method) + requestPath + body
}

// 生成签名，即 OK-ACCESS-SIGN
func HmacSha256Base64Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "Error: HmacSha256Base64Signer(): Call Write() error.", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// 生成iso格式时间字符串
func IsoTime() string {
	utcTime := time.Now().UTC()
	return utcTime.Format("2006-01-02T15:04:05.999Z")
}

// 请求验证 ========================================================================= End

// 数据转换 ========================================================================= Start
// int to string
func Int2String(arg int) string {
	return strconv.Itoa(arg)
}

// int64 to string
func Int642String(arg int64) string {
	return strconv.FormatInt(int64(arg), 10)
}

// json byte array to struct
func JsonBytes2Struct(jsonBytes []byte, result interface{}) error {
	err := json.Unmarshal(jsonBytes, result)
	return err
}

// json string to struct
func JsonString2Struct(jsonString string, result interface{}) error {
	return JsonBytes2Struct([]byte(jsonString), result)
}

// struct to json string
func Struct2JsonString(structt interface{}) (jsonString string, err error) {
	data, err := json.Marshal(structt)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 更多函数见 myfun
// 数据转换 ========================================================================= End

// 时间函数 ========================================================================= Start
// 获取当前时间戳: 1521221737.376
func EpochTime() string {
	millisecond := time.Now().UnixNano() / 1000000
	epoch := strconv.Itoa(int(millisecond))
	epochBytes := []byte(epoch)
	epoch = string(epochBytes[:10]) + "." + string(epochBytes[10:])
	return epoch
}

// 1540365300000 -> 2018-10-24 15:15:00 +0800 CST
func LongTimeToUTC8(longTime int64) time.Time {
	timeString := Int64ToString(longTime)
	sec := timeString[0:10]
	nsec := timeString[10:]
	utcTime := time.Unix(StringToInt64(sec), StringToInt64(nsec))
	return utcTime.In(time.FixedZone("CST", 8*3600))
}

// 1540365300000 -> 2018-10-24 15:15:00
func LongTimeToUTC8Format(longTime int64) string {
	return LongTimeToUTC8(longTime).Format("2006-01-02 15:04:05")
}

// "2018-11-18T16:51:55.933Z" -> 2018-11-18 16:51:55.000000933 +0000 UTC
func IsoToTime(iso string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999Z", iso)
}

// 更多函数见 mytime
// 时间函数 ========================================================================= End

// Http 函数 ======================================================================== Start
// 生成json和bin格式的body
func ParseRequestParams(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params) // data : []byte
	if err != nil {
		return "", nil, errors.New("json convert string error")
	}
	jsonBody := string(data)
	binBody := bytes.NewReader(data)
	return jsonBody, binBody, nil
}

// 设置请求头
func Headers(request *http.Request, config Config, timestamp string, sign string) {
	request.Header.Add(OK_ACCESS_KEY, config.ApiKey)
	request.Header.Add(OK_ACCESS_SIGN, sign)
	request.Header.Add(OK_ACCESS_TIMESTAMP, timestamp)
	request.Header.Add(OK_ACCESS_PASSPHRASE, config.Passphrase)
	if config.Simulated {
		request.Header.Add(X_SIMULATE_TRADING, "1")
	}
	request.Header.Add(ACCEPT, APPLICATION_JSON)
	request.Header.Add(CONTENT_TYPE, APPLICATION_JSON_UTF8)
	request.Header.Add(COOKIE, LOCALE+config.I18n)
}

// 生成一个map: {string:string}
func NewParams() map[string]string {
	return make(map[string]string)
}

/*
	  build http get request params, and order
	  eg:
	    params := make(map[string]string)
		params["bb"] = "222"
		params["aa"] = "111"
		params["cc"] = "333"
	  return string: eg: aa=111&bb=222&cc=333
*/
func BuildOrderParams(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return urlParams.Encode()
}

// 生成请求参数，版本1：
func BuildParams(requestPath string, params map[string]string) string {
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return requestPath + "?" + urlParams.Encode()
}

// 生成请求参数，版本2 ：与BuildParams()相比，只少了一个 '?'
func BuildParams2(requestPath string, params map[string]string) string {
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return requestPath + urlParams.Encode()
}

// ok. 结合上面的BuildParams2，构建一个可用的 uri
// 注意: 有时需要instId参数，且紧接着uri；有时是不需要的，不需要时，设为 ""
func BuildUri(uri string, instId string, optionalParams map[string]string) string {
	if instId != "" {
		uri = uri + "?instId=" + instId
		if len(optionalParams) > 0 {
			uri = uri + "&"
			uri = BuildParams2(uri, optionalParams)
		}
	} else if len(optionalParams) > 0 {
		uri = BuildParams(uri, optionalParams)
	}
	return uri
}

// 实现请求参数处理，生成合适的url。deepseek r1 提供。 ok
func BuildUri2(uri string, instId string, optionalParams map[string]string) string {
	// 使用strings.Builder高效构建字符串
	var builder strings.Builder
	builder.WriteString(uri)

	hasParams := false

	// 优先处理instId参数
	if instId != "" {
		builder.WriteString("?instId=")
		builder.WriteString(url.QueryEscape(instId))
		hasParams = true
	}

	// 处理可选参数
	if len(optionalParams) > 0 {
		params := url.Values{}
		for k, v := range optionalParams {
			params.Add(k, v)
		}
		encoded := params.Encode()

		if encoded != "" {
			if hasParams {
				builder.WriteByte('&') // 已有参数时追加
			} else {
				builder.WriteByte('?') // 首个参数
			}
			builder.WriteString(encoded)
			hasParams = true
		}
	}

	return builder.String()
}

func GetResponseDataJsonString(response *http.Response) string {
	return response.Header.Get(ResultDataJsonString)
}
func GetResponsePageJsonString(response *http.Response) string {
	return response.Header.Get(ResultPageJsonString)
}

/*
ternary operator biz extension
*/
func T3Ox(err error, value interface{}) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	return value, nil
}

/*
return decimalism string 9223372036854775807 -> "9223372036854775807"
*/
func Int64ToString(arg int64) string {
	return strconv.FormatInt(arg, 10)
}

func Float64ToString(arg float64, n int) string {
	return strconv.FormatFloat(arg, 'f', n, 64)
}

func IntToString(arg int) string {
	return strconv.Itoa(arg)
}

func StringToInt64(arg string) int64 {
	value, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func StringToFloat64(arg string) float64 {
	value, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func StringToInt(arg string) int {
	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	} else {
		return value
	}
}

/*
call fmt.Println(...)
*/
func FmtPrintln(flag string, info interface{}) {
	fmt.Print(flag)
	if info != nil {
		jsonString, err := Struct2JsonString(info)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jsonString)
	} else {
		fmt.Println("{}")
	}
}

func GetInstrumentIdUri(uri, instrumentId string) string {
	return strings.Replace(uri, "{instrument_id}", instrumentId, -1)
}

func GetCurrencyUri(uri, currency string) string {
	return strings.Replace(uri, "{currency}", currency, -1)
}

func GetInstrumentIdOrdersUri(uri, instrumentId string, order_client_id string) string {
	uri = strings.Replace(uri, "{instrument_id}", instrumentId, -1)
	uri = strings.Replace(uri, "{order_client_id}", order_client_id, -1)
	return uri
}

// 实现三目运算：a == b ? c : d
func T3O(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}
