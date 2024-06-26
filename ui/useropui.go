package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var UserOpContentItem = &Item{Label: "User Operation Content", Details: "Manage user operation content"}
var ExportUserOpItem = &Item{Label: "Export User Operation", Details: "Select the export format"}
var GetHashItem = &Item{Label: "Get Hash", Details: "Get the hash of the user operation with entrypoint and chainid"}

func UserOpUI() {
	items := []*Item{
		UserOpContentItem,
		ExportUserOpItem,
		GetHashItem,
		Back,
	}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
	}
	usop := new(userop.UserOp)
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case UserOpContentItem.Label:
			UserOpContentUI(usop)
		case ExportUserOpItem.Label:
			ExportUserOpUI(usop)
		case GetHashItem.Label:
			GetHashUI(usop)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

var SenderItem = &Item{Label: state.Sender, Details: "Set sender"}
var NonceItem = &Item{Label: "Nonce", Details: "Set nonce", Value: uint64(0)}
var CallDataItem = &Item{Label: "Call Data", Details: "Set Call Data"}
var CallGasLimitItem = &Item{Label: "Call Gas Limit", Details: "Set Call Gas Limit", Value: userop.DefaultCallGasLimit, DisplayValue: fmt.Sprint(userop.DefaultCallGasLimit)}
var VerificationGasLimitItem = &Item{Label: "Verification Gas Limit", Details: "Set Verification Gas Limit", Value: userop.DefaultVerificationGasLimit, DisplayValue: fmt.Sprint(userop.DefaultVerificationGasLimit)}
var PreVerificationGasItem = &Item{Label: "Pre Verification Gas", Details: "Set Pre Verification Gas", Value: userop.DefaultPreVerificationGas, DisplayValue: fmt.Sprint(userop.DefaultPreVerificationGas)}
var MaxFeePerGasItem = &Item{Label: "Max Fee Per Gas", Details: "Set Max Fee Per Gas", Value: userop.DefaultMaxFeePerGas, DisplayValue: fmt.Sprint(userop.DefaultMaxFeePerGas)}
var MaxPriorityFeePerGasItem = &Item{Label: "Max Priority Fee Per Gas", Details: "Set Max Priority Fee Per Gas", Value: userop.DefaultMaxPriorityFeePerGas, DisplayValue: fmt.Sprint(userop.DefaultMaxPriorityFeePerGas)}
var SignItem = &Item{Label: "Sign", Details: "Sign the user operation"}
var FactoryItem = &Item{Label: "Factory", Details: "Set Factory"}
var FactoryDataItem = &Item{Label: "Factory Data", Details: "Set Factory Data"}

// PaymasterItem is defined in paymasterui.go
var PaymasterDataItem = &Item{Label: "Paymaster Data", Details: "Set Paymaster Data"}
var PaymasterVerificationGasLimitItem = &Item{Label: "Paymaster Verification Gas Limit", Details: "Set Paymaster Verification Gas Limit", Value: userop.DefaultPaymasterVerificationGasLimit, DisplayValue: fmt.Sprint(userop.DefaultPaymasterVerificationGasLimit)}
var PaymasterPostOpGasLimitItem = &Item{Label: "Paymaster Post Op Gas Limit", Details: "Set Paymaster Post Op Gas Limit", Value: userop.DefaultPaymasterPostOpGasLimit, DisplayValue: fmt.Sprint(userop.DefaultPaymasterPostOpGasLimit)}
var SignatureItem = &Item{Label: "Signature", Details: "Set Signature"}

