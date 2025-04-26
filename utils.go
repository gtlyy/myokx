// 功能：提供请求验证、数据转换、时间处理等相关函数。
package myokx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gtlyy/myfun"
	"github.com/gtlyy/mytime"
)

// 设置函数别名，方便打字
var (
	ISONow            = mytime.ISONow
	ISOCSTToTs        = mytime.ISOCSTToTs
	IntToString       = myfun.IntToString
	StringToFloat64   = myfun.StringToFloat64
	Float64ToString   = myfun.Float64ToString
	JsonBytes2Struct  = myfun.JsonBytes2Struct
	Struct2JsonString = myfun.Struct2JsonString
	WriteJSONFile     = myfun.WriteJSONFile
	ReadJSONFile      = myfun.ReadJSONFile
	IfError           = myfun.IfError
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

// 请求验证 ========================================================================= End

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

// 实现请求参数处理，生成合适的url。deepseek r1 提供。 ok
func BuildUri(uri string, instId string, optionalParams map[string]string) string {
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

// Http 函数 ======================================================================== End.

// Cal Price  ======================================================================== Start:
// 根据持仓均价、收益率、方向、fee_open、fee_close，计算平仓价格
func CalCloseStrPx(p *Positions, goal, fee_open, fee_close float64) string {
	fOpen := StringToFloat64(p.AvgPx)
	var direct float64
	if p.PosSide == "long" {
		direct = 1
	} else {
		direct = -1
	}
	rf := fOpen * (goal + direct + fee_open) / (direct - fee_close)
	tickSz := GetTickSzQuick(p.InstId)
	if tickSz == "" {
		tickSz = p.AvgPx
	}
	return Float64ToString(rf, myfun.CountStrFloat(tickSz))
}

// Cal Price  ======================================================================== End.
