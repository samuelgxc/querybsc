package lib

import (
	"bufio"
	"github.com/gocraft/dbr"
	"io/ioutil"
	"os"
	"querybsc/common/db"
	"querybsc/exchange"
	"sort"
	"strconv"
	"strings"
)

//将一个字符串转换为条件.主要有三种情况
//1.固定值
//2.范围值
//3.枚举值
func NumWhere(name,exp string)(dbr.Builder,error){
	exps:=strings.Split(exp,"-")
	//斜线分隔
	if len(exps)==2{
		num1,_:=strconv.ParseFloat(exps[0],64)
		num2,_:=strconv.ParseFloat(exps[1],64)
		return dbr.And(dbr.Gte(name,num1),dbr.Lte(name,num2)),nil
	}
	exps=strings.Split(exp,",")
	if len(exps)>1{
		nums:=[]float64{}
		for _,exp:=range exps{
			//d:=new(decimal.Big)
			//d.SetString(exp)
			f,_:=strconv.ParseFloat(exp,64)

			nums=append(nums,f)
		}
		//fmt.Println("!!!!!!")
		return dbr.Eq(name,nums),nil
	}
	f,_:=strconv.ParseFloat(exp,64)
	return dbr.Eq(name,f),nil
}

func GetInAddress()[]string{

	reader := bufio.NewReader(os.Stdin)
	b,_:=ioutil.ReadAll(reader)
	//b,_=GbkToUtf8(b)
	//fmt.Println(string(b))
	//fmt.Println(b)
	lines:=strings.Split(string(b),"\n")
	addressmap:=map[string]string{}

	for _,line:=range lines{
		if strings.Contains(line,"-"){
			line=line[0:strings.Index(line,"-")]
		}
		addressmap[line]=""
	}
	var ret=[]string{}
	for k:=range addressmap{
		k=strings.ReplaceAll(k,"\r","")
		k=strings.ReplaceAll(k,"\n","")
		ret = append(ret,k)
	}
	return ret
}
func FormatAddressMap(m map[string]string)[]string{

	ret:=[]string{}
	for v,_:=range m{
		ret=append(ret,v)
	}
	//附备注
	lableMap:=map[string]string{}
	_,err:=db.Session().Select("address","lable").From("address").Where(dbr.Eq("address",ret)).Load(&lableMap)
	if err!=nil{
		panic(err)
	}
	sort.Strings(ret)
	for index,v:=range ret{
		if lable,ok:=lableMap[v];ok{
			ret[index]=v+"-"+lable
			if exchange.Have(v){
				ret[index]+="(主钱包)"
			}
		}
	}
	//fmt.Println(err)
	//fmt.Println(lableMap)
	//fmt.Println(ret)
	//os.Exit(1)
	return ret
}

func FormatAddressInt64Map(m map[string]int64)[]string{

	ret:=[]string{}
	for v,_:=range m{
		ret=append(ret,v)
	}
	//附备注
	lableMap:=map[string]string{}
	_,err:=db.Session().Select("address","lable").From("address").Where(dbr.Eq("address",ret)).Load(&lableMap)
	if err!=nil{
		panic(err)
	}
	sort.Strings(ret)
	for index,v:=range ret{
		if lable,ok:=lableMap[v];ok{
			ret[index]=v+"-"+lable
			if exchange.Have(v){
				ret[index]+="(主钱包)"
			}
			ret[index]+=" "+strconv.FormatInt(m[v],10)
		}
	}
	//fmt.Println(err)
	//fmt.Println(lableMap)
	//fmt.Println(ret)
	//os.Exit(1)
	return ret
}