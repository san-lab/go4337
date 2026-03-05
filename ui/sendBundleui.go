package ui

import (
	"fmt"
	"math/big"

	ecommon "github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/entrypoint/entrypointv6"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/rpcui"
	"github.com/san-lab/go4337/userop"
)

var BundleSignerItem = &Item[rpccalls.KeyContainer]{Label: "Signer for the bundle"}
var SendBundleItem = &Item[struct{}]{Label: "Send as bundle"}
var BeneficiaryItem = &Item[ecommon.Address]{Label: "Beneficiary"}

func SendAsBundleUI(usop *userop.UserOperation) (*ecommon.Hash, error) {
	var endpoint *state.RPCEndpoint
	var sgnr rpccalls.KeyContainer
	var beneficiary ecommon.Address
	var ok1, ok2, ok3 bool
	for {

		items := []MenuItem{EntryPointItem, rpcui.SendEndpointItem, BundleSignerItem, BeneficiaryItem}
		ok3 = BeneficiaryItem.Value != (ecommon.Address{})
		beneficiary = BeneficiaryItem.Value

		endpoint = rpcui.SendEndpointItem.Value
		ok1 = endpoint != nil
		sgnr = BundleSignerItem.Value
		ok2 = sgnr != nil
		if ok1 && ok2 && ok3 {
			items = append(items, SendBundleItem)
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
		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
		case BundleSignerItem.Label:
			GetEOASignerUI(BundleSignerItem)
		case BeneficiaryItem.Label:
			_, paddr, ok := AddressFromBookUI("Beneficiary")
			if ok {
				BeneficiaryItem.Value = *paddr
			}
		case SendBundleItem.Label:
			key := sgnr.GetECDSAKey()
			from := ecrypto.PubkeyToAddress(key.PublicKey)
			to := EntryPointItem.Value
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
			}
			gasLimit := usop.TotalGasLimit() + state.GetGasLimitOffset()
			return rpccalls.CreateAndSendTransaction(endpoint, &from, &to, big.NewInt(0), calldata, gasLimit, sgnr)

		}

	}
}

func GetEOASignerUI(it *Item[rpccalls.KeyContainer]) {
	items := []MenuItem{}
	for _, signername := range state.GetSigners() {
		s := state.GetSigner(signername)
		kc, ok := s.(rpccalls.KeyContainer)
		if ok {
			items = append(items, &Item[rpccalls.KeyContainer]{Label: s.Name(), Value: kc})
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
		it.Value = items[sel].(*Item[rpccalls.KeyContainer]).Value
	}
}
