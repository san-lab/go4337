package setauth

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/ui/abicalldata"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/nonceui"
	"github.com/san-lab/go4337/ui/signui"
	"github.com/san-lab/go4337/userop"
)

var SetCodeChainIDItem = &Item[*big.Int]{Label: "Chain ID"}

var SetCodeNonceItem = &Item[*userop.U256]{Label: "Nonce\t", Details: "Set auth nonce", Value: (*userop.U256)(new(big.Int))}
var SetCodeGasTipCapItem = &Item[*big.Int]{Label: "GasTipCap/MaxPriorityFee"}
var SetCodeGasFeeCapItem = &Item[*uint256.Int]{Label: "GasFeeCap/MaxFeePerGas"}
var SetCodeGasItem = &Item[uint64]{Label: "Gas"}
var SetCodeToItem = &Item[string]{Label: "To"}
var SetCodeValueItem = &Item[*big.Int]{Label: "Value"}
var SetCodeDataItem = &Item[[]byte]{Label: "Data"}
var SetCodeAccessListItem = &Item[any]{Label: "Access List"}
var SetCodeAuthListItem = &Item[any]{Label: "Authorization List", Display: func(v any) string {
	if al, ok := v.([]types.SetCodeAuthorization); ok {
		return fmt.Sprintf("[%d] Authorizations", len(al))
	}
	return "[]"
}}
var SetCodeSignItem = &Item[struct{}]{Label: "Sign"}

func SetAuthTxUI(satx *types.SetCodeTx) {
	if satx == nil {
		satx = &types.SetCodeTx{}
	} else {
		//Copy the values to the UI Items
		authToItems(satx)
	}

	items := []MenuItem{
		SetCodeChainIDItem,
		SetCodeSignerItem,
		SetCodeNonceItem,
		SetCodeGasTipCapItem,
		SetCodeGasFeeCapItem,
		SetCodeGasItem,
		SetCodeToItem,
		SetCodeValueItem,
		SetCodeDataItem,
		SetCodeAccessListItem,
		SetCodeAuthListItem,
		SetCodeSignItem,
		Back,
	}
	prompt := promptui.Select{Items: items, Label: "Select an item"}
	prompt.Templates = ItemTemplate
	prompt.Size = 14
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case SetCodeChainIDItem.Label:
			InputBigInt(SetCodeChainIDItem)
			if SetCodeChainIDItem.Value != nil {
				chid, _ := uint256.FromBig(SetCodeChainIDItem.Value)
				satx.ChainID = chid

			}
		case SetCodeSignerItem.Label:
			signui.SignerUI(SetCodeSignerItem)
		case SetCodeNonceItem.Label:
			nonceui.InputNonceUI(SetCodeNonceItem, SetCodeSignerItem, true)
			if SetCodeNonceItem.Value != nil {
				satx.Nonce = (*big.Int)(SetCodeNonceItem.Value).Uint64()
			}
		case SetCodeGasTipCapItem.Label:
			InputBigInt(SetCodeGasTipCapItem)
			if SetCodeGasTipCapItem.Value != nil {
				gtc, _ := uint256.FromBig(SetCodeGasTipCapItem.Value)

				satx.GasTipCap = gtc

			}
		case SetCodeGasFeeCapItem.Label:
			InputUint256(SetCodeGasFeeCapItem)
			if SetCodeGasFeeCapItem.Value != nil {
				satx.GasFeeCap = SetCodeGasFeeCapItem.Value

			}
		case SetCodeGasItem.Label:
			InputUint(SetCodeGasItem, 64)
			satx.Gas = SetCodeGasItem.Value
		case SetCodeToItem.Label:
			_, addr, ok := AddressFromBookUI("To Address")
			if ok {
				SetCodeToItem.Value = addr.String()
				satx.To = *addr
			}
		case SetCodeValueItem.Label:
			InputBigInt(SetCodeValueItem)
			if SetCodeValueItem.Value != nil {
				val, _ := uint256.FromBig(SetCodeValueItem.Value)

				satx.Value = val

			}
		case SetCodeDataItem.Label:
			calldata, err := abicalldata.PotentiallyRecursiveCallDataUI()
			if err != nil {
				fmt.Println("Error in Call Data UI:", err)
				continue
			}
			SetCodeDataItem.Value = calldata
			satx.Data = calldata
		case SetCodeAccessListItem.Label:
			fmt.Println("Not implemented yet:", sel)
		case SetCodeAuthListItem.Label:
			AuthListUI(satx)
		case SetCodeSignItem.Label:
			SignSetCodeUI(satx)
		}

	}
}

func authToItems(satx *types.SetCodeTx) {
	SetCodeChainIDItem.Value = satx.ChainID.ToBig()
	SetCodeNonceItem.Value = (*userop.U256)(big.NewInt(int64(satx.Nonce)))
	SetCodeGasTipCapItem.Value = satx.GasTipCap.ToBig()
	SetCodeGasFeeCapItem.Value = satx.GasFeeCap
	SetCodeGasItem.Value = satx.Gas
	SetCodeToItem.Value = satx.To.String()
	SetCodeValueItem.Value = satx.Value.ToBig()
	SetCodeDataItem.Value = satx.Data
	SetCodeAccessListItem.Value = satx.AccessList
	SetCodeAuthListItem.Value = satx.AuthList
	SetCodeSignItem.Display = func(_ struct{}) string { return RSVtoString(satx.R, satx.S, satx.V) }
}

