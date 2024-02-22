package tree

import (
	"fmt"
	"querybsc/common/db"
	"github.com/gocraft/dbr"
	"strconv"
)

var indexMap = map[string]string{}
var lableMap = map[string]string{}
var countMap = map[string]string{}

func Start() {

	//查询所有地址归属关系.形成图谱
	//目标节点
	db.Session().Select("address,source_address").From("address").
		Where("lable in ('直接交易所','收款地址','出款地址','币安','火币','去中心化交易所','未知机构','交互地址','交互地址3级')").Load(&indexMap)

	db.Session().Select("address,trx20_total").From("address").Load(&countMap)
	db.Session().Select("address,lable").From("address").Load(&lableMap)

	for {
		notFindup := []string{}
		for _, v := range indexMap {
			if _, ok := indexMap[v]; !ok && v != "" {
				notFindup = append(notFindup, v)
			}
		}
		//fmt.Println(notFindup)
		//如果有未找到的上级
		if len(notFindup) > 0 {
			tempMap := map[string]string{}
			db.Session().Select("address,source_address").From("address").
				Where(dbr.Eq("address", notFindup)).Load(&tempMap)
			for k, v := range tempMap {
				indexMap[k] = v
			}
		} else {
			break
		}
	}

	addresss, _ := db.Session().Select("address").From("address").Where("source_address=''").ReturnStrings()
	//--fmt.Println(addresss)
	for _, v := range addresss {
		Disp(v, "", "顶点", 0)
	}
}

//查询地址,层级,索引信息,来源备注
func Disp(address string, lstr, from_memo string, layer int) {

	fmt.Print(lstr)
	//确定与源头的关系
	//源头是否给自己交易过
	in, _ := db.Session().Select("count(*)").From("tx").
		Where("to_address =? and from_address=?", address, indexMap[address]).ReturnInt64()

	out, _ := db.Session().Select("count(*)").From("tx").
		Where("to_address =? and from_address=?", indexMap[address], address).ReturnInt64()

	if out > 0 {
		fmt.Print("<")
	} else {
	}

	if in > 0 {
		fmt.Print(">")
	} else {
	}
	fmt.Print(" ")
	inusdt, _ := db.Session().Select("FLOOR(ifnull(sum(number),0)/10000)").From("tx").
		Where("address =? and to_address =? and token_abbr='USDT'", address, address).ReturnInt64()

	fmt.Println(address+"-"+lableMap[address]+"("+countMap[address]+") 关联地址 "+strconv.Itoa(links(address))+" 累计转入USDT", inusdt, "w")

	var downlist = []string{}
	for k, v := range indexMap {
		if v == address {
			downlist = append(downlist, k)
		}
	}

	for i, v := range downlist {
		newlstr := lstr
		//fmt.Println(i+1<len(downlist),i,len(downlist))
		if i+1 < len(downlist) {
			newlstr += "  ┣ "
		} else {
			newlstr += "  ┗ "
		}
		//fmt.Println(lstr)
		rstr := []rune(newlstr)
		//fmt.Println(rstr)

		for i, v := range rstr {
			if i < len(rstr)-2 && v == '┗' {
				rstr[i] = '　'
			}
			if i < len(rstr)-2 && v == '┣' {
				rstr[i] = '┃'
			}
		}
		//fmt.Println("下属",len(downlist))
		//fmt.Println(rstr)
		Disp(v, string(rstr), "顶点", layer+1)
	}
}
func links(address string) int {
	gxaddress, _ := db.Session().Select("if(from_address='"+address+"',to_address,from_address)").From("tx").Where("(from_address in (select address from address) and to_address in (select address from address)) and address=?", address).
		ReturnStrings()
	//gxaddress
	gxmap := map[string]int64{}
	SourceMap := map[string]string{}
	db.Session().Select("address", "source_address").From("address").
		Where(dbr.Eq("address", gxaddress)).
		Load(&SourceMap)
	//累加与目前系统有关的账户的交易计数
	for _, v := range gxaddress {
		gxmap[v]++
	}
	return len(gxmap)
	//gxaddress
}
