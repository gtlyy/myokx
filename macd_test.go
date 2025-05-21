// 功能：测试macd策略
package myokx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gtlyy/mytime"
)

// 测试：CalMacd()
func TestCalMacd(t *testing.T) {
	var m MyMacdClass
	cs, err := maria.QueryClose("DOGEUSDTSWAP1H", "2021-01-01T00:00:00.000Z", mytime.ISONowCST())
	assert.True(t, err == nil)
	t.Log(len(cs))
	closeData := m.CsToFloat64(cs)
	m.Init(closeData, 12, 26, 9)
	m.CalMacd()
	t.Log(m.Hist[len(m.Hist)-3:])
}

// 测试：myMacd类及二次金叉
func TestMyMacdClass(t *testing.T) {
	fast, slow, signal := 12, 26, 9
	var m MyMacdClass

	r := maria.createData()
	N := len(r) - 33 // 这个33有点问题！！！应该是不需要减去33的。
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}

	m.Init(closeData, fast, slow, signal)
	m.CalMacd()
	t.Log(len(m.Close), len(m.Dif))
	gold := 0
	for j := 0; j < N-1; j++ {
		if m.CrossOverTwice(j) {
			gold = gold + 1
		}
	}
	t.Log(gold)
}

// 测试：myMacd UpHist ......
func TestUp3Hist(t *testing.T) {
	fast, slow, signal := 12, 26, 9
	var m MyMacdClass

	r := maria.createData()
	N := len(r)
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}
	m.Init(closeData, fast, slow, signal)
	m.CalMacd()
	t.Log(len(m.Close), len(m.Dif))
	gold := 0
	for j := 0; j < N-1; j++ {
		if m.Up3Hist(j) {
			gold = gold + 1
		}
	}
	t.Log(gold)
}

// 测试：固定fast,slow,signal && 固定盈亏比
func TestMyMacdTrade1(t *testing.T) {
	fast, slow, signal := 9, 22, 8
	var m MyMacdClass

	usdt0 := 1000.0
	usdt := usdt0
	coin := 0.0
	perF := 0.95 // 每次投入的比例
	price := 0.0
	fee := 0.0
	per := 0.01
	feeRate := 0.0005
	canBuy := true
	canSell := false
	goal := 0.01 * 1.0
	fail := -0.01 * 0.9

	r := maria.createData()
	N := len(r)
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}

	m.Init(closeData, fast, slow, signal)
	m.CalMacd()
	t.Log("Close和Dif的数量：", len(m.Close), len(m.Dif))
	t.Log("Coin本身的上涨率：", closeData[0], closeData[len(closeData)-1], Decimal(100*(closeData[len(closeData)-1]-closeData[0])/closeData[0]))
	gold, gray, grayF := 0, 0, 0
	for j := 33; j < N-1; j++ {
		if canBuy && m.CrossOverTwice(j) && usdt > 300 {
			per = Decimal(perF * usdt / m.Close[j])
			// t.Log(per)
			gold = gold + 1
			price = m.Close[j]
			fee = per * price * feeRate
			coin = coin + per
			usdt = usdt - per*price - fee
			canSell = true
			canBuy = false
			// t.Log(usdt, coin)
			// t.Log("开仓")
			// t.Log(fee)
		} else if canSell {
			if (m.Close[j]-price)/price >= goal && m.CrossUnder(j) {
				gray = gray + 1
				price = m.Close[j]
				fee = per * price * feeRate
				usdt = usdt + coin*price - fee
				coin = coin - per
				canBuy = true
				canSell = false
				// t.Log("止盈")
				// t.Log(usdt, coin)
			} else if (m.Close[j]-price)/price <= fail ||
				((m.Close[j]-price)/price <= fail*0.7 && m.CrossUnder(j)) { // 0.7 死叉，且跌到fail的70%
				grayF = grayF + 1
				price = m.Close[j]
				fee = per * price * feeRate
				usdt = usdt + coin*price - fee
				coin = coin - per
				canBuy = true
				canSell = false
				// t.Log("止损")
				// t.Log(usdt, coin)
			}
		}
	}
	t.Log("二次金叉、止盈、止损、胜率：", gold, gray, grayF, 100.0*(gray)/gold, "%")
	t.Log("参数、盈利目标、止损目标：", fast, slow, signal, goal, fail)
	t.Log("最终：", Decimal(usdt), coin, Decimal(100*(usdt+coin*closeData[len(closeData)-1]-usdt0)/usdt0), "%")

	// t.Log(TsNow())
	// mark it : TsNow --> TsNow - 200 * time_bar(15m?)
}

