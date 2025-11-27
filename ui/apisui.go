package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

var ApiCallsItem = &Item{Label: "API Calls", Details: "Call APIs"}
var ApiKeyItem = &Item{Label: "API Key", Details: "To be appended or replace the {{.}} in URL"}
var ApiURLItem = &Item{Label: "API URL", Details: "If len(api key) > 0, the key will be appended after / or put in place of {{.}}"}
var StackUpPMApiItem = &Item{Label: "StackUp Paymaster API"}
var StackUpBNApiItem = &Item{Label: "StackUp Bundler API"}
var ApiMethodItem = &Item{Label: "API Method (free-hand)"}
var StdApiItem = &Item{Label: "Standard (eth_/pm_) API methods"}

var ApiUserOpItem = &Item{Label: "User Operation"}
var ApiCallItem = &Item{Label: "Call API"}
var ZkSyncEraItem = &Item{Label: "ZkSync Era"}

func ApiCallsUI(usop *userop.UserOperation) {
	var AlchemyUIItem = &Item{Label: "Alchemy API"}
	var StackUpUIItem = &Item{Label: "StackUp API"}
	var BiconomyUIItem = &Item{Label: "Biconomy API"}
	var PimlicoUIItem = &Item{Label: "Pimlico API"}
	var CustomUIItem = &Item{Label: "Custom API"}
	var items = []*Item{AlchemyUIItem, StackUpUIItem, BiconomyUIItem, PimlicoUIItem, CustomUIItem, ZkSyncEraItem, Back}
	for {
		spr := promptui.Select{Label: "APIs", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case AlchemyUIItem.Label:
			ProvAPIUI(rpccalls.AlchemyProvider, usop)
		case StackUpUIItem.Label:
			ProvAPIUI(rpccalls.StackUpProvider, usop)
		case BiconomyUIItem.Label:
			ProvAPIUI(rpccalls.BiconomyProvider, usop)
		case PimlicoUIItem.Label:
			ProvAPIUI(rpccalls.PimlicoProvider, usop)
		case CustomUIItem.Label:
			CustomAPIUI(usop)
		case ZkSyncEraItem.Label:
			ZkSyncEraUI()
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

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
Returns a list of Entrypoint contr
*/

var eth_sendUserOperationItem = &Item{Label: "eth_sendUserOperation"}
var eth_estimateUserOperationGasItem = &Item{Label: "eth_estimateUserOperationGas"}
var eth_getUserOperationByHashItem = &Item{Label: "eth_getUserOperationByHash"}
var eth_getUserOperationReceiptItem = &Item{Label: "eth_getUserOperationReceipt"}
var eth_supportedEntryPointsItem = &Item{Label: "eth_supportedEntryPoints"}

func addETHMethods(items []*Item) []*Item {
	items = append(items, eth_sendUserOperationItem, eth_estimateUserOperationGasItem, eth_getUserOperationByHashItem, eth_getUserOperationReceiptItem, eth_supportedEntryPointsItem)
	return items
}

var rundler_maxPriorityFeePerGasItem = &Item{Label: "rundler_maxPriorityFeePerGas"}                   //Alchemy specific
var alchemy_requestPaymasterAndDataItem = &Item{Label: "alchemy_requestPaymentAndData"}               //Alchemy specific
var alchemy_requestGasAndPaymasterAndDataItem = &Item{Label: "alchemy_requestGasAndPaymasterAndData"} //Alchemy specific

var pm_getPaymasterDataItem = &Item{Label: "pm_getPaymasterData"}
var pm_getPaymasterStubDataItem = &Item{Label: "pm_getPaymasterStubData"}
var pm_sponsorUserOperationItem = &Item{Label: "pm_SponsorUserOperation"} //StackUp specific
var pm_accountsItem = &Item{Label: "pm_accounts"}                         //StackUp specific
var pm_getFeeQuoteOrDataItem = &Item{Label: "pm_getFeeQuoteOrData"}       //Biconomy

func addPMMethods(items []*Item) []*Item {
	items = append(items, pm_getPaymasterDataItem, pm_getPaymasterStubDataItem)
	return items
}

var pimlico_getTokenQuotesItem = &Item{Label: "pimlico_GetTokenQuotes"} //Pimlico specific

var APIKeyItemMap = map[string]*Item{} //&Item{Label: "API Key"}
var APIURLItemMap = map[string]*Item{} //&Item{Label: "API URL"}

var Alchemy_requestGasAndPaymasterAndDataItem = &Item{Label: "Alchemy_requestGasAndPaymasterAndData"}
var Alchemy_requestPaymasterAndDataItem = &Item{Label: "Alchemy_requestPaymasterAndData"}

func ProvAPIUI(provider string, usop *userop.UserOperation) {
	state.Log("ProvAPIUI", provider)
	var key, url string
	var ok, allOk bool
	var overrides *rpccalls.AlchemyOverrides

	if usop != nil {
		ApiUserOpItem.Value = usop
	}
	for {
		APIKeyItem, ok3 := APIKeyItemMap[provider]
		if !ok3 { //lazy initialization
			APIKeyItem = &Item{Label: provider + "API Key"}
			APIKeyItemMap[provider] = APIKeyItem
		}
		APIURLItem, ok4 := APIURLItemMap[provider]
		if !ok4 { //lazy initialization
			APIURLItem = &Item{Label: provider + "API URL", Details: "If len(api key) > 0, the key will be appended after / or put in place of {{.}}"}
			APIURLItemMap[provider] = APIURLItem

		}
		var entrypointaddress = EntryPointItem.Value.(ecommon.Address).String()
		var entrypointversion = 0
		if entrypointaddress == entrypoint.E6Address.String() {
			entrypointversion = 6
		} else if entrypointaddress == entrypoint.E7Address.String() {
			entrypointversion = 7
		} else if entrypointaddress == entrypoint.E8Address.String() {
			entrypointversion = 8
		}

		if APIKeyItem.Value != nil {
			key, ok = APIKeyItem.Value.(string)
			allOk = ok
		}
		if APIURLItem.Value != nil {
			url, ok = APIURLItem.Value.(string)
			allOk = allOk && ok
		}
		if ApiUserOpItem.Value != nil {
			usop, ok = ApiUserOpItem.Value.(*userop.UserOperation)
			allOk = allOk && ok
		} else {
			usop = nil
			allOk = false
		}
		items := []*Item{EntryPointItem, APIKeyItem, APIURLItem, ApiUserOpItem}
		if allOk {
			items = addETHMethods(items)
			items = addPMMethods(items)
			if provider == rpccalls.AlchemyProvider {
				items = append(items, rundler_maxPriorityFeePerGasItem, alchemy_requestPaymasterAndDataItem, alchemy_requestGasAndPaymasterAndDataItem)
			}
			if provider == rpccalls.StackUpProvider {
				items = append(items, pm_sponsorUserOperationItem, pm_accountsItem)
			}

		}
		items = append(items, Back)

		spr := promptui.Select{Label: provider + " API", Items: items, Templates: ItemTemplate, Size: len(items) + 3}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case EntryPointItem.Label:
			EntryPointUI()
		case ApiUserOpItem.Label:
			SelectUserOpUI(ApiUserOpItem)
		case APIKeyItem.Label:
			_, name, key, good := StringFromDictionaryUI(provider + state.ApiKeysLabel)
			if good {
				APIKeyItem.Value = key
				APIKeyItem.Details = name
			}
		case APIURLItem.Label:
			_, name, url, good := StringFromDictionaryUI(provider + state.ApiEndpointsLabel)
			if good {
				APIURLItem.Value = url
				APIURLItem.Details = name
			}
		case eth_sendUserOperationItem.Label:
			res, err := rpccalls.Eth_sendUserOperation(url, key, usop, entrypointaddress, entrypointversion, provider)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				fmt.Println("API call result:", *res)
				state.AddToDictionatyWithIndice("User Operation Hash", provider, *res)

			}

		case eth_supportedEntryPointsItem.Label:
			res, err := rpccalls.Eth_supportedEntryPoints(url, key)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				fmt.Println("Supported Entry Points:")
				for _, ep := range *res {
					fmt.Println(ep)
				}
			}
		case eth_estimateUserOperationGasItem.Label:
			res, err := rpccalls.Eth_estimateUserOperationGas(url, key, usop, entrypointaddress, entrypointversion, provider)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				IncorporateGasParametersUI(usop, res, entrypointversion, provider)
			}
		case eth_getUserOperationByHashItem.Label:
			_, _, hash, ok := StringFromDictionaryUI("User Operation Hash")
			if !ok {
				continue
			}
			uopr, err := rpccalls.Eth_getUserOperationByHash(url, key, hash, provider)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				fmt.Println("Tx Hash:", uopr.TransactionHash)
				fmt.Println("Block Number:", uopr.BlockNumber)
				fmt.Println("Entry Point:", uopr.EntryPoint)
			}

		case eth_getUserOperationReceiptItem.Label:
			_, _, hash, ok := StringFromDictionaryUI("User Operation Hash")
			if !ok {
				continue
			}
			arec, err := rpccalls.Eth_getUserOperationReceipt(url, key, hash, provider)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				fmt.Println("User Op Success:", *&arec.Success)
				fmt.Println("Tx Hash:", *&arec.Receipt.TransactionHash)
				fmt.Println("Block Number:", *&arec.Receipt.BlockNumber)
				if YesNoPromptUI("Print Receipt?") {
					pbt, _ := json.MarshalIndent(arec, "", "  ")
					fmt.Println(string(pbt))
				}
			}
		case pm_getPaymasterDataItem.Label:
			context, chainId, ok := GetPMContext(provider)
			if !ok {
				continue
			}
			pmad, err := rpccalls.PM_getPaymasterData(url, key, usop, context, entrypointaddress, chainId, provider)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				if YesNoPromptUI("Incorporate Paymaster and Data") {
					IncorporatePMandData(usop, pmad)
				}
			}
		case pm_getPaymasterStubDataItem.Label:
			context, chainId, ok := GetPMContext(provider)
			if !ok {
				continue
			}
			rpccalls.PM_getPaymasterStubData(url, key, usop, context, entrypointaddress, chainId, provider)

		case rundler_maxPriorityFeePerGasItem.Label:
			fmt.Println("Not implemented yet:", sel)
		case alchemy_requestPaymasterAndDataItem.Label:
			_, _, policyID, ok := StringFromDictionaryUI("Alchemy Policy Id")
			if !ok {
				fmt.Println("Policy ID not set")
				continue
			}
			pmad, err := rpccalls.Alchemy_requestPaymasterAndData(url, key, policyID, entrypointaddress, *usop)
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				if YesNoPromptUI("Incorporate Paymaster and Data") {
					IncorporatePMandData(usop, pmad)
				}
			}

		case alchemy_requestGasAndPaymasterAndDataItem.Label:
			if YesNoPromptUI("With Overrides?") {
				overrides = AlchemyOverridesUI(overrides)
			}
			_, _, policyID, ok1 := StringFromDictionaryUI("Alchemy Policy Id")
			if !ok1 {
				fmt.Println("Policy ID not set")
				continue
			}
			if entrypointaddress == entrypoint.EntryPointAddressV6 {
				gapad, err := rpccalls.Alchemy_requestGasAndPaymasterAndData(
					url, key, policyID, entrypointaddress,
					fmt.Sprintf("0x%x", usop.Signature), *usop, overrides)
				if err != nil {
					fmt.Println("Error making API call:", err)
				} else {
					IncorporateAlchemyGapadUI(usop, gapad)
				}
			} else {
				gapad, err := rpccalls.Alchemy_requestGasAndPaymasterAndDataV7(
					url, key, policyID, entrypointaddress,
					fmt.Sprintf("0x%x", usop.Signature), *usop, overrides)
				if err != nil {
					fmt.Println("Error making API call:", err)
				} else {
					IncorporateAlchemyGapadUIV7(usop, gapad)
				}
			}

		case pm_sponsorUserOperationItem.Label:
			res, err := rpccalls.StackUpPMPayCall(url, key, usop.ToUserOpForApiV6())
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else {
				pbt, _ := json.MarshalIndent(res, "", "  ")
				fmt.Println(string(pbt))
				if YesNoPromptUI("Incorporate Paymaster and Gas Data?") {
					rpccalls.IncorporateStackUpPMResToUserOp(usop, res)
				}
			}

		default:
			fmt.Println("Not reachable yet:", sel)

		}

	}
}

