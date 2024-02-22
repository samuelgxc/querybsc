package gephi

import (
	"flag"
	"fmt"
	"io"
	"os"
	"querybsc/common/db"
	"strconv"
)
var coin *string

func Start() {
	coin      = flag.String("c", "", "币种")
	flag.Parse()
	var Datas = []*NodeData{}
	var IdMap = map[string]string{}
	var wMap = map[string]float64{}
	//权重数据
	var coinToken string = "USDT"
	if *coin!=""{
		coinToken = *coin
	}
	fmt.Println(coinToken)
	//db.Session().SelectBySql("select address,sum(number) from tx where address=from_address and token_abbr='"+coinToken+"' group by address").Load(&wMap)
	db.Session().SelectBySql("select address,sum(number) from tx where address=from_address group by address").Load(&wMap)
	fmt.Println(wMap)
	_, err := db.Session().Select("id", "lable", "address").From("address_gephi").Load(&Datas)
	if err != nil {
		panic(err)
	}

	_,err = os.Stat("node.csv")
	if os.IsNotExist(err){
		os.Create("node.csv")
	}else{
		os.Remove("node.csv")
		os.Create("node.csv")
	}
	nodefile, err := os.OpenFile("node.csv", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	io.WriteString(nodefile, "Id,Lable,timeset,modularity_class,Weight\n")
	for _, data := range Datas {
		nodefile.WriteString(data.Id + "," + data.Address + "," + "," + data.Lable + "," + fmt.Sprintf("%1.0f", wMap[data.Address]) + "\n")
		IdMap[data.Address] = data.Id
	}

	//nodefile.Sync()
	nodefile.Close()
	//边数据
	var LDatas = []*LineData{}

	_, err = db.Session().Select("from_address", "to_address", "sum(number) usdt_sum").From("tx").
		//Where("from_address in (select address from address) and to_address in (select address from address) and token_abbr='"+coinToken+"'").
		Where("(from_address in (select address from address) or from_address in (select address from exchange)) and to_address in (select address from address)").
		GroupBy("from_address,to_address").Load(&LDatas)
	for _, data := range LDatas {
		data.FromId = IdMap[data.From]
		data.ToId = IdMap[data.To]
	}
	_,err = os.Stat("line.csv")
	if os.IsNotExist(err){
		os.Create("line.csv")
	}else{
		os.Remove("line.csv")
		os.Create("line.csv")
	}
	linefile, err := os.OpenFile("line.csv", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	io.WriteString(linefile, "Source,Target,Type,Id,usdt_sum\n")
	for id, data := range LDatas {
		linefile.WriteString(data.FromId + "," + data.ToId + "," + "," + strconv.Itoa(id) + "," + data.Sum + "\n")
	}
	linefile.Close()

}

type NodeData struct {
	Id      string `db:"id"`
	Address string `db:"address"`
	Lable   string `db:""`
}

type LineData struct {
	From   string `db:"from_address"`
	FromId string
	To     string `db:"to_address"`
	ToId   string
	Sum    string `db:"usdt_sum"`
}
