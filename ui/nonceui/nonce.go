package nonceui

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/ecsigner"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/rpcui"
	"github.com/san-lab/go4337/userop"
)

var NonceItem = &Item{Label: "Nonce\t", Details: "Set nonce", Value: uint64(0)}
var CallForNonceItem = &Item{Label: "Call for Nonce", Details: "Call for the nonce of the selected address"}

func CallForNonceUI(nit, nkit, ait *Item) {
	if ait == nil || ait.Value == nil {
		fmt.Println("No address selected")
		return
	}
	key := big.NewInt(0)
	if nkit != nil && nkit.Value != nil {
		key = nkit.Value.(*big.Int)
	}
	addr := ait.Value.(*common.Address)
	fmt.Println("Calling for nonce of address", addr.Hex())
	items := []*Item{ait, nkit, rpcui.SendEndpointItem, CallForNonceItem, Back}
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
					n, err := rpccalls.GetWalletNonce(endpoint, *addr, key)
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
	fmt.Println("here:", nit.Value)
	DirectInputItem := &Item{Label: "Direct Input", Details: "Directly input the nonce"}
	CheckOnChainItem := &Item{Label: "Check On Chain", Details: "Check the nonce on chain"}
	PureNonceItem := &Item{Label: "Nonce (pure)", Details: "Nonce w/o the key"}
	AccAddrItem := &Item{Label: "Address"}
	fullnonce := new(userop.U256)
	if nit.Value != nil {
		if simple {
			fullnonce = userop.FullNonce(big.NewInt(0), nit.Value.(*big.Int).Uint64())
		} else {
			fullnonce = (nit.Value).(*userop.U256)
		}

	}
	pureNonce := fullnonce.PureNonce()
	PureNonceItem.Value = pureNonce

	if ait != nil {
		AccAddrItem.Value = ait.Value
	}
	items := []*Item{AccAddrItem, nit}

	key := big.NewInt(0)
	NonceKeyItem := &Item{Label: "Nonce Key", Value: key}
	if !simple {
		*key = *fullnonce.Key()
		items = append(items, PureNonceItem, NonceKeyItem)
	}

	items = append(items, DirectInputItem, CheckOnChainItem)

	//DirectInputItem.Value = pureNonce // ????
	items = append(items, Back)

	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
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
		case DirectInputItem.Label:
			pureNonce, _ = InputUint(PureNonceItem, 64)
			fullnonce = userop.FullNonce(key, pureNonce)
			if simple {
				nit.Value = (*big.Int)(fullnonce)
			} else {
				nit.Value = fullnonce
			}

		case NonceKeyItem.Label:
			key, err = InputBigInt(NonceKeyItem)
			if err != nil {
				fmt.Println(err)
			}
			fullnonce = userop.FullNonce(key, pureNonce)
			nit.Value = fullnonce

		case CheckOnChainItem.Label:
			if simple {
				CallForSimpleNonceUI(nit, ait)
				PureNonceItem.Value = nit.Value.(*big.Int).Uint64()
			} else {
				CallForNonceUI(nit, NonceKeyItem, ait)
				fullnonce := nit.Value.(*userop.U256)
				PureNonceItem.Value = fullnonce.PureNonce()
				NonceKeyItem.Value = fullnonce.Key()
			}

		default:
			//fmt.Println("Not implemented yet:", sel)
		}
	}
}
