package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

//测试存储
func TestSet(t *testing.T) {
	conn := Redis()
	defer conn.Close()
	//添加数据并设置过期时间
	//方式1：
	_, err := DoEx(conn, time.Second*10, "SET", "test:k1", "v1")
	//方式2：
	//_, err := conn.Do("SET", "test:k1", "v1", "EX", 10)
	if err != nil {
		fmt.Println("存储失败， err:", err)
		return
	}
	fmt.Println("存储成功")
	//判断数据是否存在
	ok, err := redis.Bool(conn.Do("EXISTS", "test:k1"))
	if err == nil && ok {
		fmt.Println("数据存在")
	}
}

//测试获取
func TestGet(t *testing.T) {
	conn := Redis()
	defer conn.Close()

	// GET 获取为空时报错，为redigo: nil returned
	val, err := redis.String(conn.Do("GET", "test:k1"))
	if IfNil(err) {
		fmt.Println("查询结果为空， err:", err)
		return
	} else if err != nil {
		fmt.Println("查询失败， err:", err)
		return
	}
	fmt.Println(val)
}

//测试HSet
func TestHSet(t *testing.T) {
	conn := Redis()
	defer conn.Close()
	//_, err := conn.Do("HSET", "test:h1", "k1", "v1")
	_, err := conn.Do("HMSET", "test:h1", "k1", "v1", "k2", "v2", "k3", "v3")
	if err != nil {
		fmt.Println(err)
	}
	_, err = conn.Do("EXPIRE", "test:h1", 10)
	if err != nil {
		fmt.Println(err)
	}
}

//测试HGet
func TestHGet(t *testing.T) {
	conn := Redis()
	defer conn.Close()

	// HGET 获取为空时报错，为redigo: nil returned
	val, err := redis.String(conn.Do("HGET", "test:h1", "k1"))
	if IfNil(err) {
		fmt.Println("查询结果为空， err:", err)
		return
	} else if err != nil {
		fmt.Println("查询失败， err:", err)
		return
	}
	fmt.Println(val)
}

//测试HMGet
func TestHMGet(t *testing.T) {
	conn := Redis()
	defer conn.Close()

	// HGET 获取结果为长度是指定field个数的数组，为空时仅数组元素为""
	val, err := redis.Strings(conn.Do("HMGET", "test:h1", "k1", "k2"))
	if err != nil {
		fmt.Println("查询失败， err:", err)
		return
	}
	fmt.Println(len(val))
	fmt.Println(val)
}

//测试HGetAll
func TestHGetAll(t *testing.T) {
	conn := Redis()
	defer conn.Close()

	// HGETALL 获取结果为map，为空时map长度为0
	val, err := redis.StringMap(conn.Do("HGETALL", "test:h1"))
	if err != nil {
		fmt.Println("查询失败， err:", err)
		return
	}
	fmt.Println(len(val))
	fmt.Println(val)
}
