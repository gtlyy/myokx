// 功能：定义客户端数据结构，生成一个客户端。

package myokx

import "github.com/gtlyy/myfun"

type Config struct {
	Endpoint           string
	WSEndpointPublic   string
	WSEndpointPrivate  string
	WSEndpointBUSINESS string
	ApiKey             string
	SecretKey          string
	Passphrase         string
	TimeoutSecond      int
	IsPrint            bool
	I18n               string
	Simulated          bool
	Proxy              string
}

// 提供json文件，生成Config
func NewConfig(confifFile string) Config {
	var config Config
	myfun.ReadJSONFile(confifFile, &config)
	return config
}