func GetPMContext(provider string) (context, chainId interface{}, ok bool) {
	chi, ok2 := ChainIDItem.Value.(*big.Int)
	if !ok2 {
		fmt.Println("Chain ID not set")
		return nil, nil, false
	}
	switch provider {
	case rpccalls.AlchemyProvider:
		_, _, policyid, ok1 := StringFromDictionaryUI("Alchemy Policy Id")
		if !ok1 {
			fmt.Println("Policy ID not set")
			return nil, nil, false
		}

		return &rpccalls.AlchemyPMContext{PolicyId: policyid}, fmt.Sprintf("0x%x", chi), true
	case rpccalls.PimlicoProvider:
		return struct{}{}, chi, true
	default:
		fmt.Println("Not implemented yet:", provider)
		return nil, nil, false

	}
}

func CustomAPIUI(usop *userop.UserOperation) {
	fmt.Println("Custom API UI not implemented yet")
}

func IncorporateGasParametersUI(usop *userop.UserOperation, res any, endpointVersion int, provider string) {
	switch endpointVersion {
	case 6:
		if res6, ok := res.(*rpccalls.EthEstimateUserOperationGasResultV6); ok {
			IncorporateGasParametersV6UI(usop, res6)
		} else {
			fmt.Println("error casting result. entrypoint version:", endpointVersion)
			return
		}
	case 7, 8:
		if res7, ok := res.(*rpccalls.EthEstimateUserOperationGasResultV7); ok {
			IncorporateGasParametersV7UI(usop, res7)
		} else if res7b, ok := res.(*rpccalls.EthEstimateUserOperationGasResultV7Biconomy); ok {
			IncorporateGasParametersV7UI(usop,
				&rpccalls.EthEstimateUserOperationGasResultV7{
					PreVerificationGas:            fmt.Sprint(res7b.PreVerificationGas),
					VerificationGasLimit:          fmt.Sprint(res7b.VerificationGasLimit),
					PaymasterVerificationGasLimit: fmt.Sprint(res7b.PaymasterVerificationGasLimit),
					PaymasterPostOpGasLimit:       fmt.Sprint(res7b.PaymasterPostOpGasLimit),
					CallGasLimit:                  fmt.Sprint(res7b.CallGasLimit),
				},
			)
		} else {
			fmt.Println("error casting result. entrypoint version:", endpointVersion)
			return
		}

	}
}

