package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func PotentiallyRecursiveCallDataUI() (calldata []byte, err error) {
	// SetCallDataUI sets the CallDataItem value
	var CallDataDirectItem = &Item{Label: "Input as hex", Details: "Set Call Data Direct"}
	var CallDataABIItem = &Item{Label: "Input by ABI", Details: "Set Call Data using ABI"}
	var CallDataFileItem = &Item{Label: "Input by File", Details: "Set Call Data using File"}

	items := []*Item{
		CallDataDirectItem,
		CallDataABIItem,
		CallDataFileItem,
		Back,
	}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
	}
	for {
		var sel string
		_, sel, err = prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case CallDataDirectItem.Label:
			it := &Item{Label: "Input Hex", Details: "Input Hex Data"}
			err := InputBytes(it, -1)
			return it.Value.([]byte), err
		case CallDataABIItem.Label:
			it := &Item{Label: "Select ABI Methods to encode", Details: "Input ABI Data"}
			_, err := AbisUI(it)
			var retbytes []byte

			if it.Value != nil {
				retbytes = it.Value.([]byte)
			}
			return retbytes, err
		case CallDataFileItem.Label:
			bts, err := SelectFileFromFS("")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return bts, nil
		default:
			fmt.Println("Unreachable reached:", sel)
		}
	}

}