// 测试：交易2，固定12,26,9，寻找合适的盈亏比
func TestMyMacdTrade2(t *testing.T) {
	fast, slow, signal := 12, 26, 9
	var m MyMacdClass

	usdt0 := 1000.0
	usdt := usdt0
	coin := 0.0
	perF := 0.95
	price := 0.0
	fee := 0.0
	per := 0.01
	feeRate := 0.0005
	canBuy := true
	canSell := false
	gold, gray, grayF := 0, 0, 0
	// goal := 0.01 * 2.0
	// fail := -0.01 * 2.0

	goal := make([]float64, 0)
	fail := make([]float64, 0)
	for i := 0.5; i < 3.1; {
		goal = append(goal, i*0.01)
		fail = append(fail, -i*0.01)
		i = i + 0.1
	}

	r := maria.createData()
	N := len(r)
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}

	m.Init(closeData, fast, slow, signal)
	m.CalMacd()
	t.Log("Close和Dif的数量: ", len(m.Close), len(m.Dif))
	t.Log("Coin本身的上涨率: ", closeData[0], closeData[len(closeData)-1],
		Decimal(100*(closeData[len(closeData)-1]-closeData[0])/closeData[0]))

	// best := usdt
	for a := 0; a < len(goal); a++ {
		for b := 0; b < len(fail); b++ {
			// reset:
			usdt = usdt0
			coin = 0.0
			price = 0.0
			fee = 0.0
			per = 0.01
			feeRate = 0.0005
			canBuy = true
			canSell = false
			gold, gray, grayF = 0, 0, 0
			// start cal:
			for j := 33; j < N-1; j++ {
				if canBuy && m.CrossOverTwice(j) && usdt > 300 {
					per = Decimal(perF * usdt / m.Close[j])
					gold = gold + 1
					price = m.Close[j]
					fee = per * price * feeRate
					coin = coin + per
					usdt = usdt - per*price - fee
					canSell = true
					canBuy = false
					// t.Log(usdt, coin)
					// t.Log("开仓")
					// t.Log(fee)
				} else if canSell {
					if (m.Close[j]-price)/price >= goal[a] && m.CrossUnder(j) {
						gray = gray + 1
						price = m.Close[j]
						fee = per * price * feeRate
						usdt = usdt + coin*price - fee
						coin = coin - per
						canBuy = true
						canSell = false
						// t.Log("止盈")
						// t.Log(usdt, coin)
					} else if (m.Close[j]-price)/price <= fail[b] ||
						((m.Close[j]-price)/price <= fail[b]*0.7 && m.CrossUnder(j)) {
						grayF = grayF + 1
						price = m.Close[j]
						fee = per * price * feeRate
						usdt = usdt + coin*price - fee
						coin = coin - per
						canBuy = true
						canSell = false
						// t.Log("止损")
						// t.Log(usdt, coin)
					}
				}
			}
			// t.Log(gold, gray)
			// usdt = usdt + coin*price
			if usdt+coin*price > usdt0*1.00 &&
				-goal[a]/fail[b] >= 1.0 &&
				-goal[a]/fail[b] <= 1.2 {
				t.Log("二次金叉、止盈、止损、胜率：", gold, gray, grayF, 100.0*(gray)/gold, "%")
				t.Log("止盈目标、止损目标、最终：", Decimal(goal[a]*100), Decimal(fail[b]*100), Decimal(usdt), coin,
					Decimal(100*(usdt+coin*price-usdt0)/usdt0), "%")
				// best = usdt
			}
		}
	}
	t.Log("Macd参数：", fast, slow, signal)

	// t.Log(TsNow())
	// mark it : TsNow --> TsNow - 200 * time_bar(15m?)
}

