package ui

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/state"
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

func init() {
	readline.Stdout = &stderr{}
}

var PaymasterItem = &Item{Label: "Paymaster", Details: "Manage Paymaster settings"}
var UserOpItem = &Item{Label: "User Operations", Details: "Manage User Operations"}
var AbisItem = &Item{Label: "ABIs", Details: "Manage ABIs"}
var SignerItem = &Item{Label: "Signer", Details: "Manage Signer settings"}
var EntryPointItem = &Item{Label: "Entrypoint", Details: "Set Entrypoint"}
var ApiKeysItem = &Item{Label: "API's and API Keys", Details: "Manage API Access"}
var SettingsItem = &Item{Label: "Settings", Details: "Paymasters, Signers, ChainID, ..."}
var ChainCallItem = &Item{Label: "Chain Calls", Details: "Call a function on-chain"}
var DEBUGItem = &Item{Label: "DEBUG", Details: "Toggle DEBUG", DisplayValueString: fmt.Sprint(state.DEBUG)}

//var RPCEndpointsItem = &Item{Label: "RPC Endpoints", Details: "Manage RPC Endpoints"}

func RootUI() {
	items := []*Item{
		//PaymasterItem,
		SettingsItem,
		UserOpItem,
		//SignerItem,
		AbisItem,
		//ChainIDItem,
		//EntryPointItem,
		ChainCallItem,
		ApiCallsItem,
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
			AbisUI(nil)
		case ChainCallItem.Label:
			ChainCallUI()
		case ApiCallsItem.Label:
			ApiCallsUI(nil)
		case Exit.Label:
			return
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

var GasLimitOffsetItem = &Item{Label: "Gas Limit Offset", Details: "Set Gas Limit Offset to top the bundle gas limit", Value: state.GetGasLimitOffset()}

func SettingsUI() {
	items := []*Item{
		PaymasterItem,
		SignerItem,
		ChainIDItem,
		EntryPointItem,

		SendEndpointItem,
		GasLimitOffsetItem,
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
			SignerUI(SignerItem)
		case ChainIDItem.Label:
			InputUint(ChainIDItem, 64)
			state.SetChainId(ChainIDItem.Value)
		case EntryPointItem.Label:
			EntryPointUI()

		case SendEndpointItem.Label:
			RPCEndpointsUI(SendEndpointItem)
		case GasLimitOffsetItem.Label:
			InputUint(GasLimitOffsetItem, 64)
			state.SetGasLimitOffset(GasLimitOffsetItem.Value.(uint64))
		case DEBUGItem.Label:
			state.DEBUG = !state.DEBUG
			DEBUGItem.DisplayValueString = fmt.Sprint(state.DEBUG)
		case Back.Label:
			return
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func GetValue(label string, items []*Item) (interface{}, bool) {
	it, ok := GetItem(label, items)
	if !ok {
		return nil, false
	}
	return it.Value, true
}

func GetItem(label string, items []*Item) (*Item, bool) {
	for _, i := range items {
		if i.Label == label {
			return i, true
		}
	}
	return nil, false
}
