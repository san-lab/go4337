package ui

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

func init() {
	ChainIDItem.Value = state.State.ChainID

}

var ChainIDItem = &Item{Label: "Chain ID", Details: "Set Chain ID"}
var UtilsV6Item = &Item{Label: "Utils V6", Details: "Various V6 utility functions"}
var UtilsV7Item = &Item{Label: "Utils V7", Details: "Various V7 utility functions"}

var HashV7Item = &Item{Label: "Hash_v7", Details: "Get the hash of the user operation with entrypoint and chainid"}
var HashV6Item = &Item{Label: "Hash_v6", Details: "Get the V6 hash of the user operation with entrypoint and chainid"}
var PreHashV7Item = &Item{Label: "Pre Hash_v7", Details: "Hash of an encoded user operation"}
var PreHashV6Item = &Item{Label: "Pre Hash_v6", Details: "Hash of an encoded user operation"}
var EncodedBytesItem = &Item{Label: "Encoded Bytes", Details: "Encoded bytes of the user operation"}

func GetHashUI(usop *userop.UserOperation) (sig []byte, err error) {
	var SignItem = &Item{Label: "Sign", Details: "Sign the user operation"}
	prompt := promptui.Select{
		Label: "Select an option",
		//Items:     []*Item{ChainIDItem, EntryPointItem, SignerItem, SignItem, PreHashV7Item, PreHashV6Item, HashV7Item, HashV6Item, Back},
		Items:     []*Item{ChainIDItem, EntryPointItem, SignerItem, SignItem, UtilsV6Item, UtilsV7Item, Back},
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
			err = InputUint(ChainIDItem, 64)
			if err != nil {
				return
			}
			state.State.ChainID = ChainIDItem.Value.(uint64)
			fmt.Println("Chain ID set to:", state.State.ChainID)
			state.State.Save()
		case EntryPointItem.Label:
			EntryPointUI()
		case UtilsV6Item.Label:
			UtilsV6UI(usop)

		case UtilsV7Item.Label:
			UtilsV7UI(usop)

		case SignerItem.Label:
			SignerUI()
		case SignItem.Label:
			signer, ok := SignerItem.Value.(signer.Signer)
			if !ok {
				fmt.Println("Invalid Signer")
				return
			}
			ChainID := ChainIDItem.Value.(uint64)
			EntryPoint := EntryPointItem.Value.(common.Address)

			var hash [32]byte
			if EntryPoint.Hex() == entrypoint.E6Address.Hex() {
				hash, err = userop.GetUserOpHashV6(usop, EntryPoint, ChainID)
				fmt.Println("Signing using v6 hashes")
			} else {

				hash, err = userop.GetUserOpHashV7(usop.Pack(), EntryPoint, ChainID)
				fmt.Println("Signing using v7 hashes")
			}
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
			//SignatureItem.DisplayValue = ShortHex(sig, 6)
			SignatureItem.Details = hex.EncodeToString(sig[:])
			fmt.Println("Signature:", hex.EncodeToString(sig[:]))
			return
			/*
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
			*/
		default:
			return
		}
	}

}

func SetSignatureUI(userop *userop.UserOperation) (calldata []byte, err error) {

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