// 测试：5个for，查找合适的参数
// 注意：这个只是寻找best收益的参数，但有时，我们可能需要最合适的。看下一个函数。
func TestMyMacdTrade3(t *testing.T) {
	// fast, slow, signal := 12, 26, 9
	var m MyMacdClass

	usdt0 := 1000.0
	usdt := usdt0
	coin := 0.0
	perF := 0.95
	price := 0.0
	fee := 0.0
	per := 0.01
	feeRate := 0.0005
	canBuy := true
	canSell := false
	gold, gray, grayF := 0, 0, 0
	best := 0.0
	// goal := 0.01 * 2.0
	// fail := -0.01 * 2.0

	goal := make([]float64, 0)
	fail := make([]float64, 0)
	fast := make([]int, 0)
	slow := make([]int, 0)
	signal := make([]int, 0)
	for i := 0.5; i < 3.0; {
		goal = append(goal, i*0.01)
		i = i + 0.1
	}
	for i := 0.5; i < 3.0; {
		fail = append(fail, -i*0.01)
		i = i + 0.1
	}
	for i := 8; i <= 15; i++ {
		fast = append(fast, i)
	}
	for i := 20; i <= 28; i++ {
		slow = append(slow, i)
	}
	for i := 8; i <= 15; i++ {
		signal = append(signal, i)
	}

	r := maria.createData()
	N := len(r)
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}

	// m.Init(closeData, fast, slow, signal)
	// m.CalMacd()
	// t.Log(len(m.Close), len(m.Dif))
	t.Log(closeData[0], closeData[len(closeData)-1],
		Decimal(100*(closeData[len(closeData)-1]-closeData[0])/closeData[0]))

	best = usdt0
	for x := 0; x < len(fast); x++ {
		for y := 0; y < len(slow); y++ {
			for z := 0; z < len(signal); z++ {
				for a := 0; a < len(goal); a++ {
					for b := 0; b < len(fail) && -fail[b] <= goal[a]*1.5; b++ {
						// reset:
						m.Init(closeData, fast[x], slow[y], signal[z])
						m.CalMacd()
						usdt = 1000.0
						coin = 0.0
						price = 0.0
						fee = 0.0
						per = 0.01
						feeRate = 0.0005
						canBuy = true
						canSell = false
						gold, gray, grayF = 0, 0, 0

						// start cal:
						for j := 33; j < N-1; j++ {
							if canBuy && m.CrossOverTwice(j) && usdt > 300 {
								per = Decimal(perF * usdt / m.Close[j])
								gold = gold + 1
								price = m.Close[j]
								fee = per * price * feeRate
								coin = coin + per
								usdt = usdt - per*price - fee
								canSell = true
								canBuy = false
								// t.Log(usdt, coin)
								// t.Log("开仓")
								// t.Log(fee)
							} else if canSell {
								if (m.Close[j]-price)/price >= goal[a] && m.CrossUnder(j) {
									gray = gray + 1
									price = m.Close[j]
									fee = per * price * feeRate
									usdt = usdt + coin*price - fee
									coin = coin - per
									canBuy = true
									canSell = false
									// t.Log("止盈")
									// t.Log(usdt, coin)
								} else if (m.Close[j]-price)/price <= fail[b] ||
									((m.Close[j]-price)/price <= fail[b]*0.7 && m.CrossUnder(j)) {
									grayF = grayF + 1
									price = m.Close[j]
									fee = per * price * feeRate
									usdt = usdt + coin*price - fee
									coin = coin - per
									canBuy = true
									canSell = false
									// t.Log("止损")
									// t.Log(usdt, coin)
								}
							}
						}
						// t.Log(gold, gray)
						if usdt+coin*price >= best &&
							-goal[a]/fail[b] >= 1.0 &&
							-goal[a]/fail[b] <= 1.2 {
							t.Log(gold, gray, grayF, 100.0*(gray)/gold, "%")
							t.Log(fast[x], slow[y], signal[z], Decimal(goal[a]*100), Decimal(fail[b]*100),
								Decimal(usdt), coin, Decimal(100*(usdt+coin*price-usdt0)/usdt0))
							best = usdt + coin*price
						}
					}
				}
			}
		}
	}

	// t.Log(TsNow())
	// mark it : TsNow --> TsNow - 200 * time_bar(15m?)
}

