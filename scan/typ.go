package scan

import "github.com/shopspring/decimal"

type TxStru struct{
	FromAddress string `db:"from_address"`
	ToAddress   string `db:"to_address"`
	Number      decimal.Decimal `db:"number"`
}

