package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/userop"
)

// top selection menu
var InputAsHexItem = &Item{Label: "Input as HEX", Details: "Input directly as HEX"}
var EthInfVPMV7Item = &Item{Label: "EthInfinitism ValidatingPaymaster (V7)"}

// ValidatingPaymaster specific menu
var AfterItem = &Item{Label: "Valid After", Details: "After this block time", Value: 0}
var UntilItem = &Item{Label: "Valid Until", Details: "Until this block time", Value: 0}
var PaymasterSignatureItem = &Item{Label: "Paymaster Signature", Details: "Paymaster Signature", Value: []byte{}}

func SetPaymasterDataUI(it *Item, usop *userop.UserOperation) error {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{InputAsHexItem, EthInfVPMV7Item, Back},
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

		case EthInfVPMV7Item.Label:
			it.DisplayValueString = ""
			return SetEthInfVPMV7DataUI(it, usop)
		default:
			fmt.Println("Invalid selection")
			return nil
		}
	}

}

func SetEthInfVPMV7DataUI(it *Item, usop *userop.UserOperation) error {
	prompt := promptui.Select{
		Label:     "Set Validation Parameters",
		Items:     []*Item{ChainIDItem, AfterItem, UntilItem, SignerItem, PaymasterSignatureItem, Set, Back},
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
		case SignerItem.Label:
			SignerUI()

		case AfterItem.Label:
			err := InputUint(AfterItem, 48)
			if err != nil {
				fmt.Println(err)
			}

		case UntilItem.Label:
			err := InputUint(UntilItem, 48)
			if err != nil {
				fmt.Println(err)
			}

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

		}

		chainid := ChainIDItem.Value.(uint64)
		after, oka := AfterItem.Value.(uint64)
		until, oku := UntilItem.Value.(uint64)
		if oka && oku {

			_, hash, err := userop.GetPaymasterV7Hash(usop.Pack(), chainid, until, after)
			if err != nil {
				return fmt.Errorf("error hashing for paymaster: %v", err)
			}
			signer, ok := SignerItem.Value.(signer.Signer)

			if ok {
				sig, err := signer.Sign(hash)
				if err != nil {
					return fmt.Errorf("error signing for paymaster: %v", err)
				}
				PaymasterSignatureItem.Value = sig
			}

		}

	}

}
