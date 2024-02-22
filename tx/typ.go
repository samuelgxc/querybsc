package tx

import "github.com/shopspring/decimal"

type TxStru struct{
	Address     string `db:"address"`
	FromAddress string `db:"from_address"`
	ToAddress   string `db:"to_address"`
	Number      decimal.Decimal `db:"number"`
}

