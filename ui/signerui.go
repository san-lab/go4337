package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

var AddSignerItem = &Item{Label: "Add Signer", Details: "Add a new signer"}
var RemoveSignerItem = &Item{Label: "Remove Signer", Details: "Remove a signer"}

func SignerUI(signerItem *Item) {
	selectToRemove := false

	for {
		l := len(state.GetSigners())
		items := []*Item{}
		for _, sn := range state.GetSigners() {
			s := state.GetSigner(sn)
			items = append(items, &Item{Label: fmt.Sprintf("%-20s%s", s.Name()+":", s.String()), Details: "Signer of type " + s.Type()})
		}
		if !selectToRemove {
			items = append(items, AddSignerItem)
			if l > 0 {
				items = append(items, RemoveSignerItem)
			}
		}
		items = append(items, Back)

		selec := promptui.Select{
			Label:     "Manage Signers",
			Items:     items,
			Templates: ItemTemplate,
			Size:      l + 3,
		}

		i, sel, err := selec.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		name := strings.TrimSpace(strings.Split(sel, ":")[0])
		switch sel {
		case Back.Label:
			return
		case AddSignerItem.Label:
			AddSignerUI()
		case RemoveSignerItem.Label:
			selectToRemove = true
		default:
			if i < l {
				if selectToRemove {
					state.RemoveSigner(name)
					selectToRemove = false
					continue
				}

				signerItem.Value = state.GetSigner(name)
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
