package gecvs

import (
	"flag"
	"fmt"
	"io"
	"os"
	"querybsc/common/db"
	"strconv"
	"time"
)


var address *string
var direction *string
var coin *string
var lable *string
var min *string
var max *string
var event *string

func Start() {
	address   = flag.String("a", "", "地址")
	direction = flag.String("d", "", "方向")
	coin      = flag.String("c", "", "币种")
	lable     = flag.String("l", "", "标签")
	min      = flag.String("min", "", "最小金额")
	max      = flag.String("max", "", "最大金额")
	event    = flag.String("event","","事件")
	flag.Parse()

	if *address==""{
		fmt.Println("地址未填 -a")
		return
	}
	if *lable==""{
		//fmt.Println("标签未填 -l")
		//return
	}

	QueryData(*address)
}

func QueryData(address string) {
	//确认更新完主地址资料后.开始扫描后续
	where:=""
	if *direction=="in"{
		where += "to_address ='"+address+"'"
	}else if *direction=="out"{
		where += "from_address ='"+address+"'"
	}else if len(address)==42{
		//bsc.TrxGet(address,"","")
		//方向条件
		where+="address='"+address+"'"
	}

	if *lable != ""{
		//方向条件
		where+="address in (select address from address where lable='"+address+"') and address not in (select address from exchange)"
	}
	//货币类型
	if *coin!=""{
		where+=" and token_abbr='"+*coin+"'"
	}
	//fmt.Println(where)
	if min!=nil && *min!=""{
		where+=" and number>="+*min
	}
	if max!=nil &&  *max!=""{
		where+=" and number<="+*max
	}

	if *event =="transfer"{
		where += " and event_type like '%Transfer%'"
	}else if *event =="add"{
		where += " and event_type = 'Addliquid'"
	}else if *event =="remove"{
		where += " and event_type = 'Removeliquid'"
	}



	txdatas:=[]TxStru{}
	selectstmt:=db.Session().Select("from_address,to_address,token_abbr,number,transaction_id,event_type,time").From("tx").Where(where)
	//fmt.Println(selectstmt.GetSQL())
	_,err:=selectstmt.Load(&txdatas)
	if err!=nil{
		panic(err)
	}
	_,err = os.Stat("transfer.csv")
	if os.IsNotExist(err){
		os.Create("transfer.csv")
	}else{
		os.Remove("transfer.csv")
		os.Create("transfer.csv")
	}


	linefile, err := os.OpenFile("transfer.csv", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}


	io.WriteString(linefile, "from,to,token,number,hash,type,time\n")
	for _, data := range txdatas {
		atime,_ :=strconv.ParseInt(data.Time,10,64)
		t:= time.Unix(atime,0).Format("2006-01-02 15:04:05")

		linefile.WriteString(data.FromAddress + "," + data.ToAddress + "," + data.Token+ "," + data.Number + "," + data.Hash+ "," + data.Event  +"," + t  +  "\n")
	}
	linefile.Close()

}



type TxStru struct{
	FromAddress string `db:"from_address"`
	ToAddress   string `db:"to_address"`
	Number      string `db:"number"`
	Token   string `db:"token_abbr"`
	Hash   string `db:"transaction_id"`
	Event   string `db:"event_type"`
	Time   string `db:"time"`
}