func IncorporateGasParametersV7UI(usop *userop.UserOperation, res *rpccalls.EthEstimateUserOperationGasResultV7) {
	IncorporateAllItem := &Item{Label: "Incorporate All"}
	all := false
	CallGasLimitItem := &Item{Label: fmt.Sprintf("Call Gas Limit(%v) %v", usop.CallGasLimit, res.CallGasLimit), Value: false}
	PreVerificationGasItem := &Item{Label: fmt.Sprintf("Pre Verification Gas(%v) %v", usop.PreVerificationGas, res.PreVerificationGas), Value: false}
	VerificationGasLimitItem := &Item{Label: fmt.Sprintf("Verification Gas Limit(%v) %v", usop.VerificationGasLimit, res.VerificationGasLimit), Value: false}
	PaymasterVerificationGasLimitItem := &Item{Label: fmt.Sprintf("Paymaster Verification Gas Limit(%v) %v", usop.PaymasterVerificationGasLimit, res.PaymasterVerificationGasLimit), Value: false}
	PaymasterPostOpGasLimitItem := &Item{Label: fmt.Sprintf("Paymaster PostOp Gas Limit(%v) %v", usop.PaymasterPostOpGasLimit, res.PaymasterPostOpGasLimit), Value: false}
	items := []*Item{IncorporateAllItem, CallGasLimitItem, PreVerificationGasItem, VerificationGasLimitItem, PaymasterVerificationGasLimitItem, PaymasterPostOpGasLimitItem, Set, Back}
	for {
		spr := promptui.Select{Label: "Gas Parameters", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case Set.Label, IncorporateAllItem.Label:
			all = (sel == IncorporateAllItem.Label)
			if CallGasLimitItem.Value.(bool) || all {
				usop.CallGasLimit = ConvHexOrZero(res.CallGasLimit)
			}
			if PreVerificationGasItem.Value.(bool) || all {
				usop.PreVerificationGas = ConvHexOrZero(res.PreVerificationGas)
			}
			if VerificationGasLimitItem.Value.(bool) || all {
				usop.VerificationGasLimit = ConvHexOrZero(res.VerificationGasLimit)
			}
			if PaymasterPostOpGasLimitItem.Value.(bool) || all {
				usop.PaymasterPostOpGasLimit = ConvHexOrZero(res.PaymasterPostOpGasLimit)
			}
			if PaymasterVerificationGasLimitItem.Value.(bool) || all {
				usop.PaymasterVerificationGasLimit = ConvHexOrZero(res.PaymasterVerificationGasLimit)
			}
			state.Save()
			return

		default:
			it, _ := GetItem(sel, items)
			it.Value = !it.Value.(bool)

		}
	}

}

