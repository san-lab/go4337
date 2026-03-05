package ui

import (
	"fmt"
	"os"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/abicalldata"
	. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/rpcui"
	"github.com/san-lab/go4337/ui/setauth"
	"github.com/san-lab/go4337/ui/signui"
)

func init() {
	EntryPointItem.Value = entrypoint.E7Address
	ApiKeyItem.Value = state.GetApiKey("Alchemy")
	//Get rid of the bloody bell
	readline.Stdout = &stderr{}
}

type stderr struct{}

func (s *stderr) Write(b []byte) (int, error) {
	if len(b) == 1 && b[0] == 7 {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (s *stderr) Close() error {
	return os.Stderr.Close()
}

var PaymasterItem = &Item[*ecommon.Address]{Label: "Paymaster", Details: "Manage Paymaster settings"}
var UserOpItem = &Item[struct{}]{Label: "User Operations", Details: "Manage User Operations"}
var AbisItem = &Item[struct{}]{Label: "ABIs", Details: "Manage ABIs"}
var EntryPointItem = &Item[ecommon.Address]{Label: "Entrypoint", Details: "Set Entrypoint"}
var ApiKeysItem = &Item[struct{}]{Label: "API's and API Keys", Details: "Manage API Access"}
var SettingsItem = &Item[struct{}]{Label: "Settings", Details: "Paymasters, Signers, ChainID, ..."}
var ChainCallItem = &Item[struct{}]{Label: "Chain Calls", Details: "Call a function on-chain"}
var DEBUGItem = &Item[bool]{Label: "DEBUG", Details: "Toggle DEBUG", Value: state.DEBUG, Display: func(v bool) string { return fmt.Sprint(v) }}
var AddressBooksItem = &Item[struct{}]{Label: "Address Books", Details: "Manage Address Books"}

var EIP7702Item = &Item[struct{}]{Label: "EIP7702 Stuff", Details: "Engage EIP-7702"}

func RootUI() {
	items := []MenuItem{
		SettingsItem,
		UserOpItem,
		AbisItem,
		ChainCallItem,
		ApiCallsItem,
		EIP7702Item,

		Exit,
	}
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
		case SettingsItem.Label:
			SettingsUI()
		case UserOpItem.Label:
			TopUserOpUI(nil)
		case AbisItem.Label:
			abicalldata.AbisUI(nil)
		case ChainCallItem.Label:
			ChainCallUI()
		case ApiCallsItem.Label:
			ApiCallsUI(nil)
		case EIP7702Item.Label:
			setauth.AuthTxUI()
		case Exit.Label:
			return
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

var GasLimitOffsetItem = &Item[uint64]{Label: "Gas Limit Offset", Details: "Set Gas Limit Offset to top the bundle gas limit", Value: state.GetGasLimitOffset()}

func SettingsUI() {
	items := []MenuItem{
		PaymasterItem,
		SignerItem,
		ChainIDItem,
		EntryPointItem,

		rpcui.SendEndpointItem,
		GasLimitOffsetItem,
		AddressBooksItem,
		DEBUGItem,
		Back,
	}
	prompt := promptui.Select{
		Label:     "Settings",
		Items:     items,
		Templates: ItemTemplate,
		Size:      len(items),
	}
	for {
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case PaymasterItem.Label:
			PaymasterUI()
		case SignerItem.Label:
			signui.SignerUI(SignerItem)
		case ChainIDItem.Label:
			InputBigInt(ChainIDItem)
			state.SetChainId(ChainIDItem.Value)
		case EntryPointItem.Label:
			EntryPointUI()

		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
		case GasLimitOffsetItem.Label:
			InputUint64(GasLimitOffsetItem)
			state.SetGasLimitOffset(GasLimitOffsetItem.Value)
		case DEBUGItem.Label:
			state.DEBUG = !state.DEBUG
			DEBUGItem.Value = state.DEBUG
		case Back.Label:
			return
		case AddressBooksItem.Label:
			AddressBooksUI()
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func AddressBooksUI() {
	items := []MenuItem{}
	for _, ab := range state.GetAddressBooks() {
		items = append(items, &Item[struct{}]{Label: ab})
	}
	items = append(items, Back)
	prompt := promptui.Select{
		Label:     "Select an Address Book",
		Items:     items,
		Templates: ItemTemplate,
		Size:      len(items) + 2,
	}
	_, sel, _ := prompt.Run()
	if sel == Back.Label {
		return
	}
	AddressFromBookUI(sel)

}