// 测试：寻找合适的参数，与3略有区别，3是寻找best的，这个4是寻找大于某个目标值的。
func TestMyMacdTrade4(t *testing.T) {
	// fast, slow, signal := 12, 26, 9
	var m MyMacdClass

	usdt0 := 1000.0
	usdt := usdt0
	coin := 0.0
	perF := 0.95
	price := 0.0
	fee := 0.0
	per := 0.01
	feeRate := 0.0005
	canBuy := true
	canSell := false
	gold, gray, grayF := 0, 0, 0
	// goal := 0.01 * 2.0
	// fail := -0.01 * 2.0

	goal := make([]float64, 0)
	fail := make([]float64, 0)
	fast := make([]int, 0)
	slow := make([]int, 0)
	signal := make([]int, 0)
	for i := 0.5; i < 3.0; {
		goal = append(goal, i*0.01)
		i = i + 0.1
	}
	for i := 0.5; i < 3.0; {
		fail = append(fail, -i*0.01)
		i = i + 0.1
	}
	for i := 8; i <= 15; i++ {
		fast = append(fast, i)
	}
	for i := 20; i <= 28; i++ {
		slow = append(slow, i)
	}
	for i := 8; i <= 15; i++ {
		signal = append(signal, i)
	}

	r := maria.createData()
	N := len(r)
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}

	// m.Init(closeData, fast, slow, signal)
	// m.CalMacd()
	// t.Log(len(m.Close), len(m.Dif))
	t.Log(closeData[0], closeData[len(closeData)-1],
		Decimal(100*(closeData[len(closeData)-1]-closeData[0])/closeData[0]))

	for x := 0; x < len(fast); x++ {
		for y := 0; y < len(slow); y++ {
			for z := 0; z < len(signal); z++ {
				for a := 0; a < len(goal); a++ {
					for b := 0; b < len(fail) && -fail[b] <= goal[a]*1.5; b++ {
						// reset:
						m.Init(closeData, fast[x], slow[y], signal[z])
						m.CalMacd()
						usdt = 1000.0
						coin = 0.0
						price = 0.0
						fee = 0.0
						per = 0.01
						feeRate = 0.0005
						canBuy = true
						canSell = false
						gold, gray, grayF = 0, 0, 0

						// start cal:
						for j := 33; j < N-1; j++ {
							if canBuy && m.CrossOverTwice(j) && usdt > 300 {
								per = Decimal(perF * usdt / m.Close[j])
								gold = gold + 1
								price = m.Close[j]
								fee = per * price * feeRate
								coin = coin + per
								usdt = usdt - per*price - fee
								canSell = true
								canBuy = false
								// t.Log(usdt, coin)
								// t.Log("开仓")
								// t.Log(fee)
							} else if canSell {
								if (m.Close[j]-price)/price >= goal[a] && m.CrossUnder(j) {
									gray = gray + 1
									price = m.Close[j]
									fee = per * price * feeRate
									usdt = usdt + coin*price - fee
									coin = coin - per
									canBuy = true
									canSell = false
									// t.Log("止盈")
									// t.Log(usdt, coin)
								} else if (m.Close[j]-price)/price <= fail[b] ||
									((m.Close[j]-price)/price <= fail[b]*0.7 && m.CrossUnder(j)) {
									grayF = grayF + 1
									price = m.Close[j]
									fee = per * price * feeRate
									usdt = usdt + coin*price - fee
									coin = coin - per
									canBuy = true
									canSell = false
									// t.Log("止损")
									// t.Log(usdt, coin)
								}
							}
						}
						// t.Log(gold, gray)
						if usdt+coin*price >= usdt0*2.0 &&
							-goal[a]/fail[b] >= 1.0 &&
							-goal[a]/fail[b] <= 1.2 {
							t.Log(gold, gray, grayF, 100.0*(gray)/gold, "%")
							t.Log(fast[x], slow[y], signal[z], Decimal(goal[a]*100), Decimal(fail[b]*100),
								Decimal(usdt), coin, Decimal(100*(usdt+coin*price-usdt0)/usdt0))
						}
					}
				}
			}
		}
	}

	// t.Log(TsNow())
	// mark it : TsNow --> TsNow - 200 * time_bar(15m?)
}

