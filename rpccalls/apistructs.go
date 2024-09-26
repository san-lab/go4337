package rpccalls

import (
	"fmt"
)

/*
{
  "method": "eth_getUserOperationByHash",
  "params": [
    "0x10b68dfde047c51e2f4a500281a715710c7b1c0517b87ad0b48229042f61ac0e"
  ],
  "id": 1693369916,
  "jsonrpc": "2.0"
}
*/

// Works for Alchemy, StackUp,Pimlico
type UserOperationByHashResult_ASP struct {
	UserOperation struct {
		Sender               string `json:"sender"`
		Nonce                string `json:"nonce"`
		InitCode             string `json:"initCode"`
		CallData             string `json:"callData"`
		CallGasLimit         string `json:"callGasLimit"`
		VerificationGasLimit string `json:"verificationGasLimit"`
		PreVerificationGas   string `json:"preVerificationGas"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		PaymasterAndData     string `json:"paymasterAndData"`
		Signature            string `json:"signature"`
	} `json:"userOperation"`
	EntryPoint      string `json:"entryPoint"`
	BlockNumber     string `json:"blockNumber"`
	BlockHash       string `json:"blockHash"`
	TransactionHash string `json:"transactionHash"`
}

func SupportedEntryPoints(url, key string) (res *SupportedEntryPointsResult, err error) {
	ar := &APIRequest{Method: "eth_supportedEntryPoints", Params: []interface{}{}, Jsonrpc: "2.0", ID: 1}
	res = &SupportedEntryPointsResult{}
	_, err = ApiCall(url, key, ar, res)
	if err != nil {
		return nil, fmt.Errorf("error in API call: %w", err)
	}
	return res, nil
}

type SupportedEntryPointsResult []string

type UserOperationReceipt_Alchemy struct {
	UserOpHash    string       `json:"userOpHash"`
	EntryPoint    string       `json:"entryPoint"`
	Sender        string       `json:"sender"`
	Nonce         string       `json:"nonce"`
	Paymaster     string       `json:"paymaster"`
	ActualGasCost string       `json:"actualGasCost"`
	ActualGasUsed string       `json:"actualGasUsed"`
	Success       bool         `json:"success"`
	Reason        string       `json:"reason"`
	Logs          []ReceiptLog `json:"logs"`
	Receipt       struct {
		TransactionHash   string      `json:"transactionHash"`
		TransactionIndex  string      `json:"transactionIndex"`
		BlockHash         string      `json:"blockHash"`
		BlockNumber       string      `json:"blockNumber"`
		From              string      `json:"from"`
		To                string      `json:"to"`
		CumulativeGasUsed string      `json:"cumulativeGasUsed"`
		GasUsed           string      `json:"gasUsed"`
		ContractAddress   interface{} `json:"contractAddress"`
		Logs              []struct {
			Address          string   `json:"address"`
			Topics           []string `json:"topics"`
			Data             string   `json:"data"`
			BlockHash        string   `json:"blockHash"`
			BlockNumber      string   `json:"blockNumber"`
			TransactionHash  string   `json:"transactionHash"`
			TransactionIndex string   `json:"transactionIndex"`
			LogIndex         string   `json:"logIndex"`
			Removed          bool     `json:"removed"`
		} `json:"logs"`
		Status            string `json:"status"`
		LogsBloom         string `json:"logsBloom"`
		Type              string `json:"type"`
		EffectiveGasPrice string `json:"effectiveGasPrice"`
	} `json:"receipt"`
}

type ReceiptLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

type UserOperationReceipt_StackUp struct {
	UserOpHash    string `json:"userOpHash"`
	Sender        string `json:"sender"`
	Paymaster     string `json:"paymaster"`
	Nonce         string `json:"nonce"`
	Success       bool   `json:"success"`
	ActualGasCost string `json:"actualGasCost"`
	ActualGasUsed string `json:"actualGasUsed"`
	From          string `json:"from"`
	Receipt       struct {
		BlockHash         string `json:"blockHash"`
		BlockNumber       string `json:"blockNumber"`
		From              string `json:"from"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
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
		LogsBloom         string `json:"logsBloom"`
		TransactionHash   string `json:"transactionHash"`
		TransactionIndex  string `json:"transactionIndex"`
		EffectiveGasPrice string `json:"effectiveGasPrice"`
	} `json:"receipt"`
	Logs []struct {
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
}

type UserOperationReceipt_Pimlico struct {
	UserOpHash    string `json:"userOpHash"`
	EntryPoint    string `json:"entryPoint"`
	Sender        string `json:"sender"`
	Nonce         string `json:"nonce"`
	Paymaster     string `json:"paymaster"`
	ActualGasUsed string `json:"actualGasUsed"`
	ActualGasCost string `json:"actualGasCost"`
	Success       bool   `json:"success"`
	Logs          []struct {
		LogIndex         string   `json:"logIndex"`
		TransactionIndex string   `json:"transactionIndex"`
		TransactionHash  string   `json:"transactionHash"`
		BlockHash        string   `json:"blockHash"`
		BlockNumber      string   `json:"blockNumber"`
		Address          string   `json:"address"`
		Data             string   `json:"data"`
		Topics           []string `json:"topics"`
	} `json:"logs"`
	Receipt struct {
		TransactionHash   string      `json:"transactionHash"`
		TransactionIndex  string      `json:"transactionIndex"`
		BlockHash         string      `json:"blockHash"`
		BlockNumber       string      `json:"blockNumber"`
		From              string      `json:"from"`
		To                string      `json:"to"`
		CumulativeGasUsed string      `json:"cumulativeGasUsed"`
		GasUsed           string      `json:"gasUsed"`
		ContractAddress   interface{} `json:"contractAddress"`
		Logs              []struct {
			LogIndex         string   `json:"logIndex"`
			TransactionIndex string   `json:"transactionIndex"`
			TransactionHash  string   `json:"transactionHash"`
			BlockHash        string   `json:"blockHash"`
			BlockNumber      string   `json:"blockNumber"`
			Address          string   `json:"address"`
			Data             string   `json:"data"`
			Topics           []string `json:"topics"`
		} `json:"logs"`
		LogsBloom         string `json:"logsBloom"`
		Status            string `json:"status"`
		EffectiveGasPrice string `json:"effectiveGasPrice"`
	} `json:"receipt"`
}

func UserOperationReceipt(url, key string, userOpHash string) (uorec *UserOperationReceipt_StackUp, err error) {
	ar := &APIRequest{Method: "eth_getUserOperationReceipt", Params: []interface{}{userOpHash}, Jsonrpc: "2.0", ID: 1}
	_, err = ApiCall(url, key, ar, uorec)
	if err != nil {
		return nil, fmt.Errorf("could not make API call: %w", err)
	}
	return uorec, nil
}
