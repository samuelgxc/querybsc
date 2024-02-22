package track

import (
	"database/sql"
	"fmt"
	"os"
	"querybsc/bsc"
	"querybsc/common/db"
	"querybsc/interface"
	"querybsc/trx"
	"strings"
)

func Start(){
	if len(os.Args)==1{
		fmt.Println("未输入跟踪地址")
		return
	}

	address:=os.Args[1]

	addresses := strings.Split(address,",")
	if len(addresses)!=2{
		fmt.Println("输入地址不符合规范")
		return
	}

	addressQuerys := make(map[string]map[string]int)
	farAddrs := make(map[string]map[string]string)
	for i:=0;i< len(addresses);i++{
		addr := addresses[i]
		var tryGet _interface.ITryGet
		if len(address)==42{
			tryGet = new(bsc.TryGet)
		}else{
			tryGet = new(trx.TryGet)
		}
		fmt.Printf("============查询相关地址流水:%s \n==============",addr)
		tryGet.TrxGet(addr,"","地址关联追踪")
		addrs,err := db.Session().SelectBySql("select from_address as address  from tx where to_address = ? " +
			"union all select to_address as address  from tx where from_address=?",addr,addr).ReturnStrings()
		if err!=nil{
			panic(err)
		}
		fmt.Println("============查询相关线索地址流水==============")
		var mapUnionTemp = make(map[string]int)
		var mapFarTemp = make(map[string]string)
		for _,unionAddr :=range addrs{
			tryGet.TrxGet(unionAddr,"","地址关联线索地址")
			mapUnionTemp[unionAddr]=1
			fars,err := db.Session().SelectBySql("select from_address as address  from tx where to_address = ? " +
				"union all select to_address as address  from tx where from_address=?",unionAddr,unionAddr).ReturnStrings()
			if err!=nil{
				panic(err)
			}

			for _,farAddr := range fars{
				if _,find := mapFarTemp[farAddr];find {
					if !strings.Contains(mapFarTemp[farAddr],unionAddr){
						mapFarTemp[farAddr] = mapFarTemp[farAddr]+","+unionAddr
					}
				}else{
					mapFarTemp[farAddr] = unionAddr
				}
			}

		}
		farAddrs[addr] =mapFarTemp
		addressQuerys[addr]=mapUnionTemp
	}
	fmt.Println("============查询关联信息==============")
	//查询两个地址是否有直接关联
	qSid ,err := db.Session().SelectBySql("select transaction_id from tx where from_address in (?,?) and to_address in (?,?)",
		addresses[0],addresses[1],addresses[0],addresses[1]).ReturnStrings()
	if err!=nil && err!=sql.ErrNoRows{
		panic(err)
	}
	if len(qSid)>0{
		//qSidStr := strings.Join(qSid,",")
		//fmt.Printf("查询地址1：%s<==hash:%s===>查询地址2：%s \n",addresses[0],qSidStr,addresses[1])
		fmt.Printf("查询地址1：%s<=====>查询地址2：%s \n",addresses[0],addresses[1])
		return
	}
	//查询两个地址之间是否有共同地址，如果有，返回相关交易信息，如果没有，查询线索地址流水，比较两组流水中的相同地址信息
	for addr,_ := range addressQuerys[addresses[0]]{

		if _,find := addressQuerys[addresses[1]][addr];find{
			//sid1 ,err :=db.Session().SelectBySql("select transaction_id from tx " +
			//	"where from_address in (?,?) and to_address in (?,?)",addr,addresses[0],addr,addresses[0]).ReturnStrings()
			//if err!=nil{
			//	panic(err)
			//}
			//sid2 ,err :=db.Session().SelectBySql("select transaction_id from tx " +
			//	"where from_address in (?,?) and to_address in (?,?)",addr,addresses[1],addr,addresses[1]).ReturnStrings()
			//if err!=nil{
			//	panic(err)
			//}
			//sidStr1 := strings.Join(sid1,",")
			//sidStr2 := strings.Join(sid2,",")
			//fmt.Printf("查询地址1：%s<==hash:%s===>线索地址：%s<==hash:%s===>查询地址2：%s \n",
			//	addresses[0],sidStr1,addr,sidStr2,addresses[1])

			fmt.Printf("查询地址1：%s<=====>线索地址：%s<=====>查询地址2：%s \n",
				addresses[0],addr,addresses[1])
		}
	}
	for addr,_ := range farAddrs[addresses[0]]{

		if _,find := addressQuerys[addresses[1]][addr];find{
			fmt.Printf("查询地址1：%s<=====>线索地址1：%s<=====>线索地址2：%s<=====>查询地址2：%s \n",
				addresses[0],farAddrs[addresses[0]][addr],addr,addresses[1])
		}

		if _,find := farAddrs[addresses[1]][addr];find{
			fmt.Printf("查询地址1：%s<=====>线索地址1：%s<=====>线索地址2：%s<=====>线索地址3：%s<=====>查询地址2：%s \n",
				addresses[0],farAddrs[addresses[0]][addr],addr,farAddrs[addresses[1]][addr],addresses[1])
		}
	}
	for addr,_ := range farAddrs[addresses[1]]{

		if _,find := addressQuerys[addresses[0]][addr];find{
			fmt.Printf("查询地址1：%s<=====>线索地址1：%s<=====>线索地址2：%s<=====>查询地址2：%s \n",
				addresses[0],addr,farAddrs[addresses[1]][addr],addresses[1])
		}
	}

}



type AddressStru struct {
	SourceAddress string `db:"source_address"`
	Address string `db:"address"`
	Lable   string `db:"lable"`
}