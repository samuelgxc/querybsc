package cache

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"shangchenggo/conf"
	"net"
	"time"
)

var (
	redisPool *redis.Pool
)

//初始化连接redis
func init() {
	//获取redis配置
	address := net.JoinHostPort(beego.AppConfig.DefaultString(conf.REDIS_HOST, "127.0.0.1"), beego.AppConfig.DefaultString(conf.REDIS_PORT, "6379"))
	password := beego.AppConfig.DefaultString(conf.REDIS_PASSWORD, "")
	db := beego.AppConfig.DefaultInt(conf.REDIS_DB, 0)
	maxActive := beego.AppConfig.DefaultInt(conf.REDIS_MAXACTIVE, 500)
	maxIdle := beego.AppConfig.DefaultInt(conf.REDIS_MAXIDLE, 300)
	idleTimeoutTemp := beego.AppConfig.DefaultString(conf.REDIS_IDLETIMEOUT, "180") //时间单位：秒
	idleTimeout, _ := time.ParseDuration(idleTimeoutTemp + "s")
	beego.Info("--- 连接 redis ---", "address:", address, "password:", password, "db:", db)
	//连接redis并创建连接池
	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			//输入密码
			if password != "" {
				if _, err := conn.Do("AUTH", password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			//选择存贮库
			if db != 0 {
				if _, err := conn.Do("SELECT", db); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
	}

	//redis连接测试
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		beego.Error("--- 连接 redis 出错 ---", "err:", err)
	}
}

func Redis() redis.Conn {
	return redisPool.Get()
}

//Redis执行Do并设置超时时间
//  _, err = c.Do("SET", "testkey", "testvalue", "EX", "5")
func DoEx(conn redis.Conn, expireTime time.Duration, command string, args ...interface{}) (reply interface{}, err error) {
	//Redis执行Do
	if reply, err = conn.Do(command, args...); err != nil {
		return
	}
	//设置超时时间，默认单位：秒
	conn.Do("EXPIRE", args[0], int64(expireTime.Seconds()))
	return
}

//判读GET或HGET结果为空
func IfNil(err error) bool {
	if err != nil && err.Error() == "redigo: nil returned" {
		return true
	}
	return false
}

//将 对象 序列化为 string，存储到redis中
func Marshal(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

//将 redis中的string 反序列化为 对象
func Unmarshal(str string, v interface{}) error {
	err := json.Unmarshal([]byte(str), &v)
	return err
}
