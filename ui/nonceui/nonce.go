package nonceui

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/ecsigner"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/rpcui"
)

var NonceItem = &Item{Label: "Nonce\t", Details: "Set nonce", Value: uint64(0)}
var CallForNonceItem = &Item{Label: "Call for Nonce", Details: "Call for the nonce of the selected address"}

func CallForNonceUI(nit, ait *Item) {
	if ait == nil || ait.Value == nil {
		fmt.Println("No address selected")
		return
	}
	addr := ait.Value.(*common.Address)
	fmt.Println("Calling for nonce of address", addr.Hex())
	items := []*Item{ait, rpcui.SendEndpointItem, CallForNonceItem, Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Set Nonce call parameters",
		Items:     items,
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
		case Back.Label:
			return
		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
		case CallForNonceItem.Label:
			var endpoint *state.RPCEndpoint
			ok1 := false
			if rpcui.SendEndpointItem.Value != nil {
				fmt.Println("Calling Elvis")
				if endpoint, ok1 = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint); ok1 {
					n, err := rpccalls.GetWalletNonce(endpoint, *addr)
					if err != nil {
						fmt.Println(err)
						continue
					}
					nit.Value = n
					return
				} else {
					fmt.Println("No endpoint selected")
					continue
				}
			}
		}
	}
}

func CallForSimpleNonceUI(nit, ait *Item) {
	if ait == nil || ait.Value == nil {
		fmt.Println("No address selected")
		return
	}
	var addr *common.Address
	v := ait.Value
	if v != nil {
		ec, ok := v.(*ecsigner.ECSigner)
		if ok {
			addr = &(ec.SignerAddress)
		} else {
			fmt.Printf("Invalid signer type: %t\n", v)
			return
		}
	} else {
		fmt.Println("Nil signer value")
		return
	}

	fmt.Println("Calling for nonce of address", addr.Hex())
	items := []*Item{ait, rpcui.SendEndpointItem, CallForNonceItem, Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Set Nonce call parameters",
		Items:     items,
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
		case Back.Label:
			return
		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
		case CallForNonceItem.Label:
			var endpoint *state.RPCEndpoint
			ok1 := false
			if rpcui.SendEndpointItem.Value != nil {
				fmt.Println("Calling For Nonce")
				if endpoint, ok1 = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint); ok1 {
					n, err := rpccalls.GetNonce(endpoint, *addr)
					if err != nil {
						fmt.Println(err)
						continue
					}
					nit.Value = n
					return
				} else {
					fmt.Println("No endpoint selected")
					continue
				}
			}
		}
	}
}

func InputNonceUI(nit, ait *Item, simple bool) {
	DirectInputItem := &Item{Label: "Direct Input", Details: "Directly input the nonce"}
	CheckOnChainItem := &Item{Label: "Check On Chain", Details: "Check the nonce on chain"}
	items := []*Item{NonceItem, DirectInputItem, CheckOnChainItem, Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
	}
	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch sel {
	case Back.Label:
		return
	case DirectInputItem.Label:
		InputUint(nit, 64)
	case CheckOnChainItem.Label:
		if simple {
			CallForSimpleNonceUI(nit, ait)
		} else {
			CallForNonceUI(nit, ait)
		}
	default:
		fmt.Println("Not implemented yet:", sel)
	}
}
