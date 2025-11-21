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
)

var SetCodeChainIDItem = &Item{Label: "Chain ID"}

var SetCodeNonceItem = nonceui.NonceItem
var SetCodeGasTipCapItem = &Item{Label: "GasTipCap/MaxPriorityFee"}
var SetCodeGasFeeCapItem = &Item{Label: "GasFeeCap/MaxFeePerGas"}
var SetCodeGasItem = &Item{Label: "Gas"}
var SetCodeToItem = &Item{Label: "To"}
var SetCodeValueItem = &Item{Label: "Value"}
var SetCodeDataItem = &Item{Label: "Data"}
var SetCodeAccessListItem = &Item{Label: "Access List"}
var SetCodeAuthListItem = &Item{Label: "Authorization List"}
var SetCodeSignItem = &Item{Label: "Sign"}

func SetAuthTxUI(satx *types.SetCodeTx) {
	if satx == nil {
		satx = &types.SetCodeTx{}
	} else {
		//Copy the values to the UI Items
		authToItems(satx)
	}

	items := []*Item{
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
		SetCodeAuthListItem.DisplayValueString = fmt.Sprintf("[%d] Authorizations", len(satx.AuthList))

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
				chid, _ := uint256.FromBig(SetCodeChainIDItem.Value.(*big.Int))
				satx.ChainID = chid

			}
		case SetCodeSignerItem.Label:
			signui.SignerUI(SetCodeSignerItem)
		case SetCodeNonceItem.Label:
			/*
				InputUint(SetCodeNonceItem, 64)
				if SetCodeNonceItem.Value != nil {
					satx.Nonce = SetCodeNonceItem.Value.(uint64)
				}
			*/

			nonceui.InputNonceUI(SetCodeNonceItem, SetCodeSignerItem, true)
			satx.Nonce = SetCodeNonceItem.Value.(uint64)
		case SetCodeGasTipCapItem.Label:
			InputBigInt(SetCodeGasTipCapItem)
			if SetCodeGasTipCapItem.Value != nil {
				gtc, _ := uint256.FromBig(SetCodeGasTipCapItem.Value.(*big.Int))

				satx.GasTipCap = gtc

			}
		case SetCodeGasFeeCapItem.Label:
			InputUint256(SetCodeGasFeeCapItem)
			if SetCodeGasFeeCapItem.Value != nil {
				satx.GasFeeCap = SetCodeGasFeeCapItem.Value.(*uint256.Int)

			}
		case SetCodeGasItem.Label:
			InputUint(SetCodeGasItem, 64)
			if SetCodeGasItem.Value != nil {
				satx.Gas = SetCodeGasItem.Value.(uint64)
			}
		case SetCodeToItem.Label:
			_, addr, ok := AddressFromBookUI("To Address")
			if ok {
				SetCodeToItem.Value = addr.String()
				satx.To = *addr
			}
		case SetCodeValueItem.Label:
			InputBigInt(SetCodeValueItem)
			if SetCodeValueItem.Value != nil {
				val, _ := uint256.FromBig(SetCodeValueItem.Value.(*big.Int))

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
	SetCodeChainIDItem.Value = satx.ChainID
	SetCodeNonceItem.Value = satx.Nonce
	SetCodeGasTipCapItem.Value = satx.GasTipCap
	SetCodeGasFeeCapItem.Value = satx.GasFeeCap
	SetCodeGasItem.Value = satx.Gas
	SetCodeToItem.Value = satx.To
	SetCodeValueItem.Value = satx.Value
	SetCodeDataItem.Value = satx.Data
	SetCodeAccessListItem.Value = satx.AccessList
	SetCodeAuthListItem.Value = satx.AuthList
	SetCodeSignItem.DisplayValueString = RSVtoString(satx.R, satx.S, satx.V)
}

func itemsToAuth(satx *types.SetCodeTx) {
	satx.ChainID = SetCodeChainIDItem.Value.(*uint256.Int)
	satx.Nonce = SetCodeNonceItem.Value.(uint64)
	satx.GasTipCap = SetCodeGasTipCapItem.Value.(*uint256.Int)
	satx.GasFeeCap = SetCodeGasFeeCapItem.Value.(*uint256.Int)
	satx.Gas = SetCodeGasItem.Value.(uint64)
	satx.To = SetCodeToItem.Value.(common.Address)
	satx.Value = SetCodeValueItem.Value.(*uint256.Int)
	satx.Data = SetCodeDataItem.Value.([]byte)
	satx.AccessList = SetCodeAccessListItem.Value.(types.AccessList)
	satx.AuthList = SetCodeAuthListItem.Value.([]types.SetCodeAuthorization)
}

var AddAuthItem = &Item{Label: "Add Authorization"}
var RemoveAuthItem = &Item{Label: "Remove Authorization"}
var AuthImportItem = &Item{Label: "Import Authorization", Details: "Import the authorization from hex/rlp"}

func AuthListUI(satx *types.SetCodeTx) {
	if satx == nil {
		fmt.Println("No SetCodeTx")
		return
	}
	isRemoving := false

	spr := promptui.Select{Label: "Select an item"}
	spr.Templates = ItemTemplate

	for {
		authitems := []*Item{}
		for _, a := range satx.AuthList {
			authitems = append(authitems, &Item{
				Label:              fmt.Sprintf("%s_%d_%d", a.Address.Hex(), a.ChainID.Uint64(), a.Nonce),
				Value:              a,
				DisplayValueString: AuthorityString(&a),
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
			it := &Item{Label: "Input hex/rlp of the authorization"}
			err := InputBytes(it, -1)
			if err != nil {
				fmt.Println("Error importing authorization:", err)
				continue
			}
			if it.Value != nil {
				err = rlp.DecodeBytes(it.Value.([]byte), auth)
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

var SetCodeSignerItem = &Item{Label: "SetCode Signer"}

func SignSetCodeUI(satx *types.SetCodeTx) {
	if satx == nil {
		fmt.Println("No transaction to sign")
		return
	}

	spr := promptui.Select{Items: []*Item{SetCodeSignerItem, SetCodeSignItem, Back}, Label: "Select an item"}
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
			signer, ok := signerV.(signer.Signer)
			if !ok {
				fmt.Println("Invalid signer type")
				continue
			}
			nstx, err := types.SignNewTx(signer.GetKey().(*ecdsa.PrivateKey), types.NewPragueSigner(satx.ChainID.ToBig()), satx)
			if err != nil {
				fmt.Println("Error signing transaction:", err)
				continue
			}
			v, r, s := nstx.RawSignatureValues()
			satx.V, _ = uint256.FromBig(v)
			satx.R, _ = uint256.FromBig(r)
			satx.S, _ = uint256.FromBig(s)
			SetCodeSignItem.DisplayValueString = RSVtoString(satx.R, satx.S, satx.V)
			return
		}

	}
}
