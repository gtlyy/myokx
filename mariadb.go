package myokx

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"

	// "log"
	"math/rand"
	"time"

	// "mysqlgo"
	"github.com/gtlyy/mytime"

	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// struct to map[string]string
func StructToMap(s interface{}) (d map[string]string) {
	typeInfo := reflect.TypeOf(s)
	valInfo := reflect.ValueOf(s)
	num := typeInfo.NumField()
	d = make(map[string]string, num)
	for i := 0; i < num; i++ {
		key := typeInfo.Field(i).Name
		val := valInfo.Field(i).String()
		d[key] = val
	}
	return
}

// 构建 MyMariaDB 类
type MyMariaDBClass struct {
	db    *sqlx.DB
	Table string
	Sql   string
	// DbName
}

// 连接
func (mydb *MyMariaDBClass) Connect(user, passwd, ip, port, d string) error {
	uri := user + ":" + passwd + "@tcp(" + ip + ":" + port + ")" + "/" + d
	db, err := sqlx.Open("mysql", uri)
	mydb.db = db
	return err
}

// 快速初始化
func (mydb *MyMariaDBClass) Init(user, passwd, ip, port, db string) (err error) {
	return mydb.Connect(user, passwd, ip, port, db)
}

// 增
func (mydb *MyMariaDBClass) InsertDict(table string, dict map[string]string) (err error) {
	keys := make([]string, 0, len(dict))
	values := make([]string, 0, len(dict))
	updates := make([]string, 0, len(dict))
	for k, v := range dict {
		keys = append(keys, k)
		values = append(values, "\""+v+"\"")
		updates = append(updates, k+"="+"\""+v+"\"")
	}
	sql := "INSERT INTO " + table + "(" + strings.Join(keys, ",") + ")" + " VALUES(" + strings.Join(values, ",") + ")"
	sql = sql + " ON DUPLICATE KEY UPDATE " + strings.Join(updates, ",")
	// fmt.Println(sql)
	_, err = mydb.db.Exec(sql)
	return
}

// 增：0 -> 1  == old -> new
func (mydb *MyMariaDBClass) InsertStructs(table string, items interface{}) (err error) {
	valInfo := reflect.ValueOf(items)
	for i := valInfo.Len() - 1; i >= 0; i-- {
		item := valInfo.Index(i).Interface()
		dict := StructToMap(item)
		err = mydb.InsertDict(table, dict)
		// test start
		if err != nil {
			fmt.Println("In InsertStructs error:", err)
		}
		// test end.
	}
	return
}

// 删
func (mydb *MyMariaDBClass) Delete(table string, n int, direct int, orderby string) (err error) {
	sql := ""
	// del 1st : 0 to 1
	if direct == 1 {
		sql = "DELETE FROM " + table + " ORDER BY " + orderby + " LIMIT " + IntToString(n)
	} else if direct == -1 {
		sql = "DELETE FROM " + table + " ORDER BY " + orderby + " DESC LIMIT " + IntToString(n)
	}
	_, err = mydb.db.Exec(sql)
	return
}

// 删：d=1,del [0->]
func (mydb *MyMariaDBClass) DeleteNum(table string, n int, direct int, orderby string) (err error) {
	// orderby = "ts"
	err = mydb.Delete(table, n, direct, orderby)
	return
}

// 删除表
func (mydb *MyMariaDBClass) DeleteTable(table string) (err error) {
	sql := "DROP TABLE IF EXISTS " + table
	_, err = mydb.db.Exec(sql)
	return
}

// 查
func (mydb *MyMariaDBClass) Query(r interface{}, sql string) (err error) {
	return mydb.db.Select(r, sql)
}

// 查: 正向或反向(-1),查询n条数据
func (mydb *MyMariaDBClass) QueryNum(r interface{}, table string, direct int, limit int) (err error) {
	sql := `SELECT * FROM ` + table
	if direct == -1 {
		sql = sql + ` ORDER BY ts DESC` // 1 --> 0
	} else {
		sql = sql + ` ORDER BY ts ASC` // 0 -> 1
	}
	sql = sql + ` LIMIT ` + IntToString(limit)
	return mydb.Query(r, sql)
}

