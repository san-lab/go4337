package setauth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/ui/rpcui"
)

func SendSetAuthTxUI(satx *types.SetCodeTx) {
	if satx == nil {
		fmt.Println("SetAuthTx is nil")
		return
	}

	rpcOk := false
	var rpc *state.RPCEndpoint
	sendItem := &common.Item{Label: "Send SetAuthTx", Details: "Send SetAuthTx"}

	pr := promptui.Select{Label: "Sending SetAuthTx", Templates: common.ItemTemplate, Size: 10}
	//fmt.Println("Before the loop")
	for {
		rpc, rpcOk = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
		fmt.Println("RPC:", rpcOk, rpc)
		items := []*common.Item{}
		items = append(items, rpcui.SendEndpointItem)
		if rpcOk {
			sendItem.Label = fmt.Sprintf("Send SetAuthTx to %s", rpc)
			items = append(items, sendItem)
		}
		items = append(items, common.Back)
		pr.Items = items
		fmt.Println("Runnning prompt")
		_, sel, _ := pr.Run()

		switch sel {
		case common.Back.Label:
			return
		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
			fmt.Println("RPC:", rpcui.SendEndpointItem.Value)
		case sendItem.Label:

			envelope := types.NewTx(satx)

			//sbytes, _ = json.MarshalIndent(envelope, "", "  ")
			//fmt.Println("Sending SetAuthTx to", string(sbytes))

			h, err := rpccalls.SendTransaction(rpc, envelope)
			if err != nil {
				fmt.Println("Error sending SetAuthTx:", err)
				return
			}
			fmt.Println("SetAuthTx sent successfully, tx hash:", h.Hex())
			return
		}

	}

}
