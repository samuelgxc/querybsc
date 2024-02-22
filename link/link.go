package link

import (
	"flag"
	"github.com/gocraft/dbr"
	"strings"
	"querybsc/lib"
)

var lable *string
func Start(){
	lable = flag.String("l", "", "标签")
	flag.Parse()
	inaddress:=lib.GetInAddress()
	where:=dbr.And(
			//是否转账到指定标签
			//dbr.Eq("to_address"  ,),
			dbr.Expr("to_address in (select address from address where lable ='"+*lable+"')"),
			dbr.Eq("from_address",strings.Split(*lable,",")),
		)
	dbr.Select("*").From("tx").Where(where).GroupBy("address")
}