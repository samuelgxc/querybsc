package trx

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/dbr"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"querybsc/common/db"
	"querybsc/exchange"
	"strconv"
	"time"
)

type TryGet struct {}
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
	pageCount := int64(1)
	//trc20总交易长度
	trx20_total := int64(0)
	for page := int64(1); page <= pageCount; page++ {
		fmt.Println("处理", page)
		fmt.Println("https://apilist.tronscan.org/api/token_trc20/transfers?sort=-timestamp&count=true&limit=50&start=" + strconv.Itoa(int((page-1)*50)) + "&address=" + address)
		rest, err := http.Get("https://apilist.tronscan.org/api/token_trc20/transfers?sort=-timestamp&count=true&limit=50&start=" + strconv.Itoa(int((page-1)*50)) + "&relatedAddress=" + address)
		if err != nil {
			panic(err)
		}
		ret := new(Trc20Ret)
		json.NewDecoder(rest.Body).Decode(ret)
		//fmt.Println("数据长度",len(ret.Data),ret.Total)
		if page == 1 {
			//统计页数
			pageCount = ret.Total / 50
			if float64(ret.Total)/50 != 0 {
				pageCount++
			}
			trx20_total = ret.Total
			fmt.Println(address, "总记录数", ret.Total)
		}
		for _, data := range ret.TokenTransfers {
			num, _ := decimal.NewFromString(data.Quant)
			deci := data.TokenInfo.TokenDecimal
			num = num.Div(decimal.New(1,int32(deci)))
			//检查交易所标签
			lable = exchange.Get(data.ToAddress, lable)
			instStmt.Values(
				address,
				data.FromAddress,
				data.ToAddress,
				data.TokenInfo.TokenAbbr,
				num.String(),
				data.TransactionID,
				data.EventType,
				data.ContractType,
				data.FromAddressIsContract,
				data.ToAddressIsContract,
				int(data.BlockTs/1000),
				data.Block,
			)
		}
	}

	//尝试查询第一页
	pageCount = int64(1)
	//trc20总交易长度
	trx_total := int64(0)

	_, err = instStmt.Exec()
	if err != nil {
		panic(err)
	}
	_, err = tx.InsertInto("address").Columns("address", "source_address", "trx20_total", "lable").Values(address, sourceAddress, trx20_total+trx_total, lable).Exec()
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

}




func (t *TryGet)TryGetCoinAll(address string ,coins ...string){
	var lable = ""
	if address == "" {
		return
	}

	if exchange.Have(address) {
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
	trx20_total := int64(0)
	var endTime = 0
	for {
		var url = fmt.Sprintf("https://apilist.tronscan.org/api/token_trc20/transfers?sort=-timestamp&count=true&limit=50&start=0&relatedAddress=" + address+"&start_timestamp=0")
		if endTime>0{
			url = fmt.Sprintf("https://apilist.tronscan.org/api/token_trc20/transfers?sort=-timestamp&count=true&limit=50&start=0&relatedAddress=" + address+"&start_timestamp=%d&end_timestamp=%d",0,endTime*1000)
		}
		fmt.Println(url)
		rest, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		ret := new(Trc20Ret)
		json.NewDecoder(rest.Body).Decode(ret)
		//if endTime>0 && ret.Total==10000{
		//	fmt.Println(fmt.Sprintf("当前时间数量大"))
		//	endTime = endTime-3600
		//	continue
		//}

		if ret.Total==0{
			fmt.Println(fmt.Println("查询结束"))
			body, _ := io.ReadAll(rest.Body)
			fmt.Println(string(body))
			//if len(string(body))==0{
			//	fmt.Println("重启")
			//	continue
			//}
			break
		}
		//fmt.Println("数据长度",len(ret.Data),ret.Total)
		for _, data := range ret.TokenTransfers {
			if endTime==0{
				var exists =0
				db.Session().QueryRow("select time from tx where address=?  order by time asc limit 1",address).Scan(&exists)
				if exists>0 {
					fmt.Println(exists)
					endTime = exists-1
					break
				}
			}
			num, _ := decimal.NewFromString(data.Quant)
			deci := data.TokenInfo.TokenDecimal
			num = num.Div(decimal.New(1,int32(deci)))
			//检查交易所标签
			lable = exchange.Get(data.ToAddress, lable)
			instStmt.Values(
				address,
				data.FromAddress,
				data.ToAddress,
				data.TokenInfo.TokenAbbr,
				num.String(),
				data.TransactionID,
				data.EventType,
				data.ContractType,
				data.FromAddressIsContract,
				data.ToAddressIsContract,
				int(data.BlockTs/1000),
				data.Block,
			)
			trx20_total++
			endTime = int(data.BlockTs/1000-1)
			if len(instStmt.Value)>1000{
				_, err = instStmt.Exec()
				if err != nil {
					panic(err)
				}
				fmt.Println("更新部分批量数据")
				tx.Commit()
				tx, err = db.Session().Begin()
				if err != nil {
					panic(tx)
				}
				instStmt = tx.InsertInto("tx").Columns(
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
			}
		}
	}

	endTime = 0
	for {
		var url = fmt.Sprintf("https://apilist.tronscanapi.com/api/transaction?sort=-timestamp&count=true&limit=50&start=0&address="+address+"&type=1")
		if endTime>0{
			url = fmt.Sprintf("https://apilist.tronscanapi.com/api/transaction?sort=-timestamp&count=true&limit=50&start=0&address="+address+"&type=1&start_timestamp=%d&end_timestamp=%d",0,endTime*1000)
		}
		fmt.Println(url)
		rest, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		ret := new(TrcRet)
		json.NewDecoder(rest.Body).Decode(ret)
		//if endTime>0 && ret.Total==10000{
		//	fmt.Println(fmt.Sprintf("当前时间数量大"))
		//	endTime = endTime-3600
		//	continue
		//}

		if ret.Total==0{
			fmt.Println(fmt.Println("查询结束"))
			body, _ := io.ReadAll(rest.Body)
			fmt.Println(string(body))
			//if len(string(body))==0{
			//	fmt.Println("重启")
			//	continue
			//}
			break
		}
		fmt.Println("数据长度",len(ret.TokenTransfers),ret.Total)
		for _, data := range ret.TokenTransfers {
			if endTime==0{
				var exists =0
				db.Session().QueryRow("select time from tx where address=?  and token_abbr='TRX' order by time asc limit 1",address).Scan(&exists)
				if exists>0 {
					fmt.Println(exists)
					endTime = exists-1
					break
				}else{
					endTime = int(time.Now().Unix())
				}
			}

			trx20_total++
			num, _ := decimal.NewFromString(data.Quant)
			deci := data.TokenInfo.TokenDecimal
			num = num.Div(decimal.New(1,int32(deci)))
			//检查交易所标签
			lable = exchange.Get(data.ToAddress, lable)
			instStmt.Values(
				address,
				data.FromAddress,
				data.ToAddress,
				"TRX",
				num.String(),
				data.TransactionID,
				"Trx Transfer",
				data.ContractType,
				0,
				0,
				int(data.BlockTs/1000),
				data.Block,
			)
			endTime = int(data.BlockTs/1000-1)
			if len(instStmt.Value)>1000{
				_, err = instStmt.Exec()
				if err != nil {
					panic(err)
				}
				fmt.Println("更新部分批量数据")
				tx.Commit()
				tx, err = db.Session().Begin()
				if err != nil {
					panic(tx)
				}
				instStmt = tx.InsertInto("tx").Columns(
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
			}
		}
	}
	if len(instStmt.Value)>0{
		_, err = instStmt.Exec()
		if err != nil {
			panic(err)
		}
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

