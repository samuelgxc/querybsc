package trx

import "encoding/json"

type Trc20Ret struct {
	ContractInfo   struct{} `json:"contractInfo"`
	RangeTotal     int64    `json:"rangeTotal"`
	TokenTransfers []struct {
		ApprovalAmount        string `json:"approval_amount"`
		Block                 int64  `json:"block"`
		BlockTs               int64  `json:"block_ts"`
		Confirmed             bool   `json:"confirmed"`
		ContractRet           string `json:"contractRet"`
		ContractAddress       string `json:"contract_address"`
		ContractType          string `json:"contract_type"`
		EventType             string `json:"event_type"`
		FinalResult           string `json:"finalResult"`
		FromAddressIsContract bool   `json:"fromAddressIsContract"`
		FromAddress           string `json:"from_address"`
		Quant                 string `json:"quant"`
		Revert                bool   `json:"revert"`
		ToAddressIsContract   bool   `json:"toAddressIsContract"`
		ToAddress             string `json:"to_address"`
		TokenInfo             struct {
			TokenAbbr    string `json:"tokenAbbr"`
			TokenCanShow int64  `json:"tokenCanShow"`
			TokenDecimal int64  `json:"tokenDecimal"`
			TokenID      string `json:"tokenId"`
			TokenLogo    string `json:"tokenLogo"`
			TokenName    string `json:"tokenName"`
			TokenType    string `json:"tokenType"`
			Vip          bool   `json:"vip"`
		} `json:"tokenInfo"`
		TransactionID string `json:"transaction_id"`
	} `json:"token_transfers"`
	Total int64 `json:"total"`
}


type TrcRet struct {
	ContractInfo   struct{} `json:"contractInfo"`
	RangeTotal     int64    `json:"rangeTotal"`
	TokenTransfers []struct {
		Block                 int64  `json:"block"`
		BlockTs               int64  `json:"timestamp"`
		ContractAddress       string `json:"contract_address"`
		FromAddress           string `json:"ownerAddress"`
		Quant                 string `json:"amount"`
		ToAddress             string `json:"toAddress"`
		ContractType          string `json:"tokenType"`
		TokenInfo             struct {
			TokenAbbr    string `json:"tokenAbbr"`
			TokenCanShow int64  `json:"tokenCanShow"`
			TokenDecimal int64  `json:"tokenDecimal"`
			TokenID      string `json:"tokenId"`
			TokenLogo    string `json:"tokenLogo"`
			TokenName    string `json:"tokenName"`
			TokenType    string `json:"tokenType"`
			Vip          bool   `json:"vip"`
		} `json:"tokenInfo"`
		TransactionID string `json:"hash"`
	} `json:"data"`
	Total int64 `json:"total"`
}


type transferRet struct {
	ContractInfo struct{} `json:"contractInfo"`
	ContractMap  json.RawMessage `json:"contractMap"`
	Data []struct {
		Amount      int64  `json:"amount"`
		Block       int64  `json:"block"`
		Confirmed   bool   `json:"confirmed"`
		ContractRet string `json:"contractRet"`
		Data        string `json:"data"`
		ID          string `json:"id"`
		Revert      bool   `json:"revert"`
		Timestamp   int64  `json:"timestamp"`
		TokenInfo   struct {
			TokenAbbr    string `json:"tokenAbbr"`
			TokenCanShow int64  `json:"tokenCanShow"`
			TokenDecimal int64  `json:"tokenDecimal"`
			TokenID      string `json:"tokenId"`
			TokenLogo    string `json:"tokenLogo"`
			TokenName    string `json:"tokenName"`
			TokenType    string `json:"tokenType"`
			Vip          bool   `json:"vip"`
		} `json:"tokenInfo"`
		TokenName           string `json:"tokenName"`
		TransactionHash     string `json:"transactionHash"`
		TransferFromAddress string `json:"transferFromAddress"`
		TransferToAddress   string `json:"transferToAddress"`
	} `json:"data"`
	RangeTotal int64 `json:"rangeTotal"`
	Total      int64 `json:"total"`
}
