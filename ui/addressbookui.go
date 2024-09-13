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
var RenameAddressItem = &Item{Label: "Rename an Address"}
var FromAnotherBookItem = &Item{Label: "From another Address Book"}

// Returns selcted name, address and a bool indicating if the selection was successful
func AddressFromBookUI(label string) (string, *common.Address, bool) {
	selectToRemove := false
	selectToRename := false
	abook, _ := state.GetAddressBook(label)
	normalLabel := "Select a " + label
	removeLabel := "Select a " + label + " to remove"
	renameLabel := "Select a " + label + " to rename"
	currentLabel := normalLabel
	for {
		items := []*Item{}
		for _, key := range abook.Keys() {
			items = append(items, &Item{Label: fmt.Sprintf("%-25s", key), Details: "Select this " + label, Value: (*abook)[key]})

		}
		if !selectToRemove && !selectToRename {
			items = append(items, FromAnotherBookItem, AddAddressItem, RemoveAddressItem, RenameAddressItem)
		}
		items = append(items, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     currentLabel,
			Items:     items,
			Templates: ItemTemplate,
		}
		prompt.Size = len(items)
		if prompt.Size > 25 {
			prompt.Size = 25
		}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return "", nil, false
		}
		name := strings.TrimSpace(strings.Split(sel, ":")[0])

		switch sel {

		case Back.Label:
			return "", nil, false
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
		case RenameAddressItem.Label:
			selectToRename = true
			currentLabel = renameLabel
		case FromAnotherBookItem.Label:
			otherbook, oname, addr, ok := AddressFromAllBooksUI()
			if ok {
				abook.Add(fmt.Sprintf("%s_%s", otherbook, oname), addr)
				return oname, addr, true
			}
		default:
			val, ok := GetValue(sel, items)
			if !ok || val == nil {
				fmt.Println("Invalid selection: ", sel)
				return "", nil, false

			}
			if selectToRemove {
				abook.RemoveByName(name)
				selectToRemove = false
				currentLabel = normalLabel

			} else if selectToRename {
				RenameItem := &Item{Label: fmt.Sprintf("New name for >>%s<<", name)}
				err := InputNewStringUI(RenameItem)
				if err != nil {
					fmt.Println(err)
					continue
				}
				nname := RenameItem.Value.(string)
				abook.Rename(name, nname)
				selectToRename = false
				currentLabel = normalLabel
			} else {
				return name, val.(*common.Address), true
			}
		}
	}

}

// Returns the selected addressbook name, selected name and address, and a bool indicating if the selection was successful
func AddressFromAllBooksUI() (string, string, *common.Address, bool) {
	items := []*Item{}
	for _, name := range state.GetAddressBooks() {
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
		return "", "", nil, false
	}
	switch sel {
	case Back.Label:
		return "", "", nil, false
	default:
		oname, addr, ok := AddressFromBookUI(sel)
		return sel, oname, addr, ok
	}
}
