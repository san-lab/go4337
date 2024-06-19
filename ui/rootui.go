package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/manifoldco/promptui"
)

const EntryPointAddressV6 = "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"
const EntryPointAddressV7 = "0x0000000071727De22E5E9d8BAf0edAc6f37da032"

func init() {
	EntryPointItem.DisplayValue = EntryPointAddressV6
	bt, err := hex.DecodeString(EntryPointAddressV6[2:])
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

func RootUI() {
	items := []*Item{
		PaymasterItem,
		UserOpItem,
		SignerItem,
		AbisItem,
		ChainIDItem,
		EntryPointItem,
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
		case PaymasterItem.Label:
			PaymasterUI()
		case UserOpItem.Label:
			UserOpUI()
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
