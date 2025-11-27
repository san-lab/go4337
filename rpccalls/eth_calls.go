package rpccalls

import (
	"fmt"
	"log"

	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

// providers
const AlchemyProvider = "alchemy"
const StackUpProvider = "stackup"
const PimlicoProvider = "pimlico"
const BiconomyProvider = "biconomy"

func Eth_sendUserOperation(url, key string, usop *userop.UserOperation, entrypoint string, entrypointVersion int, provider string) (*string, error) {
	var usopdata interface{}
	switch entrypointVersion {
	case 6:
		usopdata = usop.ToUserOpForApiV6()
	case 7, 8:
		usopdata = usop.ToUserOpForApiV78(provider)
	default:
		return nil, fmt.Errorf("Unsupported entrypoint version: %d", entrypointVersion)
	}

	ar := &APIRequest{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_sendUserOperation",
		Params:  []interface{}{usopdata, entrypoint},
	}
	result := new(string)
	_, err := ApiCall(url, key, ar, result)
	if err != nil {
		return nil, fmt.Errorf("API Call error: %w", err)
	}
	return result, nil

}

func Eth_estimateUserOperationGas(url, key string, usop *userop.UserOperation, entrypoint string, entrypointVersion int, provider string) (any, error) {
	if state.DEBUG {
		log.Println(entrypoint, entrypointVersion)
	}
	var usopdata, result interface{}
	switch entrypointVersion {
	case 6:
		usopdata = usop.ToUserOpForApiV6()
		result = new(EthEstimateUserOperationGasResultV6)
	case 7, 8:
		usopdata = usop.ToUserOpForApiV78(provider)
		result = new(EthEstimateUserOperationGasResultV7)
		if provider == BiconomyProvider {
			result = new(EthEstimateUserOperationGasResultV7Biconomy)
		}
	default:
		return nil, fmt.Errorf("Unsupported entrypoint version: %d", entrypointVersion)
	}

	ar := &APIRequest{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_estimateUserOperationGas",
		Params:  []interface{}{usopdata, entrypoint},
	}

	_, err := ApiCall(url, key, ar, result)
	if err != nil {
		return nil, fmt.Errorf("API Call error: %w", err)
	}

	/*
		switch provider {
		case AlchemyProvider, PimlicoProvider:
			finalResult.CallGasLimit, _ = strconv.ParseUint(result.(*AlchemyEstimateGasCostResponse).CallGasLimit[2:], 16, 64)
			finalResult.VerificationGasLimit, _ = strconv.ParseUint(result.(*AlchemyEstimateGasCostResponse).VerificationGasLimit[2:], 16, 64)
			finalResult.PreVerificationGas, _ = strconv.ParseUint(result.(*AlchemyEstimateGasCostResponse).PreVerificationGas[2:], 16, 64)
		default:
			finalResult = result.(*EthEstimateUserOperationGasResultV6)
		}
	*/

	return result, nil

}

type EthEstimateUserOperationGasResultV6 struct {
	CallGasLimit         uint64 `json:"callGasLimit"`
	VerificationGasLimit uint64 `json:"verificationGasLimit"`
	PreVerificationGas   uint64 `json:"preVerificationGas"`
	ValidUntil           string `json:"validUntil"`
	ValidAfter           string `json:"validAfter"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
}

type EthEstimateUserOperationGasResultV7 struct {
	PreVerificationGas            string `json:"preVerificationGas"`
	VerificationGasLimit          string `json:"verificationGasLimit"`
	CallGasLimit                  string `json:"callGasLimit"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"`
}

type EthEstimateUserOperationGasResultV7Biconomy struct {
	PreVerificationGas            uint64 `json:"preVerificationGas"`
	VerificationGasLimit          uint64 `json:"verificationGasLimit"`
	CallGasLimit                  uint64 `json:"callGasLimit"`
	PaymasterVerificationGasLimit uint64 `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       uint64 `json:"paymasterPostOpGasLimit"`
}

type AlchemyEstimateGasCostResponse struct {
	PreVerificationGas   string `json:"preVerificationGas"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
}

func Eth_getUserOperationByHash(url, key string, hash string, provider string) (*UserOperationByHashResult_ASP, error) {
	ar := &APIRequest{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getUserOperationByHash",
		Params:  []interface{}{hash},
	}
	if provider == BiconomyProvider {
		bres := new(UserOperationByHashResult_Biconomy)
		_, err := ApiCall(url, key, ar, bres)
		if err != nil {
			return nil, fmt.Errorf("API Call error: %w", err)
		}
		return RepackageUSOPResult(bres)
	}

	result := new(UserOperationByHashResult_ASP)
	_, err := ApiCall(url, key, ar, result)
	if err != nil {
		return nil, fmt.Errorf("API Call error: %w", err)
	}
	return result, nil

}

func Eth_getUserOperationReceipt(url, key string, hash string, provider string) (*UserOperationReceipt_Alchemy, error) {
	//fmt.Println("eth_getUserOperationReceipt not implemented")
	ar := &APIRequest{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getUserOperationReceipt",
		Params:  []interface{}{hash},
	}

	if provider == BiconomyProvider {
		bres := new(UserOperationReceipt_Biconomy)
		_, err := ApiCall(url, key, ar, bres)
		if err != nil {
			return nil, fmt.Errorf("API Call error: %w", err)
		}
		return RepackageUSOPReceiptResult(bres)
	}
	arec := new(UserOperationReceipt_Alchemy)
	_, err := ApiCall(url, key, ar, arec)
	return arec, err

}

func Eth_supportedEntryPoints(url, key string) (*[]string, error) {
	ar := &APIRequest{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_supportedEntryPoints",
		Params:  []interface{}{},
	}
	result := &[]string{}
	_, err := ApiCall(url, key, ar, result)
	if err != nil {
		return nil, fmt.Errorf("API Call error: %w", err)
	}
	return result, nil

}
