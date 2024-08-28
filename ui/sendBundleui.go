package ui

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/entrypoint/entrypointv6"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var SendEndpointItem = &Item{Label: "Blockchain RPC Endpoint"}
var BundleSignerItem = &Item{Label: "Signer for the bundle"}
var SendBundleItem = &Item{Label: "Send as bundle"}
var BeneficiaryItem = &Item{Label: "Beneficiary"}

func SendAsBundleUI(usop *userop.UserOperation) (*common.Hash, error) {
	var endpoint *state.RPCEndpoint
	var signer rpccalls.KeyContainer
	var beneficiary common.Address
	var ok1, ok2, ok3 bool
	for {

		items := []*Item{EntryPointItem, SendEndpointItem, BundleSignerItem, BeneficiaryItem}
		if BeneficiaryItem.Value != nil {
			beneficiary, ok3 = BeneficiaryItem.Value.(common.Address)
		}

		if SendEndpointItem.Value != nil && BundleSignerItem.Value != nil {

			endpoint, ok1 = SendEndpointItem.Value.(*state.RPCEndpoint)
			signer, ok2 = BundleSignerItem.Value.(rpccalls.KeyContainer)
			if ok1 && ok2 && ok3 {
				items = append(items, SendBundleItem)
			}
		}
		items = append(items, Back)
		spr := promptui.Select{Label: "Send UserOp as a Bundle", Items: items, Templates: ItemTemplate, Size: 10}

		_, sel, err := spr.Run()
		if err != nil {
			return nil, fmt.Errorf("could not run prompt: %v", err)
		}
		switch sel {
		case Back.Label:
			return nil, nil
		case SendEndpointItem.Label:
			RPCEndpointsUI(SendEndpointItem)
		case BundleSignerItem.Label:
			GetEOASignerUI(BundleSignerItem)
		case BeneficiaryItem.Label:
			_, paddr, ok := AddressFromBookUI("Beneficiary")
			if ok {
				BeneficiaryItem.Value = *paddr
			}
		case SendBundleItem.Label:
			key := signer.GetKey()
			from := ecrypto.PubkeyToAddress(key.PublicKey)
			to := EntryPointItem.Value.(common.Address)
			fmt.Println("from:", from.Hex())
			fmt.Println("to:", to.Hex())
			fmt.Println("endpoint:", endpoint.Name)
			var calldata []byte
			if EntryPointItem.Value == entrypoint.E7Address {
				fmt.Println("entrypoint v7")
				eabi, err := state.GetABI(state.EntrypointV7)
				if err != nil {
					fmt.Println("could not get abi:", err)
					return nil, err
				}
				calldata, err = eabi.ABI.Pack("handleOps", []userop.PackedUserOp{*usop.Pack()}, beneficiary)
				if err != nil {
					fmt.Println("could not pack data:", err)
					return nil, err
				}
				//fmt.Println("calldata:", hex.EncodeToString(bt))

			}
			if EntryPointItem.Value == entrypoint.E6Address {
				fmt.Println("entrypoint v6")
				eabi, err := state.GetABI(state.EntrypointV6)
				if err != nil {
					fmt.Println("could not get abi:", err)
					return nil, err
				}
				calldata, err = eabi.ABI.Pack("handleOps", []entrypointv6.UserOperation{usop.MarshalV6UserOp()}, beneficiary)
				if err != nil {
					fmt.Println("could not pack data:", err)
					return nil, err
				}
				//fmt.Println("calldata:", hex.EncodeToString(bt))
			}
			gasLimit := usop.TotalGasLimit() + state.GetGasLimitOffset()
			//fmt.Println("gasLimit:", gasLimit, len(calldata))
			//return nil, fmt.Errorf("not implemented")
			return rpccalls.CreateAndSendTransaction(endpoint, from, to, big.NewInt(0), calldata, gasLimit, signer)

		}

	}
}

func GetEOASignerUI(it *Item) {
	items := []*Item{}
	for _, signer := range state.GetSigners() {
		_, ok := signer.(rpccalls.KeyContainer)
		if ok {
			items = append(items, &Item{Label: signer.Name(), Value: signer})
		}

	}
	items = append(items, Back)
	prompt := promptui.Select{Label: "Select a Signer", Items: items, Templates: ItemTemplate, Size: 10}
	sel, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	//Careful. This assumes there is only BACK item attached to the KeyContainers list
	if sel < len(items)-1 {
		it.Value = items[sel].Value
	}
}
