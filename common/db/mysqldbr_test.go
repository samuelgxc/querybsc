package db

import (
	"fmt"
	"github.com/gocraft/dbr"
	"shangchenggo/common/page"
	"testing"
	"time"
)

type config struct {
	id    int64
	key   string
	value string
}

func TestDbr(t *testing.T) {
	sess := Session()
	stmt := sess.Select("`name`", "data").From("config").Where(dbr.Eq("id", []int64{68,69,70,71,72,73,74,75,76}))
	data, err := page.NewPaginator(stmt, 1, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}
}

func TestSelectMap(t *testing.T) {
	var maps []dbr.Maps
	sess := Session()
	//方式1：where条件 = 的值可以是一个数组
	num, err := sess.Select("`name`", "data").From("config").Maps(&maps)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num, err)
		for _, term := range maps {
			fmt.Println(term["name"], ":", term["data"])
		}
	}
}

func TestSelect(t *testing.T) {
	sess := Session()
	//方式1：where条件 = 的值可以是一个数组
	rows, err := sess.Select("key", "value").From("config").Where("id > ? AND id < ?", 1, 100).Rows()
	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			var key, value string
			rows.Scan(&key, &value)
			fmt.Println(key, value)
		}
	}
	//方式2：where条件 eq 的值可以是一个数组;可以LoadOne只查一条数据
	var configs []config
	where := dbr.And(dbr.Lt("id", 100), dbr.Gt("id", 1))
	_, err = sess.Select("key", "value").From("config").Where(where).Load(&configs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(configs)
	}
}

func TestInsert(t *testing.T) {
	sess := Session()
	//方式1：
	sess.InsertInto("config").Columns("key", "value").Values("akey", "avalue").Exec()
	//方式2：
	conf := config{
		key:   "bkey",
		value: "bvalue",
	}
	sess.InsertInto("config").Columns("key", "value").Record(&conf).Exec()
	//方式3：插入同时查询
	sess.InsertInto("config").Columns("key", "value").Record(&conf).Returning("id").Load(&conf.id)
	fmt.Println(conf.id)
}

func TestUpdate(t *testing.T) {
	sess := Session()
	//方式1：
	sess.Update("config").Set("value", "aavalue").Where("key = akey").Exec()
	//方式2：
	setMap := map[string]interface{}{
		"value": "bbvalue",
	}
	sess.Update("config").SetMap(setMap).Where("key = bkey").Exec()
	//方式3：
	num, err := sess.Update("config").Set("deleted", 1).Set("value", dbr.Expr("value + ?", "10.123")).Where("id = 8").Exec()
	fmt.Println(err)
	if num != nil {
		fmt.Println(num.RowsAffected())
	}
}

func TestDb1(t *testing.T) {
	sess := Session()
	//测试查不到数据时的返回情况
	rows, err := sess.Select("*").From("user_symbol").Where("id = 1").Rows()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("aaa")
		if rows.Next() {
			fmt.Println("bbb")
		}
	}

	fmt.Println("----------")
	var data []int64
	count, err := sess.Select("1").From("user_symbol").Where("id = 1").Load(&data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(count)
		fmt.Println(data)
	}

	fmt.Println("----------")
	var data2 int64
	err = sess.Select("SUM(id)").From("user_symbol").Where("id = 1").LoadOne(&data2)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err == dbr.ErrNotFound)
	} else {
		fmt.Println(data2)
	}
}

func TestDb2(t *testing.T) {
	sess := Session()

	//setMap := UserSymbol{
	//	//Id:0,
	//	User_id:   1,
	//	Symbol_id: 1,
	//	W_time:    time.Now().Unix(),
	//}
	//ct, err := sess.InsertInto("user_symbol").Columns("user_id", "symbol_id", "w_time").Record(setMap).Exec()

	setMap := map[string]interface{}{
		//"id":        0,
		"user_id":   1,
		"symbol_id": 1,
		"w_time":    time.Now().Unix(),
	}
	ct, err := sess.InsertInto("user_symbol").Map(setMap).Exec()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ct.LastInsertId())
	}
}

func TestDb3(t *testing.T) {
	sess := Session()
	var data []int64
	count, err := sess.Select("id").From("user_symbol").Where("user_id = 2").
		Where(dbr.Expr("symbol_id in (select symbol_id from user_symbol where id in ?)", []int64{10, 11})).
		//Where("symbol_id in (select symbol_id from user_symbol where id in ?)", []int64{10, 11}).
		Load(&data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(count)
		fmt.Println(data)
	}
}

func TestDb4(t *testing.T) {
	sess := Session()
	data := UserSymbol{}
	where := dbr.Eq("user_id", 2)
	where = dbr.And(nil, where, dbr.Eq("symbol_id", 8))
	err := sess.Select("*").From("user_symbol").Where(where).LoadOne(&data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
		fmt.Println(data.User_id)
		fmt.Println(data.CreateTime)
	}
}

func TestDb5(t *testing.T) {
	sess := Session()
	data := []UserSymbol{}
	num, err := sess.Select("*").From("user_symbol").Where("id = 5").ShowSql().Load(&data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num)
		for _, rs := range data {
			fmt.Println(rs)
		}
	}
}

type UserSymbol struct {
	Id         int64
	User_id    int64               //下划线法
	SymbolId   int64               //驼峰法
	CreateTime int64 `db:"w_time"` //注解法
}
