package myokx

// 行情信息结构，待完善
type BaseTickerResult struct {
	InstId   string `json:"instId"`
	InstType string `json:"instType"`
	Last     string `json:"last"`
	LastSz   string `json:"lastSz"`
	Ts       string `json:"ts"`
	AskPx    string `json:"askPx"`
	AskSz    string `json:"askSz"`
	BidPx    string `json:"bidPx"`
	BidSz    string `json:"bidSz"`
}
type TickerResult struct {
	ApiCodeMsg
	Data []BaseTickerResult `json:"data"`
}

// 深度数据
type BooksData struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Ts   string     `json:"ts"`
}
type Books struct {
	ApiCodeMsg
	Data []BooksData `json:"data"`
}

// Kline:
// 该格式主要用于数据库查询和插入kline
type KlineDataS struct {
	Ts          string `json:"ts" db:"ts"`
	O           string `json:"o" db:"o"`
	H           string `json:"h" db:"h"`
	L           string `json:"l" db:"l"`
	C           string `json:"c" db:"c"`
	Vol         string `json:"vol" db:"vol"`
	VolCcy      string `json:"volCcy" db:"volCcy"`
	VolCcyQuote string `json:"volCcyQuote" db:"volCcyQuote"`
	Confirm     string `json:"confirm" db:"confirm"`
	Tstocst     string `json:"tstocst" db:"tstocst"` // 这个我自己加上去的，方便看时间
}

type KlineData [9]string

// type of kline:
type Klines struct {
	ApiCodeMsg
	Data []KlineData `json:"data"`
}

// Kline: End.
