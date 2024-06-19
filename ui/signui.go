package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/userop"
)

var ChainIDItem = &Item{Label: "Chain ID", Details: "Set Chain ID", Value: uint64(1), DisplayValue: "1"}
var HashItem = &Item{Label: "Hash", Details: "Get the hash of the user operation with entrypoint and chainid"}
var PreHashItem = &Item{Label: "Pre Hash", Details: "Hash of an encoded user operation"}
var EncodedBytesItem = &Item{Label: "Encoded Bytes", Details: "Encoded bytes of the user operation"}

func GetHashUI(usop *userop.UserOp) (sig []byte, err error) {
	var SignItem = &Item{Label: "Sign", Details: "Sign the user operation"}
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{ChainIDItem, EntryPointItem, SignerItem, SignItem, PreHashItem, HashItem, Back},
		Templates: ItemTemplate,
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
			err = InputUint(ChainIDItem, 64)
			if err != nil {
				return
			}
		case EntryPointItem.Label:
			EntryPointUI()
		case PreHashItem.Label:
			h, err := userop.GetUsOpLibPrehash(usop.Pack())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("PreHash:", hex.EncodeToString(h[:]))
			}
		case EncodedBytesItem.Label:
			encoded, err := userop.GetUserOpBytesToHash(usop.Pack(), common.BytesToAddress(EntryPointItem.Value.([]byte)), ChainIDItem.Value.(uint64))
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Encoded Bytes:", hex.EncodeToString(encoded[:]))
			}

		case SignerItem.Label:
			SignerUI()
		case SignItem.Label:
			signer, ok := SignerItem.Value.(signer.Signer)
			if !ok {
				fmt.Println("Invalid Signer")
				return
			}
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.([]byte)
			var hash [32]byte
			hash, err = userop.GetUserOpHash(usop.Pack(), common.BytesToAddress(EntryPoint), ChainID)
			if err != nil {
				fmt.Println("error hashing:", err)
				return
			}
			sig, err = signer.Sign(hash[:])
			if err != nil {
				fmt.Println("error signing:", err)
				return
			}
			usop.Signature = sig
			SignatureItem.Value = sig
			SignatureItem.DisplayValue = ShortHex(sig, 6)
			SignatureItem.Details = hex.EncodeToString(sig[:])
			fmt.Println("Signature:", hex.EncodeToString(sig[:]))
			return
		case HashItem.Label:
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.([]byte)
			var hash [32]byte
			hash, err = userop.GetUserOpHash(usop.Pack(), common.BytesToAddress(EntryPoint), ChainID)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Hash:", hex.EncodeToString(hash[:]))
			}

			return
		default:
			return
		}
	}

}

func SetSignatureUI(userop *userop.UserOp) (calldata []byte, err error) {

	var SignatureItemDirectItem = &Item{Label: "Input signature as hex", Details: "Input signature directly as hex"}
	var UseSignerItem = &Item{Label: "Calculate using a Signer", Details: "Set Call Data using ABI"}

	items := []*Item{
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
			it := &Item{Label: "Input Hex", Details: "Input Hex Data"}
			err := InputBytes(it, -1)
			if err != nil {
				userop.Signature = it.Value.([]byte)
			}
			return it.Value.([]byte), err
		case UseSignerItem.Label:
			return GetHashUI(userop)
		default:
			fmt.Println("Unreachable reached:", sel)
		}
	}

}
