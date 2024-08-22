package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/userop"
)

var SelectUserOpItem = &Item{Label: "Select User Operation", Details: "Select a user operation"}
var CreateUserOpItem = &Item{Label: "Create User Operation", Details: "Create a new user operation"}
var CloneUserOpItem = &Item{Label: "Clone User Operation", Details: "Clone a user operation"}

func TopUserOpUI(callbackItem *Item) {

	usOpItem := &Item{}
	workItem := &Item{}
	deleteItem := &Item{}
	for {
		items := []*Item{}
		if usOpItem.Value != nil {
			if usop, ok := usOpItem.Value.(*userop.UserOperation); ok {
				if callbackItem != nil {
					callbackItem.Label = "Select this UserOperation"
					callbackItem.Value = usop
					items = append(items, callbackItem)
				}

				workItem = &Item{Label: "Work on User Operation: " + usOpItem.Label, Details: fmt.Sprintf("%s/%v", usop.Sender, usop.Nonce)}
				deleteItem = &Item{Label: "Delete User Operation: " + usOpItem.Label, Details: "Delete this user operation"}
				items = append(items, workItem, deleteItem)

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
			usop := usOpItem.Value.(*userop.UserOperation) //Has been checked when generating ui, and there should be no concurrency, so it is safe
			UserOpUI(usop)
		case deleteItem.Label:
			state.RemoveUserOp(usOpItem.Label)
			return
		case callbackItem.Label:
			return
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func InputNonceUI(nit, ait *Item) {
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
		CallForNonceUI(nit, ait)
	default:
		fmt.Println("Not implemented yet:", sel)
	}
}

var CallForNonceItem = &Item{Label: "Call for Nonce", Details: "Call for the nonce of the selected address"}

func CallForNonceUI(nit, ait *Item) {
	if ait == nil || ait.Value == nil {
		fmt.Println("No address selected")
		return
	}
	addr := ait.Value.(*common.Address)
	items := []*Item{ait, SendEndpointItem, CallForNonceItem, Back}
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
		case SendEndpointItem.Label:
			RPCEndpointsUI(SendEndpointItem)
		case CallForNonceItem.Label:
			var endpoint *state.RPCEndpoint
			ok1 := false
			if SendEndpointItem.Value != nil {
				fmt.Println("Calling Elvis")
				if endpoint, ok1 = SendEndpointItem.Value.(*state.RPCEndpoint); ok1 {
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

func CloneUserOpUI(topIt *Item) {
	it := &Item{}
	SelectUserOpUI(it)
	if it.Value == nil {
		return
	}
	usop, ok := it.Value.(*userop.UserOperation)
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
		if _, ok := state.GetUserOps()[name]; !ok {
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
	nusop := clone.(*userop.UserOperation)
	nusop.Nonce++
	topIt.Value = clone
	topIt.Label = newname
	state.AddUserOp(newname, nusop)

}

func SelectUserOpUI(topit *Item) {
	items := []*Item{}
	keys := make([]string, 0, len(state.GetUserOps()))
	for k := range state.GetUserOps() {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		fmt.Println("No UserOps available")
		return
	}
	sort.Strings(keys)
	for _, name := range keys {
		uop, _ := state.GetUserOp(name)
		items = append(items, &Item{Label: name, Details: fmt.Sprintf("Sender: %s, Nonce: %d", uop.Sender, uop.Nonce)})
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
	topit.Value, _ = state.GetUserOp(sel)
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
	if _, ok := state.GetUserOps()[name]; ok {
		fmt.Println("UserOp already exists")
		return
	}
	nuop := userop.NewUserOperationWithDefaults()
	state.AddUserOp(name, nuop)
	if err != nil {
		fmt.Println(err)
		return
	}
	topIt.Label = name
	topIt.Value = nuop
	UserOpContentUI(topIt)

}

var UserOpContentItem = &Item{Label: "User Operation Content", Details: "Manage user operation content"}
var ExportUserOpItem = &Item{Label: "Export User Operation", Details: "Select the export format"}
var GetHashItem = &Item{Label: "Hashes and signatures", Details: "Get the hash of the user operation with entrypoint and chainid"}
var SendAsBundleItem = &Item{Label: "Send as Bundle", Details: "Send the user operation as a bundle"}

func UserOpUI(usop *userop.UserOperation) {
	items := []*Item{
		UserOpContentItem,
		ExportUserOpItem,
		GetHashItem,
		SendAsBundleItem,
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
		case SendAsBundleItem.Label:
			h, err := SendAsBundleUI(usop)
			if err != nil {
				fmt.Println(err)
			} else if h != nil {
				fmt.Println("Transaction sent with hash:", h.Hex())
			} else {
				fmt.Println("Transaction not sent")
			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

var SenderItem = &Item{Label: state.Sender, Details: "Set sender"}
var NonceItem = &Item{Label: "Nonce", Details: "Set nonce", Value: uint64(0)}
var CallDataItem = &Item{Label: "Call Data", Details: "Set Call Data"}
var CallGasLimitItem = &Item{Label: "Call Gas Limit", Details: "Set Call Gas Limit"}
var VerificationGasLimitItem = &Item{Label: "Verification Gas Limit", Details: "Set Verification Gas Limit"}
var PreVerificationGasItem = &Item{Label: "Pre Verification Gas", Details: "Set Pre Verification Gas"}
var MaxFeePerGasItem = &Item{Label: "Max Fee Per Gas", Details: "Set Max Fee Per Gas"}
var MaxPriorityFeePerGasItem = &Item{Label: "Max Priority Fee Per Gas", Details: "Set Max Priority Fee Per Gas"}
var SignItem = &Item{Label: "Sign", Details: "Sign the user operation"}
var FactoryItem = &Item{Label: "Factory", Details: "Set Factory"}
var FactoryDataItem = &Item{Label: "Factory Data", Details: "Set Factory Data"}

// PaymasterItem is defined in paymasterui.go
var PaymasterDataItem = &Item{Label: "Paymaster Data", Details: "Set Paymaster Data"}
var PaymasterVerificationGasLimitItem = &Item{Label: "Paymaster Verification Gas Limit", Details: "Set Paymaster Verification Gas Limit", Value: userop.DefaultPaymasterVerificationGasLimit}
var PaymasterPostOpGasLimitItem = &Item{Label: "Paymaster Post Op Gas Limit", Details: "Set Paymaster Post Op Gas Limit", Value: userop.DefaultPaymasterPostOpGasLimit}
var SignatureItem = &Item{Label: "Signature", Details: "Set Signature"}

func UserOpContentUI(topIt *Item) {
	var usop *userop.UserOperation

	ok := false
	usop, ok = topIt.Value.(*userop.UserOperation)
	if !ok || usop == nil {
		fmt.Println("Invalid UserOp passed to UserOpContentUI")
		return
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
		Size:      22,
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
		case NonceItem.Label:
			InputNonceUI(NonceItem, SenderItem)
		case CallGasLimitItem.Label, VerificationGasLimitItem.Label,
			PreVerificationGasItem.Label, MaxFeePerGasItem.Label, MaxPriorityFeePerGasItem.Label,
			PaymasterVerificationGasLimitItem.Label, PaymasterPostOpGasLimitItem.Label:
			it, _ := GetItem(sel, items)
			InputUint(it, 64)
			copyValuesToUserOp(usop)
		case CallDataItem.Label, FactoryDataItem.Label:
			it, _ := GetItem(sel, items)
			caldat, err := PotentiallyRecursiveCallDataUI()
			if err != nil {
				fmt.Println(err)
			} else {
				it.Value = caldat
				usop.CallData = caldat

			}
		//InputBytes(it)
		case PaymasterDataItem.Label:
			it, _ := GetItem(sel, items)
			SetPaymasterDataUI(it, usop)
			if it.Value != nil {
				usop.PaymasterData = it.Value.([]byte)
			}

			//copyValuesToUserOp(usop)
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
				usop.Sender = SenderItem.Value.(*common.Address)

			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func copyFromUseropToItems(uop *userop.UserOperation) {
	SenderItem.Value = uop.Sender
	NonceItem.Value = uop.Nonce
	CallDataItem.Value = uop.CallData
	FactoryItem.Value = uop.Factory
	FactoryDataItem.Value = uop.FactoryData
	CallGasLimitItem.Value = uop.CallGasLimit
	VerificationGasLimitItem.Value = uop.VerificationGasLimit
	PreVerificationGasItem.Value = uop.PreVerificationGas
	MaxFeePerGasItem.Value = uop.MaxFeePerGas
	MaxPriorityFeePerGasItem.Value = uop.MaxPriorityFeePerGas
	PaymasterItem.Value = uop.Paymaster
	PaymasterDataItem.Value = uop.PaymasterData
	PaymasterVerificationGasLimitItem.Value = uop.PaymasterVerificationGasLimit
	PaymasterPostOpGasLimitItem.Value = uop.PaymasterPostOpGasLimit
	SignatureItem.Value = uop.Signature
}

func copyValuesToUserOp(uop *userop.UserOperation) {
	defer state.Save()
	if SenderItem.Value != nil {
		uop.Sender = SenderItem.Value.(*common.Address)
	}
	uop.Nonce = NonceItem.Value.(uint64)
	if CallDataItem.Value != nil {
		uop.CallData = CallDataItem.Value.([]byte)
	}
	if FactoryItem.Value != nil {
		uop.Factory = FactoryItem.Value.(*common.Address)
	}
	if FactoryDataItem.Value != nil {
		uop.FactoryData = FactoryDataItem.Value.([]byte)
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

var ExportAsUOPJSONItem = &Item{Label: "Export as JSON", Details: "Export as JSON"}
var ExportAsRemixTupleV7Item = &Item{Label: "Export as Remix Tuple (V7)", Details: "Export as Remix Tuple"}
var ExportAsRemixTupleV6Item = &Item{Label: "Export as Remix Tuple (V6)", Details: "Export as Remix Tuple"}
var ExportAsCurlToAlchemyItem = &Item{Label: "Export as Curl to Alchemy's Bundler", Details: "Export as Curl to Alchemy's Bundler"}
var ExportAsCurlToEntryItem = &Item{Label: "Export as Curl to Entrypoint", Details: "Export as Curl to Entrypoint"}

func ExportUserOpUI(uop *userop.UserOperation) {
	items := []*Item{ExportAsUOPJSONItem, ExportAsRemixTupleV6Item, ExportAsRemixTupleV7Item, ExportAsCurlToAlchemyItem, ExportAsCurlToEntryItem, Back}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Select an option",
		Items:     items,
		Templates: ItemTemplate,
		Size:      len(items),
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
			fmt.Println(uop.Pack().MarshalRemix()) // "Pack" only makes sense for v7+ userops
		case ExportAsCurlToAlchemyItem.Label:
			if ApiKeyItem == nil || ApiKeyItem.Value == "" {
				fmt.Println("Set Alchemy's API Key in Settings first!")
				continue
			}
			ExportAsAlchemy(uop)
		case ExportAsCurlToEntryItem.Label:
			ExportAsCurl(uop)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func ExportAsJSON(uop *userop.UserOperation) {
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

func ExportAsCurl(uop *userop.UserOperation) {
	fmt.Println("Not implemented yet")
}

func ExportAsAlchemy(uop *userop.UserOperation) {
	apiKey := fmt.Sprintf("%s", ApiKeyItem.Value)
	head := `curl --request POST --url https://eth-sepolia.g.alchemy.com/v2/` + apiKey + ` --header 'accept: application/json' --data ' { "id": 1, "jsonrpc": "2.0", "method": "eth_sendUserOperation", "params":`
	fmt.Printf("%s [ %s, \"%s\" ] }'\n", head, uop.MarshalAlchemy(), EntryPointItem.Value) // newline or the string gets cut
}
