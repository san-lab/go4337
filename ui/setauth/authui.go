package setauth

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/nonceui"
	"github.com/san-lab/go4337/ui/signui"
	"github.com/san-lab/go4337/userop"
)

var AuthChainIDItem = &common.Item[*big.Int]{Label: "Chain ID", Details: "The chain ID must be the same as TX or 0"}
var AuthAddressItem = &common.Item[string]{Label: "Address", Details: "The addrss to instrument"}
var AuthNonceItem = &common.Item[*userop.U256]{Label: "Nonce", Details: "Current nonce at the authority ", Value: (*userop.U256)(new(big.Int))}
var AuthSignItem = &common.Item[struct{}]{Label: "Sign", Details: "Sign the authorization"}
var AuthExportItem = &common.Item[struct{}]{Label: "Export as rlp or json", Details: "Export the authorization as hex of the rlp encoded data or as json"}
var AuthResetItem = &common.Item[struct{}]{Label: "Remove"}

func AuthUI(auth *types.SetCodeAuthorization) *types.SetCodeAuthorization {
	if auth == nil {
		auth = new(types.SetCodeAuthorization)
		AuthChainIDItem.Value = nil
		AuthAddressItem.Value = ""
		AuthNonceItem.Value = (*userop.U256)(new(big.Int))
		AuthSignItem.Display = nil
		AuthSignerItem.Details = ""
	} else {
		AuthChainIDItem.Value = auth.ChainID.ToBig()
		AuthAddressItem.Value = auth.Address.String()
		AuthNonceItem.Value = (*userop.U256)(big.NewInt(int64(auth.Nonce)))
		AuthSignItem.Display = func(_ struct{}) string {
			return RSVtoString(&auth.R, &auth.S, uint256.NewInt(uint64(auth.V)))
		}
		AuthSignerItem.Details = AuthorityString(auth)
	}

	authitems := []common.MenuItem{
		AuthChainIDItem,
		AuthAddressItem,
		AuthSignerItem,
		AuthNonceItem,
		AuthSignItem,
		common.Set,
		AuthResetItem,
		AuthExportItem,
		common.Back,
	}
	spr := promptui.Select{Items: authitems, Label: "Select an item"}
	spr.Templates = common.ItemTemplate
	spr.Size = 10
	for {
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("Prompt failed", err)
			return nil
		}
		switch sel {
		case common.Back.Label:
			return nil
		case common.Set.Label:
			return auth
		case AuthExportItem.Label:
			ExoprtAuthUI(auth)
			return nil

		case AuthChainIDItem.Label:
			common.InputBigInt(AuthChainIDItem)
			if AuthChainIDItem.Value != nil {
				chid, _ := uint256.FromBig(AuthChainIDItem.Value)
				auth.ChainID = *chid

			}
		case AuthAddressItem.Label:
			_, addr, ok := common.AddressFromBookUI("Target Address")
			if ok {
				AuthAddressItem.Value = addr.String()
				auth.Address = *addr
			}
		case AuthSignerItem.Label:
			signui.SignerUI(AuthSignerItem)

		case AuthNonceItem.Label:
			nonceui.InputNonceUI(AuthNonceItem, AuthSignerItem, true)
			if AuthNonceItem.Value != nil {
				auth.Nonce = (*big.Int)(AuthNonceItem.Value).Uint64()
			}
		case AuthSignItem.Label:
			SignAuthorizationUI(auth)
		case AuthResetItem.Label:
			return nil
		}

	}

}

var AuthSignerItem = &common.Item[signer.Signer]{Label: "Auth Signer", Details: "The signer to use for signing the authorization"}

func SignAuthorizationUI(auth *types.SetCodeAuthorization) {
	if auth == nil {
		fmt.Println("No authorization to sign")
		return
	}
	spr := promptui.Select{Items: []common.MenuItem{AuthSignerItem, AuthSignItem, common.Back}, Label: "Select an item"}
	spr.Templates = common.ItemTemplate
	for {
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("Prompt failed", err)
			return
		}
		switch sel {
		case common.Back.Label:
			return
		case AuthSignerItem.Label:
			signui.SignerUI(AuthSignerItem)

		case AuthSignItem.Label:
			if AuthSignerItem.Value == nil {
				fmt.Println("No signer selected")
			} else {

				fmt.Println("Signing with", AuthSignerItem.Value)
				signerV := AuthSignerItem.Value
				a2, err := types.SignSetCode(signerV.GetKey().(*ecdsa.PrivateKey), *auth)
				if err != nil {
					fmt.Println("Failed to sign", err)
				} else {
					auth.R = a2.R
					auth.S = a2.S
					auth.V = uint8(a2.V)

					AuthSignItem.Display = func(_ struct{}) string {
						return RSVtoString(&a2.R, &a2.S, uint256.NewInt(uint64(a2.V)))
					}
					return
				}

			}

		}
	}

}

func AuthorityString(auth *types.SetCodeAuthorization) string {
	if auth == nil {
		return "Not Set"
	}
	authority, err := auth.Authority()
	if err != nil {
		return "Failed to get authority"
	}
	return "authority: " + authority.Hex()
}

func RSVtoString(r, s, v *uint256.Int) string {
	if r == nil || s == nil || v == nil {
		return "Not Set"
	}
	return fmt.Sprintf("R: %s, S: %s, V: %v", common.ShortHex(r.Bytes(), 12), common.ShortHex(s.Bytes(), 12), v)
}

func ExoprtAuthUI(auth *types.SetCodeAuthorization) {
	exportAsHex := &common.Item[struct{}]{Label: "Export as hex", Details: "Export the authorization as hex of the rlp encoded data"}
	exportAsJSON := &common.Item[struct{}]{Label: "Export as JSON", Details: "Export the authorization as JSON"}
	selector := promptui.Select{Items: []common.MenuItem{exportAsHex, exportAsJSON, common.Back}, Label: "Select an item"}
	selector.Templates = common.ItemTemplate
	selector.Size = 4
	_, sel, err := selector.Run()
	if err != nil {
		fmt.Println("Prompt failed", err)
		return
	}
	switch sel {
	case common.Back.Label:
		return
	case exportAsHex.Label:
		msg := ExportAuthAsHex(auth)
		fmt.Println("Authorization in hex:", msg)
	case exportAsJSON.Label:
		msg := ExportAuthAsJSON(auth)
		fmt.Println("Authorization in JSON:", msg)
	}
}

func ExportAuthAsHex(auth *types.SetCodeAuthorization) string {
	if auth == nil {
		return "Not Set"
	}
	// Marshal the authorization to RLP
	rlpData, err := rlp.EncodeToBytes(auth)
	if err != nil {
		return "Failed to encode authorization"
	}
	// Convert the RLP data to hex
	hexData := fmt.Sprintf("0x%x", rlpData)
	return hexData

}

func ExportAuthAsJSON(auth *types.SetCodeAuthorization) string {
	if auth == nil {
		return "Not Set"
	}
	// Marshal the authorization to JSON
	jsonData, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return "Failed to encode authorization"
	}
	// Convert the JSON data to string
	jsonString := string(jsonData)
	return jsonString
}
