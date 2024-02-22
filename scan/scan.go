package scan

import (
	"flag"
	"fmt"
	"querybsc/bsc"
	"querybsc/common/db"
	_interface "querybsc/interface"
	"querybsc/trx"
)
//扫描地址信息
var address *string
var direction *string
var coin *string
var lable *string
var min *string
var max *string
var nums *string
var count *string

//机构地址信息
var exchange map[string]string
func Start(){
	address   = flag.String("a", "", "地址")
	direction = flag.String("d", "", "方向")
	coin      = flag.String("c", "", "币种")
	lable     = flag.String("l", "", "标签")
	nums      = flag.String("nums", "", "最大金额")
	count     = flag.String("count", "", "个人汇总记录数")

	flag.Parse()
	if *direction!="in" && *direction!="out"{
		fmt.Println("方向不正确 -d")
		return
	}
	if *coin==""{
		fmt.Println("币种未确定 -c")
		return
	}
	if *address==""{
		fmt.Println("地址未填 -a")
		return
	}
	if *lable==""{
		//fmt.Println("标签未填 -l")
		//return
	}

	//地址长度如果为34.则是根据地址做跟踪
	findForAddress(*address)
	//加载查询主地址的信息
	return

}

func findForAddress(address string){
	var tryGet _interface.ITryGet
	if len(address)==42{
		tryGet = new(bsc.TryGet)
	}else{
		tryGet = new(trx.TryGet)
	}

	//确认更新完主地址资料后.开始扫描后续
	where:=""
	if len(address)==42{
		tryGet.TrxGet(address,"","")
		//方向条件
		where+="address='"+address+"'"
	}else if len(address)==34{
		tryGet.TrxGet(address,"","")
		//方向条件
		where+="address='"+address+"'"
	}else{
		//方向条件
		where+="address in (select address from address where lable='"+address+"') and address not in (select address from exchange)"
	}

	//货币类型
	where+=" and token_abbr='"+*coin+"'"
	//fmt.Println(where)
	if min!=nil && *min!=""{
		where+=" and number>="+*min
	}
	if max!=nil &&  *max!=""{
		where+=" and number<="+*max
	}
	if nums!=nil && *nums!=""{
		where+=" and number in ("+*nums+")"
	}
	txdatas:=[]TxStru{}
	selectstmt:=db.Session().Select("*").From("tx").Where(where)
	//fmt.Println(selectstmt.GetSQL())
	_,err:=selectstmt.Load(&txdatas)
	if err!=nil{
		panic(err)
	}
	//fmt.Println(txdatas)
	addressmap:=map[string]string{}
	var countmap=map[string]int64{}
	for _,req:=range txdatas{

		var scanaddress string
		var soureAddress string
		if *direction=="in"{
			scanaddress=req.FromAddress
			soureAddress=req.ToAddress
		}else{
			scanaddress=req.ToAddress
			soureAddress=req.FromAddress
		}
		//如果本身是交易所地址,则忽略迭代
		if _,ok:=addressmap[scanaddress];ok{
			continue
		}
		//如果有COUNT参数.则是要汇总自己的符合条件的交易记录数
		if *count!=""{
			countmap[soureAddress]++;
		}else {
			var tryGetS _interface.ITryGet
			if len(address)==42{
				tryGetS = new(bsc.TryGet)
			}else{
				tryGetS = new(trx.TryGet)
			}
			addressmap[scanaddress]=tryGetS.TrxGet(scanaddress,soureAddress,*lable)
		}

		//fmt.Println(scanaddress,soureAddress)

	}
	//寻找到如下地址
	for k,v:=range addressmap{
		fmt.Println(k,v)
	}
}