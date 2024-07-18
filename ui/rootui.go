package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
)

func init() {
	EntryPointItem.DisplayValue = entrypoint.EntryPointAddressV6
	bt, err := hex.DecodeString(entrypoint.EntryPointAddressV6[2:])
	if err != nil {
		fmt.Println(err)
	}
	EntryPointItem.Value = bt
}

var PaymasterItem = &Item{Label: "Paymaster", Details: "Manage paymaster settings"}
var UserOpItem = &Item{Label: "User Operations", Details: "Manage user operations"}
var AbisItem = &Item{Label: "ABIs", Details: "Manage ABIs"}
var SignerItem = &Item{Label: "Signer", Details: "Manage signer settings"}
var EntryPointItem = &Item{Label: "Entry Point", Details: "Set entry point"}

var SettingsItem = &Item{Label: "Settings", Details: "Paymasters, Signers, ChainID, ..."}

func RootUI() {
	items := []*Item{
		//PaymasterItem,
		SettingsItem,
		UserOpItem,
		//SignerItem,
		AbisItem,
		//ChainIDItem,
		//EntryPointItem,
		Exit,
	}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case SettingsItem.Label:
			SettingsUI()
		case PaymasterItem.Label:
			PaymasterUI()
		case UserOpItem.Label:
			TopUserOpUI()
		case SignerItem.Label:
			SignerUI()
		case AbisItem.Label:
			AbisUI(nil)
		case EntryPointItem.Label:
			EntryPointUI()
		case Exit.Label:
			return
		case ChainIDItem.Label:
			InputUint(ChainIDItem, 64)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func SettingsUI() {
	items := []*Item{
		PaymasterItem,
		SignerItem,
		ChainIDItem,
		EntryPointItem,
		Back,
	}
	prompt := promptui.Select{
		Label:     "Settings",
		Items:     items,
		Templates: ItemTemplate,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case PaymasterItem.Label:
			PaymasterUI()
		case SignerItem.Label:
			SignerUI()
		case ChainIDItem.Label:
			InputUint(ChainIDItem, 64)
		case EntryPointItem.Label:
			EntryPointUI()
		case Back.Label:
			return
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func GetValue(label string, items []*Item) (interface{}, bool) {
	it, ok := GetItem(label, items)
	if !ok {
		return nil, false
	}
	return it.Value, true
}

func GetItem(label string, items []*Item) (*Item, bool) {
	for _, i := range items {
		if i.Label == label {
			return i, true
		}
	}
	return nil, false
}