func IncorporateGasParametersV6UI(usop *userop.UserOperation, res *rpccalls.EthEstimateUserOperationGasResultV6) {
	IncorporateAllItem := &Item{Label: "Incorporate All"}
	all := false
	CallGasLimitItem := &Item{Label: fmt.Sprintf("Call Gas Limit(%v) %v", usop.CallGasLimit, res.CallGasLimit), Value: false}
	PreVerificationGasItem := &Item{Label: fmt.Sprintf("Pre Verification Gas(%v) %v", usop.PreVerificationGas, res.PreVerificationGas), Value: false}
	VerificationGasLimitItem := &Item{Label: fmt.Sprintf("Verification Gas Limit(%v) %v", usop.VerificationGasLimit, res.VerificationGasLimit), Value: false}
	ValidUntilItem := &Item{Label: fmt.Sprintf("Valid Until %v", res.ValidUntil), Value: false}
	ValidAfterItem := &Item{Label: fmt.Sprintf("Valid After %v", res.ValidAfter), Value: false}
	MaxPriorityFeePerGasItem := &Item{Label: fmt.Sprintf("Max Priority Fee Per Gas(%v) %v", usop.MaxPriorityFeePerGas, res.MaxPriorityFeePerGas), Value: false}
	MaxFeePerGasItem := &Item{Label: fmt.Sprintf("Max Fee Per Gas(%v) %v", usop.MaxFeePerGas, res.MaxFeePerGas), Value: false}
	items := []*Item{IncorporateAllItem, CallGasLimitItem, PreVerificationGasItem, VerificationGasLimitItem, ValidUntilItem, ValidAfterItem, MaxPriorityFeePerGasItem, MaxFeePerGasItem, Set, Back}
	for {
		spr := promptui.Select{Label: "Gas Parameters", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case Set.Label, IncorporateAllItem.Label:
			all = (sel == IncorporateAllItem.Label)
			if CallGasLimitItem.Value.(bool) || all {
				usop.CallGasLimit = res.CallGasLimit
			}
			if PreVerificationGasItem.Value.(bool) || all {
				usop.PreVerificationGas = res.PreVerificationGas
			}
			if VerificationGasLimitItem.Value.(bool) || all {
				usop.VerificationGasLimit = res.VerificationGasLimit
			}
			if MaxPriorityFeePerGasItem.Value.(bool) || all {
				usop.MaxPriorityFeePerGas = ConvHexOrZero(res.MaxPriorityFeePerGas)
			}
			if MaxFeePerGasItem.Value.(bool) || all {
				usop.MaxFeePerGas = ConvHexOrZero(res.MaxFeePerGas)
			}
			state.Save()
			return

		default:
			it, _ := GetItem(sel, items)
			it.Value = !it.Value.(bool)

		}
	}

}

