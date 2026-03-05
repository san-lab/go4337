package rpcui

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
)

var SendEndpointItem = &common.Item[*state.RPCEndpoint]{Label: "Blockchain RPC Endpoint"}

// Maintain a map ChainID -> RPC Endpoint
func RPCEndpointsUI(it *common.Item[*state.RPCEndpoint]) {
	deleting := false
	AddRPCEnpointItem := &common.Item[struct{}]{Label: "Add RPC Endpoint", Details: "Add a new RPC Endpoint"}
	RemoveRPCEnpointItem := &common.Item[struct{}]{Label: "Remove RPC Endpoint", Details: "Select an RPC Endpoint to remove"}

	for {
		items := []common.MenuItem{}
		for _, rpc := range state.GetRPCEndpoints() {
			items = append(items, &common.Item[string]{Label: fmt.Sprintf("%s/%v", rpc.Name, rpc.ChainId), Value: rpc.URL})
		}
		if !deleting {
			items = append(items, AddRPCEnpointItem, RemoveRPCEnpointItem)
		}
		items = append(items, common.Back)
		prompt := promptui.Select{Label: "RPC Endpoints", Items: items, Templates: common.ItemTemplate, Size: 10}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case common.Back.Label:
			return
		case RemoveRPCEnpointItem.Label:
			deleting = true
			continue
		case AddRPCEnpointItem.Label:
			AddRPCEnpointUI()
			return
		default:
			name := strings.Split(sel, "/")[0]
			if deleting {
				label := fmt.Sprintf("Are you sure you want to delete RPC Endpoint %s (yes/no)?", sel)

				if common.YesNoPromptUI(label) {
					state.RemoveRPCEndpoint(name)
				}
				deleting = false
			} else {
				v := EditRPCEnpointUI(name)
				if it != nil {
					it.Value = v
					return
				}
			}

		}
	}

}

func AddRPCEnpointUI() {
	nameItem := &common.Item[string]{Label: "RPC Endpoint Name", Details: "Enter a name for the RPC Endpoint"}
	err := common.InputNewStringUI(nameItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := nameItem.Value
	if name == "" {
		fmt.Println("Invalid value")
		return
	}
	urlItem := &common.Item[string]{Label: "RPC Endpoint URL", Details: "Enter the URL for the RPC Endpoint"}
	err = common.InputNewStringUI(urlItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	url := urlItem.Value
	if url == "" {
		fmt.Println("Invalid value")
		return
	}
	chainIdItem := &common.Item[*big.Int]{Label: "Chain ID", Details: "Enter the Chain ID for the RPC Endpoint"}
	_, err = common.InputBigInt(chainIdItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	chainId := chainIdItem.Value

	state.AddRPCEndpoint(name, url, chainId)

}

// TODO Error handling
func EditRPCEnpointUI(name string) *state.RPCEndpoint {
	return state.GetRPCEndpoint(name)
}
