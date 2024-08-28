package ui

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
)

func EntryPointUI() {
	var EntryPointV6Item = &Item{Label: "V6 at " + entrypoint.EntryPointAddressV6,
		Details: fmt.Sprintf("Use the standard V6 entrypoint at: %s", entrypoint.E6Address)}
	var EntryPointV7Item = &Item{Label: "V7 at " + entrypoint.EntryPointAddressV7,
		Details: fmt.Sprintf("Use the standard V7 entrypoint at: %s", entrypoint.E7Address)}
	var CustomEntryPointItem = &Item{Label: "Custom", Details: "Use a custom entrypoint"}
	items := []*Item{
		EntryPointV6Item,
		EntryPointV7Item,
		CustomEntryPointItem,
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
		entrypoint.E6Address = common.HexToAddress(entrypoint.EntryPointAddressV6)
		EntryPointItem.Value = entrypoint.E6Address
	case EntryPointV7Item.Label:
		entrypoint.E7Address = common.HexToAddress(entrypoint.EntryPointAddressV7)
		EntryPointItem.Value = entrypoint.E7Address
	case CustomEntryPointItem.Label:
		CustomEntryPointUI()
	default:
		fmt.Println("Unreachable reached:", sel)
	}

}

func CustomEntryPointUI() {
	//Input custom address

	_, cadd, ok := AddressFromBookUI("Custom Entry Point")
	if !ok {
		fmt.Println("Invalid address")
		return
	}

	//Select v6 or v7
	CustomV6Item := &Item{Label: "V6", Details: "Use the V6 spec"}
	CustomV7Item := &Item{Label: "V7", Details: "Use the V7 spec"}
	items := []*Item{
		CustomV6Item, CustomV7Item, Back,
	}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Specification version:",
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
	case CustomV6Item.Label:
		entrypoint.E6Address = *cadd
		EntryPointItem.Value = entrypoint.E6Address
	case CustomV7Item.Label:
		entrypoint.E7Address = *cadd
		EntryPointItem.Value = entrypoint.E7Address
	default:
		fmt.Println("Unreachable reached:", sel)
	}

}
