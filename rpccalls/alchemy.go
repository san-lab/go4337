package rpccalls

import (
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

// USerOpV6 - standard

type AlchemyV7UserOp struct {
	Sender                        string `json:"sender"`
	Nonce                         string `json:"nonce"`
	Factory                       string `json:"factory"`
	FactoryData                   string `json:"factoryData"`
	CallData                      string `json:"callData"`
	CallGasLimit                  string `json:"callGasLimit"`
	VerificationGasLimit          string `json:"verificationGasLimit"`
	PreVerificationGas            string `json:"preVerificationGas"`
	MaxFeePerGas                  string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          string `json:"maxPriorityFeePerGas"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"`
	Signature                     string `json:"signature"`
	Paymaster                     string `json:"paymaster"`
	PaymasterData                 string `json:"paymasterData"`
}

/*

eth_sendUserOperation
Submits a user operation to a Bundler. If the request is successful, the endpoint will return a user operation hash that the caller can use to look up the status of the user operation. If it fails, or another error occurs, an error code and description will be returned.
eth_estimateUserOperationGas
Estimates the gas values for a user operation. It returns the preVerificationGas, verificationGasLimit, and callGasLimit values associated with the provided user operation.
eth_getUserOperationByHash
Returns a user operation based on the given user operation hash. It returns the user operation along with extra information including what block/transaction it was included in. If the operation has not yet been included, it will return null.
eth_getUserOperationReceipt
Returns a user operation receipt ( metadata associated with the given user operation ) based on the given user operation hash. It returns null if the user operation has not yet been included.
rundler_maxPriorityFeePerGas
Returns a fee per gas that is an estimate of how much users should set as a priority fee in UOs for Rundler endpoints.
eth_supportedEntryPoints
Returns a list of Entrypoint contract addresses supported by the bunder endpoints.

*/

func Alchemy_requestGasAndPaymasterAndData(url, key, policyID, entrypoint, dummysignature string,
	usop userop.UserOperation, overrides *AlchemyOverrides) (*AlchemyGasAndPaymasterDataResult, error) {
	ar := &APIRequest{
		ID:      4338,
		Jsonrpc: "2.0",
		Method:  "alchemy_requestGasAndPaymasterAndData",
		Params:  []interface{}{AlchemyReqGasAndPMandDataParams{policyID, entrypoint, dummysignature, usop.ToUserOpForApiV6(), overrides}},
	}
	state.Log("Alchemy Overrides:", overrides)
	agapad := &AlchemyGasAndPaymasterDataResult{}
	_, err := ApiCall(url, key, ar, agapad)
	return agapad, err
}

type AlchemyGasAndPaymasterDataResult struct {
	PaymasterAndData     string `json:"paymasterAndData"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	PreVerificationGas   string `json:"preVerificationGas"`
}

type AlchemyReqGasAndPMandDataParams struct {
	PolicyId       string                 `json:"policyId"`
	EntryPoint     string                 `json:"entryPoint"`
	DummySignature string                 `json:"dummySignature"`
	UserOperation  *userop.UserOpForApiV6 `json:"userOperation"`
	Overrides      *AlchemyOverrides      `json:"overrides"`
}

func Alchemy_requestPaymasterAndData(url, key, policyID, entrypoint string,
	usop userop.UserOperation) (*PMandDataResult, error) {
	ar := &APIRequest{
		ID:      4337,
		Jsonrpc: "2.0",
		Method:  "alchemy_requestPaymasterAndData",
		Params:  []interface{}{AlchemyReqPMandDatParams{policyID, entrypoint, usop.ToUserOpForApiV6()}},
	}
	pmad := &PMandDataResult{}
	_, err := ApiCall(url, key, ar, pmad)
	return pmad, err

}

type AlchemyReqPMandDatParams struct {
	PolicyId      string                 `json:"policyId"`
	EntryPoint    string                 `json:"entryPoint"`
	UserOperation *userop.UserOpForApiV6 `json:"userOperation"`
}

type AlchemyOverrides struct {
	/*
			{
		  "maxFeePerGas": "hex string" | { "multiplier": number },
		  "maxPriorityFeePerGas": "hex string" | { "multiplier": number },
		  "callGasLimit": "hex string" | { "multiplier": number },
		  "verificationGasLimit": "hex string" | { "multiplier": number },
		  "preVerificationGas": "hex string" | { "multiplier": number },
		}
	*/
	MaxFeePerGas         interface{} `json:"maxFeePerGas"`
	MaxPriorityFeePerGas interface{} `json:"maxPriorityFeePerGas"`
	CallGasLimit         interface{} `json:"callGasLimit"`
	VerificationGasLimit interface{} `json:"verificationGasLimit"`
	PreVerificationGas   interface{} `json:"preVerificationGas"`
}

type AlchemyOverrideMultiplier struct {
	Multiplier float64 `json:"multiplier"`
}
