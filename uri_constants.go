package myokx

const (
	// 交易账户 account
	// 获取持仓信息
	ACCOUNT_POSITIONS = "/api/v5/account/positions"
	// 可交易产品基础信息
	ACCOUNT_INSTRUMENTS = "/api/v5/account/instruments"
	// 余额
	ACCOUNT_BALANCE = "/api/v5/account/balance"
	// 交易流水：七天
	ACCOUNT_BILLS = "/api/v5/account/bills"
	// 交易流水：三个月
	ACCOUNT_BILLS_ARCHIVE = "/api/v5/account/bills-archive"
	// 账户配置
	ACCOUNT_CONFIG = "/api/v5/account/config"
	// 设置账户交易模式
	SET_POSITIONMODE = "/api/v5/account/set-position-mode"
	// 获取杠杆倍数
	ACCOUNT_LEVERINFO = "/api/v5/account/leverage-info"
	// 设置杠杠
	SET_LEVER = "/api/v5/account/set-leverage"
	// 获取最大可下单数量
	Max_SIZE = "/api/v5/account/max-size"
	// 获取最大可用余额/保证金
	Max_Avail_SIZE = "/api/v5/account/max-avail-size"
	// 获取交易手续费率
	TRADE_FEE = "/api/v5/account/trade-fee"

	// 资金账户 asset
	// 获取平台所有币种列表
	ASSET_CURRENCIES = "/api/v5/asset/currencies"
	// 获取资金账户余额
	ASSET_BALANCES = "/api/v5/asset/balances"

	// 公共数据 public
	// 服务器时间
	OKEX_TIME_URI = "/api/v5/public/time"
	// 可交易产品基础信息
	PUBLIC_INSTRUMENTS = "/api/v5/public/instruments"
	// 当前资金费率
	SWAP_FUNDING_RATE = "/api/v5/public/funding-rate"
	// 历史资金费率
	SWAP_HISTORY_FUNDING_RATE = "/api/v5/public/funding-rate-history"

	// 行情数据 market
	// 获取单个产品行情信息
	MARKET_TICKER = "/api/v5/market/ticker"
	// Kline
	MARKET_CANDLES         = "/api/v5/market/candles"
	MARKET_CANDLES_HISTORY = "/api/v5/market/history-candles"
	// books
	MARKET_BOOKS = "/api/v5/market/books"

	// 下单、查询订单
	TRADE_ORDER = "/api/v5/trade/order"
	// 获取未成交订单
	TRADE_PENDING_ORDER = "/api/v5/trade/orders-pending"
	// 撤单
	TRADE_CANCEL_ORDER = "/api/v5/trade/cancel-order"
)
