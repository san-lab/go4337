package ui

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddPaymasterItem = &Item{Label: "Add Paymaster", Details: "Add a new paymaster address"}

func PaymasterUI() (*common.Address, error) {
	for {
		items := []*Item{}
		for _, s := range state.State.Paymasters {
			items = append(items, &Item{Label: s.String(), Details: "Select this Paymaster", Value: s})

		}
		items = append(items, AddPaymasterItem, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Select a Paymaster",
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		switch sel {
		case Back.Label:
			return nil, fmt.Errorf("Back")
		case AddPaymasterItem.Label:
			pymas, err := InputNewAddressUI("Add Paymester")
			if err != nil {
				fmt.Println(err)
			} else {
				state.State.Paymasters = append(state.State.Paymasters, pymas)
				state.State.Save()
			}

		default:
			val, ok := GetValue(sel, items)
			if !ok {
				return nil, fmt.Errorf("Invalid selection")
			}
			addr := val.(*common.Address)

			PaymasterItem.Value = addr
			PaymasterItem.DisplayValue = addr.String()
			fmt.Println("Selected Paymaster:", addr.String())
			return addr, nil
		}
	}
	return nil, nil
}
