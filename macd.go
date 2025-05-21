package myokx

import (
	// . "myrabbitmq"
	"fmt"
	"strconv"

	tb "github.com/markcheno/go-talib"
	// . "mytime"
	// "reflect"
	// "testing"
	// "time"
	// "github.com/stretchr/testify/assert"
)

// 构建 MyMacdClass 类
type MyMacdClass struct {
	Date   []string // ts or str???
	Close  []float64
	Dif    []float64
	Dea    []float64
	Hist   []float64
	Fast   int
	Slow   int
	Signal int
}

// 初始化
func (mymacd *MyMacdClass) Init(closeData []float64, fast int, slow int, signal int) {
	mymacd.Fast, mymacd.Slow, mymacd.Signal = fast, slow, signal
	mymacd.Close = make([]float64, len(closeData))
	copy(mymacd.Close, closeData)
}

// 添加数据    add date???
func (mymacd *MyMacdClass) Append(close1 float64) {
	mymacd.Close = append(mymacd.Close, close1)
}

// 计算Macd
func (mymacd *MyMacdClass) CalMacd() {
	mymacd.Dif, mymacd.Dea, mymacd.Hist = tb.Macd(mymacd.Close, mymacd.Fast, mymacd.Slow, mymacd.Signal)
	for i, v := range mymacd.Hist {
		mymacd.Hist[i] = 2 * v
	}
}

// close数组约定 old ---> new
// n 表示从new往old方向的位移
// pos 标识确切位置 0,1,2,3...N-1
// 判断金叉
func (mymacd *MyMacdClass) CrossOver(pos int) bool {
	N := len(mymacd.Dif)
	if N < 3 || pos < 1 || pos > N-1 {
		return false
	}
	return mymacd.Dif[pos-1] <= mymacd.Dea[pos-1] && mymacd.Dif[pos] > mymacd.Dea[pos]
}

// 判断死叉
func (mymacd *MyMacdClass) CrossUnder(pos int) bool {
	N := len(mymacd.Dif)
	if N < 3 || pos < 1 || pos > N-1 {
		return false
	}
	return mymacd.Dif[pos-1] >= mymacd.Dea[pos-1] && mymacd.Dif[pos] < mymacd.Dea[pos]
}

// 判断是否0轴上方
func (mymacd *MyMacdClass) OverZero(pos int) bool {
	N := len(mymacd.Dif)
	if N < 3 || pos < 1 || pos > N-1 {
		return false
	}
	return mymacd.Dif[pos] > 0 && mymacd.Dea[pos] > 0
}

// 判断是否0轴下方
func (mymacd *MyMacdClass) UnderZero(pos int) bool {
	N := len(mymacd.Dif)
	if N < 3 || pos < 1 || pos > N-1 {
		return false
	}
	return mymacd.Dif[pos] < 0 && mymacd.Dea[pos] < 0
}

// 判断是否上涨 Hist
func (mymacd *MyMacdClass) UpHist(pos int) bool {
	N := len(mymacd.Dif)
	if N < 4 || pos < 1 || pos > N-1 {
		return false
	}
	return mymacd.Hist[pos] > mymacd.Hist[pos-1]
}

// 判断是否连续上涨3个Hist
func (mymacd *MyMacdClass) Up3Hist(pos int) bool {
	N := len(mymacd.Dif)
	if N < 4 || pos < 1 || pos > N-1 {
		return false
	}
	return mymacd.UpHist(pos) && mymacd.UpHist(pos-1) && mymacd.UpHist(pos-2) && !mymacd.UpHist(pos-3)
}

// 判断低位二次金叉
func (mymacd *MyMacdClass) CrossOverTwice(pos int) bool {
	// 如当前不是金叉，或不在0轴之下，返回
	if !mymacd.CrossOver(pos) || !mymacd.UnderZero(pos) {
		return false
	}
	// N := len(mymacd.Dif)
	flag := false
	// 回溯，判断此前是否发生过金叉
	for i := pos - 1; i > 1; i-- {
		// 如中途不在0轴之下，返回
		if !mymacd.UnderZero(i) {
			break
		}
		if mymacd.CrossOver(i) {
			flag = true
			break
		}
	}
	return flag
}

// 判断高位二次死叉
func (mymacd *MyMacdClass) CrossUnderTwice(pos int) bool {
	// 如当前不是死叉，或不在0轴之上，返回
	if !mymacd.CrossUnder(pos) || !mymacd.OverZero(pos) {
		return false
	}
	flag := false
	// 回溯，判断此前是否发生过死叉
	for i := pos - 1; i > 0; i-- {
		// 如中途不在0轴之上，返回
		if !mymacd.OverZero(i) {
			break
		}
		if mymacd.CrossUnder(i) {
			flag = true
			break
		}
	}
	return flag
}

// 将 cs:[]string --->  closeData:[]float64
func (mymacd *MyMacdClass) CsToFloat64(cs []string) []float64 {
	closeData := make([]float64, len(cs))
	for i, v := range cs {
		closeData[i] = StringToFloat64(v)
	}
	return closeData
}

// todo:背离，看Candlestick  my_macd.py

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
