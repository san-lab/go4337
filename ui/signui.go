package ui

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/signui"
	"github.com/san-lab/go4337/userop"
)

func init() {
	ChainIDItem.Value = state.GetChainId()
}

var ChainIDItem = &Item[*big.Int]{Label: "Chain ID", Details: "Set Chain ID"}
var UtilsV6Item = &Item[struct{}]{Label: "Utils V6", Details: "Various V6 utility functions"}
var UtilsV7Item = &Item[struct{}]{Label: "Utils V7", Details: "Various V7 utility functions"}

var HashV7Item = &Item[struct{}]{Label: "Hash_v7", Details: "Get the hash of the user operation with entrypoint and chainid"}
var HashV6Item = &Item[struct{}]{Label: "Hash_v6", Details: "Get the V6 hash of the user operation with entrypoint and chainid"}
var PreHashV7Item = &Item[struct{}]{Label: "Pre Hash_v7", Details: "Hash of an encoded user operation"}
var PreHashV6Item = &Item[struct{}]{Label: "Pre Hash_v6", Details: "Hash of an encoded user operation"}
var EncodedBytesItem = &Item[struct{}]{Label: "Encoded Bytes", Details: "Encoded bytes of the user operation"}

// SignerItem is the signer used for signing user operations.
var SignerItem = &Item[signer.Signer]{Label: "Signer", Details: "Manage Signer settings"}

func GetHashUI(usop *userop.UserOperation) (sig []byte, err error) {
	var SignItem = &Item[struct{}]{Label: "Sign", Details: "Sign the user operation"}
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []MenuItem{ChainIDItem, EntryPointItem, SignerItem, SignItem, UtilsV6Item, UtilsV7Item, Back},
		Templates: ItemTemplate,
		Size:      10,
	}
	var sel string
	for {
		_, sel, err = prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case ChainIDItem.Label:
			_, err = InputBigInt(ChainIDItem)
			if err != nil {
				return
			}
			state.SetChainId(ChainIDItem.Value)
			fmt.Println("Chain ID set to:", state.GetChainId())
		case EntryPointItem.Label:
			EntryPointUI()
		case UtilsV6Item.Label:
			UtilsV6UI(usop)
		case UtilsV7Item.Label:
			UtilsV7UI(usop)
		case SignerItem.Label:
			signui.SignerUI(SignerItem)
		case SignItem.Label:
			asigner := SignerItem.Value
			if asigner == nil {
				fmt.Println("Invalid Signer")
				return
			}
			ChainID := ChainIDItem.Value
			EntryPoint := EntryPointItem.Value
			sig, err = asigner.SignUserop(usop, ChainID, EntryPoint)
			if err != nil {
				fmt.Println("error signing:", err)
				return
			}
			if sig[64] < 27 {
				sig[64] += 27
			}
			usop.Signature = sig
			SignatureItem.Value = sig
			SignatureItem.Details = hex.EncodeToString(sig[:])
			fmt.Println("Signature:", hex.EncodeToString(sig[:]))
			return
		}
	}
}

func SetSignatureUI(userop *userop.UserOperation) (calldata []byte, err error) {

	var SignatureItemDirectItem = &Item[struct{}]{Label: "Input signature as hex", Details: "Input signature directly as hex"}
	var UseSignerItem = &Item[struct{}]{Label: "Calculate using a Signer", Details: "Set Call Data using ABI"}

	items := []MenuItem{
		SignatureItemDirectItem,
		UseSignerItem,
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
		case SignatureItemDirectItem.Label:
			it := &Item[[]byte]{Label: "Input Hex", Details: "Input Hex Data"}
			err := InputBytes(it, -1)
			if err != nil {
				userop.Signature = it.Value
			}
			return it.Value, err
		case UseSignerItem.Label:
			return GetHashUI(userop)
		default:
			fmt.Println("Unreachable reached:", sel)
		}
	}
}