func UserOpContentUI(usop *userop.UserOp) {

	items := []*Item{
		SenderItem,
		NonceItem,
		FactoryItem,
		FactoryDataItem,
		CallDataItem,
		CallGasLimitItem,
		VerificationGasLimitItem,
		PreVerificationGasItem,
		MaxFeePerGasItem,
		MaxPriorityFeePerGasItem,
		PaymasterItem,
		PaymasterDataItem,
		PaymasterVerificationGasLimitItem,
		PaymasterPostOpGasLimitItem,
		SignatureItem,
		Back,
	}
	copyFromUseropToItems(usop)
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
		Size:      15,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			copyValuesToUserOp(usop)

			return
		case NonceItem.Label, CallGasLimitItem.Label, VerificationGasLimitItem.Label,
			PreVerificationGasItem.Label, MaxFeePerGasItem.Label, MaxPriorityFeePerGasItem.Label,
			PaymasterVerificationGasLimitItem.Label, PaymasterPostOpGasLimitItem.Label:
			it, _ := GetItem(sel, items)
			err = InputUint(it, 64)
			if err != nil {
				copyValuesToUserOp(usop)
			}
		case CallDataItem.Label:
			it, _ := GetItem(sel, items)
			caldat, err := PotentiallyRecursiveCallDataUI()
			if err != nil {
				fmt.Println(err)
			} else {
				SetCallDataValue(caldat, it)
			}
		//InputBytes(it)
		case FactoryDataItem.Label, PaymasterDataItem.Label:
			it, _ := GetItem(sel, items)
			err = InputBytes(it, -1)
			if err != nil {
				copyValuesToUserOp(usop)
			}
		case SignatureItem.Label:
			copyValuesToUserOp(usop)
			ret, err := SetSignatureUI(usop)
			if err == nil {
				SignItem.Value = ret
				SignItem.DisplayValue = ShortHex(ret, 7)
				SignItem.Details = hex.EncodeToString(ret)
			} else {
				fmt.Println(err)
			}

		case SenderItem.Label, PaymasterItem.Label, FactoryItem.Label:
			it, _ := GetItem(sel, items)
			addr, ok := AddressFromBookUI(sel)
			if ok {

				it.Value = addr
				it.DisplayValue = addr.String()
			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func copyFromUseropToItems(uop *userop.UserOp) {
	if uop.Sender != nil {
		SenderItem.Value = uop.Sender
		SenderItem.DisplayValue = uop.Sender.String()
	}
	NonceItem.Value = uop.Nonce
	if uop.CallData != nil {
		SetCallDataValue(uop.CallData, CallDataItem)
	}
	CallGasLimitItem.Value = uop.CallGasLimit
	VerificationGasLimitItem.Value = uop.VerificationGasLimit
	PreVerificationGasItem.Value = uop.PreVerificationGas
	MaxFeePerGasItem.Value = uop.MaxFeePerGas
	MaxPriorityFeePerGasItem.Value = uop.MaxPriorityFeePerGas
	if uop.Paymaster != nil {
		PaymasterItem.Value = uop.Paymaster
		PaymasterItem.DisplayValue = uop.Paymaster.String()
	}
	if uop.PaymasterData != nil {
		SetCallDataValue(uop.PaymasterData, PaymasterDataItem)
	}
	PaymasterVerificationGasLimitItem.Value = uop.PaymasterVerificationGasLimit
	PaymasterPostOpGasLimitItem.Value = uop.PaymasterPostOpGasLimit
	if uop.Signature != nil {
		SignatureItem.Value = uop.Signature
		SignatureItem.DisplayValue = hex.EncodeToString(uop.Signature)
	}
}

func copyValuesToUserOp(uop *userop.UserOp) {
	if SenderItem.Value != nil {
		uop.Sender = SenderItem.Value.(*common.Address)
	}
	uop.Nonce = NonceItem.Value.(uint64)
	if CallDataItem.Value != nil {
		uop.CallData = CallDataItem.Value.([]byte)
	}
	uop.CallGasLimit = CallGasLimitItem.Value.(uint64)
	uop.VerificationGasLimit = VerificationGasLimitItem.Value.(uint64)
	uop.PreVerificationGas = PreVerificationGasItem.Value.(uint64)
	uop.MaxFeePerGas = MaxFeePerGasItem.Value.(uint64)
	uop.MaxPriorityFeePerGas = MaxPriorityFeePerGasItem.Value.(uint64)
	if PaymasterItem.Value != nil {
		uop.Paymaster = PaymasterItem.Value.(*common.Address)
	}
	if PaymasterDataItem.Value != nil {
		uop.PaymasterData = PaymasterDataItem.Value.([]byte)
	}
	uop.PaymasterVerificationGasLimit = PaymasterVerificationGasLimitItem.Value.(uint64)
	uop.PaymasterPostOpGasLimit = PaymasterPostOpGasLimitItem.Value.(uint64)
	if SignatureItem.Value != nil {
		uop.Signature = SignatureItem.Value.([]byte)
	}
}

func SetCallDataValue(data []byte, item *Item) {
	item.Value = data
	item.Details = fmt.Sprintf("Call Data: %s", hex.EncodeToString(data))
	if len(data) < 16 {
		item.DisplayValue = hex.EncodeToString(data)
		return
	}
	item.DisplayValue = ShortHex(data, 6)
}

var ExportAsUOPJSONItem = &Item{Label: "Export as JSON", Details: "Export as JSON"}
var ExportAsRemixTupleItem = &Item{Label: "Export as Remix Tuple", Details: "Export as Remix Tuple"}
var ExportAsCurlToEntryItem = &Item{Label: "Export as Curl to Entrypoint", Details: "Export as Curl to Endpoint"}

func ExportUserOpUI(uop *userop.UserOp) {
	items := []*Item{ExportAsUOPJSONItem, ExportAsRemixTupleItem, ExportAsCurlToEntryItem, Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case ExportAsUOPJSONItem.Label:
			ExportAsJSON(uop)
		case ExportAsRemixTupleItem.Label:
			ExportToRemixTuple(uop)
		case ExportAsCurlToEntryItem.Label:
			ExportAsCurl(uop)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func ExportAsJSON(uop *userop.UserOp) {
	bt, err := json.MarshalIndent(uop, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bt))
}

//Sample
/*
["0xCf6290218F6F970657c475E5BFb98Edf45085495",0,"0x","0x42",
"0x0000000000000000000000000000000000000000000000000000000000008233",33332,
"0x0000000000000000000000000000000000000000000000000000000000008235","0x",
"0xfc78e0bcb3b9e4a294d0bfaccebe57111b053679f19dfc31b18486f94a52709129bd1a4f0b446384d851da2ae8076a7aaf56ac3fafae65810602efc3c30efd321B"]
*/

func ExportToRemixTuple(uop *userop.UserOp) {
	fmt.Println(uop.Pack().MarshalRemix())
}

func ExportAsCurl(uop *userop.UserOp) {
	fmt.Println("Not implemented yet")
}
