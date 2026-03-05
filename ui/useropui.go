package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/abiutil"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/abicalldata"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/nonceui"
	"github.com/san-lab/go4337/ui/setauth"
	"github.com/san-lab/go4337/userop"
)

var SelectUserOpItem = &Item[struct{}]{Label: "Select User Operation", Details: "Select a user operation"}
var CreateUserOpItem = &Item[struct{}]{Label: "Create User Operation", Details: "Create a new user operation"}
var CloneUserOpItem = &Item[struct{}]{Label: "Clone User Operation", Details: "Clone a user operation"}
var DeleteUserOpItem = &Item[struct{}]{Label: "Delete User Operation", Details: "Delete a user operation"}

func TopUserOpUI(usOpItem *Item[*userop.UserOperation]) {
	if usOpItem == nil {
		usOpItem = &Item[*userop.UserOperation]{}
	}
	workItem := &Item[struct{}]{}

	for {
		items := []MenuItem{}
		if usOpItem.Value != nil {
			usop := usOpItem.Value
			workItem = &Item[struct{}]{Label: "Work on User Operation: " + usOpItem.Label, Details: fmt.Sprintf("%s/%v", usop.Sender, usop.Nonce)}
			items = append(items, workItem)
		}
		items = append(items,
			SelectUserOpItem,
			CreateUserOpItem,
			CloneUserOpItem,
			DeleteUserOpItem,
		)
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
			usop := usOpItem.Value
			UserOpUI(usop)
		case DeleteUserOpItem.Label:
			DeleteUserOpUI(usOpItem)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func DeleteUserOpUI(topIt *Item[*userop.UserOperation]) {
	it := &Item[*userop.UserOperation]{}
	SelectUserOpUI(it)
	if it.Value == nil {
		return
	}
	usop := it.Value
	for k, v := range state.GetUserOps() {
		if v == usop {
			if YesNoPromptUI(fmt.Sprintf("Delete UserOp %s?", k)) {

				delete(state.GetUserOps(), k)
				state.Save()
				topIt.Value = nil
				return
			} else {
				fmt.Println("UserOp not deleted")
				return

			}
		}
		fmt.Println("UserOp not found in state")
	}

}

func CloneUserOpUI(topIt *Item[*userop.UserOperation]) {
	it := &Item[*userop.UserOperation]{}
	SelectUserOpUI(it)
	if it.Value == nil {
		return
	}
	usop := it.Value
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
	nusop.Nonce.Increment()
	topIt.Value = nusop
	topIt.Label = newname
	state.AddUserOp(newname, nusop)

}

func SelectUserOpUI(topit *Item[*userop.UserOperation]) {
	items := []MenuItem{}
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
		items = append(items, &Item[string]{Label: name, Details: fmt.Sprintf("Sender: %s, Nonce: %s", uop.Sender, uop.Nonce)})
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

func CreateUserOpUI(topIt *Item[*userop.UserOperation]) {
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

var UserOpContentItem = &Item[struct{}]{Label: "User Operation Content", Details: "Manage user operation content"}
var ExportUserOpItem = &Item[struct{}]{Label: "Export User Operation", Details: "Select the export format"}
var GetHashItem = &Item[struct{}]{Label: "Hashes and signatures", Details: "Get the hash of the user operation with entrypoint and chainid"}
var SendAsBundleItem = &Item[struct{}]{Label: "Send as Bundle", Details: "Send the user operation as a bundle"}

func UserOpUI(usop *userop.UserOperation) {
	var uopname string
	for k, v := range state.GetUserOps() {
		if v == usop {
			uopname = k
			break
		}
	}
	if len(uopname) == 0 {
		fmt.Println("UserOp not found in state")
		return
	}

	items := []MenuItem{
		UserOpContentItem,
		ExportUserOpItem,
		GetHashItem,
		SendAsBundleItem,
		ApiCallsItem,
		Back,
	}
	// Create a new select prompt
	prompt := promptui.Select{
		Label:     "Working with User Operation: " + uopname,
		Items:     items,
		Templates: ItemTemplate,
		Size:      10,
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
			it := &Item[*userop.UserOperation]{Value: usop}
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
		case ApiCallsItem.Label:
			ApiCallsUI(usop)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

var SenderItem = &Item[*common.Address]{Label: state.Sender, Details: "Set sender"}

var CallDataItem = &Item[[]byte]{Label: "Call Data", Details: "Set Call Data"}
var CallGasLimitItem = &Item[*big.Int]{Label: "Call Gas Limit\t ", Details: "Set Call Gas Limit"}
var VerificationGasLimitItem = &Item[*big.Int]{Label: "Verification Gas Limit", Details: "Set Verification Gas Limit"}
var PreVerificationGasItem = &Item[*big.Int]{Label: "Pre Verification Gas   ", Details: "Set Pre Verification Gas"}
var MaxFeePerGasItem = &Item[*big.Int]{Label: "Max Fee Per Gas\t ", Details: "Set Max Fee Per Gas"}
var MaxPriorityFeePerGasItem = &Item[*big.Int]{Label: "Max Priority Fee Per Gas", Details: "Set Max Priority Fee Per Gas"}

var FactoryItem = &Item[*common.Address]{Label: "Factory", Details: "Set Factory"}
var FactoryDataItem = &Item[[]byte]{Label: "Factory Data", Details: "Set Factory Data"}

// PaymasterItem is defined in paymasterui.go
var PaymasterDataItem = &Item[[]byte]{Label: "Paymaster Data", Details: "Set Paymaster Data"}
var PaymasterVerificationGasLimitItem = &Item[*big.Int]{Label: "Paymaster Verif. Gas Limit", Details: "Set Paymaster Verification Gas Limit", Value: userop.DefaultPaymasterVerificationGasLimit}
var PaymasterPostOpGasLimitItem = &Item[*big.Int]{Label: "Paymaster Post Op Gas Limit", Details: "Set Paymaster Post Op Gas Limit", Value: userop.DefaultPaymasterPostOpGasLimit}
var SignatureItem = &Item[[]byte]{Label: "Signature", Details: "Set Signature"}
var EIP7702AuthItem = &Item[*types.SetCodeAuthorization]{Label: "EIP7702Authorization", Display: func(v *types.SetCodeAuthorization) string {
	return setauth.AuthorityString(v)
}}

func UserOpContentUI(topIt *Item[*userop.UserOperation]) {
	var usop *userop.UserOperation

	usop = topIt.Value
	if usop == nil {
		fmt.Println("Invalid UserOp passed to UserOpContentUI")
		return
	}

	items := []MenuItem{
		SenderItem,
		nonceui.NonceItem,
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
		EIP7702AuthItem,
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
		case nonceui.NonceItem.Label:
			nonceui.InputNonceUI(nonceui.NonceItem, SenderItem, false)
		case CallGasLimitItem.Label, VerificationGasLimitItem.Label,
			PreVerificationGasItem.Label, MaxFeePerGasItem.Label, MaxPriorityFeePerGasItem.Label,
			PaymasterVerificationGasLimitItem.Label, PaymasterPostOpGasLimitItem.Label:
			typedIt, ok := GetTypedItem[*big.Int](sel, items)
			if ok {
				InputBigInt(typedIt)
			}
			copyValuesToUserOp(usop)
		case CallDataItem.Label, FactoryDataItem.Label:
			typedIt, ok := GetTypedItem[[]byte](sel, items)
			caldat, err := abicalldata.PotentiallyRecursiveCallDataUI()
			if err != nil {
				fmt.Println(err)
			} else if ok {
				typedIt.Value = caldat
				usop.CallData = caldat

			}
		case PaymasterDataItem.Label:
			SetPaymasterDataUI(PaymasterDataItem, usop)
			if PaymasterDataItem.Value != nil {
				usop.PaymasterData = PaymasterDataItem.Value
			}
		case SignatureItem.Label:
			copyValuesToUserOp(usop)
			ret, err := SetSignatureUI(usop)
			if err == nil {
				SignatureItem.Value = ret
				SignatureItem.Details = hex.EncodeToString(ret)
			} else {
				fmt.Println(err)
			}

		case SenderItem.Label, PaymasterItem.Label, FactoryItem.Label:
			_, addr, ok := AddressFromBookUI(sel)
			if ok {
				switch sel {
				case SenderItem.Label:
					SenderItem.Value = addr
					usop.Sender = addr
				case PaymasterItem.Label:
					PaymasterItem.Value = addr
					usop.Paymaster = addr
				case FactoryItem.Label:
					FactoryItem.Value = addr
					usop.Factory = addr
				}
			}
		case EIP7702AuthItem.Label:
			usop.EIP7702Auth = setauth.AuthUI(usop.EIP7702Auth)
			EIP7702AuthItem.Value = usop.EIP7702Auth

		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func copyFromUseropToItems(uop *userop.UserOperation) {
	SenderItem.Value = uop.Sender
	nonceui.NonceItem.Value = uop.Nonce
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
	EIP7702AuthItem.Value = uop.EIP7702Auth
}

func copyValuesToUserOp(uop *userop.UserOperation) {
	defer state.Save()
	if SenderItem.Value != nil {
		uop.Sender = SenderItem.Value
	}
	uop.Nonce = nonceui.NonceItem.Value
	if CallDataItem.Value != nil {
		uop.CallData = CallDataItem.Value
	}
	if FactoryItem.Value != nil {
		uop.Factory = FactoryItem.Value
	}
	if FactoryDataItem.Value != nil {
		uop.FactoryData = FactoryDataItem.Value
	}
	uop.CallGasLimit = CallGasLimitItem.Value
	uop.VerificationGasLimit = VerificationGasLimitItem.Value
	uop.PreVerificationGas = PreVerificationGasItem.Value
	uop.MaxFeePerGas = MaxFeePerGasItem.Value
	uop.MaxPriorityFeePerGas = MaxPriorityFeePerGasItem.Value
	if PaymasterItem.Value != nil {
		uop.Paymaster = PaymasterItem.Value
	}
	if PaymasterDataItem.Value != nil {
		uop.PaymasterData = PaymasterDataItem.Value
	}
	uop.PaymasterVerificationGasLimit = PaymasterVerificationGasLimitItem.Value
	uop.PaymasterPostOpGasLimit = PaymasterPostOpGasLimitItem.Value
	if SignatureItem.Value != nil {
		uop.Signature = SignatureItem.Value
	}
	if EIP7702AuthItem.Value != nil {
		uop.EIP7702Auth = EIP7702AuthItem.Value
	}
}

var ExportAsUOPJSONItem = &Item[struct{}]{Label: "Export as JSON", Details: "Export as JSON"}
var ExportAsRemixTupleV7Item = &Item[struct{}]{Label: "Export as Remix Tuple (V7)", Details: "Export as Remix Tuple"}
var ExportAsRemixTupleV6Item = &Item[struct{}]{Label: "Export as Remix Tuple (V6)", Details: "Export as Remix Tuple"}
var ExportAsCurlToAlchemyItem = &Item[struct{}]{Label: "Export as Curl to Alchemy's Bundler", Details: "Export as Curl to Alchemy's Bundler"}
var ExportAsCurlToEntryItem = &Item[struct{}]{Label: "Export as Curl to Entrypoint", Details: "Export as Curl to Entrypoint"}

func ExportUserOpUI(uop *userop.UserOperation) {
	items := []MenuItem{ExportAsUOPJSONItem, ExportAsRemixTupleV6Item, ExportAsRemixTupleV7Item, ExportAsCurlToAlchemyItem, ExportAsCurlToEntryItem, Back}
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

func ExportAsCurl(uop *userop.UserOperation) {
	u6 := uop.MarshalV6UserOp()
	bt, _ := json.MarshalIndent(u6, "", "  ")
	fmt.Println(string(bt))
}

func ExportAsAlchemy(uop *userop.UserOperation) {
	apiKey := ApiKeyItem.Value
	head := `curl --request POST --url https://eth-sepolia.g.alchemy.com/v2/` + apiKey + ` --header 'accept: application/json' --data ' { "id": 1, "jsonrpc": "2.0", "method": "eth_sendUserOperation", "params":`
	fmt.Printf("%s [ %s, \"%s\" ] }'\n", head, uop.MarshalAlchemy(), EntryPointItem.Value) // newline or the string gets cut
}
