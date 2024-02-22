package tx

import (
	"flag"
	"fmt"
	"github.com/gocraft/dbr"
	"querybsc/common/db"
	"querybsc/lib"
	"strconv"
	"strings"
)
//扫描地址信息
var address *string
var direction *string
var coin *string
var lable *string
var min *string
var max *string
var num *string
var count *string
//机构地址信息
var exchange map[string]string
func Start(){
	direction = flag.String("d", "", "方向")
	coin = flag.String("c", "", "币种")
	num = flag.String("num", "", "最小金额")
	lable = flag.String("l", "", "过滤标签")
	count = flag.String("count", "", "统计个人符合条件交易数量")
	flag.Parse()
	if *direction!="in" && *direction!="out" && *direction!="inout"{
		fmt.Println("方向不正确 -d")		
		return
	}

	if *coin==""{
		*coin = "Bnb"
	}
	//fmt.Println(2222)
	//地址长度如果为34.则是根据地址做跟踪
		findForAddress()
	//加载查询主地址的信息
	return
}
func findForAddress(){
	//fmt.Println(lib.GetInAddress())
	var where = dbr.And(
		dbr.Eq("token_abbr",*coin),
		)
	inaddress:=lib.GetInAddress()


	if *direction=="in"{
		where      = dbr.And(where,dbr.Eq("address",inaddress),dbr.Eq("to_address",inaddress))
	}
	if *direction=="out"{
		where      = dbr.And(where,dbr.Eq("address",inaddress),dbr.Eq("from_address",inaddress))
	}
	if *direction=="inout" {
		where      = dbr.And(where,dbr.Eq("address",inaddress),dbr.Or(dbr.Eq("to_address",inaddress),dbr.Eq("from_address",inaddress)))
	}

	if *num!=""{
		numwhere,_:=lib.NumWhere("number",*num)
		where      = dbr.And(where,numwhere)
	}
	txdatas:=[]TxStru{}
	selectstmt:=db.Session().Select("*").From("tx").Where(where)
	//fmt.Println(selectstmt.GetSQL())
	//fmt.Println(selectstmt.GetSQL())
	_,err:=selectstmt.Load(&txdatas)
	if err!=nil{
		panic(err)
	}
	//fmt.Println(txdatas)
	addressmap:=map[string]string{}
	countmap :=map[string]int64{}
	inmap :=map[string]int64{}
	outmap:=map[string]int64{}
	for _,req:=range txdatas{

		var scanaddress string
		var soureAddress string
		//var soureAddress string

		if req.Address==req.ToAddress{
			inmap[req.FromAddress]=1
			scanaddress=req.FromAddress
			soureAddress=req.ToAddress
		}else{
			outmap[req.ToAddress]=1
			scanaddress=req.ToAddress
			soureAddress=req.FromAddress
		}


		if *count!=""{
			countmap[soureAddress]++;
			addressmap[soureAddress]=""
		}else{
			addressmap[scanaddress]=""
		}
	}
	if *direction=="inout" {
		for k:=range addressmap{
			_,ok1:=inmap[k]
			_,ok2:=outmap[k]
			if !ok1 || !ok2{
				delete(addressmap,k)
			}
		}
	}
	if *count != "" {
		countarr:=strings.Split(*count,"-")
		//存在最小值判定
		var min int64
		var max int64 =-1
		if countarr[0]!=""{
			min,_=strconv.ParseInt(countarr[0],10,0)
		}
		if len(countarr)==2{
			if countarr[1]!=""{
				max,_=strconv.ParseInt(countarr[1],10,0)
			}
		}else{
			//如果只有一个数字.那么数量则固定
			max=min
		}
		for k:=range addressmap{
			val:=countmap[k]
			max=max
			if val<min{
				delete(addressmap,k)
				continue
			}
			if max>1 && val>max{
				delete(addressmap,k)
				continue
			}
		}
	}
	lables:=strings.Split(*lable,",")

	//寻找到如下地址
	for _,v:=range lib.FormatAddressMap(addressmap){
		if v==""{
			continue
		}
		if *lable==""{
			fmt.Println(v+"\r")
		}else{
			for _,iflable:=range lables {
				if strings.Contains(v,iflable){
					fmt.Println(v+"\r")
				}
			}
		}
	}
}