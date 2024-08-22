package ui

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
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
			ViewCallUI()

		case TxCallItem.Label:
			TxCallUI()

			return
		}
	}

}

var TargetContractItem = &Item{Label: "Target Contract", Details: "Select the target address"}
var CallDataItem2 = &Item{Label: "Call Data", Details: "Enter the call data"}
var MakeTheCallItem = &Item{Label: "Call", Details: "Call with the selected values"}

func ViewCallUI() {
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     []*Item{TargetContractItem, SendEndpointItem, CallDataItem2, MakeTheCallItem, Back},
		Templates: ItemTemplate,
	}
	addressOk := false
	calldatOk := false
	if TargetContractItem.Value != nil {
		_, addressOk = TargetContractItem.Value.(*common.Address)
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
		case SendEndpointItem.Label:
			RPCEndpointsUI(SendEndpointItem)
		case TargetContractItem.Label:
			TargetContractItem.Value, addressOk = AddressFromBookUI("Target Contract")

		case CallDataItem2.Label:
			CallDataItem2.Value, err = PotentiallyRecursiveCallDataUI()
			calldatOk = err == nil

		case MakeTheCallItem.Label:
			if addressOk && calldatOk {
				fmt.Println("Calling contract...")
				res, err := rpccalls.CallContract(SendEndpointItem.Value.(*state.RPCEndpoint), &common.Address{},
					TargetContractItem.Value.(*common.Address), CallDataItem2.Value.([]byte))
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Result:", string(res), len(res), ShortHex(res, 32))
				}
				return
			} else {
				fmt.Printf("Target address ok: %v, Call data ok: %v\n", addressOk, calldatOk)
			}
		}
	}

}

func TxCallUI() {
	fmt.Println("Not implemented yet") // TODO
}
