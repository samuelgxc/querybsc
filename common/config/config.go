package config

import (
	"fmt"
	"github.com/ericlagergren/decimal"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"querybsc/common/db"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var configVal = map[string]*conData{}

type Bool func(v ...bool) bool
type Int64 func(v ...int64) int64
type Uint64 func(v ...uint64) uint64
type String func(v ...string) string
type Big func(v ...*decimal.Big) *decimal.Big

type conData struct {
	raw  string      //数据库原始配置信息
	val  interface{} //转换后实际信息
	typ  reflect.Type
	lock sync.RWMutex
	isup bool //是否被更新过
}

func init() {
	loadConfig()
	loadVal(&Val)
}

func Init() {}

func Open(name string) bool {
	return false
}

func Update() {
	//打开数据库
	//判断数据库中conVer是否一致。不一致则把数据库中的最新配置信息同步一次。
	//更新数据中的conVer以及
}

func SetValue(key string, value interface{}) {
	sourceType := reflect.TypeOf(configVal[key].val)
	paramsType := reflect.TypeOf(value)
	if sourceType == paramsType {
		configVal[key].lock.RLock()
		defer configVal[key].lock.RUnlock()
		configVal[key].val = value
	}
}

func GetValue(key string) interface{} {
	return getCon(key)
}

func loadConfig() {
	rows, err := db.Session().Select("`name`", "`data`").From("config").Rows()
	if err != nil {
		log.Printf("--- 加载 appConfig 配置出错，err: %v ---", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var key, value string
		rows.Scan(&key, &value)
		key = strings.ToLower(key)
		configVal[key] = &conData{raw: value}
	}
}

func loadVal(obj interface{}) {
	//反射配置对象
	objType := reflect.TypeOf(obj)
	objVal := reflect.ValueOf(obj)
	if objType.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
		objType = objVal.Type()
	}
	for i := 0; i < objVal.NumField(); i++ {
		//属性类型
		fieldType := objType.Field(i)
		//属性值（方法）
		fieldVal := objVal.Field(i)
		//创建配置属性锁
		key := strings.ToLower(fieldType.Name)
		if _, ok := configVal[key]; !ok {
			panic("配置信息 [" + key + "] 未找到")
		}
		//创建配置属性方法
		switch fieldType.Type.Name() {
		case "Bool":
			fieldVal.Set(makeBool(key))
		case "Uint64":
			fieldVal.Set(makeUint64(key))
		case "Int64":
			fieldVal.Set(makeInt64(key))
		case "String":
			fieldVal.Set(makeString(key))
		case "Big":
			fieldVal.Set(makeBig(key))
		default:
			panic("未知配置类型" + fieldType.Type.Name())
		}
	}
}

func getCon(key string) interface{} {
	configVal[key].lock.RLock()
	defer configVal[key].lock.RUnlock()
	return configVal[key].val
}

func makeBool(key string) reflect.Value {
	var err error
	configVal[key].val, err = strconv.ParseBool(configVal[key].raw)
	if err != nil {
		panic(fmt.Sprintf("转换Int64配置项%s出错%s", key, err.Error()))
	}
	return reflect.ValueOf(func(v ...bool) bool {
		return getCon(key).(bool)
	})
}

func makeUint64(key string) reflect.Value {
	var err error
	configVal[key].val, err = strconv.ParseUint(configVal[key].raw, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("转换Int64配置项%s出错%s", key, err.Error()))
	}
	return reflect.ValueOf(func(v ...uint64) uint64 {
		return getCon(key).(uint64)
	})
}

func makeInt64(key string) reflect.Value {
	var err error
	configVal[key].val, err = strconv.ParseInt(configVal[key].raw, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("转换Int64配置项%s出错%s", key, err.Error()))
	}
	return reflect.ValueOf(func(v ...int64) int64 {
		return getCon(key).(int64)
	})
}

func makeString(key string) reflect.Value {
	configVal[key].val = configVal[key].raw
	return reflect.ValueOf(func(v ...string) string {
		return getCon(key).(string)
	})
}

func makeBig(key string) reflect.Value {
	var ok bool
	configVal[key].val, ok = decimal.New(0, 0).SetString(configVal[key].raw)
	if !ok {
		panic(fmt.Sprintf("转换Int64配置项%s出错%s", key, configVal[key].raw))
	}
	return reflect.ValueOf(func(v ...*decimal.Big) *decimal.Big {
		return decimal.New(0, 0).Set(getCon(key).(*decimal.Big))
	})
}
