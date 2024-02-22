package address

import (
	"flag"
	"fmt"
	"github.com/gocraft/dbr"
	"querybsc/common/db"
	"querybsc/lib"
)

var find  *string
var lable *string
var trx20_total *string
var memo *bool
func Start(){
	lable = flag.String("l", "", "标签")
	find  = flag.String("f", "", "查询")
	trx20_total= flag.String("trx20_total", "", "TRC20累计交易量")
	flag.Parse()

	var where dbr.Builder
	//通过标签查找
	if *lable!=""{
		where=dbr.Like("lable","%"+*lable+"%")
	}

	//查询地址
	if *find!=""{
		where=dbr.Like("address","%"+*find+"%")
	}
	if *trx20_total!=""{
		newwhere,_:=lib.NumWhere("trx20_total",*trx20_total)
		where=dbr.And(where,newwhere)
	}

	//fmt.Println(db.Session().Select("CONCAT(address,'-',lable)").From("address").Where(where).GetSQL())
	addresss,err:=db.Session().Select("CONCAT(address,'-',lable,'(',trx20_total,')')").From("address").Where(where).OrderAsc("id").ReturnStrings()
	//addresss,err:=db.Session().Select("CONCAT(address,'-',lable)").From("lable").Where(where).ReturnStrings()

	if err!=nil{

		fmt.Println(err.Error())
		panic(err)
	}

	for _,data:=range addresss{
		fmt.Println(data+"\r")
	}
	//fmt.Println(len(*find))
	if len(addresss)==0 && (len(*find)==42||len(*find)==34){
		fmt.Println(*find+"-未找到的地址")
	}
}