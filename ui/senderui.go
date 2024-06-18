package ui

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddSenderItem = &Item{Label: "Add Sender", Details: "Add a new sender"}

func SenderUI() (*common.Address, error) {

	for {
		items := []*Item{}
		for _, s := range state.State.Senders {
			items = append(items, &Item{Label: s.String(), Details: "Select this sender", Value: s})

		}
		items = append(items, AddSenderItem, Back)
		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Select a sender",
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		if sel == Back.Label {
			return nil, fmt.Errorf("Back")
		}
		if sel == AddSenderItem.Label {
			sender, err := InputNewAddressUI("Add Sender")
			if err != nil {
				fmt.Println(err)
			} else {
				state.State.Senders = append(state.State.Senders, sender)
				state.State.Save()
				continue
			}
		} else {
			val, ok := GetValue(sel, items)
			if !ok || val == nil {
				return nil, fmt.Errorf("Invalid selection")
			}
			return val.(*common.Address), nil
		}
	}

}
