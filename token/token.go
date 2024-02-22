package token

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

const (
	ContractAddress = "0x2aa504586d6cab3c59fa629f74c586d78b93a025"
	ApiKey = "J2DF58A68NKCXTI2QNUZAHT3YNKVUKYAW9"
)

var action  *string
var address  *string

type BscRet struct {
	Status	int64	`json:"status"`
	Message string  `json:"message"`
	Result  string 	`json:"result"`
}

func Start(){
	action = flag.String("a", "", "方法")
	address = flag.String("addr", "", "地址")
	flag.Parse()

	//发币的部署地址，可以通过 https://coinmarketcap.com，查询匹配合约地址，不建议自动拉取，因为不精准
	//币的一些信息，

	var url = "https://api.bscscan.com/api"
	url += "?apikey=" + ApiKey
	url += "&contractaddress=" + ContractAddress

	switch *action {
	//返回代币的当前流通量
	case "tokenCsupply":
		url += "&module=stats"
		url += "&action=tokenCsupply"
	//返回代币的当前总量
	case "tokensupply":
		url += "&module=stats"
		url += "&action=tokensupply"
	//返回某地址的代币余额
	case "tokenbalance":
		url += "&module=account"
		url += "&action=tokenbalance"
		url += "&address=" + *address
		url += "&tag=latest"
	//返回代币持有地址和持有数量
	case "tokenholderlist":
		url += "&module=token"
		url += "&action=tokenholderlist"
		url += "&page=1"
		url += "&offset=10"
	//返回代币信息
	case "tokeninfo":
		url += "&module=token"
		url += "&action=tokeninfo"
	}

	url = "https://bscscan.com/token/" + ContractAddress

	rest, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	ret := new(BscRet)
	json.NewDecoder(rest.Body).Decode(ret)
	fmt.Println(ret)

	//持币前十，https://bscscan.com/token/0x2aa504586d6cab3c59fa629f74c586d78b93a025#balances

	//LP的信息，https://bscscan.com/token/0x2aa504586d6cab3c59fa629f74c586d78b93a025#tokenTrade

	//LP持币前十，根据ContractAddress得到hash，查询当时logs，发现添加的交易对 0x3de032d5d11c94d2d79dba0c34d7851ffaa05dd8
	//----------获取此交易对的持币前十，https://bscscan.com/token/0x3de032d5d11c94d2d79dba0c34d7851ffaa05dd8#balances

	//币合约的owner之类的这种，明面上的owner是0地址，还有holder、fundAddress等地址，不确定是哪个

	//这几天，恶补了一番区块链相关的知识，看了bscscan相关的接口，看了coinmarketcap.com相关的接口，又看了pancakeswap相关的接口，真是活到老学到老！

}