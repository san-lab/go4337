package rpccalls

import "fmt"

type UserOperationByHashResult_Biconomy struct {
	Sender               string `json:"sender"`
	Nonce                int    `json:"nonce"`
	InitCode             string `json:"initCode"`
	CallData             string `json:"callData"`
	CallGasLimit         int    `json:"callGasLimit"`
	PreVerificationGas   int    `json:"preVerificationGas"`
	VerificationGasLimit int    `json:"verificationGasLimit"`
	PaymasterAndData     string `json:"paymasterAndData"`
	MaxFeePerGas         int64  `json:"maxFeePerGas"`
	MaxPriorityFeePerGas int    `json:"maxPriorityFeePerGas"`
	Signature            string `json:"signature"`
	EntryPoint           string `json:"entryPoint"`
	TransactionHash      string `json:"transactionHash"`
	BlockNumber          int    `json:"blockNumber"`
	BlockHash            string `json:"blockHash"`
}

func RepackageUSOPResult(ausopr *UserOperationByHashResult_Biconomy) (*UserOperationByHashResult_ASP, error) {
	if ausopr == nil {
		return nil, fmt.Errorf("No result in the input")
	}
	res := &UserOperationByHashResult_ASP{}
	res.BlockHash = ausopr.BlockHash
	res.BlockNumber = fmt.Sprintf("0x%x", ausopr.BlockNumber)
	res.EntryPoint = ausopr.EntryPoint
	res.TransactionHash = ausopr.TransactionHash
	res.UserOperation.Sender = ausopr.Sender
	res.UserOperation.InitCode = ausopr.InitCode
	res.UserOperation.CallData = ausopr.CallData
	res.UserOperation.CallGasLimit = fmt.Sprintf("0x%x", ausopr.CallGasLimit)
	res.UserOperation.VerificationGasLimit = fmt.Sprintf("0x%x", ausopr.VerificationGasLimit)
	res.UserOperation.PreVerificationGas = fmt.Sprintf("0x%x", ausopr.PreVerificationGas)
	res.UserOperation.MaxFeePerGas = fmt.Sprintf("0x%x", ausopr.MaxFeePerGas)
	res.UserOperation.MaxPriorityFeePerGas = fmt.Sprintf("0x%x", ausopr.MaxPriorityFeePerGas)
	res.UserOperation.PaymasterAndData = ausopr.PaymasterAndData
	res.UserOperation.Signature = ausopr.Signature
	return res, nil

}

type UserOperationReceipt_Biconomy struct {
	UserOpHash    string `json:"userOpHash"`
	EntryPoint    string `json:"entryPoint"`
	Sender        string `json:"sender"`
	Nonce         int    `json:"nonce"`
	Success       string `json:"success"`
	Paymaster     string `json:"paymaster"`
	ActualGasCost int64  `json:"actualGasCost"`
	ActualGasUsed int    `json:"actualGasUsed"`
	Reason        string `json:"reason"`
	Logs          []struct {
		Address          string   `json:"address"`
		Topics           []string `json:"topics"`
		Data             string   `json:"data"`
		BlockNumber      string   `json:"blockNumber"`
		TransactionHash  string   `json:"transactionHash"`
		TransactionIndex string   `json:"transactionIndex"`
		BlockHash        string   `json:"blockHash"`
		LogIndex         string   `json:"logIndex"`
		Removed          bool     `json:"removed"`
	} `json:"logs"`
	Receipt struct {
		BlockHash         string      `json:"blockHash"`
		BlockNumber       string      `json:"blockNumber"`
		ContractAddress   interface{} `json:"contractAddress"`
		CumulativeGasUsed string      `json:"cumulativeGasUsed"`
		EffectiveGasPrice string      `json:"effectiveGasPrice"`
		From              string      `json:"from"`
		GasUsed           string      `json:"gasUsed"`
		Logs              []struct {
			Address          string   `json:"address"`
			Topics           []string `json:"topics"`
			Data             string   `json:"data"`
			BlockNumber      string   `json:"blockNumber"`
			TransactionHash  string   `json:"transactionHash"`
			TransactionIndex string   `json:"transactionIndex"`
			BlockHash        string   `json:"blockHash"`
			LogIndex         string   `json:"logIndex"`
			Removed          bool     `json:"removed"`
		} `json:"logs"`
		LogsBloom        string `json:"logsBloom"`
		Status           string `json:"status"`
		To               string `json:"to"`
		TransactionHash  string `json:"transactionHash"`
		TransactionIndex string `json:"transactionIndex"`
		Type             string `json:"type"`
	} `json:"receipt"`
}

func RepackageUSOPReceiptResult(ausopr *UserOperationReceipt_Biconomy) (*UserOperationReceipt_Alchemy, error) {
	if ausopr == nil {
		return nil, fmt.Errorf("No result in the input")
	}
	res := &UserOperationReceipt_Alchemy{}
	res.UserOpHash = ausopr.UserOpHash
	res.EntryPoint = ausopr.EntryPoint
	res.Sender = ausopr.Sender
	res.Nonce = fmt.Sprintf("0x%x", ausopr.Nonce)
	res.Success = (ausopr.Success == "true")
	res.Paymaster = ausopr.Paymaster
	res.ActualGasCost = fmt.Sprintf("0x%x", ausopr.ActualGasCost)
	res.ActualGasUsed = fmt.Sprintf("0x%x", ausopr.ActualGasUsed)
	res.Reason = ausopr.Reason
	for _, log := range ausopr.Logs {
		res.Logs = append(res.Logs, ReceiptLog{
			Address:          log.Address,
			Topics:           log.Topics,
			Data:             log.Data,
			BlockNumber:      log.BlockNumber,
			TransactionHash:  log.TransactionHash,
			TransactionIndex: log.TransactionIndex,
			BlockHash:        log.BlockHash,
			LogIndex:         log.LogIndex,
			Removed:          log.Removed,
		})
	}
	res.Receipt.BlockHash = ausopr.Receipt.BlockHash
	res.Receipt.BlockNumber = ausopr.Receipt.BlockNumber
	res.Receipt.ContractAddress = ausopr.Receipt.ContractAddress
	res.Receipt.CumulativeGasUsed = ausopr.Receipt.CumulativeGasUsed
	res.Receipt.EffectiveGasPrice = ausopr.Receipt.EffectiveGasPrice
	res.Receipt.From = ausopr.Receipt.From
	res.Receipt.GasUsed = ausopr.Receipt.GasUsed
	res.Receipt.LogsBloom = ausopr.Receipt.LogsBloom
	res.Receipt.Status = ausopr.Receipt.Status
	res.Receipt.To = ausopr.Receipt.To
	res.Receipt.TransactionHash = ausopr.Receipt.TransactionHash
	res.Receipt.TransactionIndex = ausopr.Receipt.TransactionIndex
	res.Receipt.Type = ausopr.Receipt.Type
	for _, log := range ausopr.Receipt.Logs {
		res.Receipt.Logs = append(res.Receipt.Logs, ReceiptLog{
			Address:          log.Address,
			Topics:           log.Topics,
			Data:             log.Data,
			BlockNumber:      log.BlockNumber,
			TransactionHash:  log.TransactionHash,
			TransactionIndex: log.TransactionIndex,
			BlockHash:        log.BlockHash,
			LogIndex:         log.LogIndex,
			Removed:          log.Removed,
		})
	}

	return res, nil

}
