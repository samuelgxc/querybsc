package bsc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/shopspring/decimal"
	"net/http"
	url2 "net/url"
	"querybsc/common/db"
	"querybsc/exchange"
	"strconv"
	"strings"
)

var Api = "QQZCVDSJYREXE8FRJEZWRFPDJ51VSCEHZ5"
var Url = "api.bscscan.com"
var coin = "0x55d398326f99059fF775485246999027B3197955"
//var bscApi = "QQZCVDSJYREXE8FRJEZWRFPDJ51VSCEHZ5"
//var bscUrl = "api.bscscan.com"
//var coin = "0x55d398326f99059fF775485246999027B3197955"
//var ethApi = "YSVTMJI57ANG2VCNNAHWDFPQUFWX7CPTRA"
//var ethUrl = "api.etherscan.io"
//var coin = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
//var proxyURL = "http://127.0.0.1:7890"
var proxyURL = ""

//https://api.bscscan.com/api?module=logs&action=getLogs&fromBlock=1&toBlock=4993832&address=0xbc35B8a2EF3aE2f1390FA3C36259dEb8d963b320&topic0=0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925&apikey=QQZCVDSJYREXE8FRJEZWRFPDJ51VSCEHZ5

type TryGet struct {}

func HttpGetFromProxy(url string) (*http.Response,error) {
	if proxyURL==""{
		return http.Get(url)
	}
	req,_ := http.NewRequest("GET",url,nil)
	proxy,err := url2.Parse(proxyURL)
	if err != nil {
		return nil,err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy : http.ProxyURL(proxy),},}
	return client.Do(req)
}
//登记地址信息
//地址
//查询来源
func (t *TryGet)TrxGet(address string, sourceAddress string, lable string) string {
	if address == "" {
		return ""
	}
	//查询的地址本身是一个机构.(可能是充币地址)
	if exchange.Have(address) {
		return exchange.Get(address, "") + "-出币"
	}

	address = strings.ToLower(address)
	//fmt.Println(db.Session().Select("ifnull(lable,'')").From("address").Where("address=?",address).GetSQL())
	var lable2 string
	lable2, err := db.Session().Select("lable").From("address").Where("address=?", address).ReturnString()
	if err != nil && err != dbr.ErrNotFound {
		panic(err)
	}
	if err == nil {
		return lable2
	}

	tx, err := db.Session().Begin()
	if err != nil {
		panic(tx)
	}
	if sourceAddress == "" {
		//尝试找到一条来源记录
		sourceAddress, err = db.Session().Select("from_address").From("tx").Where("to_address=?", address).ReturnString()
	}
	if sourceAddress == "" {
		//	//尝试找到一条来源记录
		sourceAddress, err = db.Session().Select("to_address").From("tx").Where("from_address=?", address).ReturnString()
	}
	//id
	//from_address
	//to_address
	//direction
	//tokenAbbr
	//number
	//transaction_id
	//event_type
	//contract_type
	instStmt := tx.InsertInto("tx").Columns(
		"address",
		"from_address",
		"to_address",
		"token_abbr",
		"number",
		"transaction_id",
		"event_type",
		"contract_type",
		"FromAddressIsContract",
		"ToAddressIsContract",
		"Time",
		"block",
	)

	//尝试查询第一页
	pageCount := int64(2)
	//trc20总交易长度
	trx20_total := int64(0)
	for page := int64(1); page <= pageCount; page++ {
		var url = "https://"+Url+"/api?module=account&action=txlist&address="+address+"&startblock=0&endblock=99999999&page="+strconv.Itoa(int(page))+"&offset=5000&sort=desc&apikey="+Api
		fmt.Println(url)
		rest, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		ret := new(BscRet)
		json.NewDecoder(rest.Body).Decode(ret)
		//fmt.Println("数据长度",len(ret.Data),ret.Total)
		for _, data := range ret.TokenTransfers {
			if data.Input!="0x"{
				continue
			}
			num, _ := decimal.NewFromString(data.Quant)
			num = num.Div(decimal.New(1,18))
			//fmt.Println(num.String())
			blockTs,_ := strconv.ParseInt(data.BlockTs,10,64)
			if num.Cmp(decimal.New(1,-3))<0{
				continue
			}
			//fmt.Println(data.FromAddress)
			//检查交易所标签
			//lable = exchange.Get(data.FromAddress,lable)
			lable = exchange.Get(data.ToAddress,lable)
			instStmt.Values(
				address,
				data.FromAddress,
				data.ToAddress,
				"Bnb",
				""+num.String()+"",
				data.TransactionID,
				"Erc10 Transfer",
				"0",
				"0",
				"0",
				int(blockTs),
				data.Block,
			)
			trx20_total++
		}

	}

	//获取默认usdt信息
	pageCount = int64(2)
	//coin := "0x55d398326f99059fF775485246999027B3197955"
	var queryExists = false
	//trc20总交易长度
	for page := int64(1); page <= pageCount; page++ {
		var url = "https://"+Url+"/api?module=account&action=tokentx&contractaddress="+coin+"&address="+address+"&page="+strconv.Itoa(int(page))+"&offset=5000&startblock=0&endblock=999999999&sort=desc&apikey="+Api
		fmt.Println(url)
		rest, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		ret := new(Bsc20Ret)
		json.NewDecoder(rest.Body).Decode(ret)
		//fmt.Println("数据长度",len(ret.Data),ret.Total)
		for _, data := range ret.TokenTransfers {
			if !queryExists{
				var exists =0
				err = db.Session().QueryRow("select 1 from tx where address=? and token_abbr=? limit 1",address,data.TokenName).Scan(&exists)
				if err!=nil && err!= sql.ErrNoRows{
					fmt.Println(err.Error())
					break
				}
				if exists>0{
					break
				}
				queryExists = true
			}

			//lable = exchange.Get(data.FromAddress,lable)
			lable = exchange.Get(data.ToAddress,lable)
			num, _ := decimal.NewFromString(data.Quant)
			deci,_ := strconv.ParseInt(data.TokenDecimal,10,64)
			num = num.Div(decimal.New(1,int32(deci)))
			fmt.Println(num.String())
			blockTs,_ := strconv.ParseInt(data.BlockTs,10,64)
			//检查交易所标签
			instStmt.Values(
				address,
				data.FromAddress,
				data.ToAddress,
				data.TokenName,
				num.String(),
				data.TransactionID,
				"Erc20 Transfer",
				"0",
				"0",
				"0",
				int(blockTs),
				data.Block,
			)
			trx20_total++
		}
	}
	if len(instStmt.Value)>0{
		_, err = instStmt.Exec()
		if err != nil {
			panic(err)
		}
	}

	_, err = tx.InsertInto("address").Columns("address", "source_address", "trx20_total", "lable").Values(address, sourceAddress, trx20_total, lable).Exec()
	if err != nil {
		panic(err)
	}


	if lable != "" {
		_, err = tx.InsertInto("lable").Columns("address", "lable").Values(address, lable).Exec()
		if err != nil {
			panic(err)
		}
	}



	tx.Commit()
	return lable
}


