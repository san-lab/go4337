package ui

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
)

var ViewCallItem = &Item{Label: "View Call", Details: "Call a view/pure function"}
var TxCallItem = &Item{Label: "Tx Call", Details: "Send a transaction"}

func ChainCallUI() {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{ViewCallItem, TxCallItem, Back},
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

		case ViewCallItem.Label:
			//ViewCallUI()
			TxCallUI(false)
		case TxCallItem.Label:
			TxCallUI(true)

			return
		}
	}

}

var TargetContractItem = &Item{Label: "Target Contract", Details: "Select the target address"}
var CallDataViewItem = &Item{Label: "CallData for View call", Details: "Enter the call data"}
var CallDataTxItem = &Item{Label: "CallData for Tx call", Details: "Enter the call data"}
var MakeTheCallItem = &Item{Label: "Make call", Details: "Call with the selected values"}

var CallSignerItem = &Item{Label: "Signer/Sender", Details: "Select the signer"}
var CallValueItem = &Item{Label: "Value", Details: "Enter the value to send", Value: big.NewInt(0)}
var GasLimitItem = &Item{Label: "Gas Limit", Details: "Enter the gas limit", Value: uint64(200000)}
var retArgs abi.Arguments

func TxCallUI(transactional bool) {
	items := []*Item{TargetContractItem, SendEndpointItem}
	if transactional {
		items = append(items, CallDataTxItem, CallSignerItem, CallValueItem, GasLimitItem)
	} else {
		items = append(items, CallDataViewItem) // Just to keep the calldata around when switching to/from transactional
	}
	items = append(items, MakeTheCallItem, Back)
	prompt := promptui.Select{
		Label:     "Set the Call parameters",
		Items:     items,
		Templates: ItemTemplate,
	}
	addressOk := false
	calldatOk := false
	signerOk := false
	if TargetContractItem.Value != nil {
		_, addressOk = TargetContractItem.Value.(*common.Address)
	}
	if CallDataViewItem.Value != nil {
		_, calldatOk = CallDataViewItem.Value.([]byte)
	}
	if CallSignerItem.Value != nil {
		_, signerOk = CallSignerItem.Value.(rpccalls.KeyContainer)
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
		case CallSignerItem.Label:
			GetEOASignerUI(CallSignerItem)
			_, signerOk = CallSignerItem.Value.(rpccalls.KeyContainer)
		case SendEndpointItem.Label:
			RPCEndpointsUI(SendEndpointItem)
		case TargetContractItem.Label:
			TargetContractItem.Value, addressOk = AddressFromBookUI("Target Contract")

		case CallDataViewItem.Label:
			CallDataViewItem.Value, err = PotentiallyRecursiveCallDataUI()
			if err == nil {
				calldatOk = true
				retArgs = currentReturnType
			} else {
				fmt.Println(err)
				calldatOk = false
			}

		case CallDataTxItem.Label:
			CallDataTxItem.Value, err = PotentiallyRecursiveCallDataUI()
			calldatOk = err == nil
		case MakeTheCallItem.Label:
			if !transactional {
				if addressOk && calldatOk {
					fmt.Println("Calling contract...")
					ret, res, err := rpccalls.CallContract(SendEndpointItem.Value.(*state.RPCEndpoint), &common.Address{},
						TargetContractItem.Value.(*common.Address), CallDataViewItem.Value.([]byte), retArgs)
					if err != nil && err != rpccalls.ErrRetParse {
						fmt.Println(err)
					} else {
						if err == rpccalls.ErrRetParse {
							fmt.Println("Could not parse return values")
						} else {
							for i, v := range ret {
								fmt.Printf("Return value %d: %v\n", i, v)
							}
						}
						fmt.Printf("Raw response: %s\n", ShortHex(res, 256))
					}
					return
				} else {
					fmt.Printf("Target address ok: %v, Call data ok: %v\n", addressOk, calldatOk)
				}
			} else {
				if addressOk && calldatOk && signerOk {
					signer := CallSignerItem.Value.(rpccalls.KeyContainer)
					key := signer.GetKey()
					from := ecrypto.PubkeyToAddress(key.PublicKey)
					to := TargetContractItem.Value.(*common.Address)
					value := CallValueItem.Value.(*big.Int)
					calldata := CallDataTxItem.Value.([]byte)
					h, err := rpccalls.CreateAndSendTransaction(SendEndpointItem.Value.(*state.RPCEndpoint),
						from,
						*to,
						value,
						calldata,
						GasLimitItem.Value.(uint64),
						signer)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Transaction sent, hash:", h.Hex())
					}
					return

				} else {
					fmt.Printf("Target address ok: %v, Call data ok: %v, Signer ok: %v\n", addressOk, calldatOk, signerOk)
				}
			}
		}
	}
}
