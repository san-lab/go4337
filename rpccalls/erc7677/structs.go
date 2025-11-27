package erc7677

import "github.com/ethereum/go-ethereum/common"

const GetPaymasterStubData_jsonV6 = `{"method":"pm_getPaymasterStubData","params":[{
      "sender": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
      "nonce": "0x2a",
      "initCode": "0x",
      "callData": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
      "callGasLimit": "0x0",
      "verificationGasLimit": "0x0",
      "preVerificationGas": "0x0",
      "maxFeePerGas": "0x0",
      "maxPriorityFeePerGas": "0x0"
    }, "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789", "0x14A34",
    {
      "policyId": "631528b0-d444-4a9b-a575-40dd3aa4a13a"
    }],"jsonrpc":"2.0", "id": 1}`

type Sponsor struct {
	name string
	icon string
}
type PaymasterStubDataResV6 struct {
	Sponsor          Sponsor `json:"sponsor,omitempty"`
	PaymasterAndData string  `json:"paymasterAndData"` //hex
}

const GetPaymasterStubData_jsonV7 = `{"method":"pm_getPaymasterStubData","params":[{
      "sender": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
      "nonce": "0x2a",
      "factory": "0x0BA5ED0c6AA8c49038F819E587E2633c4A9F428a",
      "factoryData": "0x",
      "callData": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
      "callGasLimit": "0x0",
      "verificationGasLimit": "0x0",
      "preVerificationGas": "0x0",
      "maxFeePerGas": "0x0",
      "maxPriorityFeePerGas": "0x0",
      "paymasterVerificationGasLimit": "0x0",
      "paymasterPostOpGasLimit": "0x0"
    }, "0x0000000071727De22E5E9d8BAf0edAc6f37da032", "0x14A34",
    {
      "policyId": "631528b0-d444-4a9b-a575-40dd3aa4a13a" 
    }],"jsonrpc":"2.0", "id": 1}`

type PaymasterStubDataResV7 struct {
	Paymaster                     common.Address
	PaymasterData                 string
	PaymasterVerificationGasLimit string
	PaymasterPostOpGasLimit       string
	Sponsor                       Sponsor
}
