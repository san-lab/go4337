package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddSignerItem = &Item{Label: "Add Signer", Details: "Add a new signer"}

func SignerUI() {
	for {
		items := []*Item{}
		for _, s := range state.State.Signers {
			items = append(items, &Item{Label: s.String(), Details: "Signer of type " + s.Type()})
		}
		items = append(items, AddSignerItem, Back)

		selec := promptui.Select{
			Label:     "Manage Signers",
			Items:     items,
			Templates: ItemTemplate,
		}

		i, sel, err := selec.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case AddSignerItem.Label:
			AddSignerUI()
		default:
			if i < len(state.State.Signers) {
				SignerItem.Value = state.State.Signers[i]
				//SignerItem.DisplayValue = state.State.Signers[i].String()
				return
			}
			fmt.Println("Unreachable reached:", sel)
		}
	}

}

func AddSignerUI() {
	items := []*Item{}
	for k, _ := range state.SignerTypes {
		items = append(items, &Item{Label: k, Details: "Add a new signer of type " + k})
	}
	items = append(items, Back)
	for {
		Label := "Select new Signer Type"
		if len(items) == 1 {
			Label = "No Signer Types available"
		}
		selec := promptui.Select{
			Label:     Label,
			Items:     items,
			Templates: ItemTemplate,
		}

		_, sel, err := selec.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return
		default:
			if state.SignerTypes[sel] != nil {
				//Call the type specific creation function
				err := state.SignerTypes[sel]()
				if err != nil {
					fmt.Println(err)
				}
				return
			} else {
				fmt.Println("Unreachable reached:", sel)
			}
		}
	}
}
