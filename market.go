package myokx

import (
	"log"
	"time"

	"github.com/gtlyy/mytime"
)

// 获取指定类型的行情信息
func (client *Client) GetTickerType(instType string) (r TickerResult, err error) {
	uri := MARKET_TICKERS + "?instType=" + instType
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return TickerResult{}, err
	}
	return
}

// 获取单个产品行情信息
func (client *Client) GetTicker(instId string) (r TickerResult, err error) {
	uri := MARKET_TICKER + "?instId=" + instId
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return TickerResult{}, err
	}
	return
}

// 获取产品深度 books
func (client *Client) GetBooks(params map[string]string) (r Books, err error) {
	uri := BuildUri(MARKET_BOOKS, "", params)
	if _, err = client.Request(GET, uri, nil, &r); err != nil {
		return Books{}, err
	}
	return
}

// Klines ================================================================================= Start:
// 获取K线数据 (end, start)  0->1  new->old  ；最多可获取最近1440条数据。
func (client *Client) GetKlines(instId string, params map[string]string) (r Klines, err error) {
	uri := BuildUri(MARKET_CANDLES, instId, params)
	_, err = client.Request(GET, uri, nil, &r)
	return
}

// 获取K线历史数据。这个不行的，要看下面的Plus
func (client *Client) GetKlinesHistory(instId string, params map[string]string) (r Klines, err error) {
	uri := BuildUri(MARKET_CANDLES_HISTORY, instId, params)
	_, err = client.Request(GET, uri, nil, &r)
	return
}

// Time to ms. 15m to xxxxxxxxxx
func TimeToMs(s string) (r int64) {
	last := s[len(s)-1:]
	switch last {
	case "m":
		r = StringToInt64(s[:len(s)-1]) * 60 * 1000
	case "H":
		r = StringToInt64(s[:len(s)-1]) * 60 * 60 * 1000
	case "D":
		r = StringToInt64(s[:len(s)-1]) * 24 * 60 * 60 * 1000
	}
	return
}

// 功能：将给定的时间段，切分成每段100个k线。但是，每段的开始部分，回退了1s
func TimeSplit1(end string, start string, bar string) [][2]string {
	r := make([][2]string, 0, 100) // 这个100没有多大意义
	ms := TimeToMs(bar)

	e := StringToInt64(end)
	var s int64
	for {
		s = e - ms*100 // ms * 100 ：时间间隔 100条数据
		if s <= StringToInt64(start) {
			r = append(r, [2]string{Int64ToString(e), Int64ToString(StringToInt64(start) - 1000)}) // -1000 ??? right!
			break
		} else {
			r = append(r, [2]string{Int64ToString(e), Int64ToString(s - 1000)}) // -1000：???
			e = s
		}
	}
	return r
}

// 功能：将给定的时间段，切分成每段100个k线；不过，与上面的不同，这是整齐的分段。
func TimeSplit2(end string, start string, bar string) [][2]string {
	r := make([][2]string, 0, 100) // 这个100没有多大意义
	ms := TimeToMs(bar)

	e := StringToInt64(end)
	var s int64
	for {
		s = e - ms*100 // ms * 100 ：时间间隔 100条数据
		if s <= StringToInt64(start) {
			r = append(r, [2]string{Int64ToString(e), Int64ToString(StringToInt64(start))})
			break
		} else {
			r = append(r, [2]string{Int64ToString(e), Int64ToString(s)})
			e = s
		}
	}
	return r
}

// 获取K线数据 history 加强版1 (End,start]，且不限时间段，不限数量。 0 --> 1 == new to old, but not include the kline(now).
func (client *Client) GetKlinesHistoryPlus1(instId string, params map[string]string) ([]KlineData, error) {
	start := params["before"]
	end := params["after"]

	bar := params["bar"]
	tsArray := TimeSplit1(end, start, bar)

	ks := make([]KlineData, 0, 1440) // 1440 no meanings.
	var err error

	for i := 0; i < len(tsArray); i++ {
		var r Klines
		params["before"] = tsArray[i][1]
		params["after"] = tsArray[i][0]
		uri := BuildUri(MARKET_CANDLES_HISTORY, instId, params)
		_, err = client.Request(GET, uri, nil, &r)
		if (i+1)%20 == 0 {
			time.Sleep(time.Second * 2)
		}
		ks = append(ks, r.Data...)
	}
	return ks, err
}

