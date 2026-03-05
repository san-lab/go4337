package nonceui

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/ecsigner"
	"github.com/san-lab/go4337/rpccalls"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/rpcui"
	"github.com/san-lab/go4337/userop"
)

var NonceItem = &Item[*userop.U256]{Label: "Nonce\t", Details: "Set nonce", Value: (*userop.U256)(new(big.Int))}
var CallForNonceItem = &Item[struct{}]{Label: "Call for Nonce", Details: "Call for the nonce of the selected address"}

func CallForNonceUI(nit *Item[*userop.U256], nkit *Item[*big.Int], ait *Item[*common.Address]) {
	if ait == nil || ait.Value == nil {
		fmt.Println("No address selected")
		return
	}
	key := big.NewInt(0)
	if nkit != nil && nkit.Value != nil {
		key = nkit.Value
	}
	addr := ait.Value
	fmt.Println("Calling for nonce of address", addr.Hex())
	items := []MenuItem{ait, nkit, rpcui.SendEndpointItem, CallForNonceItem, Back}
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
		case nkit.Label:
			newKey, err := InputBigInt(nkit)
			if err == nil {
				key = newKey
			}
		case CallForNonceItem.Label:
			endpoint := rpcui.SendEndpointItem.Value
			if endpoint != nil {
				fmt.Println("Calling Elvis")
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

func CallForSimpleNonceUI(nit *Item[*big.Int], ait *Item[any]) {
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
	items := []MenuItem{ait, rpcui.SendEndpointItem, CallForNonceItem, Back}
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
			endpoint := rpcui.SendEndpointItem.Value
			if endpoint != nil {
				fmt.Println("Calling For Nonce")
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

func InputNonceUI(nit *Item[*userop.U256], ait MenuItem, simple bool) {
	fmt.Println("here:", nit.Value)
	DirectInputItem := &Item[struct{}]{Label: "Direct Input", Details: "Directly input the nonce"}
	CheckOnChainItem := &Item[struct{}]{Label: "Check On Chain", Details: "Check the nonce on chain"}
	PureNonceItem := &Item[uint64]{Label: "Nonce (pure)", Details: "Nonce w/o the key"}
	AccAddrItem := &Item[any]{Label: "Address"}
	fullnonce := new(userop.U256)
	if nit.Value != nil {
		if simple {
			fullnonce = userop.FullNonce(big.NewInt(0), (*big.Int)(nit.Value).Uint64())
		} else {
			fullnonce = nit.Value
		}

	}
	pureNonce := fullnonce.PureNonce()
	PureNonceItem.Value = pureNonce

	if ait != nil {
		// Try to extract address value via reflection for AccAddrItem display
		AccAddrItem.Value = getMenuItemValue(ait)
	}
	items := []MenuItem{AccAddrItem, nit}

	key := big.NewInt(0)
	NonceKeyItem := &Item[*big.Int]{Label: "Nonce Key", Value: key}
	if !simple {
		*key = *fullnonce.Key()
		items = append(items, PureNonceItem, NonceKeyItem)
	}

	items = append(items, DirectInputItem, CheckOnChainItem)

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
				nit.Value = (*userop.U256)(big.NewInt(0).Set((*big.Int)(fullnonce)))
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
				// Use a *big.Int item for simple nonce
				simpleNonceItem := &Item[*big.Int]{Label: "Simple Nonce"}
				// We need an any-typed item for the signer
				anyAit := &Item[any]{Label: AccAddrItem.Label, Value: AccAddrItem.Value}
				CallForSimpleNonceUI(simpleNonceItem, anyAit)
				if simpleNonceItem.Value != nil {
					nit.Value = (*userop.U256)(simpleNonceItem.Value)
					PureNonceItem.Value = simpleNonceItem.Value.Uint64()
				}
			} else {
				// Extract *common.Address from ait
				addrItem := extractAddressItem(ait)
				if addrItem != nil {
					CallForNonceUI(nit, NonceKeyItem, addrItem)
					if nit.Value != nil {
						pureNonce = nit.Value.PureNonce()
						PureNonceItem.Value = pureNonce
						NonceKeyItem.Value = nit.Value.Key()
					}
				} else {
					fmt.Println("Cannot extract address for nonce call")
				}
			}

		default:
			//fmt.Println("Not implemented yet:", sel)
		}
	}
}

// getMenuItemValue extracts the underlying value from a MenuItem via reflection
func getMenuItemValue(m MenuItem) any {
	if m == nil {
		return nil
	}
	import_reflect := func() interface{} {
		return nil
	}
	_ = import_reflect
	// We rely on the fact that MenuItem wraps *Item[T]
	// Use type switch for common cases
	switch v := m.(type) {
	case *Item[*common.Address]:
		return v.Value
	case *Item[any]:
		return v.Value
	default:
		return nil
	}
}

// extractAddressItem tries to convert a MenuItem to *Item[*common.Address]
func extractAddressItem(m MenuItem) *Item[*common.Address] {
	if m == nil {
		return nil
	}
	switch v := m.(type) {
	case *Item[*common.Address]:
		return v
	default:
		return nil
	}
}
