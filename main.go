package main

import (
	"os"
	"querybsc/address"
	_ "querybsc/common/db"
	"querybsc/disp"
	"querybsc/exchange"
	"querybsc/gecvs"
	"querybsc/gephi"
	"querybsc/lable"
	"querybsc/scan"
	"querybsc/track"
	"querybsc/tree"
	"querybsc/tx"
)

func main() {
	//初始化交易所信息
	exchange.Start()
	if len(os.Args) > 1 {
		mod := os.Args[1]
		os.Args = append(os.Args[0:1], os.Args[2:]...)
		switch mod {
		case "scan":
			scan.Start()
		case "address":
			address.Start()
		case "disp":
			disp.Start()
		case "tx":
			tx.Start()
		case "gecvs":
			gecvs.Start()
		case "lable":
			lable.Start()
		case "分析":
			track.Start()
		case "link":
			track.Start()
		case "tree":
			tree.Start()
		case "gephi":
			gephi.Start()
		}
	}
}
