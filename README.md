# myokx
Learn how to trade use okx api.
学会如何使用okx的api进行交易。

## 交易之前，记住如下血的教训：
- 交易第一要义：生存。
- 交易第二要义：无对冲，不交易。
- 交易第三要义：如赌博，不交易。

## Install
```
go mod init
go mod tidy
```

## config.json如下：
```
{
  "Endpoint": "https://www.okx.com",
  "WSEndpointPublic": "wss://wseea.okx.com:8443/ws/v5/public",
  "WSEndpointPrivate": "wss://wseea.okx.com:8443/ws/v5/private",
  "WSEndpointBusiness": "wss://wseea.okx.com:8443/ws/v5/business",
  "ApiKey": "your ApiKey",
  "SecretKey": "your SecretKey",
  "Passphrase": "your Passphrase",
  "TimeoutSecond": 0,
  "IsPrint": false,
  "I18n": "",
  "Simulated": false,
  "Proxy": "127.0.0.1:1080"
}
```

## tag v0.0.1
go test -v -run TestGetAccountBalance

## Licence
MIT License - see LICENSE for more details
