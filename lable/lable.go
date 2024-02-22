package lable

import (
	"fmt"
	"github.com/gocraft/dbr"
	"os"
	"querybsc/bsc"
	"querybsc/common/db"
	_interface "querybsc/interface"
	"querybsc/lib"
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
var count *int64
func Start(){
	//fmt.Println(os.Args)
	if len(os.Args)==1 {
		fmt.Println("未指定标签")
		return
	}
	lable:=os.Args[1]

	var coin =""
	if len(os.Args)>2 && os.Args[2]!="-r"{
		coin = os.Args[2]
	}


	//地址长度如果为34.则是根据地址做跟踪
	//findForAddress(*address)
	addresss:=lib.GetInAddress()
	for _,address:=range addresss{
		var tryGet _interface.ITryGet
		if len(address)==42{
			tryGet = new(bsc.TryGet)
		}else{
			tryGet = new(trx.TryGet)
		}
		if len(os.Args)>3 && os.Args[3]=="-all"{
			tryGet.TryGetCoinAll(address,coin)
		}else if coin!=""{
			tryGet.TryGetCoin(address,coin)
		}else{
			tryGet.TrxGet(address,"",lable)
		}

		//更新lp的event_type
		if len(os.Args)>3 && os.Args[3]=="-lp"{
			var addIds []int64
			row1,_ := db.Session().Query("select id from tx t  WHERE t.address = ? AND EXISTS (SELECT 1 FROM tx t2 WHERE "+
				" address = ? and t.transaction_id = t2.transaction_id and t2.token_abbr!=t.token_abbr and t2.from_address = t.from_address and t2.from_address != ? )",address,address,address)
			for row1.Next(){
				var id int64
				row1.Scan(&id)
				addIds = append(addIds, id)
			}
			db.Session().Update("tx").Set("event_type","Addliquid").Where(dbr.Eq("id",addIds)).Exec()

			var removeIds []int64
			row2,_ := db.Session().Query("select id from tx t  WHERE t.address = ? AND EXISTS (SELECT 1 FROM tx t2 WHERE "+
				" address = ? and t.transaction_id = t2.transaction_id and t2.token_abbr!=t.token_abbr and t2.from_address = t.from_address and t2.to_address != ? )",address,address,address)
			for row2.Next(){
				var id int64
				row2.Scan(&id)
				removeIds = append(removeIds, id)
			}
			db.Session().Update("tx").Set("event_type","Removeliquid").Where(dbr.Eq("id",removeIds)).Exec()
		}
	}
	//替换标签
	if len(os.Args)>2 && os.Args[2]=="-r"{
		db.Session().Update("address").Set("lable",lable).Where(dbr.Eq("address",addresss)).Exec()
	}






	return
}
func findForAddress(address string) {
}