// 测试：Up3Hist策略的交易表现
func TestUp3HistStrategy(t *testing.T) {
	fast, slow, signal := 12, 26, 9
	var m MyMacdClass

	// 初始化参数
	usdt0 := 1000.0
	usdt := usdt0
	coin := 0.0
	positionOpen := false
	entryPrice := 0.0
	feeRate := 0.0005
	goal := 0.01 * 0.6
	fail := -0.01 * 0.5

	// 获取测试数据
	r := maria.createData()
	N := len(r)
	closeData := make([]float64, N)
	for i := 0; i < N; i++ {
		closeData[i] = StringToFloat64(r[i].C)
	}

	// 初始化MACD计算
	m.Init(closeData, fast, slow, signal)
	m.CalMacd()

	t.Log("========= Up3Hist策略测试 =========")
	t.Logf("数据总数 | Close: %d, Hist: %d", len(m.Close), len(m.Hist))
	t.Logf("标的涨幅: %.2f%%", 100*(closeData[len(closeData)-1]-closeData[0])/closeData[0])

	// 交易统计
	tradeCount := 0
	winCount := 0
	lossCount := 0

	// 遍历所有数据点
	for j := 3; j < len(m.Hist); j++ { // 需要至少4个点来判断连续3次上涨
		currentPrice := m.Close[j]

		// 开仓条件：出现Up3Hist信号且无持仓
		if !positionOpen && m.Up3Hist(j) && m.OverZero(j) {
			// 计算买入数量（95%资金使用率）
			positionSize := (usdt * 0.95) / currentPrice
			fee := positionSize * currentPrice * feeRate

			// 执行买入
			usdt -= positionSize*currentPrice + fee
			coin += positionSize
			entryPrice = currentPrice
			positionOpen = true
			tradeCount++

			t.Logf("[%s] 开仓 | 价格: %.4f | 数量: %.4f", r[j].Tstocst, entryPrice, positionSize)

			// 平仓条件：持仓时出现Hist下跌
		} else if positionOpen && ((currentPrice-entryPrice)/entryPrice >= goal || ((currentPrice-entryPrice)/entryPrice <= fail || ((currentPrice-entryPrice)/entryPrice <= fail*0.7 && m.CrossUnder(j)))) {
			// 计算卖出金额
			fee := coin * currentPrice * feeRate

			// 执行卖出
			usdt += coin*currentPrice - fee
			coin = 0
			positionOpen = false

			// 统计盈亏
			if currentPrice > entryPrice {
				winCount++
			} else {
				lossCount++
			}

			t.Logf("[%s] 平仓 | 价格: %.4f | 盈利: %.2f%%",
				r[j].Tstocst,
				currentPrice,
				100*(currentPrice-entryPrice)/entryPrice)
		}
	}

	// 最终统计
	finalValue := usdt + coin*m.Close[len(m.Close)-1]
	roi := 100 * (finalValue - usdt0) / usdt0

	t.Log("\n========= 最终结果 =========")
	t.Logf("初始资金: %.2f", usdt0)
	t.Logf("最终价值: %.2f", finalValue)
	t.Logf("总收益率: %.2f%%", roi)
	t.Logf("交易次数: %d (胜率 %.1f%%)", tradeCount, 100*float64(winCount)/float64(winCount+lossCount))
	t.Logf("盈利次数: %d", winCount)
	t.Logf("亏损次数: %d", lossCount)
	t.Log("============================")

	// 添加断言验证基本逻辑
	assert.Greater(t, len(m.Hist), 100, "应有足够的历史数据")
	if tradeCount > 0 {
		assert.NotZero(t, winCount+lossCount, "应完成完整的交易")
	}
}
