package rpccalls

import (
	"github.com/san-lab/go4337/userop"
)

func PM_getPaymasterStubData(url, key string, usop *userop.UserOperation, context interface{},
	entrypoint string, chainId interface{}, provider string) ([]byte, error) {
	ar := &APIRequest{
		ID:      17,
		Jsonrpc: "2.0",
		Method:  "pm_getPaymasterStubData",
		Params:  []interface{}{usop.ToUserOpForApiV6(), entrypoint, chainId, context},
	}
	//pmad := &PMandDataResult{}
	bt, err := ApiCall(url, key, ar, nil)
	return bt, err
}

func PM_getPaymasterData(url, key string, usop *userop.UserOperation, context interface{}, entrypoint string,
	chainId interface{}, provider string) (*PMandDataResult, error) {
	ar := &APIRequest{
		ID:      17,
		Jsonrpc: "2.0",
		Method:  "pm_getPaymasterData",
		Params:  []interface{}{usop.ToUserOpForApiV6(), entrypoint, chainId, context},
	}
	pmad := &PMandDataResult{}
	_, err := ApiCall(url, key, ar, pmad)
	return pmad, err
}

type PMandDataResult struct {
	PaymasterAndData string `json:"paymasterAndData"`
}

type AlchemyPMContext struct {
	PolicyId string `json:"policyId"`
}

type BiconomyPMContext struct {
	Mode               string `json:"mode"`
	CalculateGasLimits bool   `json:"calculateGasLimits"`
	ExpiryDuration     int    `json:"expiryDuration"`
	SponsorshipInfo    struct {
		WebhookData struct {
		} `json:"webhookData"`
		SmartAccountInfo struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"smartAccountInfo"`
	} `json:"sponsorshipInfo"`
}

type PimlicoPMTokenContext struct {
	Token string `json:"token"`
}

/* Biconomy request and response
{
  "method": "pm_getPaymasterData",
  "params": [
        {
      "sender": "0x5A927A01a32cE02AF4B438b7848BceBc52C8Ea3e",
      "nonce": "0x1a",
      "initCode": "0x",
      "callData": "0x0000189a0000000000000000000000001758f42af7026fbbb559dc60ece0de3ef81f665e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002440d097c30000000000000000000000005a927a01a32ce02af4b438b7848bcebc52c8ea3e00000000000000000000000000000000000000000000000000000000",
      "maxFeePerGas": "0x05f5e100",
      "maxPriorityFeePerGas": "0",
      "verificationGasLimit": "0x01536f",
      "callGasLimit": "0x012da1",
      "preVerificationGas": "0x0bcb43",
    }, // unsigned user operation
    "0x5ff137d4b0fdcd49dca30c7cf57e578a026d2789", // entry point address
    "0x2105", // chain id
        {
      "mode": "SPONSORED",
      "calculateGasLimits": true,
      "expiryDuration": 300 // duration (secs) for which the generate paymasterAndData will be valid. Default duration is 300 secs.
      "sponsorshipInfo": {
        "webhookData": {},
        "smartAccountInfo": {
            "name": "BICONOMY",
            "version": "2.0.0"
        }
      }
    }
  ],
  "id": 1693369916,
  "jsonrpc": "2.0"
}

{
  "jsonrpc": "2.0",
  "id": 1693369916,
  "result": {
    "paymasterAndData": "0xabc..."
  }
}
*/

/* Pimlico
 ----V6 Request-----
 {
    "jsonrpc": "2.0",
    "method": "pm_getPaymasterData",
    "params": [
        {
            "sender":"0xb341FEAFaF71b09089d03B7D114599f8F491EE45",
            "nonce":"0x0",
            "initCode":"0x5de4839a76cf55d0c90e2061ef4386d962E15ae3296601cd0000000000000000000000000da6a956b9488ed4dd761e59f52fdc6c8068e6b5000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000084d1f57894000000000000000000000000d9ab5096a832b9ce79914329daee236f8eea039000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000014375cd3E53E18f65672E9d0Eb6AD174511b0BF98100000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "callData":"0x51945447",
            "callGasLimit":"0x115b5c0",
            "verificationGasLimit":"0x249f0",
            "preVerificationGas":"0xeb11",
            "maxPriorityFeePerGas":"0x12a05f200",
            "maxFeePerGas":"0x5b08082fa"
        },
        "0x0000000071727De22E5E9d8BAf0edAc6f37da032",
        "0x1",
        {
            "token": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
        }
    ],
    "id": 4337
}

----V6 Response-----
{
  "jsonrpc": "2.0",
  "result": {
      "paymasterAndData": "0x0000000000000039cd5e8aE05257CE51C473ddd101000066d1a1a4000000000000036cbd53842c5426634e7929541ec2318f3dcf7e0000000000000000000000000000c350000000000000000000000000000000000000000000000000000000009666598f0b846603deb0a8e59b78ba3dce9c3466394ccf07795d38ecf7925dfe12c07a022c27bb199099fa54de2f5e3e87dd9c581df52e9d3d199166a31124cc1227a9921b",
      "preVerificationGas": "0x350f7",
      "verificationGasLimit": "0x501ab",
      "callGasLimit": "0x212df"
  },
  "id": 4337
}

	----V7 Request-----
	{
    "jsonrpc": "2.0",
    "method": "pm_getPaymasterData",
    "params": [
        {
            "sender": "0x5a6b47F4131bf1feAFA56A05573314BcF44C9149",
            "nonce": "0x845ADB2C711129D4F3966735ED98A9F09FC4CE5700000000000000000000",
            "factory": "0xd703aaE79538628d27099B8c4f621bE4CCd142d5",
            "factoryData": "0xc5265d5d000000000000000000000000aac5d4240af87249b3f71bc8e4a2cae074a3e4190000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001243c3b752b01845ADb2C711129d4f3966735eD98a9F09fC4cE570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000014375d883Cb4afb913aC35c4B394468C4bC73d77C40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "callData": "0xe9ae5c53",
            "callGasLimit": "0x13880",
            "verificationGasLimit": "0x60B01",
            "preVerificationGas": "0xD3E3",
            "maxPriorityFeePerGas": "0x3B9ACA00",
            "maxFeePerGas": "0x7A5CF70D5",
            "paymaster": null,
            "paymasterVerificationGasLimit": null,
            "paymasterPostOpGasLimit": null,
            "paymasterData": null
        },
        "0x0000000071727De22E5E9d8BAf0edAc6f37da032",
        "0x1",
        {
            "token": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
        }
    ],
    "id": 4337
}
	----V7 Response-----
	{
  "jsonrpc": "2.0",
  "id": 4337,
  "result": {
      "paymaster": "0x0000000000000039cd5e8aE05257CE51C473ddd1",
      "paymasterData": "0x01000066d1a1a4000000000000036cbd53842c5426634e7929541ec2318f3dcf7e0000000000000000000000000000c350000000000000000000000000000000000000000000000000000000009666598f0b846603deb0a8e59b78ba3dce9c3466394ccf07795d38ecf7925dfe12c07a022c27bb199099fa54de2f5e3e87dd9c581df52e9d3d199166a31124cc1227a9921b"
    }
}

*/
