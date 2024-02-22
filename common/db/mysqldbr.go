package db

import (
	"context"
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"querybsc/conf"
	"time"
)

var (
	conn     *dbr.Connection
	connRead *dbr.Connection
)

func init() {
	//写数据库连接
	config := conf.GetDbConfig()
	//beego.Info("--- 连接 mysql 主库（写库） ---", "config:", config)
	var err error
	conn, err = dbr.Open("mysql", config, nil)
	if err != nil {
		beego.Error("--- 连接 mysql 主库（写库）出错 ---", "err:", err)
		return
	}
	maxOpenConns := beego.AppConfig.DefaultInt(conf.MYSQL_MAXOPENCONNS, 500)
	maxIdleConns := beego.AppConfig.DefaultInt(conf.MYSQL_MAXIDLECONNS, 300)
	connMaxLifetime, _ := time.ParseDuration(beego.AppConfig.DefaultString(conf.MYSQL_CONNMAXLIFETIME, "180") + "s") //时间单位：秒
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetConnMaxLifetime(connMaxLifetime)

	//读数据库连接（若没有单独配置读库，则与写库相同）
	//config = conf.GetDbReadConfig()
	//beego.Info("--- 连接 mysql 从库（读库） ---", "config:", config)
	//connRead, err = dbr.Open("mysql", config, nil)
	//if err != nil {
	//	beego.Error("--- 连接 mysql 从库（读库）出错 ---", "err:", err)
	//	connRead = conn
	//	return
	//}
	//connRead.SetMaxOpenConns(maxOpenConns)
	//connRead.SetMaxIdleConns(maxIdleConns)
	//connRead.SetConnMaxLifetime(connMaxLifetime)
}

func Session() *dbr.Session {
	return conn.NewSession(nil)
}

func SessionRead() *dbr.Session {
	return connRead.NewSession(nil)
}

func GetTx(db *dbr.Session, level sql.IsolationLevel) (*dbr.Tx, error) {
	return db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: level,
	})
}

//判断是否是 *dbr.Tx
func IsTx(sess dbr.SessionRunner) (*dbr.Tx, bool) {
	tx, ok := sess.(*dbr.Tx)
	return tx, ok
}
