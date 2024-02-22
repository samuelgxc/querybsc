package bsc

type Bsc20Ret struct {
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	TokenTransfers []struct {
		Block                 string  `json:"blockNumber"`
		BlockTs               string  `json:"timeStamp"`
		ContractAddress       string `json:"contractAddress"`
		FromAddress           string `json:"from"`
		Quant                 string `json:"value"`
		ToAddress             string `json:"to"`
		TokenDecimal string  `json:"tokenDecimal"`
		TokenName    string `json:"tokenSymbol"`
		TransactionID string `json:"hash"`
	} `json:"result"`
}


type BscRet struct {
	Status      string 	   `json:"status"`
	Message     string    `json:"message"`
	TokenTransfers []struct {
		Block                 string  `json:"blockNumber"`
		BlockTs               string  `json:"timeStamp"`
		ContractAddress       string `json:"contractAddress"`
		FromAddress           string `json:"from"`
		Quant                 string `json:"value"`
		ToAddress             string `json:"to"`
		Input				  string `json:"input"`
		TransactionID string `json:"hash"`
	} `json:"result"`
}
