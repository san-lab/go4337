package ui

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

var ApiCallsItem = &Item[struct{}]{Label: "API Calls", Details: "Call APIs"}
var ApiKeyItem = &Item[string]{Label: "API Key", Details: "To be appended or replace the {{.}} in URL"}
var ApiURLItem = &Item[string]{Label: "API URL", Details: "If len(api key) > 0, the key will be appended after / or put in place of {{.}}"}
var StackUpPMApiItem = &Item[struct{}]{Label: "StackUp Paymaster API"}
var StackUpBNApiItem = &Item[struct{}]{Label: "StackUp Bundler API"}
var ApiMethodItem = &Item[struct{}]{Label: "API Method (free-hand)"}
var StdApiItem = &Item[struct{}]{Label: "Standard (eth_/pm_) API methods"}

var ApiUserOpItem = &Item[*userop.UserOperation]{Label: "User Operation"}
var ApiCallItem = &Item[struct{}]{Label: "Call API"}
var ZkSyncEraItem = &Item[struct{}]{Label: "ZkSync Era"}

func ApiCallsUI(usop *userop.UserOperation) {
	var AlchemyUIItem = &Item[struct{}]{Label: "Alchemy API"}
	var StackUpUIItem = &Item[struct{}]{Label: "StackUp API"}
	var BiconomyUIItem = &Item[struct{}]{Label: "Biconomy API"}
	var PimlicoUIItem = &Item[struct{}]{Label: "Pimlico API"}
	var CustomUIItem = &Item[struct{}]{Label: "Custom API"}
	var items = []MenuItem{AlchemyUIItem, StackUpUIItem, BiconomyUIItem, PimlicoUIItem, CustomUIItem, ZkSyncEraItem, Back}
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

var eth_sendUserOperationItem = &Item[struct{}]{Label: "eth_sendUserOperation"}
var eth_estimateUserOperationGasItem = &Item[struct{}]{Label: "eth_estimateUserOperationGas"}
var eth_getUserOperationByHashItem = &Item[struct{}]{Label: "eth_getUserOperationByHash"}
var eth_getUserOperationReceiptItem = &Item[struct{}]{Label: "eth_getUserOperationReceipt"}
var eth_supportedEntryPointsItem = &Item[struct{}]{Label: "eth_supportedEntryPoints"}

func addETHMethods(items []MenuItem) []MenuItem {
	items = append(items, eth_sendUserOperationItem, eth_estimateUserOperationGasItem, eth_getUserOperationByHashItem, eth_getUserOperationReceiptItem, eth_supportedEntryPointsItem)
	return items
}

var rundler_maxPriorityFeePerGasItem = &Item[struct{}]{Label: "rundler_maxPriorityFeePerGas"}                   //Alchemy specific
var alchemy_requestPaymasterAndDataItem = &Item[struct{}]{Label: "alchemy_requestPaymentAndData"}               //Alchemy specific
var alchemy_requestGasAndPaymasterAndDataItem = &Item[struct{}]{Label: "alchemy_requestGasAndPaymasterAndData"} //Alchemy specific

var pm_getPaymasterDataItem = &Item[struct{}]{Label: "pm_getPaymasterData"}
var pm_getPaymasterStubDataItem = &Item[struct{}]{Label: "pm_getPaymasterStubData"}
var pm_sponsorUserOperationItem = &Item[struct{}]{Label: "pm_SponsorUserOperation"} //StackUp specific
var pm_accountsItem = &Item[struct{}]{Label: "pm_accounts"}                         //StackUp specific
var pm_getFeeQuoteOrDataItem = &Item[struct{}]{Label: "pm_getFeeQuoteOrData"}       //Biconomy

func addPMMethods(items []MenuItem) []MenuItem {
	items = append(items, pm_getPaymasterDataItem, pm_getPaymasterStubDataItem)
	return items
}

var pimlico_getTokenQuotesItem = &Item[struct{}]{Label: "pimlico_GetTokenQuotes"} //Pimlico specific

var APIKeyItemMap = map[string]*Item[string]{}
var APIURLItemMap = map[string]*Item[string]{}

var Alchemy_requestGasAndPaymasterAndDataItem = &Item[struct{}]{Label: "Alchemy_requestGasAndPaymasterAndData"}
var Alchemy_requestPaymasterAndDataItem = &Item[struct{}]{Label: "Alchemy_requestPaymasterAndData"}

func ProvAPIUI(provider string, usop *userop.UserOperation) {
	state.Log("ProvAPIUI", provider)
	var key, url string
	var allOk bool
	var overrides *rpccalls.AlchemyOverrides

	if usop != nil {
		ApiUserOpItem.Value = usop
	}
	for {
		APIKeyItem, ok3 := APIKeyItemMap[provider]
		if !ok3 { //lazy initialization
			APIKeyItem = &Item[string]{Label: provider + "API Key"}
			APIKeyItemMap[provider] = APIKeyItem
		}
		APIURLItem, ok4 := APIURLItemMap[provider]
		if !ok4 { //lazy initialization
			APIURLItem = &Item[string]{Label: provider + "API URL", Details: "If len(api key) > 0, the key will be appended after / or put in place of {{.}}"}
			APIURLItemMap[provider] = APIURLItem

		}
		var entrypointaddress = EntryPointItem.Value.String()
		var entrypointversion = 0
		if entrypointaddress == entrypoint.E6Address.String() {
			entrypointversion = 6
		} else if entrypointaddress == entrypoint.E7Address.String() {
			entrypointversion = 7
		} else if entrypointaddress == entrypoint.E8Address.String() {
			entrypointversion = 8
		}

		key = APIKeyItem.Value
		url = APIURLItem.Value
		allOk = url != ""

		usop = ApiUserOpItem.Value
		if usop == nil {
			allOk = false
		}

		items := []MenuItem{EntryPointItem, APIKeyItem, APIURLItem, ApiUserOpItem}
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
			_, name, keyVal, good := StringFromDictionaryUI(provider + state.ApiKeysLabel)
			if good {
				APIKeyItem.Value = keyVal
				APIKeyItem.Details = name
			}
		case APIURLItem.Label:
			_, name, urlVal, good := StringFromDictionaryUI(provider + state.ApiEndpointsLabel)
			if good {
				APIURLItem.Value = urlVal
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
				IncorporateRecommendedValuesUI(usop, res)
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
				IncorporateRecommendedValuesUI(usop, pmad)
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
				IncorporateRecommendedValuesUI(usop, pmad)
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
					IncorporateRecommendedValuesUI(usop, gapad)
				}
			} else {
				gapad, err := rpccalls.Alchemy_requestGasAndPaymasterAndDataV7(
					url, key, policyID, entrypointaddress,
					fmt.Sprintf("0x%x", usop.Signature), *usop, overrides)
				if err != nil {
					fmt.Println("Error making API call:", err)
				} else {
					IncorporateRecommendedValuesUI(usop, gapad)
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
	chi := ChainIDItem.Value
	if chi == nil {
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
	case rpccalls.BiconomyProvider:
		return rpccalls.GetBiconomyContext(), chi, true
	default:
		fmt.Println("Not implemented yet:", provider)
		return nil, nil, false

	}
}

func CustomAPIUI(usop *userop.UserOperation) {
	fmt.Println("Custom API UI not implemented yet")
}

func AlchemyOverridesUI(overrides *rpccalls.AlchemyOverrides) *rpccalls.AlchemyOverrides {
	if overrides == nil {
		overrides = &rpccalls.AlchemyOverrides{}
	}
	MaxFeePerGasItem := &Item[any]{Label: "Max Fee Per Gas", Value: overrides.MaxFeePerGas}
	MaxPriorityFeePerGasItem := &Item[any]{Label: "Max Priority Fee Per Gas", Value: overrides.MaxPriorityFeePerGas}
	CallGasLimitItem := &Item[any]{Label: "Call Gas Limit", Value: overrides.CallGasLimit}
	VerificationGasLimitItem := &Item[any]{Label: "Verification Gas Limit", Value: overrides.VerificationGasLimit}
	PreVerificationGasItem := &Item[any]{Label: "Pre Verification Gas", Value: overrides.PreVerificationGas}
	items := []MenuItem{MaxFeePerGasItem, MaxPriorityFeePerGasItem, CallGasLimitItem, VerificationGasLimitItem, PreVerificationGasItem, Set, Back}
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
			it, ok := GetItem(sel, items)
			if ok {
				SetOverrideValue(sel, it.(*Item[any]))
			}

		}
	}
}

var isHex = regexp.MustCompile("^0x[0-9a-fA-F]*$")

func SetOverrideValue(sel string, it *Item[any]) {
	HexItem := &Item[struct{}]{Label: "As Absolute Value"}
	MultiplierItem := &Item[struct{}]{Label: "As Multiplier"}
	items := []MenuItem{HexItem, MultiplierItem, Back}
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
		u64Item := &Item[uint64]{Label: it.Label}
		v, err := InputUint(u64Item, 64)
		if err != nil {
			fmt.Println("Error getting value:", err)
			return
		}
		it.Value = fmt.Sprintf("0x%x", v)
	case MultiplierItem.Label:
		nit := &Item[float64]{Label: "Multiplier Value", Value: 1}
		err := InputFloatUI(nit)
		if err != nil {
			fmt.Println("Error getting multiplier:", err)
			return
		}
		it.Value = &rpccalls.AlchemyOverrideMultiplier{Multiplier: nit.Value}
	}
}

// Needed for type compatibility
var _ = (*big.Int)(nil)
var _ = ecommon.Address{}
