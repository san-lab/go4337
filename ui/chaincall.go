package ui

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/abicalldata"
	. "github.com/san-lab/go4337/ui/common"
	rpcui "github.com/san-lab/go4337/ui/rpcui"
)

var ViewCallItem = &Item[struct{}]{Label: "View Call", Details: "Call a view/pure function"}
var LegacyTxCallItem = &Item[struct{}]{Label: "Legacy Tx Call", Details: "Send a legacy transaction"}

// var TxCallItem = &Item[struct{}]{Label: "NewTx Call", Details: "Send an EIP2718 transaction"}
var DeployViaCreateItem = &Item[struct{}]{Label: "Deploy via Create", Details: "Deploy a contract via CREATE. TO will be left empty"}
var GenericRPCCallITem = &Item[struct{}]{Label: "Generic RPC Call", Details: "Call a generic RPC function"}

func ChainCallUI() {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []MenuItem{GenericRPCCallITem, ViewCallItem /*TxCallItem,*/, LegacyTxCallItem, DeployViaCreateItem, Back},
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
			TxCallUI(false, false)
		case LegacyTxCallItem.Label:
			TxCallUI(true, false)
		case GenericRPCCallITem.Label:
			GenericRPCCallUI()
		case DeployViaCreateItem.Label:
			TxCallUI(true, true)
			return
		}
	}

}

var TargetContractItem = &Item[*ecommon.Address]{Label: "Target Address ('to')", Details: "Select the target address"}
var CallDataViewItem = &Item[[]byte]{Label: "CallData for View call", Details: "Enter the call data"}
var CallDataTxItem = &Item[[]byte]{Label: "CallData for Tx call", Details: "Enter the call data"}
var MakeTheCallItem = &Item[struct{}]{Label: "Make call", Details: "Call with the selected values"}

var CallSignerItem = &Item[rpccalls.KeyContainer]{Label: "Signer/Sender", Details: "Select the signer"}
var CallValueItem = &Item[*big.Int]{Label: "Value", Details: "Enter the value to send", Value: big.NewInt(0)}
var GasLimitItem = &Item[uint64]{Label: "Gas Limit", Details: "Enter the gas limit", Value: uint64(200000)}
var retArgs abi.Arguments

var UtilCallItem = &Item[struct{}]{Label: "Web3Api Call", Details: "Call a Web3 API function"}

var FromAddress = ecommon.HexToAddress("0xaab05558448C8a9597287Db9F61e2d751645B12a")

func TxCallUI(transactional bool, deploy bool) {

	addressOk := TargetContractItem.Value != nil
	calldatOk := CallDataViewItem.Value != nil
	signerOk := !transactional // Only if transactional, we need a signer
	if transactional {
		signerOk = CallSignerItem.Value != nil
	}
	rpcOk := rpcui.SendEndpointItem.Value != nil

	for {
		items := []MenuItem{TargetContractItem, rpcui.SendEndpointItem}
		if transactional {
			items = append(items, CallDataTxItem, CallSignerItem, CallValueItem, GasLimitItem)
		} else {
			items = append(items, CallDataViewItem) // Just to keep the calldata around when switching to/from transactional
		}
		if rpcOk && (addressOk || deploy) {
			if calldatOk && signerOk {
				items = append(items, MakeTheCallItem)
			}
		}
		items = append(items, Back)
		prompt := promptui.Select{
			Label:     "Set the Call parameters",
			Items:     items,
			Templates: ItemTemplate,
		}
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
			signerOk = CallSignerItem.Value != nil
		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
			rpcOk = rpcui.SendEndpointItem.Value != nil
		case TargetContractItem.Label:
			_, addr, ok := AddressFromBookUI("Target Contract")
			if ok {
				TargetContractItem.Value = addr
				addressOk = true
			}
		case GasLimitItem.Label:
			InputUint64(GasLimitItem)

		case CallDataViewItem.Label:
			calldata, err2 := abicalldata.PotentiallyRecursiveCallDataUI()
			if err2 == nil {
				CallDataViewItem.Value = calldata
				calldatOk = true
				retArgs = abicalldata.CurrentReturnType
			} else {
				fmt.Println(err2)
				calldatOk = false
			}

		case CallDataTxItem.Label:
			calldata, err2 := abicalldata.PotentiallyRecursiveCallDataUI()
			if err2 == nil {
				CallDataTxItem.Value = calldata
				calldatOk = true
			}
		case MakeTheCallItem.Label:
			if !transactional {
				fmt.Println("Calling contract...")
				ret, res, err := rpccalls.CallContract(rpcui.SendEndpointItem.Value,
					&FromAddress,
					TargetContractItem.Value, CallDataViewItem.Value, retArgs)
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
				sgnr := CallSignerItem.Value
				key := sgnr.GetECDSAKey()
				from := ecrypto.PubkeyToAddress(key.PublicKey)
				var to *ecommon.Address
				if deploy {
					to = nil
				} else {
					to = TargetContractItem.Value
				}
				value := CallValueItem.Value
				calldata := CallDataTxItem.Value
				h, err := rpccalls.CreateAndSendTransaction(rpcui.SendEndpointItem.Value,
					&from,
					to,
					value,
					calldata,
					GasLimitItem.Value,
					sgnr)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Transaction sent, hash:", h.Hex())
				}
				return
			}
		}
	}
}

func GenericRPCCallUI() {

	addressOk := TargetContractItem.Value != nil
	rpcOk := rpcui.SendEndpointItem.Value != nil

	for {
		items := []MenuItem{TargetContractItem, rpcui.SendEndpointItem}
		if rpcOk && addressOk {
			items = append(items, UtilCallItem)
		}
		items = append(items, Back)
		prompt := promptui.Select{
			Label:     "Set the Call parameters",
			Items:     items,
			Templates: ItemTemplate,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case TargetContractItem.Label:
			_, addr, ok := AddressFromBookUI("Target Contract")
			if ok {
				TargetContractItem.Value = addr
				addressOk = true
			}

		case Back.Label:
			return

		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
			rpcOk = rpcui.SendEndpointItem.Value != nil
		case UtilCallItem.Label:
			fmt.Println("Calling utility function...")
			if rpcOk && addressOk {
				UtilCallUI(rpcui.SendEndpointItem.Value, *TargetContractItem.Value)
			} else {
				fmt.Println("Need RPC endpoint and target address")
			}

		}
	}
}

func UtilCallUI(endpoint *state.RPCEndpoint, account ecommon.Address) {
	PendingNonceItem := &Item[struct{}]{Label: "Get Pending Nonce", Details: "Get pending nonce at address"}
	NonceItem := &Item[struct{}]{Label: "Get Nonce", Details: "Get nonce at address"}
	BalanceItem := &Item[struct{}]{Label: "Get Balance", Details: "Get balance at address"}
	Items := []MenuItem{PendingNonceItem, NonceItem, BalanceItem, Back}
	prompt := promptui.Select{
		Label:     "Select action",
		Items:     Items,
		Templates: ItemTemplate,
		Size:      10,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case PendingNonceItem.Label:
			nonce, err := rpccalls.GetPendingNonce(endpoint, account)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Pending Nonce at %s: %v\n", account, nonce)
			}
		case NonceItem.Label:
			nonce, err := rpccalls.GetNonce(endpoint, account)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Nonce at %s: %v\n", account, nonce)
			}
		case BalanceItem.Label:
			balance, err := rpccalls.GetBalance(endpoint, account)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Balance at %s: %v\n", account, balance)
			}
		case Back.Label:
			return
		}
	}

}