func IncorporateAlchemyGapadUI(usop *userop.UserOperation, gapad *rpccalls.AlchemyGasAndPaymasterDataResultV6) {
	IncorporateAllItem := &Item{Label: "Incorporate All"}
	PaymasterDataItem := &Item{Label: "Paymaster and PM Data", Value: false}
	CallGasLimitItem := &Item{Label: fmt.Sprintf("Call Gas Limit(%v/0x%x) %s", usop.CallGasLimit, usop.CallGasLimit, gapad.CallGasLimit), Value: false}
	VerificationGasLimitItem := &Item{Label: fmt.Sprintf("Verification Gas Limit(%v/0x%x) %s", usop.VerificationGasLimit, usop.VerificationGasLimit, gapad.VerificationGasLimit), Value: false}
	MaxPriorityFeePerGasItem := &Item{Label: fmt.Sprintf("Max Priority Fee Per Gas(%v/0x%x) %s", usop.MaxPriorityFeePerGas, usop.MaxPriorityFeePerGas, gapad.MaxPriorityFeePerGas), Value: false}
	MaxFeePerGasItem := &Item{Label: fmt.Sprintf("Max Fee Per Gas(%v/0x%x) %s", usop.MaxFeePerGas, usop.MaxFeePerGas, gapad.MaxFeePerGas), Value: false}
	PreVerificationGasItem := &Item{Label: fmt.Sprintf("Pre Verification Gas(%v/0x%x) %s", usop.PreVerificationGas, usop.PreVerificationGas, gapad.PreVerificationGas), Value: false}

	items := []*Item{IncorporateAllItem, PaymasterDataItem, CallGasLimitItem, VerificationGasLimitItem, MaxPriorityFeePerGasItem,
		MaxFeePerGasItem, PreVerificationGasItem, Set, Back}
	for {
		spr := promptui.Select{Label: "Select parameters to incorporate", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case Set.Label, IncorporateAllItem.Label:
			all := (sel == IncorporateAllItem.Label)
			if PaymasterDataItem.Value.(bool) || all {
				IncorporatePMandData(usop, &rpccalls.PMandDataResult{PaymasterAndData: gapad.PaymasterAndData})

			}
			if CallGasLimitItem.Value.(bool) || all {
				usop.CallGasLimit = ConvHexOrZero(gapad.CallGasLimit)
			}
			if VerificationGasLimitItem.Value.(bool) || all {
				usop.VerificationGasLimit = ConvHexOrZero(gapad.VerificationGasLimit)
			}
			if MaxPriorityFeePerGasItem.Value.(bool) || all {
				usop.MaxPriorityFeePerGas = ConvHexOrZero(gapad.MaxPriorityFeePerGas)
			}
			if MaxFeePerGasItem.Value.(bool) || all {
				usop.MaxFeePerGas = ConvHexOrZero(gapad.MaxFeePerGas)
			}
			if PreVerificationGasItem.Value.(bool) || all {
				usop.PreVerificationGas = ConvHexOrZero(gapad.PreVerificationGas)
			}
			state.Save()
			return

		default:
			it, _ := GetItem(sel, items)
			it.Value = !it.Value.(bool)

		}
	}
}

