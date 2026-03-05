package abicalldata

import (
	"fmt"

	"github.com/manifoldco/promptui"
	. "github.com/san-lab/go4337/ui/common"
)

func PotentiallyRecursiveCallDataUI() (calldata []byte, err error) {
	// SetCallDataUI sets the CallDataItem value
	var CallDataDirectItem = &Item[struct{}]{Label: "Input as hex", Details: "Set Call Data Direct"}
	var CallDataABIItem = &Item[struct{}]{Label: "Input by ABI", Details: "Set Call Data using ABI"}
	var CallDataFileItem = &Item[struct{}]{Label: "Input by File", Details: "Set Call Data using File"}

	items := []MenuItem{
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
			it := &Item[[]byte]{Label: "Input Hex", Details: "Input Hex Data"}
			err := InputBytes(it, -1)
			return it.Value, err
		case CallDataABIItem.Label:
			it := &Item[[]byte]{Label: "Select ABI Methods to encode", Details: "Input ABI Data"}
			_, err := AbisUI(it)
			return it.Value, err
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