func itemsToAuth(satx *types.SetCodeTx) {
	if SetCodeChainIDItem.Value != nil {
		chid, _ := uint256.FromBig(SetCodeChainIDItem.Value)
		satx.ChainID = chid
	}
	if SetCodeNonceItem.Value != nil {
		satx.Nonce = (*big.Int)(SetCodeNonceItem.Value).Uint64()
	}
	if SetCodeGasTipCapItem.Value != nil {
		gtc, _ := uint256.FromBig(SetCodeGasTipCapItem.Value)
		satx.GasTipCap = gtc
	}
	satx.GasFeeCap = SetCodeGasFeeCapItem.Value
	satx.Gas = SetCodeGasItem.Value
	satx.To = common.HexToAddress(SetCodeToItem.Value)
	if SetCodeValueItem.Value != nil {
		val, _ := uint256.FromBig(SetCodeValueItem.Value)
		satx.Value = val
	}
	satx.Data = SetCodeDataItem.Value
	if al, ok := SetCodeAccessListItem.Value.(types.AccessList); ok {
		satx.AccessList = al
	}
	if al, ok := SetCodeAuthListItem.Value.([]types.SetCodeAuthorization); ok {
		satx.AuthList = al
	}
}

var AddAuthItem = &Item[struct{}]{Label: "Add Authorization"}
var RemoveAuthItem = &Item[struct{}]{Label: "Remove Authorization"}
var AuthImportItem = &Item[struct{}]{Label: "Import Authorization", Details: "Import the authorization from hex/rlp"}

func AuthListUI(satx *types.SetCodeTx) {
	if satx == nil {
		fmt.Println("No SetCodeTx")
		return
	}
	isRemoving := false

	spr := promptui.Select{Label: "Select an item"}
	spr.Templates = ItemTemplate

	for {
		authitems := []MenuItem{}
		for _, a := range satx.AuthList {
			aCopy := a
			authitems = append(authitems, &Item[types.SetCodeAuthorization]{
				Label:   fmt.Sprintf("%s_%d_%d", a.Address.Hex(), a.ChainID.Uint64(), a.Nonce),
				Value:   aCopy,
				Display: func(v types.SetCodeAuthorization) string { return AuthorityString(&v) },
			})
		}
		if !isRemoving {
			authitems = append(authitems, AddAuthItem, RemoveAuthItem, AuthImportItem)
		}
		authitems = append(authitems, Back)
		spr.Items = authitems
		idx, sel, err := spr.Run()
		if err != nil {
			fmt.Println("Prompt failed", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case AddAuthItem.Label:
			auth := &types.SetCodeAuthorization{}
			auth = AuthUI(auth)
			if auth != nil {
				satx.AuthList = append(satx.AuthList, *auth)
			}
		case AuthImportItem.Label:
			auth := &types.SetCodeAuthorization{}
			it := &Item[[]byte]{Label: "Input hex/rlp of the authorization"}
			err := InputBytes(it, -1)
			if err != nil {
				fmt.Println("Error importing authorization:", err)
				continue
			}
			if it.Value != nil {
				err = rlp.DecodeBytes(it.Value, auth)
				if err != nil {
					fmt.Println("Error importing authorization:", err)
					continue
				}
				satx.AuthList = append(satx.AuthList, *auth)
			}
		case RemoveAuthItem.Label:
			isRemoving = true
			continue
		default:
			if idx >= len(satx.AuthList) || idx < 0 {
				fmt.Println("Invalid index")
				continue
			}
			auth := satx.AuthList[idx]
			if isRemoving {
				satx.AuthList = append(satx.AuthList[:idx], satx.AuthList[idx+1:]...)
				isRemoving = false
			} else {
				fmt.Println("Editing Authorization", idx)
				AuthUI(&auth)
				satx.AuthList[idx] = auth

			}
		}
	}

}

var SetCodeSignerItem = &Item[signer.Signer]{Label: "SetCode Signer"}

func SignSetCodeUI(satx *types.SetCodeTx) {
	if satx == nil {
		fmt.Println("No transaction to sign")
		return
	}

	spr := promptui.Select{Items: []MenuItem{SetCodeSignerItem, SetCodeSignItem, Back}, Label: "Select an item"}
	spr.Templates = ItemTemplate
	for {
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("Prompt failed", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case SetCodeSignerItem.Label:
			signui.SignerUI(SetCodeSignerItem)
		case SetCodeSignItem.Label:
			fmt.Println("Signing transaction...")
			signerV := SetCodeSignerItem.Value
			if signerV == nil {
				fmt.Println("No signer selected")
				continue
			}
			nstx, err := types.SignNewTx(signerV.GetKey().(*ecdsa.PrivateKey), types.NewPragueSigner(satx.ChainID.ToBig()), satx)
			if err != nil {
				fmt.Println("Error signing transaction:", err)
				continue
			}
			v, r, s := nstx.RawSignatureValues()
			satx.V, _ = uint256.FromBig(v)
			satx.R, _ = uint256.FromBig(r)
			satx.S, _ = uint256.FromBig(s)
			SetCodeSignItem.Display = func(_ struct{}) string { return RSVtoString(satx.R, satx.S, satx.V) }
			return
		}

	}
}