func IncorporateAlchemyGapadUIV7(usop *userop.UserOperation, gapad *rpccalls.AlchemyGasAndPaymasterDataResultV7) {
	IncorporateAllItem := &Item{Label: "Incorporate All"}
	PaymasterItem := &Item{Label: "Paymaster", Value: false}
	PaymasterDataItem := &Item{Label: "Paymaster Data", Value: false}
	CallGasLimitItem := &Item{Label: fmt.Sprintf("Call Gas Limit(%v/0x%x) %s", usop.CallGasLimit, usop.CallGasLimit, gapad.CallGasLimit), Value: false}
	VerificationGasLimitItem := &Item{Label: fmt.Sprintf("Verification Gas Limit(%v/0x%x) %s", usop.VerificationGasLimit, usop.VerificationGasLimit, gapad.VerificationGasLimit), Value: false}
	MaxPriorityFeePerGasItem := &Item{Label: fmt.Sprintf("Max Priority Fee Per Gas(%v/0x%x) %s", usop.MaxPriorityFeePerGas, usop.MaxPriorityFeePerGas, gapad.MaxPriorityFeePerGas), Value: false}
	MaxFeePerGasItem := &Item{Label: fmt.Sprintf("Max Fee Per Gas(%v/0x%x) %s", usop.MaxFeePerGas, usop.MaxFeePerGas, gapad.MaxFeePerGas), Value: false}
	PreVerificationGasItem := &Item{Label: fmt.Sprintf("Pre Verification Gas(%v/0x%x) %s", usop.PreVerificationGas, usop.PreVerificationGas, gapad.PreVerificationGas), Value: false}
	PaymasterVerificationGasLimitItem := &Item{Label: fmt.Sprintf("Paymaster Verification Gas Limit(%v/0x%x) %s", usop.PaymasterVerificationGasLimit, usop.PaymasterVerificationGasLimit, gapad.PaymasterVerificationGasLimit), Value: false}

	items := []*Item{IncorporateAllItem, PaymasterVerificationGasLimitItem, PaymasterItem, PaymasterDataItem, CallGasLimitItem, VerificationGasLimitItem, MaxPriorityFeePerGasItem,
		MaxFeePerGasItem, PreVerificationGasItem, Set, Back}
	for {
		spr := promptui.Select{Label: "Select parameters to incorporate", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case Set.Label, IncorporateAllItem.Label:
			all := (sel == IncorporateAllItem.Label)
			if PaymasterDataItem.Value.(bool) || all {
				fmt.Println("incorporating all:", all)
				p := ecommon.HexToAddress(gapad.Paymaster)
				usop.Paymaster = &p
				usop.PaymasterData, _ = hex.DecodeString(strings.TrimPrefix(gapad.PaymasterData, "0x"))
			}
			if CallGasLimitItem.Value.(bool) || all {
				usop.CallGasLimit = ConvHexOrZero(gapad.CallGasLimit)
			}
			if VerificationGasLimitItem.Value.(bool) || all {
				usop.VerificationGasLimit = ConvHexOrZero(gapad.VerificationGasLimit)
			}
			if MaxPriorityFeePerGasItem.Value.(bool) || all {
				usop.MaxPriorityFeePerGas = ConvHexOrZero(gapad.MaxPriorityFeePerGas)
			}
			if MaxFeePerGasItem.Value.(bool) || all {
				usop.MaxFeePerGas = ConvHexOrZero(gapad.MaxFeePerGas)
			}
			if PreVerificationGasItem.Value.(bool) || all {
				usop.PreVerificationGas = ConvHexOrZero(gapad.PreVerificationGas)
			}

			usop.PaymasterVerificationGasLimit = ConvHexOrZero(gapad.PaymasterVerificationGasLimit)
			usop.PaymasterPostOpGasLimit = ConvHexOrZero(gapad.PaymasterPostOpGasLimit)
			state.Save()
			return

		default:
			it, _ := GetItem(sel, items)
			it.Value = !it.Value.(bool)

		}
	}
}

