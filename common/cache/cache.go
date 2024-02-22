package cache

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/ericlagergren/decimal"
	"shangchenggo/conf"
	"shangchenggo/tools"
)

var beegoCache cache.Cache

func init() {
	config := `{"key":"cache","conn":"` + beego.AppConfig.DefaultString(conf.REDIS_HOST, "127.0.0.1") + `:` + beego.AppConfig.DefaultString(conf.REDIS_PORT, "6379") + `","dbNum":"` + beego.AppConfig.DefaultString(conf.REDIS_DB, "0") + `","password":"` + beego.AppConfig.DefaultString(conf.REDIS_PASSWORD, "") + `"}`
	beego.Info("--- 连接 redis cache ---", "config:", config)
	var err error
	beegoCache, err = cache.NewCache("redis", config)
	if err != nil {
		beego.Error("--- 连接 redis cache 出错 ---", "err:", err)
	}
}

func NewCache() *Cache {
	c := &Cache{}
	c.Cache = beegoCache
	return c
}

type Cache struct {
	cache.Cache
}

func (c *Cache) GetBool(key string) bool {
	value := c.Get(key)
	return tools.ToBool(value)
}

func (c *Cache) GetInt(key string) int {
	value := c.Get(key)
	return tools.ToInt(value)
}

func (c *Cache) GetInt64(key string) int64 {
	value := c.Get(key)
	return tools.ToInt64(value)
}

func (c *Cache) GetFloat64(key string) float64 {
	value := c.Get(key)
	return tools.ToFloat64(value)
}

func (c *Cache) GetString(key string) string {
	value := c.Get(key)
	return tools.ToString(value)
}

func (c *Cache) GetBig(key string) *decimal.Big {
	value := c.Get(key)
	return tools.ToBig(value)
}