func (t *TryGet)TryGetCoin(address string ,coins ...string){
	var lable = ""
	if address == "" {
		return
	}
	if exchange.Have(address) {
		return
	}

	address = strings.ToLower(address)
	if len(coins)==0{
		return
	}
	tx, err := db.Session().Begin()
	if err != nil {
		panic(tx)
	}
	instStmt := tx.InsertInto("tx").Columns(
		"address",
		"from_address",
		"to_address",
		"token_abbr",
		"number",
		"transaction_id",
		"event_type",
		"contract_type",
		"FromAddressIsContract",
		"ToAddressIsContract",
		"Time",
		"block",
	)
	var coinlist = coins[0]
	coinlists := strings.Split(coinlist,",")

	trx20_total := int64(0)
	for i:=0;i< len(coinlists);i++{
		var coin = coinlists[i]
		//尝试查询第一页
		pageCount := int64(2)
		//trc20总交易长度
		for page := int64(1); page <= pageCount; page++ {
			var url = "https://"+Url+"/api?module=account&action=tokentx&contractaddress="+coin+"&address="+address+"&page="+strconv.Itoa(int(page))+"&offset=5000&startblock=0&endblock=999999999&sort=desc&apikey="+Api
			fmt.Println(url)
			rest, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			ret := new(Bsc20Ret)
			json.NewDecoder(rest.Body).Decode(ret)
			//fmt.Println("数据长度",len(ret.Data),ret.Total)
			for _, data := range ret.TokenTransfers {
				var exists =0
				db.Session().QueryRow("select 1 from tx where address=? and token_abbr=? limit 1",address,data.TokenName).Scan(&exists)
				if exists>0{
					break
				}
				//lable = exchange.Get(data.FromAddress,lable)
				lable = exchange.Get(data.ToAddress,lable)
				num, _ := decimal.NewFromString(data.Quant)
				deci,_ := strconv.ParseInt(data.TokenDecimal,10,64)
				num = num.Div(decimal.New(1,int32(deci)))
				blockTs,_ := strconv.ParseInt(data.BlockTs,10,64)
				//检查交易所标签
				instStmt.Values(
					address,
					data.FromAddress,
					data.ToAddress,
					data.TokenName,
					num.String(),
					data.TransactionID,
					"Erc20 Transfer",
					"0",
					"0",
					"0",
					int(blockTs),
					data.Block,
				)
				trx20_total++
			}
		}
	}



	if len(instStmt.Value)>0{
		_, err = instStmt.Exec()
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	_, err = tx.UpdateBySql("update address set trx20_total =trx20_total+ ? where address = ?",trx20_total,address).Exec()
	if err != nil {
		panic(err)
	}

	if lable!=""{
		_, err = tx.UpdateBySql("update address set lable= ? where address = ?",lable,address).Exec()
		if err != nil {
			panic(err)
		}

		_, err =  tx.UpdateBySql("update lable set lable= ? where address = ?",lable,address).Exec()
		if err != nil {
			panic(err)
		}
	}

	tx.Commit()



}




func (t *TryGet)TryGetCoinAll(address string ,coins ...string){
	var lable = ""
	if address == "" {
		return
	}
	if exchange.Have(address) {
		return
	}

	address = strings.ToLower(address)
	if len(coins)==0{
		return
	}
	tx, err := db.Session().Begin()
	if err != nil {
		panic(tx)
	}
	instStmt := tx.InsertInto("tx").Columns(
		"address",
		"from_address",
		"to_address",
		"token_abbr",
		"number",
		"transaction_id",
		"event_type",
		"contract_type",
		"FromAddressIsContract",
		"ToAddressIsContract",
		"Time",
		"block",
	)
	var coinlist = coins[0]
	coinlists := strings.Split(coinlist,",")
	trx20_total := int64(0)
	for i:=0;i< len(coinlists);i++{
		var coin = coinlists[i]
		var queryBlock int64 = 999999999
		//trc20总交易长度
		for {
			var url = fmt.Sprintf("https://"+Url+"/api?module=account&action=tokentx&contractaddress=%s" +
				"&address=%s&page=%d&offset=5000&startblock=0&endblock=%d&sort=desc&apikey="+Api,coin,address,1,queryBlock)
			fmt.Println(url)
			rest, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			ret := new(Bsc20Ret)
			json.NewDecoder(rest.Body).Decode(ret)
			if ret.Status!="1"{
				fmt.Println(ret.Message)
				break
			}
			if len(ret.TokenTransfers)==0{
				break
			}
			//fmt.Println("数据长度",len(ret.Data),ret.Total)
			for _, data := range ret.TokenTransfers {
				if queryBlock==999999999{
					var exists =0
					db.Session().QueryRow("select block from tx where address=? and token_abbr=? order by block asc limit 1",address,data.TokenName).Scan(&exists)
					if exists>0 {
						queryBlock = int64(exists)-1
						break
					}
				}

				lable = exchange.Get(data.ToAddress,lable)
				num, _ := decimal.NewFromString(data.Quant)
				deci,_ := strconv.ParseInt(data.TokenDecimal,10,64)
				num = num.Div(decimal.New(1,int32(deci)))
				blockTs,_ := strconv.ParseInt(data.BlockTs,10,64)
				//检查交易所标签
				instStmt.Values(
					address,
					data.FromAddress,
					data.ToAddress,
					data.TokenName,
					num.String(),
					data.TransactionID,
					"Erc20 Transfer",
					"0",
					"0",
					"0",
					int(blockTs),
					data.Block,
				)
				trx20_total++
				nBlock,_ := strconv.ParseInt(data.Block,10,64)
				queryBlock = nBlock-1
			}
		}
	}
	if len(instStmt.Value)>0{
		_, err = instStmt.Exec()
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	_, err = tx.UpdateBySql("update address set trx20_total =trx20_total+ ? where address = ?",trx20_total,address).Exec()
	if err != nil {
		panic(err)
	}

	if lable!=""{
		_, err = tx.UpdateBySql("update address set lable= ? where address = ?",lable,address).Exec()
		if err != nil {
			panic(err)
		}

		_, err =  tx.UpdateBySql("update lable set lable= ? where address = ?",lable,address).Exec()
		if err != nil {
			panic(err)
		}
	}
	tx.Commit()
}

