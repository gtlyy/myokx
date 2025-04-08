package myokx

// 返回结果必包含这两个基本信息
type ApiCodeMsg struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// 交易账户 account： (对应，下面还有资金账户 asset )
// Account ========================================================= Start:
// 获取当前账户可交易产品的信息列表
type AccountInstrumentsData struct {
	InstType     string `json:"instType"`     // 产品类型
	InstId       string `json:"instId"`       // 产品ID
	Uly          string `json:"uly"`          // 标的指数
	InstFamily   string `json:"instFamly"`    // 交易品种
	BaseCcy      string `json:"baseCcy"`      // 交易货币币种
	QuoteCcy     string `json:"quoteCcy"`     // 计价货币币种
	SettleCcy    string `json:"settleCcy"`    // 盈亏结算和保证金币种
	CtVal        string `json:"ctVal"`        // 合约面值
	CtMult       string `json:"ctMult"`       // 合约乘数
	CtValCcy     string `json:"ctValCcy"`     // 合约面值计价币种
	OptType      string `json:"optType"`      // 期权类型
	Stk          string `json:"stk"`          // 行权价格
	ListTime     string `json:"listTime"`     // 上线时间
	ExpTime      string `json:"expTime"`      // 下线时间
	Lever        string `json:"lever"`        // 最大杠杆倍数
	TickSz       string `json:"tickSz"`       // 下单价格精度
	LotSz        string `json:"lotSz"`        // 下单数量精度
	MinSz        string `json:"minSz"`        // 最小下单数量
	CtType       string `json:"ctType"`       // 合约类型
	State        string `json:"state"`        // 产品状态
	RuleType     string `json:"ruleType"`     // 交易规则类型
	MaxLmtSz     string `json:"maxLmtSz"`     // 限价单的单笔最大委托数量
	MaxMktSz     string `json:"maxMktSz"`     // 市价单的单笔最大委托数量
	MaxLmtAmt    string `json:"maxLmtAmt"`    // 限价单的单笔最大美元价值
	MaxMktAmt    string `json:"maxMktAmt"`    // 市价单的单笔最大美元价值
	MaxTwapSz    string `json:"maxTwapSz"`    // 时间加权单的单笔最大委托数量
	MaxIcebergSz string `json:"maxIcebergSz"` // 冰山委托的单笔最大委托数量
	MaxTriggerSz string `json:"maxTriggerSz"` // 计划委托委托的单笔最大委托数量
	MaxStopSz    string `json:"maxStopSz"`    // 止盈止损市价委托的单笔最大委托数量
}
type AccountInstrumentsResult struct {
	ApiCodeMsg
	Data []AccountInstrumentsData `json:"data"`
}

// 账户余额
type AccountBalanceData struct {
	UTime       string                      `json:"uTime"`       // 账户信息的更新时间
	TotalEq     string                      `json:"totalEq"`     // 美金层面权益
	IsoEq       string                      `json:"isoEq"`       // 美金层面逐仓仓位权益
	AdjEq       string                      `json:"adjEq"`       // 美金层面有效保证金
	OrdFroz     string                      `json:"ordFroz"`     // 美金层面全仓挂单占用保证金
	Imr         string                      `json:"imr"`         // 美金层面占用保证金
	Mmr         string                      `json:"mmr"`         // 美金层面维持保证金
	BorrowFroz  string                      `json:"borrowFroz"`  // 账户美金层面潜在借币占用保证金率
	MgnRatio    string                      `json:"mgnRatio"`    // 美金层面保证金率
	NotionalUsd string                      `json:"notionalUsd"` // 仓位美金价值
	Upl         string                      `json:"upl"`         // 账户层面未实现盈亏
	Details     []AccountBalanceDataDetails `json:"details"`     // 各币种资产详细信息
}
type AccountBalanceDataDetails struct {
	Ccy        string `json:"ccy"`        // 币种
	Eq         string `json:"eq"`         // 币种总权益
	CashBal    string `json:"cashBal"`    // 币种余额
	UTime      string `json:"uTime"`      // 时间戳
	IsoEq      string `json:"isoEq"`      // 币种逐仓仓位权益
	AvailEq    string `json:"availEq"`    // 可用保证金
	DisEq      string `json:"disEq"`      // 美金层面币种折算权益
	FixedBal   string `json:"fixedBal"`   // 抄底宝、逃顶宝功能的币种冻结金额
	AvailBal   string `json:"availBal"`   // 可用余额
	FronzenBal string `json:"fronzenBal"` // 币种占用余额
	OrdFrozen  string `json:"ordFrozen"`  // 挂单冻结数量
	Liab       string `json:"liab"`       // 币种负债额
	Upl        string `json:"upl"`        // 未实现盈亏
	UplLiab    string `json:"uplLiab"`    // 由于仓位未实现亏损导致的负债
	CrossLiab  string `json:"crossLiab"`  // 币种全仓负债额
	IsoLiab    string `json:"isoLiab"`    // 币种逐仓负债额
	MgnRatio   string `json:"mgnRatio"`   // 币种全仓保证金率
	Imr        string `json:"imr"`        // 币种维度全仓占用保证金
	Mmr        string `json:"mmr"`        // 币种维度全仓维持保证金
	Interest   string `json:"interest"`   // 应扣未扣利息
	// 未完待续。。。。。。
}
type AccountBalanceResult struct {
	ApiCodeMsg
	Data []AccountBalanceData `json:"data"`
}

