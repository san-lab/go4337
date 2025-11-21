package ui

import (
	"fmt"
	"math/big"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/signui"
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
		Size:      10,
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
			copyValuesToUserOp(usop)
			return SetEthInfVPMV7DataUI(it, usop)
		default:
			fmt.Println("Invalid selection")
			return nil
		}
	}

}

var PaymasterSignerItem = &Item{Label: "Paymaster Signer", Details: "Signer for Paymaster"}

func SetEthInfVPMV7DataUI(it *Item, usop *userop.UserOperation) error {
	prompt := promptui.Select{
		Label:     "Set Validation Parameters",
		Items:     []*Item{ChainIDItem, AfterItem, UntilItem, PaymasterSignerItem, PaymasterSignatureItem, Set, Back},
		Templates: ItemTemplate,
		Size:      10,
	}
	var after, until uint64
	if len(usop.PaymasterData) > 64 {
		a := new(big.Int).SetBytes(usop.PaymasterData[:32])
		b := new(big.Int).SetBytes(usop.PaymasterData[32:64])
		AfterItem.Value = a.Uint64()
		UntilItem.Value = b.Uint64()
	}
	for {

		_, sel, err := prompt.Run()
		if err != nil {
			return err
		}
		switch sel {
		case Back.Label:
			return nil
		case PaymasterSignerItem.Label:
			signui.SignerUI(PaymasterSignerItem)

		case AfterItem.Label:
			after, err = InputUint(AfterItem, 48)
			if err != nil {
				fmt.Println(err)
			}

		case UntilItem.Label:
			until, err = InputUint(UntilItem, 48)
			if err != nil {
				fmt.Println(err)
			}

		case Set.Label:
			pmd := []byte{}

			a := userop.PackUints(0, after)
			u := userop.PackUints(0, until)
			pmd = append(pmd, a[:]...)
			pmd = append(pmd, u[:]...)

			sig, ok := PaymasterSignatureItem.Value.([]byte)
			if ok {
				pmd = append(pmd, sig...)
			}
			if len(pmd) > 0 {
				it.Value = pmd

			}

			return nil

		}

		chainid := ChainIDItem.Value.(*big.Int)

		_, hash, err := userop.GetPaymasterV7Hash(usop.Pack(), chainid, until, after)
		if err != nil {
			return fmt.Errorf("error hashing for paymaster: %v", err)
		}
		asigner, ok := PaymasterSignerItem.Value.(signer.Signer)

		if ok {
			mhash := signer.ToEthSignedMessageHash(hash)
			sig, err := asigner.SignHash(mhash[:])
			if err != nil {
				return fmt.Errorf("error signing for paymaster: %v", err)
			}
			PaymasterSignatureItem.Value = sig
		} else {
			fmt.Println("No signer set")
		}

	}

}

/*

Alchemy's Paymaster endpoint returns:

Paymaster Address: 0xc03aac639bb21233e0139381970328db8bceeb67
Gas Limit?:        000000000000000066b9c1b700000000
Gas PostOp?:       00000000000000000000000000000000
Signature:         fc9686b707d33d7d4fe9c34d5a420e810de94a8540dd4a4ce3d5b60152423ef44300824fae70351502982671c70b6d23cdd703617a53aca8b70b56fec78e2fe31c

Our paymaster:

Paymaster Address: 0xd8c9e8f357810a5e702764f626d4a380ca462fb1
Gas Limit:         00000000000000000000000000989680
Gas PostOp:        000000000000000000000000000f4240
ValidUntil:        0000000000000000000000000000000000000000000000000000000000000000
ValidAfter:        0000000000000000000000000000000000000000000000000000000000000000
Signature:         d525443c573b55b4b365fc46e6ddba6daa7babc41ac23d84e0d9e064a52741b03ae907e8a8de6608f057c07bc6f770b96b15e631ea170f893aeac5a31c97f0111b

Maybe their Bundler expects the same fields...
*/
