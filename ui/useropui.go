package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var SelectUserOpItem = &Item{Label: "Select User Operation", Details: "Select a user operation"}
var CreateUserOpItem = &Item{Label: "Create User Operation", Details: "Create a new user operation"}
var CloneUserOpItem = &Item{Label: "Clone User Operation", Details: "Clone a user operation"}

func TopUserOpUI() {

	usOpItem := &Item{}
	workItem := &Item{}
	for {
		items := []*Item{}
		if usOpItem.Value != nil {
			if usop, ok := usOpItem.Value.(*userop.UserOp); ok {
				workItem = &Item{Label: "Work on User Operation: " + usOpItem.Label, Details: fmt.Sprintf("%s/%v", usop.Sender.String(), usop.Nonce)}
				items = append(items, workItem)
			}
		}
		items = append(items, []*Item{
			SelectUserOpItem,
			CreateUserOpItem,
			CloneUserOpItem,
		}...)
		//Select a userOp, create a new one, or go back

		items = append(items, Back)

		// Create a new select prompt
		prompt := promptui.Select{
			Label:     "Select an option",
			Items:     items,
			Templates: ItemTemplate,
			Size:      10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case SelectUserOpItem.Label:

			SelectUserOpUI(usOpItem)
		case CreateUserOpItem.Label:
			CreateUserOpUI(usOpItem)
		case CloneUserOpItem.Label:
			CloneUserOpUI(usOpItem)
		case workItem.Label:
			usop := usOpItem.Value.(*userop.UserOp) //Has been checked when generating ui, and there should be no concurrency, so it is safe
			UserOpUI(usop)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func CloneUserOpUI(topIt *Item) {
	it := &Item{}
	SelectUserOpUI(it)
	if it.Value == nil {
		return
	}
	usop, ok := it.Value.(*userop.UserOp)
	if !ok {
		fmt.Println("Invalid UserOp selected. This should be impossible...")
		return
	}
	//Prompt for a new name
	newname := ""
	for {
		prompt := promptui.Prompt{
			Label: "Enter new UserOp Name",
		}
		name, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		newname = name
		if _, ok := state.State.UserOps[name]; !ok {
			break
		}
		fmt.Println("UserOp already exists")

	}
	if len(newname) == 0 {
		return
	}
	clone, err := abiutil.CloneStruct(usop)
	if err != nil {
		fmt.Printf("Error cloning UserOp: %v\n", err)
		return
	}

	topIt.Value = clone
	topIt.Label = newname
	state.State.UserOps[newname] = clone.(*userop.UserOp)
	state.State.Save()
}

func SelectUserOpUI(topit *Item) {
	items := []*Item{}
	keys := make([]string, 0, len(state.State.UserOps))
	for k := range state.State.UserOps {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		fmt.Println("No UserOps available")
		return
	}
	sort.Strings(keys)
	for _, name := range keys {
		uop := state.State.UserOps[name]
		items = append(items, &Item{Label: name, Details: fmt.Sprintf("Sender: %s, Nonce: %d", uop.Sender.String(), uop.Nonce)})
	}
	items = append(items, Back)
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select UserOperation",
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
	}
	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	if sel == Back.Label {
		return
	}
	topit.Value = state.State.UserOps[sel]
	topit.Label = sel

}

func CreateUserOpUI(topIt *Item) {
	//prompt for name
	prompt := promptui.Prompt{
		Label: "Enter UserOp Name",
	}
	name, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, ok := state.State.UserOps[name]; ok {
		fmt.Println("UserOp already exists")
		return
	}
	topIt.Label = name
	UserOpContentUI(topIt)
	state.State.UserOps[name] = topIt.Value.(*userop.UserOp)
	state.State.Save()

}

var UserOpContentItem = &Item{Label: "User Operation Content", Details: "Manage user operation content"}
var ExportUserOpItem = &Item{Label: "Export User Operation", Details: "Select the export format"}
var GetHashItem = &Item{Label: "Hashes and signatures", Details: "Get the hash of the user operation with entrypoint and chainid"}

func UserOpUI(usop *userop.UserOp) {
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

	for {
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}
		switch sel {
		case Back.Label:
			return
		case UserOpContentItem.Label:
			it := &Item{Value: usop}
			UserOpContentUI(it)
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
var CallGasLimitItem = &Item{Label: "Call Gas Limit", Details: "Set Call Gas Limit", Value: userop.DefaultCallGasLimit}
var VerificationGasLimitItem = &Item{Label: "Verification Gas Limit", Details: "Set Verification Gas Limit", Value: userop.DefaultVerificationGasLimit}
var PreVerificationGasItem = &Item{Label: "Pre Verification Gas", Details: "Set Pre Verification Gas", Value: userop.DefaultPreVerificationGas}
var MaxFeePerGasItem = &Item{Label: "Max Fee Per Gas", Details: "Set Max Fee Per Gas", Value: userop.DefaultMaxFeePerGas}
var MaxPriorityFeePerGasItem = &Item{Label: "Max Priority Fee Per Gas", Details: "Set Max Priority Fee Per Gas", Value: userop.DefaultMaxPriorityFeePerGas}
var SignItem = &Item{Label: "Sign", Details: "Sign the user operation"}
var FactoryItem = &Item{Label: "Factory", Details: "Set Factory"}
var FactoryDataItem = &Item{Label: "Factory Data", Details: "Set Factory Data"}

// PaymasterItem is defined in paymasterui.go
var PaymasterDataItem = &Item{Label: "Paymaster Data", Details: "Set Paymaster Data"}
var PaymasterVerificationGasLimitItem = &Item{Label: "Paymaster Verification Gas Limit", Details: "Set Paymaster Verification Gas Limit", Value: userop.DefaultPaymasterVerificationGasLimit}
var PaymasterPostOpGasLimitItem = &Item{Label: "Paymaster Post Op Gas Limit", Details: "Set Paymaster Post Op Gas Limit", Value: userop.DefaultPaymasterPostOpGasLimit}
var SignatureItem = &Item{Label: "Signature", Details: "Set Signature"}

func UserOpContentUI(topIt *Item) {
	var usop *userop.UserOp
	if topIt.Value == nil {
		usop = new(userop.UserOp)
	} else {
		ok := false
		usop, ok = topIt.Value.(*userop.UserOp)
		if !ok {
			fmt.Println("Invalid UserOp passed to UserOpContentUI")
			return
		}
	}

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
			topIt.Value = usop
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
				SignatureItem.Value = ret
				//SignItem.DisplayValue = ShortHex(ret, 7)
				SignatureItem.Details = hex.EncodeToString(ret)
			} else {
				fmt.Println(err)
			}

		case SenderItem.Label, PaymasterItem.Label, FactoryItem.Label:
			it, _ := GetItem(sel, items)
			addr, ok := AddressFromBookUI(sel)
			if ok {

				it.Value = addr

			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func copyFromUseropToItems(uop *userop.UserOp) {
	if uop.Sender != nil {
		SenderItem.Value = uop.Sender

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

	}
	if uop.PaymasterData != nil {
		SetCallDataValue(uop.PaymasterData, PaymasterDataItem)
	}
	PaymasterVerificationGasLimitItem.Value = uop.PaymasterVerificationGasLimit
	PaymasterPostOpGasLimitItem.Value = uop.PaymasterPostOpGasLimit
	//if uop.Signature != nil {
	SignatureItem.Value = uop.Signature

	//}
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

}

var ExportAsUOPJSONItem = &Item{Label: "Export as JSON", Details: "Export as JSON"}
var ExportAsRemixTupleV7Item = &Item{Label: "Export as Remix Tuple (V7)", Details: "Export as Remix Tuple"}
var ExportAsRemixTupleV6Item = &Item{Label: "Export as Remix Tuple (V6)", Details: "Export as Remix Tuple"}
var ExportAsCurlToEntryItem = &Item{Label: "Export as Curl to Entrypoint", Details: "Export as Curl to Endpoint"}

func ExportUserOpUI(uop *userop.UserOp) {
	items := []*Item{ExportAsUOPJSONItem, ExportAsRemixTupleV6Item, ExportAsRemixTupleV7Item, ExportAsCurlToEntryItem, Back}
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
		case ExportAsRemixTupleV6Item.Label:
			fmt.Println(uop.MarshalRemixV6())
		case ExportAsRemixTupleV7Item.Label:
			fmt.Println(uop.Pack().MarshalRemix()) //"Pack" only makes sence for v7+ userops
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

func ExportAsCurl(uop *userop.UserOp) {
	fmt.Println("Not implemented yet")
}