// 账户流水
type Bills struct {
	InstId   string `json:"instId"`
	InstType string `json:"instType"`
	MgnMode  string `json:"mgnMode"`
	Ts       string `json:"ts"`
	Type     string `json:"type"`
	SubType  string `json:"subType"`
	Sz       string `json:"sz"`
	Pnl      string `json:"pnl"`    // 收益
	Fee      string `json:"fee"`    // 手续费
	BillId   string `json:"billId"` // 账单ID
}
type BillsResult struct {
	ApiCodeMsg
	Data []Bills `json:"data"`
}

// 持仓信息
type Positions struct {
	InstId   string `json:"instId"`
	InstType string `json:"instType"`
	MgnMode  string `json:"mgnMode"`
	PosId    string `json:"posId"`
	PosSide  string `json:"posSide"`
	Pos      string `json:"pos"`
	AvailPos string `json:"availPos"`
	AvgPx    string `json:"avgPx"`
	Upl      string `json:"upl"`
	UplRatio string `json:"uplRatio"`
	Lever    string `json:"lever"`
	Last     string `json:"last"`
	CTime    string `json:"cTime"`
	UTime    string `json:"uTime"`
	Imr      string `json:"imr"`
	Mmr      string `json:"mmr"`
	Margin   string `json:"margin"`
	MgnRatio string `json:"mgnRatio"`
}
type PositionsResult struct {
	ApiCodeMsg
	Data []Positions `json:"data"`
}

// 账户配置
type AccConfig struct {
	Uid             string `json:"uid"`
	AcctLv          string `json:"acctLv"`
	PosMode         string `json:"posMode"`
	AutoLoan        bool   `json:"autoLoan"`
	GreeksType      string `json:"greeksType"`
	Level           string `json:"level"`
	Label           string `json:"label"`
	Ip              string `json:"ip"`
	Perm            string `json:"perm"`
	LiquidationGear string `json:"liquidationGear"`
}
type ConfigResult struct {
	ApiCodeMsg
	Data []AccConfig `json:"data"`
}

// 设置持仓模式时提交的参数
type PositionModeParam struct {
	PosMode string `json:"posMode"`
}

// 查询杠杠倍数
type LeverInfo struct {
	InstId  string `json:"instId"`
	MgnMode string `json:"mgnMode"`
	PosSide string `json:"posSide"`
	Lever   string `json:"lever"`
}
type LeverInfoResult struct {
	ApiCodeMsg
	Data []LeverInfo `json:"data"`
}

// 设置杠杠时提交的参数
type LeverParam struct {
	InstId  string `json:"instId"`
	Ccy     string `json:"ccy"`
	Lever   string `json:"lever"`
	MgnMode string `json:"mgnMode"` // isolated, cross
	PosSide string `json:"posSide"` // long, short
}

// 设置杠杠返回的数据结构
type BaseLeverResult struct {
	Lever   string `json:"lever"`
	MgnMode string `json:"mgnMode"` // isolated, cross
	InstId  string `json:"instId"`
	PosSide string `json:"posSide"` // long, short
}
type LeverResult struct {
	ApiCodeMsg
	Data []BaseLeverResult `json:"data"`
}

// 获取最大可下单数量：对应下单时的 sz 字段
type MaxSize struct {
	InstId  string `json:"instId"`
	Ccy     string `json:"ccy"`
	MaxBuy  string `json:"maxBuy"`
	MaxSell string `json:"maxSell"`
}
type MaxSizeResult struct {
	ApiCodeMsg
	Data []MaxSize `json:"data"`
}

// 获取最大可用余额/保证金
type MaxAvailSize struct {
	InstId    string `json:"instId"`
	AvailBuy  string `json:"availBuy"`
	AvailSell string `json:"availSell"`
}
type MaxAvailSizeResult struct {
	ApiCodeMsg
	Data []MaxAvailSize `json:"data"`
}

// 获取当前账户交易手续费率
type TradeFee struct {
	Level    string `json:"level"`
	Maker    string `json:"maker"`
	Taker    string `json:"taker"`
	MakerU   string `json:"makerU"` // USDT 合约
	TakerU   string `json:"takerU"`
	Delivery string `json:"delivery"` // 交割
}
type TradeFeeResult struct {
	ApiCodeMsg
	Data []TradeFee `json:"data"`
}

// Account ========================================================= End.

// 资金账户：
// Asset =========================================================== Start:
// 获取币种列表
type Currencies struct {
	Ccy         string `json:"ccy"`
	Name        string `json:"name"`
	Chain       string `json:"chain"`
	CanDep      bool   `json:"canDep"`
	CanWd       bool   `json:"canWd"`
	CanInternal bool   `json:"canInternal"`
	MinWd       string `json:"minWd"`
	MinFee      string `json:"minFee"`
	MaxFee      string `json:"maxFee"`
}
type CurrenciesResult struct {
	ApiCodeMsg
	Data []Currencies `json:"data"`
}

// 获取资金账户余额
type AssetBalances struct {
	Ccy       string `json:"ccy"`
	Bal       string `json:"bal"`
	FrozenBal string `json:"frozenBal"`
	AvailBal  string `json:"availBal"`
}
type AssetBalancesResult struct {
	ApiCodeMsg
	Data []AssetBalances `json:"data"`
}

// Asset =========================================================== End:
