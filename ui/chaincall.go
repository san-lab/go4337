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
	"github.com/san-lab/go4337/ui/common"
	rpcui "github.com/san-lab/go4337/ui/rpcui"
)

var ViewCallItem = &common.Item{Label: "View Call", Details: "Call a view/pure function"}
var LegacyTxCallItem = &common.Item{Label: "Legacy Tx Call", Details: "Send a legacy transaction"}

// var TxCallItem = &common.Item{Label: "NewTx Call", Details: "Send an EIP2718 transaction"}
var DeployViaCreateItem = &common.Item{Label: "Deploy via Create", Details: "Deploy a contract via CREATE. TO will be left empty"}
var GenericRPCCallITem = &common.Item{Label: "Generic RPC Call", Details: "Call a generic RPC function"}

func ChainCallUI() {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*common.Item{GenericRPCCallITem, ViewCallItem /*TxCallItem,*/, LegacyTxCallItem, DeployViaCreateItem, common.Back},
		Templates: common.ItemTemplate,
	}

	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
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

var TargetContractItem = &common.Item{Label: "Target Address ('to')", Details: "Select the target address"}
var CallDataViewItem = &common.Item{Label: "CallData for View call", Details: "Enter the call data"}
var CallDataTxItem = &common.Item{Label: "CallData for Tx call", Details: "Enter the call data"}
var MakeTheCallItem = &common.Item{Label: "Make call", Details: "Call with the selected values"}

var CallSignerItem = &common.Item{Label: "Signer/Sender", Details: "Select the signer"}
var CallValueItem = &common.Item{Label: "Value", Details: "Enter the value to send", Value: big.NewInt(0)}
var GasLimitItem = &common.Item{Label: "Gas Limit", Details: "Enter the gas limit", Value: uint64(200000)}
var retArgs abi.Arguments

var UtilCallItem = &common.Item{Label: "Web3Api Call", Details: "Call a Web3 API function"}

var FromAddress = ecommon.HexToAddress("0xaab05558448C8a9597287Db9F61e2d751645B12a")

func TxCallUI(transactional bool, deploy bool) {

	addressOk := false
	calldatOk := false
	signerOk := !transactional // Only if transactional, we need a signer
	rpcOk := false
	if TargetContractItem.Value != nil {
		_, addressOk = TargetContractItem.Value.(*ecommon.Address)
	}
	if CallDataViewItem.Value != nil {
		_, calldatOk = CallDataViewItem.Value.([]byte)
	}
	if CallSignerItem.Value != nil {
		_, signerOk = CallSignerItem.Value.(rpccalls.KeyContainer)
	}
	if rpcui.SendEndpointItem.Value != nil {
		_, rpcOk = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
	}
	for {
		items := []*common.Item{TargetContractItem, rpcui.SendEndpointItem}
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
		items = append(items, common.Back)
		prompt := promptui.Select{
			Label:     "Set the Call parameters",
			Items:     items,
			Templates: common.ItemTemplate,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
			return

		case CallSignerItem.Label:
			GetEOASignerUI(CallSignerItem)
			_, signerOk = CallSignerItem.Value.(rpccalls.KeyContainer)
		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
			_, rpcOk = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
		case TargetContractItem.Label:
			_, TargetContractItem.Value, addressOk = common.AddressFromBookUI("Target Contract")
		case GasLimitItem.Label:
			common.InputUint(GasLimitItem, 64)

		case CallDataViewItem.Label:
			CallDataViewItem.Value, err = abicalldata.PotentiallyRecursiveCallDataUI()
			if err == nil {
				calldatOk = true
				retArgs = abicalldata.CurrentReturnType
			} else {
				fmt.Println(err)
				calldatOk = false
			}

		case CallDataTxItem.Label:
			CallDataTxItem.Value, err = abicalldata.PotentiallyRecursiveCallDataUI()
			calldatOk = err == nil
		case MakeTheCallItem.Label:
			if !transactional {

				fmt.Println("Calling contract...")
				ret, res, err := rpccalls.CallContract(rpcui.SendEndpointItem.Value.(*state.RPCEndpoint),
					&FromAddress,
					TargetContractItem.Value.(*ecommon.Address), CallDataViewItem.Value.([]byte), retArgs)
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
					fmt.Printf("Raw response: %s\n", common.ShortHex(res, 256))
				}
				return

			} else {

				signer := CallSignerItem.Value.(rpccalls.KeyContainer)
				key := signer.GetECDSAKey()
				from := ecrypto.PubkeyToAddress(key.PublicKey)
				var to *ecommon.Address
				if deploy {
					to = nil
				} else {
					to = TargetContractItem.Value.(*ecommon.Address)
				}
				value := CallValueItem.Value.(*big.Int)
				calldata := CallDataTxItem.Value.([]byte)
				h, err := rpccalls.CreateAndSendTransaction(rpcui.SendEndpointItem.Value.(*state.RPCEndpoint),
					&from,
					to,
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

			}
		}
	}
}

func GenericRPCCallUI() {

	addressOk := false
	rpcOk := false
	if TargetContractItem.Value != nil {
		_, addressOk = TargetContractItem.Value.(*ecommon.Address)
	}

	if rpcui.SendEndpointItem.Value != nil {
		_, rpcOk = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
	}
	for {
		items := []*common.Item{TargetContractItem, rpcui.SendEndpointItem}
		if rpcOk && addressOk {
			items = append(items, UtilCallItem)

		}
		items = append(items, common.Back)
		prompt := promptui.Select{
			Label:     "Set the Call parameters",
			Items:     items,
			Templates: common.ItemTemplate,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case TargetContractItem.Label:
			_, TargetContractItem.Value, addressOk = common.AddressFromBookUI("Target Contract")

		case common.Back.Label:
			return

		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
			_, rpcOk = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
		case UtilCallItem.Label:
			fmt.Println("Calling utility function...")
			if rpcOk && addressOk {
				UtilCallUI(rpcui.SendEndpointItem.Value.(*state.RPCEndpoint), *TargetContractItem.Value.(*ecommon.Address))
			} else {
				fmt.Println("Need RPC endpoint and target address")
			}

		}
	}
}

func UtilCallUI(endpoint *state.RPCEndpoint, account ecommon.Address) {
	PendingNonceItem := &common.Item{Label: "Get Pending Nonce", Details: "Get pending nonce at address"}
	NonceItem := &common.Item{Label: "Get Nonce", Details: "Get nonce at address"}
	BalanceItem := &common.Item{Label: "Get Balance", Details: "Get balance at address"}
	Items := []*common.Item{PendingNonceItem, NonceItem, BalanceItem, common.Back}
	prompt := promptui.Select{
		Label:     "Select action",
		Items:     Items,
		Templates: common.ItemTemplate,
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
		case common.Back.Label:
			return
		}
	}

}
