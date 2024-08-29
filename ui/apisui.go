package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var ApiCallsItem = &Item{Label: "API Calls", Details: "Call APIs"}
var ApiKeyItem = &Item{Label: "API Key"}
var ApiURLItem = &Item{Label: "API URL"}
var StackUpPMApiItem = &Item{Label: "StackUp Paymaster API"}
var StackUpBNApiItem = &Item{Label: "StackUp Bundler API"}
var ApiMethodItem = &Item{Label: "API Method (free-hand)"}

var ApiUserOpItem = &Item{Label: "User Operation"}
var ApiCallItem = &Item{Label: "Call API"}

func ApiKeysUI(usop *userop.UserOperation) {
	var key, url, mtemplate string
	var ok, allOk bool
	var keyAndUrl bool
	if usop != nil {
		ApiUserOpItem.Value = usop
	}
	for {
		if ApiKeyItem.Value != nil {
			key, ok = ApiKeyItem.Value.(string)
			allOk = ok
		}
		if ApiURLItem.Value != nil {
			url, ok = ApiURLItem.Value.(string)
			allOk = allOk && ok
			keyAndUrl = allOk
		}
		if ApiMethodItem.Value != nil {
			mtemplate, ok = ApiMethodItem.Value.(string)
			allOk = allOk && ok
		}
		if ApiUserOpItem.Value != nil {
			usop, ok = ApiUserOpItem.Value.(*userop.UserOperation)
			allOk = allOk && ok
			keyAndUrl = keyAndUrl && ok
		}
		items := []*Item{ApiKeyItem, ApiURLItem, ApiUserOpItem}
		if keyAndUrl {
			items = append(items, StackUpPMApiItem, StackUpBNApiItem, ApiMethodItem)
		}
		if allOk {
			items = append(items, ApiCallItem)
		}
		items = append(items, Back)

		spr := promptui.Select{Label: "API Keys and URL", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case StackUpPMApiItem.Label:
			res, aerr, err := rpccalls.StackUpPMPayCall(url, key, usop.ToUserOpForApiV6())
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else if aerr != nil {
				fmt.Printf("Error from API call: %v, %s, %s\n", aerr.Code, aerr.Message, string(aerr.Data))
			} else {

				fmt.Printf("API call result:\n CallGasLimit: %s (0x%x)\n PreVerificationGas: %s (0x%x)\n VerificationGasLimit: %s (0x%x)\n PaymasterAndData: %s\n",
					res.CallGasLimit, usop.CallGasLimit,
					res.PreVerificationGas, usop.PreVerificationGas,
					res.VerificationGasLimit, usop.VerificationGasLimit,
					res.PaymasterAndData)
				if YesNoPromptUI("Incorporate to the UserOp?") {
					rpccalls.IncorporateStackUpPMResToUserOp(usop, res)
					state.Save()
				}
			}

		case ApiKeyItem.Label:
			_, name, key, good := StringFromDictionaryUI(state.ApiKeysLabel)
			if good {
				ApiKeyItem.Value = key
				ApiKeyItem.Details = name
			}
		case ApiURLItem.Label:
			_, name, url, good := StringFromDictionaryUI(state.ApiEndpointsLabel)
			if good {
				ApiURLItem.Value = url
				ApiURLItem.Details = name
			}
		case ApiMethodItem.Label:
			_, name, mtemplate, good := StringFromDictionaryUI(state.MethodTemplatesLabel)
			if good {
				ApiMethodItem.Value = mtemplate
				ApiMethodItem.Details = name
			}
		case ApiUserOpItem.Label:
			SelectUserOpUI(ApiUserOpItem)
		case ApiCallItem.Label:
			ret, aerr, err := rpccalls.ApiFreeHandCall(url, key, mtemplate, usop.ToUserOpForApiV6())
			if err != nil {
				fmt.Println("Error making API call:", err)
			} else if aerr != nil {
				fmt.Printf("Error from API call: %v, %s, %s\n", aerr.Code, aerr.Message, string(aerr.Data))

			} else {
				fmt.Println("API call result:", string(ret))

			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}

	}
}
