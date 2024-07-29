package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/userop"
)

var RequiredPrefundItem = &Item{Label: "Required Prefund", Details: "Get the required prefund for the user operation"}

func UtilsV7UI(usop *userop.UserOperation) {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{PreHashV7Item, HashV7Item, RequiredPrefundItem, Back},
		Templates: ItemTemplate,
	}

	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return

		case PreHashV7Item.Label:
			h, err := userop.GetUsOpLibPrehash(usop.Pack())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("PreHash:", hex.EncodeToString(h[:]))
			}

		case HashV7Item.Label:
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.(common.Address)
			var hash [32]byte
			hash, err = userop.GetUserOpHash(usop.Pack(), EntryPoint, ChainID)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Hash:", hex.EncodeToString(hash[:]))
			}

			return
		case RequiredPrefundItem.Label:
			val := fmt.Sprint(usop.GetRequiredPrefund())
			if len(val) > 18 {
				val = val[:len(val)-18] + "." + val[len(val)-18:] + "ETH"
			}
			fmt.Printf("Required Prefund: %s\n", val)
		default:
			return
		}
	}

}

func UtilsV6UI(usop *userop.UserOperation) {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{PreHashV6Item, HashV6Item, RequiredPrefundItem, Back},
		Templates: ItemTemplate,
	}

	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return

		case PreHashV6Item.Label:
			h, err := userop.GetUsOpLibPrehash(usop.Pack())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("PreHash:", hex.EncodeToString(h[:]))
			}

		case HashV6Item.Label:
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.(common.Address)
			var hash [32]byte
			hash, err = userop.GetUserOpHash(usop.Pack(), EntryPoint, ChainID)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Hash:", hex.EncodeToString(hash[:]))
			}

			return
		case RequiredPrefundItem.Label:
			fmt.Println("Required Prefund:", usop.GetRequiredPrefund())
		default:
			return
		}
	}

}