func IncorporatePMandData(usop *userop.UserOperation, pmad *rpccalls.PMandDataResult) {
	state.Log("IncorporatePMandData", pmad)
	if len(pmad.PaymasterAndData) > 0 {
		pmadHex := pmad.PaymasterAndData
		if len(pmadHex) < 42 {
			fmt.Println("Paymaster and Data too short")
			return
		}
		pmadbt, _ := hex.DecodeString(pmadHex[42:])
		paddr := ecommon.HexToAddress(pmadHex[:42])
		usop.Paymaster = &paddr
		usop.PaymasterData = pmadbt
	} else {
		pa := ecommon.HexToAddress(pmad.Paymaster)
		usop.Paymaster = &pa
		usop.PaymasterData, _ = hex.DecodeString(strings.TrimPrefix(pmad.PaymasterData, "0x"))
	}

}

func AlchemyOverridesUI(overrides *rpccalls.AlchemyOverrides) *rpccalls.AlchemyOverrides {
	if overrides == nil {
		overrides = &rpccalls.AlchemyOverrides{}
	}
	MaxFeePerGasItem := &Item{Label: "Max Fee Per Gas", Value: overrides.MaxFeePerGas}
	MaxPriorityFeePerGasItem := &Item{Label: "Max Priority Fee Per Gas", Value: overrides.MaxPriorityFeePerGas}
	CallGasLimitItem := &Item{Label: "Call Gas Limit", Value: overrides.CallGasLimit}
	VerificationGasLimitItem := &Item{Label: "Verification Gas Limit", Value: overrides.VerificationGasLimit}
	PreVerificationGasItem := &Item{Label: "Pre Verification Gas", Value: overrides.PreVerificationGas}
	items := []*Item{MaxFeePerGasItem, MaxPriorityFeePerGasItem, CallGasLimitItem, VerificationGasLimitItem, PreVerificationGasItem, Set, Back}
	for {
		spr := promptui.Select{Label: "Select parameters to override", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return overrides
		}
		switch sel {
		case Back.Label:
			return overrides
		case Set.Label:
			if MaxFeePerGasItem.Value != nil {
				overrides.MaxFeePerGas = MaxFeePerGasItem.Value
			} else {
				overrides.MaxFeePerGas = &rpccalls.AlchemyOverrideMultiplier{Multiplier: 1}
			}
			if MaxPriorityFeePerGasItem.Value != nil {
				overrides.MaxPriorityFeePerGas = MaxPriorityFeePerGasItem.Value
			} else {
				overrides.MaxPriorityFeePerGas = &rpccalls.AlchemyOverrideMultiplier{Multiplier: 1}
			}
			if CallGasLimitItem.Value != nil {
				overrides.CallGasLimit = CallGasLimitItem.Value
			} else {
				overrides.CallGasLimit = &rpccalls.AlchemyOverrideMultiplier{Multiplier: 1}
			}
			if VerificationGasLimitItem.Value != nil {
				overrides.VerificationGasLimit = VerificationGasLimitItem.Value
			} else {
				overrides.VerificationGasLimit = &rpccalls.AlchemyOverrideMultiplier{Multiplier: 1}
			}
			if PreVerificationGasItem.Value != nil {
				overrides.PreVerificationGas = PreVerificationGasItem.Value
			} else {
				overrides.PreVerificationGas = &rpccalls.AlchemyOverrideMultiplier{Multiplier: 1}
			}

			return overrides

		default:
			it, _ := GetItem(sel, items)
			SetOverrideValue(sel, it)

		}
	}
}

