package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/userop"
)

var RequiredPrefundItem = &Item{Label: "Required Prefund", Details: "Get the required prefund for the user operation"}

var PaymasterV7HashItem = &Item{Label: "Paymaster V7 Hash", Details: "Get the hash as cheked in VerifyingPaymasterV7"}

func UtilsV7UI(usop *userop.UserOperation) {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{PreHashV7Item, HashV7Item, PaymasterV7HashItem, RequiredPrefundItem, Back},
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
			hash, err = userop.GetUserOpHashV7(usop.Pack(), EntryPoint, ChainID)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Hash:", hex.EncodeToString(hash[:]))
			}

			return
		case PaymasterV7HashItem.Label:
			ChainID := ChainIDItem.Value.(uint64)
			vafterItem := &Item{Label: "Valid After", Details: "Valid After"}
			a, err := InputUint(vafterItem, 48)
			if err != nil {
				fmt.Println(err)
				return
			}
			vuntilItem := &Item{Label: "Valid Until", Details: "Valid Until"}
			u, err := InputUint(vuntilItem, 48)
			if err != nil {
				fmt.Println(err)
				return
			}

			var hash []byte
			bts, hash, err := userop.GetPaymasterV7Hash(usop.Pack(), ChainID, a, u)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Paymaster V7 Hash:", hex.EncodeToString(hash[:]))
				fmt.Println("Paymaster V7 Hash Bytes:", hex.EncodeToString(bts))

			}
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
			h, err := userop.GetUsOpLibPrehashV6(usop)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("PreHash:", hex.EncodeToString(h[:]))
			}

		case HashV6Item.Label:
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.(common.Address)
			var hash [32]byte
			hash, err = userop.GetUserOpHashV6(usop, EntryPoint, ChainID)
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
