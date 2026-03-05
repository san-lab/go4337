package signui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
)

var AddSignerItem = &common.Item[struct{}]{Label: "Add Signer", Details: "Add a new signer"}
var RemoveSignerItem = &common.Item[struct{}]{Label: "Remove Signer", Details: "Remove a signer"}

func SignerUI(signerItem *common.Item[signer.Signer]) {
	selectToRemove := false

	for {
		l := len(state.GetSigners())
		items := []common.MenuItem{}
		for _, sn := range state.GetSigners() {
			s := state.GetSigner(sn)
			items = append(items, &common.Item[signer.Signer]{Label: fmt.Sprintf("%-20s%s", s.Name()+":", s.String()), Details: "Signer of type " + s.Type(), Value: s})
		}
		if !selectToRemove {
			items = append(items, AddSignerItem)
			if l > 0 {
				items = append(items, RemoveSignerItem)
			}
		}
		items = append(items, common.Back)

		selec := promptui.Select{
			Label:     "Manage Signers",
			Items:     items,
			Templates: common.ItemTemplate,
			Size:      l + 3,
		}

		i, sel, err := selec.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		name := strings.TrimSpace(strings.Split(sel, ":")[0])
		switch sel {
		case common.Back.Label:
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
				return
			}
			fmt.Println("Unreachable reached:", sel)
		}
	}

}

func AddSignerUI() {
	items := []common.MenuItem{}
	for k := range state.SignerTypes {
		items = append(items, &common.Item[struct{}]{Label: k, Details: "Add a new signer of type " + k})
	}
	items = append(items, common.Back)
	for {
		Label := "Select new Signer Type"
		if len(items) == 1 {
			Label = "No Signer Types available"
		}
		selec := promptui.Select{
			Label:     Label,
			Items:     items,
			Templates: common.ItemTemplate,
		}

		_, sel, err := selec.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
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
