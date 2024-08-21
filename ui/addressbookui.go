package ui

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddAddressItem = &Item{Label: "Add a new Address"}
var RemoveAddressItem = &Item{Label: "Remove an Address"}

func AddressFromBookUI(label string) (*common.Address, bool) {
	selectToRemove := false
	abook, ok := state.GetAddressBook(label)
	if !ok {
		state.State.AddressBooks[label] = &state.AddressBook{}
	}
	normalLabel := "Select a " + label
	removeLabel := "Select a " + label + " to remove"
	currentLabel := normalLabel
	for {
		items := []*Item{}
		for _, s := range *abook {
			items = append(items, &Item{Label: s.String(), Details: "Select this " + label, Value: s,
				DisplayValueString: " "})

		}
		if !selectToRemove {
			items = append(items, AddAddressItem, RemoveAddressItem)
		}
		items = append(items, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     currentLabel,
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return nil, false
		}

		switch sel {

		case Back.Label:
			return nil, false
		case AddAddressItem.Label:

			naddrs, err := InputNewAddressUI("Add a new " + label)
			if err != nil {
				fmt.Println(err)
			} else {
				abook.Add(naddrs)

				continue
			}
		case RemoveAddressItem.Label:
			selectToRemove = true
			currentLabel = removeLabel
		default:
			val, ok := GetValue(sel, items)
			if !ok || val == nil {
				fmt.Println("Invalid selection: ", sel)
				return nil, false

			}
			if selectToRemove {
				abook.Remove(val.(*common.Address))
				selectToRemove = false
				currentLabel = normalLabel
			} else {
				return val.(*common.Address), true
			}
		}
	}

}
