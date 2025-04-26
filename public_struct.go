package myokx

// 系统时间
type Times struct {
	Ts string `json:"ts"`
}
type ServerTime struct {
	Code string  `json:"code"`
	Data []Times `json:"data"`
	Msg  string  `json:"msg"`
}

// 资金费率
type FundingRateData struct {
	InstType        string `json:"instType" db:"instType"`
	InstId          string `json:"instId" db:"instId"`
	FundingRate     string `json:"fundingRate" db:"fundingRate"`
	FundingTime     string `json:"fundingTime" db:"fundingTime"`
	RealizedRate    string `json:"realizedRate" db:"realizedRate"`       // 获取历史资金费率返回这个
	NextFundingRate string `json:"nextFundingRate" db:"nextFundingRate"` // 获取当前资金费率返回这个
}
type FundingRate struct {
	ApiCodeMsg
	Data []FundingRateData `json:"data"`
}
