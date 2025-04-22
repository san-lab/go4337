package rpcui

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
)

var SendEndpointItem = &common.Item{Label: "Blockchain RPC Endpoint"}

// Maintain a map ChainID -> RPC Endpoint
func RPCEndpointsUI(it *common.Item) {
	deleting := false
	AddRPCEnpointItem := &common.Item{Label: "Add RPC Endpoint", Details: "Add a new RPC Endpoint"}
	RemoveRPCEnpointItem := &common.Item{Label: "Remove RPC Endpoint", Details: "Select an RPC Endpoint to remove"}

	for {
		items := []*common.Item{}
		for _, rpc := range state.GetRPCEndpoints() {
			items = append(items, &common.Item{Label: fmt.Sprintf("%s/%v", rpc.Name, rpc.ChainId), Value: rpc.URL})
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
	tItem := &common.Item{Label: "RPC Endpoint Name", Details: "Enter a name for the RPC Endpoint"}
	err := common.InputNewStringUI(tItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	name, ok := tItem.Value.(string)
	if !ok {
		fmt.Println("Invalid value")
		return
	}
	tItem.Label = "RPC Endpoint URL"
	tItem.Details = "Enter the URL for the RPC Endpoint"
	tItem.Value = ""
	err = common.InputNewStringUI(tItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	url, ok := tItem.Value.(string)
	if !ok {
		fmt.Println("Invalid value")
		return
	}
	tItem.Label = "Chain ID"
	tItem.Details = "Enter the Chain ID for the RPC Endpoint"
	tItem.Value = nil
	err = common.InputBigInt(tItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	chainId, ok := tItem.Value.(*big.Int)

	state.AddRPCEndpoint(name, url, chainId)

}

// TODO Error handling
func EditRPCEnpointUI(name string) *state.RPCEndpoint {
	return state.GetRPCEndpoint(name)
}
