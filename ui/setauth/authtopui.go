package setauth

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
	. "github.com/san-lab/go4337/ui/common"
)

// select menu with options: select/create//clone/export/import setAuthTx
var SelectSetAuthTxItem = &Item[struct{}]{Label: "Select SetAuthTx", Details: "Select SetAuthTx"}
var CreateSetAuthTxItem = &Item[struct{}]{Label: "Create SetAuthTx", Details: "Create SetAuthTx"}
var DeleteSetAuthTxItem = &Item[struct{}]{Label: "Delete SetAuthTx", Details: "Delete SetAuthTx"}
var CloneSetAuthTxItem = &Item[struct{}]{Label: "Clone SetAuthTx", Details: "Clone SetAuthTx"}
var ExportSetAuthTxItem = &Item[struct{}]{Label: "Export SetAuthTx", Details: "Export SetAuthTx"}
var ImportSetAuthTxItem = &Item[struct{}]{Label: "Import SetAuthTx", Details: "Import SetAuthTx"}
var ContentSetAuthTxItem = &Item[struct{}]{Label: "Content SetAuthTx", Details: "Content SetAuthTx"}
var SendSetAuthTxItem = &Item[struct{}]{Label: "Send SetAuthTx", Details: "Send SetAuthTx"}

var satxName string

func AuthTxUI() {

	// Create a new select prompt
	prompt := promptui.Select{
		Templates: ItemTemplate,
		Size:      10,
	}
	for {
		var selectedSetAuthTx *types.SetCodeTx
		var ok bool
		items := []MenuItem{}
		if selectedSetAuthTx, ok = state.GetSetAuthTx(satxName); ok {
			items = append(items, ContentSetAuthTxItem, DeleteSetAuthTxItem, CloneSetAuthTxItem, ExportSetAuthTxItem, SendSetAuthTxItem)
			prompt.Label = fmt.Sprintf("Working with tx To: %s at Nonce: %d", selectedSetAuthTx.To.Hex(), selectedSetAuthTx.Nonce)
		} else {

		}
		items = append(items,
			SelectSetAuthTxItem,
			CreateSetAuthTxItem,
			ImportSetAuthTxItem,
			Back,
		)
		prompt.Items = items

		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case SelectSetAuthTxItem.Label:
			SelectSetAuthTxUI()
		case CreateSetAuthTxItem.Label:
			CreateSetAuthTxUI()
		case DeleteSetAuthTxItem.Label:
			DeleteSetAuthTxUI()
		case CloneSetAuthTxItem.Label:
			CloneSetAuthTxUI()
		case ExportSetAuthTxItem.Label:
			ExportSetAuthTxUI(selectedSetAuthTx)
		case ImportSetAuthTxItem.Label:
			ImportSetAuthTxUI()
		case ContentSetAuthTxItem.Label:
			SetAuthTxUI(selectedSetAuthTx)
		case SendSetAuthTxItem.Label:
			SendSetAuthTxUI(selectedSetAuthTx)

		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func SelectSetAuthTxUI() {
	items := []MenuItem{}
	satxNames := state.ListSetAuthTxs()
	for _, sn := range satxNames {
		items = append(items, &Item[string]{
			Label: sn,
			Value: sn,
		})
	}
	items = append(items, Back)
	if len(satxNames) == 0 {
		fmt.Println("No SetAuthTxs found")
		return
	}
	prompt := promptui.Select{
		Label:     "Select SetAuthTx",
		Items:     items,
		Templates: ItemTemplate,
	}
	_, sel, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	if sel == Back.Label {
		return
	}
	_, ok := state.GetSetAuthTx(sel)
	if ok {
		satxName = sel
	}

}

func CreateSetAuthTxUI() {
	satx := &types.SetCodeTx{}
	SetAuthTxUI(satx)
	name := fmt.Sprintf("satx%03d", len(state.ListSetAuthTxs()))
	state.AddSetAuthTx(name, satx)
	satxName = name
}

func DeleteSetAuthTxUI() {
	state.RemoveSetAuthTx(satxName)
	state.Save()
	satxName = ""
}

func CloneSetAuthTxUI() {
	// CloneSetAuthTxUI clones a SetAuthTx
}

func ExportSetAuthTxUI(satx *types.SetCodeTx) {
	//RLP encode the SetAuthTx
	tx := types.NewTx(satx)
	bt, err := tx.MarshalBinary()
	if err != nil {
		fmt.Println("could not marshal SetAuthTx:", err)
		return
	}
	fmt.Printf("\n%x\n", bt)

	// ExportSetAuthTxUI exports a SetAuthTx

}

func ImportSetAuthTxUI() {
	it := &Item[[]byte]{Label: "Input Hex", Details: "Input Hex Data"}
	err := InputBytes(it, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	bts := it.Value
	if len(bts) == 0 || bts[0] != 0x04 {
		fmt.Println("Invalid SetAuthTx")
		return
	}
	satx := &types.SetCodeTx{}

	err = rlp.DecodeBytes(bts[1:], satx)
	if err != nil {
		fmt.Println(err)
		return
	}
	jb, err := json.MarshalIndent(satx, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jb))

	nameit := &Item[string]{Label: "Set Name", Details: "Set Name for the imported SetAuth Transaction"}
	err = InputNewStringUI(nameit)
	if err != nil {
		fmt.Println(err)
		return
	}
	state.AddSetAuthTx(nameit.Value, satx)
	state.Save()
	satxName = nameit.Value

}
