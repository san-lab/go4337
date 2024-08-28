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
var ApiMethodItem = &Item{Label: "API Method"}

var ApiUserOpItem = &Item{Label: "User Operation"}
var ApiCallItem = &Item{Label: "Call API"}

func ApiKeysUI(usop *userop.UserOperation) {
	var key, url, mtemplate string
	var ok, allOk bool
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
		}
		if ApiMethodItem.Value != nil {
			mtemplate, ok = ApiMethodItem.Value.(string)
			allOk = allOk && ok
		}
		if ApiUserOpItem.Value != nil {
			usop, ok = ApiUserOpItem.Value.(*userop.UserOperation)
			allOk = allOk && ok
		}
		items := []*Item{ApiKeyItem, ApiURLItem, ApiMethodItem, ApiUserOpItem}
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
			ret, err := rpccalls.ApiCall(url, key, mtemplate, usop)
			if err != nil {
				fmt.Println("Error in API call:", err)
			} else {
				fmt.Println("API call result:", string(ret))
			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}

	}
}
