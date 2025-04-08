// 功能：实现了与API交互的客户端核心，统一了请求入口方法 Request()。

package myokx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type Client struct {
	Config     Config
	HttpClient *http.Client
}

// 提供Config，生成Client
func NewClient(config Config) *Client {
	var client Client
	client.Config = config
	timeout := config.TimeoutSecond
	if timeout <= 0 {
		timeout = 30
	}

	if config.Proxy != "" { // Proxy Start:
		dialer, err := proxy.SOCKS5("tcp", config.Proxy, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}
		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
		client.HttpClient = &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: httpTransport,
		} // Proxy End.
	} else {
		client.HttpClient = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	return &client
}

// 发送请求，然后返回响应
func (client *Client) Request(method string, requestPath string, params, result interface{}) (response *http.Response, err error) {
	config := client.Config

	// 处理url，确保不以'/'结尾。
	endpoint := config.Endpoint
	if strings.HasSuffix(config.Endpoint, "/") {
		endpoint = config.Endpoint[0 : len(config.Endpoint)-1]
	}
	url := endpoint + requestPath

	// 生成json和bin格式的body
	var jsonBody string
	var binBody = bytes.NewReader(make([]byte, 0))
	if params != nil {
		jsonBody, binBody, err = ParseRequestParams(params)
		if err != nil {
			return response, err
		}
	}

	// 生成签名
	timestamp := IsoTime()
	preHash := PreHashString(timestamp, method, requestPath, jsonBody)
	sign, err := HmacSha256Base64Signer(preHash, config.SecretKey)
	if err != nil {
		return response, err
	}

	// 设置请求头
	request, err := http.NewRequest(method, url, binBody)
	if err != nil {
		return response, err
	}
	Headers(request, config, timestamp, sign)

	// 打印请求信息
	if config.IsPrint {
		printRequest(config, request, jsonBody, preHash)
	}

	// 发送请求，得到响应信息
	response, err = client.HttpClient.Do(request)
	if err != nil {
		return response, err
	}
	defer response.Body.Close()

	// 处理响应信息，如下：
	status := response.StatusCode
	message := response.Status
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return response, err
	}

	if config.IsPrint {
		printResponse(status, message, body)
	}

	responseBodyString := string(body)
	response.Header.Add(ResultDataJsonString, responseBodyString)

	// 将响应信息转为需要的结构体，存入result。
	if status >= 200 && status < 300 {
		if body != nil && result != nil {
			err := JsonBytes2Struct(body, result)
			if err != nil {
				return response, err
			}
		}
		return response, nil
	} else if status >= 400 || status <= 500 {
		errMsg := "Http error(400~500) result: status=" + IntToString(status) + ", message=" + message + ", body=" + responseBodyString
		fmt.Println(errMsg)
		if body != nil {
			err := errors.New(errMsg)
			return response, err
		}
	} else {
		fmt.Println("Http error result: status=" + IntToString(status) + ", message=" + message + ", body=" + responseBodyString)
		return response, errors.New(message)
	}
	return response, nil
}

// 打印请求信息
func printRequest(config Config, request *http.Request, body string, preHash string) {
	if config.SecretKey != "" {
		fmt.Println("  Secret-Key: " + config.SecretKey)
	}
	fmt.Println("  Request(" + IsoTime() + "):")
	fmt.Println("\tUrl: " + request.URL.String())
	fmt.Println("\tMethod: " + strings.ToUpper(request.Method))
	if len(request.Header) > 0 {
		fmt.Println("\tHeaders: ")
		for k, v := range request.Header {
			if strings.Contains(k, "Ok-") {
				k = strings.ToUpper(k)
			}
			fmt.Println("\t\t" + k + ": " + v[0])
		}
	}
	fmt.Println("\tBody: " + body)
	if preHash != "" {
		fmt.Println("  PreHash: " + preHash)
	}
}

// 打印服务器响应信息
func printResponse(status int, message string, body []byte) {
	fmt.Println("  Response(" + IsoTime() + "):")
	statusString := strconv.Itoa(status)
	message = strings.Replace(message, statusString, "", -1)
	message = strings.Trim(message, " ")
	fmt.Println("\tStatus: " + statusString)
	fmt.Println("\tMessage: " + message)
	fmt.Println("\tBody: " + string(body))
}