// 查：根据时间范围 [a, b] 查询
func (mydb *MyMariaDBClass) QueryStartEnd(r interface{}, table string, start string, end string) (err error) {
	end1 := mytime.ISOCSTToTs(end)
	start1 := mytime.ISOCSTToTs(start)
	sql := `SELECT * FROM ` + table
	sql = sql + ` WHERE ts<=` + end1 + ` AND ts>=` + start1
	sql = sql + ` ORDER BY ts ASC` // 0 -> 1
	// sql = sql + ` ORDER BY ts DESC` // 1 -> 0
	return mydb.Query(r, sql)
}

// 获取close价格，以用于计算macd等。
func (mydb *MyMariaDBClass) QueryClose(table string, start string, end string) ([]string, error) {
	var r []KlineDataS
	err := mydb.QueryStartEnd(&r, table, start, end)
	if err != nil {
		fmt.Println("Error in QueryClose().")
	}
	cs := make([]string, 0, len(r))
	for _, v := range r {
		cs = append(cs, v.C)
	}
	return cs, err
}

// 查：随机n条
func (mydb *MyMariaDBClass) QueryRand(r interface{}, table string, n int) (err error) {
	// 随机获取起始数据条
	query1 := "SELECT ts FROM " + table + " ORDER BY RAND() LIMIT 1"
	err = mydb.Query(r, query1)
	IfError("In QueryRand 1.", err)
	// log.Println(r)

	// 获取起始数据条的 ts 值
	rValue := reflect.ValueOf(r)
	rValue = reflect.Indirect(rValue) // 获取指向实际值的指针
	startTs := rValue.Index(0).FieldByName("Ts").String()

	// 清空 r 切片
	rValue.SetLen(0)

	// 然后，获取连续 n 条
	query2 := "SELECT * FROM " + table + " WHERE ts >= " + startTs + " LIMIT " + IntToString(n)
	err = mydb.Query(r, query2)
	IfError("In QueryRand 2.", err)
	// log.Println(r)
	return err
}

// 查询表是否存在
func (mydb *MyMariaDBClass) CheckTableExists(table string) bool {
	sql_check_table := "SHOW TABLES LIKE" + "'" + table + "'"
	var r []string
	err := mydb.Query(&r, sql_check_table)
	IfError("In CheckTableExists():", err)
	return len(r) > 0
}