// 获取K线数据 history 加强版2，且不限时间段，不限数量。 flag=0, (end, start) ; flag=1, (end, start] ; flag=2, [end, start]
func (client *Client) GetKlinesHistoryPlus2(instId string, params map[string]string, flag int) ([]KlineData, error) {
	start := params["before"]
	end := params["after"]
	bar := params["bar"]
	tsArray := TimeSplit2(end, start, bar)

	ks := make([]KlineData, 0, 1440) // 1440 no meanings.
	var err error

	for i := 0; i < len(tsArray); i++ {
		var r Klines
		if i != 0 || flag == 2 {
			params["after"] = mytime.TimeToTs(mytime.TsToTime(tsArray[i][0]).Add(2 * time.Second)) // end
		} else {
			params["after"] = tsArray[i][0]
		}

		if i != len(tsArray)-1 || (i == len(tsArray)-1 && flag != 0) {
			params["before"] = mytime.TimeToTs(mytime.TsToTime(tsArray[i][1]).Add(-2 * time.Second)) // start
		} else {
			params["before"] = tsArray[i][1]
		}

		uri := BuildUri(MARKET_CANDLES_HISTORY, instId, params)
		_, err = client.Request(GET, uri, nil, &r)
		if (i+1)%20 == 0 {
			time.Sleep(time.Second * 2)
		}
		ks = append(ks, r.Data...)
	}
	return ks, err
}

// 计算需要等待到下一个完整K线生成的时间（精确到毫秒级对齐）
func calcNextBarWaitDuration(bar string) time.Duration {
	now := time.Now().UTC()
	var nextTime time.Time

	switch bar {
	case "1m":
		base := 1 * time.Minute
		truncated := now.Truncate(base) // 对齐到最近的1分钟整数倍（向下）
		nextTime = truncated.Add(base)  // 下一个周期的开始时间
	case "15m":
		// 原理：将当前时间对齐到最近的15分钟整数倍（向下取整），然后加15分钟得到下一个周期起点
		// 示例：
		// 当前时间 12:07:25 --> Truncate(15m)=12:00:00 --> Add(15m)=12:15:00
		// 当前时间 12:15:00 --> Truncate(15m)=12:15:00 --> Add(15m)=12:30:00
		base := 15 * time.Minute
		truncated := now.Truncate(base) // 对齐到最近的15分钟整数倍（向下）
		nextTime = truncated.Add(base)  // 下一个周期的开始时间

	case "1H":
		// 对齐到整小时并加1小时
		// 示例：
		// 14:25:30 --> 14:00:00 + 1h = 15:00:00
		base := time.Hour
		truncated := now.Truncate(base)
		nextTime = truncated.Add(base)

	case "4H":
		base := 4 * time.Hour
		truncated := now.Truncate(base)
		nextTime = truncated.Add(base)

	case "1D":
		base := 24 * time.Hour
		truncated := now.Truncate(base)
		nextTime = truncated.Add(base)

	default:
		panic("unsupported bar type: " + bar)
	}

	// 计算时间差并添加缓冲（交易所通常需要3-5秒生成K线）
	waitDuration := nextTime.Sub(now)

	// 添加安全缓冲（根据经验值调整）：
	// - 15m/1H 等短周期：+5秒
	// - 4H/1D 等长周期：+10秒
	buffer := 5 * time.Second
	if bar == "4H" || bar == "1D" {
		buffer = 10 * time.Second
	}

	// 最终需要等待的时间 = 到下一个周期的时间差 + 缓冲时间
	return waitDuration + buffer
}

// 实时获取k线，并发送到通道ch。这个函数主要是为了方便地写入数据库。
func (client *Client) GetKlinesSync(id string, p map[string]string, ch chan []KlineData) {
	lastTs := p["before"]
	// 获取初始k线数据
	r, err := client.GetKlinesHistoryPlus2(id, p, 2)
	if err != nil {
		log.Printf("GetKlinesHistoryPlus2 error: %v", err)
		return
	}
	// 逆序发送初始k线数据
	// log.Println("In GetKlinesSync():", len(r))
	if len(r) > 0 {
		for i := len(r) - 1; i >= 0; i-- {
			ch <- []KlineData{r[i]}
		}
		// 顺便设置 lastTs
		lastTs = r[0][0]
	}

	// 实时监测k线变化
	for {
		t1 := calcNextBarWaitDuration(p["bar"])
		time.Sleep(t1)

		flag := false
		for !flag {
			p2 := NewParams()
			p2["after"] = Int64ToString(time.Now().UnixNano()/1000000 + 10000000)
			p2["before"] = lastTs
			r2, err2 := client.GetKlines(id, p2)
			if err2 != nil {
				log.Printf("In GetKlinesSync(): GetKlines error: %v", err2)
				return
			}
			// 处理(lastTs, nowTs)的k线，未完成的k线（其实就是第0条）忽略。
			for i := len(r2.Data) - 1; i >= 0; i-- {
				if r2.Data[i][8] == "1" {
					flag = true
					lastTs = r2.Data[i][0]
					ch <- []KlineData{r2.Data[i]}
				} else {
					time.Sleep(2 * time.Second)
					break
				}
			}
		}
	}
}

// Klines ================================================================================= End.
