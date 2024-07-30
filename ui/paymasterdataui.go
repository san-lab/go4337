package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/userop"
)

// top selection menu
var InputAsHexItem = &Item{Label: "Input as HEX", Details: "Input directly as HEX"}
var EthInfVPMItem = &Item{Label: "EthInfinitism ValidatingPaymaster"}

// ValidatingPaymaster specific menu
var AfterItem = &Item{Label: "Valid After", Details: "After this block time", Value: 0}
var UntilItem = &Item{Label: "Valid Until", Details: "Until this block time", Value: 0}
var PaymasterSignatureItem = &Item{Label: "Paymaster Signature", Details: "Paymaster Signature", Value: []byte{}}

func SetPaymasterDataUI(it *Item) error {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{InputAsHexItem, EthInfVPMItem, Back},
		Templates: ItemTemplate,
	}

	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return err
		}
		switch sel {
		case Back.Label:
			return nil

		case InputAsHexItem.Label:
			err := InputBytes(it, -1)
			if err != nil {
				return err
			}
			it.DisplayValueString = ""
			return nil

		case EthInfVPMItem.Label:
			it.DisplayValueString = ""
			return SetEthInfVPMDataUI(it)
		default:
			fmt.Println("Invalid selection")
			return nil
		}
	}

}

func SetEthInfVPMDataUI(it *Item) error {
	prompt := promptui.Select{
		Label:     "Set Validation Parameters",
		Items:     []*Item{AfterItem, UntilItem, PaymasterSignatureItem, Set, Back},
		Templates: ItemTemplate,
	}

	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return err
		}
		switch sel {
		case Back.Label:
			return nil

		case AfterItem.Label:
			err := InputUint(AfterItem, 48)
			if err != nil {
				return err
			}

		case UntilItem.Label:
			err := InputUint(UntilItem, 48)
			if err != nil {
				return err
			}

		case PaymasterSignatureItem.Label:
			PaymasterSignatureItem.Value = make([]byte, 64)

		case Set.Label:
			pmd := []byte{}

			after, oka := AfterItem.Value.(uint64)
			until, oku := UntilItem.Value.(uint64)
			if oka && oku {
				a := userop.PackUints(0, after)
				u := userop.PackUints(0, until)
				pmd = append(pmd, a[:]...)
				pmd = append(pmd, u[:]...)
			}
			sig, ok := PaymasterSignatureItem.Value.([]byte)
			if ok {
				pmd = append(pmd, sig...)
			}
			if len(pmd) > 0 {
				it.Value = pmd

			}

			return nil
		default:
			fmt.Println("Invalid selection")
			return nil
		}
	}

}