// 创建表
func (mydb *MyMariaDBClass) CreateTable(table string) error {
	createSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			ts varchar(20) NOT NULL,
			o varchar(20) NOT NULL,
			h varchar(20) NOT NULL,
			l varchar(20) NOT NULL,
			c varchar(20) NOT NULL,
			vol varchar(45) NOT NULL,
			volCcy varchar(45) NOT NULL,
			volCcyQuote varchar(45) NOT NULL,
			confirm varchar(20) NOT NULL,
			tstocst varchar(45) NOT NULL,
			PRIMARY KEY (ts)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci
	`, table)

	// return mydb.db.Select(nil, createSQL)
	// 使用 Exec 方法执行
	_, err := mydb.db.Exec(createSQL)
	return err
}

// get(from server) --> save(to maria)   [a, b)
// id, bar, start, end, table, history:是否使用历史接口
// history is false: 获取当前变化中的kline，confirm=0，这个不宜插入到数据库
func (mydb *MyMariaDBClass) InsertUseIdAndBar(c *Client, id string, bar string, start string, end string, table string, history bool) error {
	p := NewParams()

	if end == "-1" { // update to now.
		p["after"] = mytime.TsNow() // end
	} else {
		p["after"] = mytime.ISOCSTToTs(end)
	}

	if start == "-1" { // auto update
		var r0 []KlineDataS
		mydb.QueryNum(&r0, table, -1, 1)
		if len(r0) <= 0 {
			fmt.Println("Empty......")
			p["before"] = mytime.ISOCSTToTs("2021-01-01T00:00:00Z") //空表的话，默认从2021年开始
		} else {
			p["before"] = r0[0].Ts // start
		}
	} else {
		p["before"] = mytime.ISOCSTToTs(start)
	}

	p["bar"] = bar
	var r []KlineData
	var err error
	if history {
		r, err = c.GetKlinesHistoryPlus1(id, p)
		if err != nil {
			fmt.Println("Error: In InsertUseIdAndBar(): 1")
			return err
		}
	} else {
		_r, err1 := c.GetKlines(id, p)
		if err1 != nil {
			fmt.Println("Error: In InsertUseIdAndBar(): 1")
			return err1
		}
		r = _r.Data
	}

	// log.Println("Insert klines: ", len(r))
	if len(r) == 0 {
		return err
	}

	var ks []KlineDataS
	var k KlineDataS
	ks = make([]KlineDataS, 0, len(r))
	for i := 0; i < len(r); i++ {
		k.Ts = r[i][0]                      //ts
		k.O = r[i][1]                       //open
		k.H = r[i][2]                       //high
		k.L = r[i][3]                       //low
		k.C = r[i][4]                       //close
		k.Vol = r[i][5]                     //vol
		k.VolCcy = r[i][6]                  //volCcy
		k.VolCcyQuote = r[i][7]             // volCcyQuote
		k.Confirm = r[i][8]                 // confirm
		k.Tstocst = mytime.TsToISOCST(k.Ts) // tstocst
		// test start
		s1, _ := Struct2JsonString(k)
		fmt.Println(s1)
		// test end.
		ks = append(ks, k)
	}
	err = mydb.InsertStructs(table, ks)
	return err
}

// id := "DOGE-USDT-SWAP" bar := "1H"  --> table: DOGEUSDTSWAP1H
func IdAndBarToTable(id, bar string) string {
	return strings.Replace(id+bar, "-", "", -1)
}

// 生成测试数据
func (mydb *MyMariaDBClass) createData() []KlineDataS {
	end := "2025-03-17T00:00:00.000Z"
	// end := ISONowCST()
	start := "2025-01-01T00:00:00.000Z"
	id := "DOGE-USDT"
	bar := "15m"
	table := IdAndBarToTable(id, bar)
	var r []KlineDataS
	err := mydb.QueryStartEnd(&r, table, start, end)
	IfError("Error in createData2()", err)
	return r
}

// 生成测试数据 for game
func (mydb *MyMariaDBClass) CreateGameData2() (r []KlineDataS) {
	id := "DOGE-USDT"
	bar := "15m"
	table := IdAndBarToTable(id, bar)
	n := 720
	err := mydb.QueryRand(&r, table, n)
	IfError("Error in createGameData2()", err)
	return r
}

// 获取csv文件内容
func GetLinesCsv(filename string, passTitle bool) ([][]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取 CSV 文件并存储到 result 中,跳过第一行标题
	result, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 如果文件不为空,则返回除第一行外的所有行
	if len(result) > 0 {
		if passTitle {
			return result[1:], nil
		} else {
			return result, nil
		}
	}

	return nil, nil
}

// 生成测试数据 for tradegame  混合大A和加密币。3：返回 股票名称
func (mydb *MyMariaDBClass) CreateTradeGameData3(a bool, b bool, bar string) (r []KlineDataS, stock string, name string) {
	// fmt.Println(a, b, bar)
	// A
	d := make(map[string]string)
	idsA := make([]string, 0)
	stockss, err := GetLinesCsv("stocks.csv", true)
	IfError("Error: In CreateTradeGameData3(), Call GetLinesCsv(): ", err)
	for _, v := range stockss {
		stock1 := v[0]
		d[stock1] = v[1]
		// name := v[1]
		// date1 := v[2]
		idsA = append(idsA, stock1)
	}

	idsB := []string{"DOGE-USDT-SWAP", "ETC-USDT-SWAP", "BTC-USDT-SWAP", "KAITO-USDT-SWAP", "TRUMP-USDT-SWAP"}
	idsAll := make([]string, 0)
	if a && b {
		idsAll = append(idsA, idsB...)
	} else if a && !b {
		idsAll = idsA
	} else if !a && b {
		idsAll = idsB
	}

	// 使用随机数生成索引来选择一个随机的 id
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(idsAll))
	id := idsAll[randomIndex]
	stock = id

	// Big A or coin
	suffix := "SH" // 上证
	if id[0] >= '0' && id[0] <= '9' {
		bar = "1D"
		name = d[stock]
		if id[0] == '0' || id[0] == '3' {
			suffix = "SZ" // 深证
		}
	} else {
		suffix = "" // coin
	}

	table := IdAndBarToTable(stock+suffix, bar)
	n := 720
	err = mydb.QueryRand(&r, table, n)
	IfError("Error in CreateTradeGameData3()", err)

	return
}

/*
SET @a = (SELECT ts FROM DOGEUSDTSWAP1H ORDER BY RAND() LIMIT 1);
select @a;
SELECT * FROM DOGEUSDTSWAP1H WHERE  ts >= (select @a) limit 10;
*/

// InsertUseIdAndBar2 ： 重构 InsertUseIdAndBar
// get(from server) --> 检查数据 --> save(to maria)   [a, b)
// id, bar, start, end, table, history:是否使用历史接口
// history is false: 获取当前变化中的kline，confirm=0，这个不宜插入到数据库
func (mydb *MyMariaDBClass) InsertUseIdAndBar2(c *Client, id string, bar string, start string, end string, table string, history bool) error {
	p := NewParams()

	// 检查表，没有就新建一个
	if !mydb.CheckTableExists(table) {
		mydb.CreateTable(table)
	}

	// 查询最新的ts
	var r0 []KlineDataS
	bEmpty := false
	newestTs := ""
	mydb.QueryNum(&r0, table, -1, 1)
	if len(r0) <= 0 {
		bEmpty = true
	} else {
		newestTs = r0[0].Ts // start
	}

	// 准备参数：截止时间
	if end == "-1" { // update to now.
		p["after"] = mytime.TsNow() // end
	} else {
		p["after"] = mytime.ISOCSTToTs(end)
	}

	// 准备参数：开始时间
	if start == "-1" { // auto update
		if bEmpty {
			fmt.Println("Empty......")
			p["before"] = mytime.ISOCSTToTs("2021-01-01T00:00:00Z") //空表的话，默认从2021年开始
		} else {
			p["before"] = r0[0].Ts // start
		}
	} else {
		p["before"] = mytime.ISOCSTToTs(start)
	}

	p["bar"] = bar
	var r []KlineData
	var err error
	if history {
		// r, err = c.GetKlinesHistoryPlus(id, p)
		r, err = c.GetKlinesHistoryPlus1(id, p)
		if err != nil {
			fmt.Println("Error: In InsertUseIdAndBar(): 1")
			return err
		}
	} else {
		_r, err1 := c.GetKlines(id, p)
		if err1 != nil {
			fmt.Println("Error: In InsertUseIdAndBar(): 1")
			return err1
		}
		r = _r.Data
	}

	if len(r) == 0 {
		return err
	}

	var ks []KlineDataS
	var k KlineDataS
	ks = make([]KlineDataS, 0, len(r))
	for i := 0; i < len(r); i++ {
		k.Ts = r[i][0] //ts
		// 更新的数据才写入数据库。
		if k.Ts <= newestTs {
			continue
		}
		k.O = r[i][1]                       //open
		k.H = r[i][2]                       //high
		k.L = r[i][3]                       //low
		k.C = r[i][4]                       //close
		k.Vol = r[i][5]                     //vol
		k.VolCcy = r[i][6]                  //volCcy
		k.VolCcyQuote = r[i][7]             // volCcyQuote
		k.Confirm = r[i][8]                 // confirm
		k.Tstocst = mytime.TsToISOCST(k.Ts) // tstocst
		// test start
		// s1, _ := Struct2JsonString(k)
		// fmt.Println(s1)
		// test end.
		ks = append(ks, k)
	}

	log.Println("Insert klines: ", len(ks))
	err = mydb.InsertStructs(table, ks)
	return err
}

// 实时同步K线数据（有sleep，阻塞式函数，最好goroutine运行）
func (maria *MyMariaDBClass) InsertUseIdAndBarSync(c *Client, id string, bar string, start string) {
	table := IdAndBarToTable(id, bar)
	// 首次同步两次确保数据完整
	maria.InsertUseIdAndBar2(c, id, bar, start, "-1", table, true)
	maria.InsertUseIdAndBar2(c, id, bar, "-1", "-1", table, true)

	for {
		// 计算需要等待的时间
		sleepDuration := calcNextBarWaitDuration(bar)

		// 等待到下一个K线周期
		time.Sleep(sleepDuration)

		// 执行同步
		maria.InsertUseIdAndBar2(c, id, bar, "-1", "-1", table, true)
	}
}

// 实时同步K线数据（有sleep，使用goroutine运行，并根据最新的GetKlinsSync函数进行修改。）
func (mydb *MyMariaDBClass) InsertUseIdAndBarSyncCh(c *Client, id string, bar string, start string) {
	// 表，没有就新建一个。
	table := IdAndBarToTable(id, bar)
	if !mydb.CheckTableExists(table) {
		mydb.CreateTable(table)
	}

	// 准备参数：开始时间
	p := NewParams()
	var r0 []KlineDataS
	mydb.QueryNum(&r0, table, -1, 1)
	// log.Println("len(r0)=", len(r0))
	if start == "-1" {
		if len(r0) <= 0 {
			fmt.Println("Empty......")
			p["before"] = mytime.ISOCSTToTs("2021-01-01T00:00:00Z") //空表的话，默认从2021年开始
		} else {
			p["before"] = r0[0].Ts
		}
	} else {
		p["before"] = mytime.ISOCSTToTs(start)
	}

	log.Println("start time:", mytime.TsToISOCST(p["before"]))

	// 准备参数：截止时间、bar
	p["after"] = mytime.TsNow()
	p["bar"] = bar

	// 准备通道，用于获取k线
	ch := make(chan []KlineData, 10)

	// 用于处理获取的k线
	go func() {
		for klineData := range ch {
			var ks []KlineDataS
			ks = make([]KlineDataS, 0)
			for _, kline := range klineData {
				// 写入数据库：
				var k KlineDataS
				k.Ts = kline[0] //ts
				// 跳过已有数据
				if len(r0) > 0 {
					if k.Ts <= r0[0].Ts {
						continue
					}
				}
				k.O = kline[1]                      //open
				k.H = kline[2]                      //high
				k.L = kline[3]                      //low
				k.C = kline[4]                      //close
				k.Vol = kline[5]                    //vol
				k.VolCcy = kline[6]                 //volCcy
				k.VolCcyQuote = kline[7]            // volCcyQuote
				k.Confirm = kline[8]                // confirm
				k.Tstocst = mytime.TsToISOCST(k.Ts) // tstocst
				fmt.Println("to insert:", kline[0], mytime.TsToISOCST(kline[0]))
				ks = append(ks, k)
			}
			if len(ks) > 0 {
				log.Println("Insert klines: ", len(ks))
				mydb.InsertStructs(table, ks)
			}
		}
	}()

	// 去获取k线
	go c.GetKlinesSync(id, p, ch)

	// 处理 Ctrl + C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	fmt.Println("Received interrupt signal, exiting...")
}
