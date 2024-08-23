package ui

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddAddressItem = &Item{Label: "Add a new Address"}
var RemoveAddressItem = &Item{Label: "Remove an Address"}
var FromAnotherBookItem = &Item{Label: "From another Address Book"}

func AddressFromBookUI(label string) (*common.Address, bool) {
	selectToRemove := false
	abook, _ := state.GetAddressBook(label)
	normalLabel := "Select a " + label
	removeLabel := "Select a " + label + " to remove"
	currentLabel := normalLabel
	for {
		items := []*Item{}
		for name, addr := range *abook {
			items = append(items, &Item{Label: fmt.Sprintf("%-25s", name), Details: "Select this " + label, Value: addr})

		}
		if !selectToRemove {
			items = append(items, FromAnotherBookItem, AddAddressItem, RemoveAddressItem)
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
		name := strings.TrimSpace(strings.Split(sel, ":")[0])

		switch sel {

		case Back.Label:
			return nil, false
		case AddAddressItem.Label:

			nname, naddrs, err := InputNewAddressUI("Add a new " + label)
			if err != nil {
				fmt.Println(err)
			} else {
				abook.Add(nname, naddrs)

				continue
			}
		case RemoveAddressItem.Label:
			selectToRemove = true
			currentLabel = removeLabel
		case FromAnotherBookItem.Label:
			otherbook, addr, ok := AddressFromAllBooksUI()
			if ok {
				abook.Add(fmt.Sprintf("%s_%d", otherbook, len(*abook)), addr)
				return addr, true
			}
		default:
			val, ok := GetValue(sel, items)
			if !ok || val == nil {
				fmt.Println("Invalid selection: ", sel)
				return nil, false

			}
			if selectToRemove {
				abook.RemoveByName(name)
				selectToRemove = false
				currentLabel = normalLabel
			} else {
				return val.(*common.Address), true
			}
		}
	}

}

func AddressFromAllBooksUI() (string, *common.Address, bool) {
	items := []*Item{}
	for name := range state.GetAddressBooks() {
		items = append(items, &Item{Label: name, Details: "Select an Address from " + name})
	}
	items = append(items, Back)
	prompt := promptui.Select{
		Label:     "Select an Address Book",
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
	}
	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return "", nil, false
	}
	switch sel {
	case Back.Label:
		return "", nil, false
	default:
		addr, ok := AddressFromBookUI(sel)
		return sel, addr, ok
	}
}
