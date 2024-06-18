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
var SignItem = &Item{Label: "Sign", Details: "Sign the user operation"}

func GetHashUI(usop *userop.UserOp) {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{ChainIDItem, EntryPointItem, PreHashItem, HashItem, SignerItem, SignItem, Back},
		Templates: ItemTemplate,
	}

	for {
		_, sel, err := prompt.Run()
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
			signer, ok := SignItem.Value.(signer.Signer)
			if !ok {
				fmt.Println("Invalid Signer")
				return
			}
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.([]byte)
			hash, err := userop.GetUserOpHash(usop.Pack(), common.BytesToAddress(EntryPoint), ChainID)
			if err != nil {
				fmt.Println("error hashing:", err)
				return
			}
			sig, err := signer.Sign(hash[:])
			if err != nil {
				fmt.Println("error signing:", err)
				return
			}
			usop.Signature = sig
			SignatureItem.Value = sig
			SignatureItem.DisplayValue = ShortHex(sig, 6)
			fmt.Println("Signature:", hex.EncodeToString(sig[:]))
		case HashItem.Label:
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.([]byte)
			hash, err := userop.GetUserOpHash(usop.Pack(), common.BytesToAddress(EntryPoint), ChainID)
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
