package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/manifoldco/promptui"
)

var EntryPointV6Item = &Item{Label: "V6", Details: fmt.Sprintf("Use the V6 entrypoint at: %s", EntryPointAddressV6)}
var EntryPointV7Item = &Item{Label: "V7", Details: fmt.Sprintf("Use the V7 entrypoint at: %s", EntryPointAddressV7)}

func EntryPointUI() {
	items := []*Item{
		EntryPointV6Item,
		EntryPointV7Item,
		Back,
	}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
	}

	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch sel {
	case Back.Label:
		return
	case EntryPointV6Item.Label:
		EntryPointItem.DisplayValue = EntryPointAddressV6
		bt, _ := hex.DecodeString(EntryPointAddressV6)
		EntryPointItem.Value = bt
	case EntryPointV7Item.Label:
		EntryPointItem.DisplayValue = EntryPointAddressV7
		bt, _ := hex.DecodeString(EntryPointAddressV7)
		EntryPointItem.Value = bt
	default:
		fmt.Println("Unreachable reached:", sel)
	}

}