var isHex = regexp.MustCompile("^0x[0-9a-fA-F]*$")

func SetOverrideValue(sel string, it *Item) {
	HexItem := &Item{Label: "As Absolute Value"}
	MultiplierItem := &Item{Label: "As Multiplier"}
	items := []*Item{HexItem, MultiplierItem, Back}
	selector := promptui.Select{Label: "Select Value Type for " + sel, Items: items, Templates: ItemTemplate, Size: 10}
	_, sel2, err := selector.Run()
	if err != nil {
		fmt.Println("could not run prompt:", err)
		return
	}
	switch sel2 {
	case Back.Label:
		return
	case HexItem.Label:
		v, err := InputUint(it, 64)
		if err != nil {
			fmt.Println("Error getting value:", err)
			return
		}
		it.Value = fmt.Sprintf("0x%x", v)
	case MultiplierItem.Label:
		nit := &Item{Label: "Multiplier Value", Value: 1}
		err := InputFloatUI(nit)
		if err != nil {
			fmt.Println("Error getting multiplier:", err)
			return
		}
		f, ok := nit.Value.(float64)
		if !ok {
			fmt.Println("Not a valid multiplier")
			return
		}
		it.Value = &rpccalls.AlchemyOverrideMultiplier{Multiplier: f}
	}
}
