package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/abicalldata"
	. "github.com/san-lab/go4337/ui/common"
	rpcui "github.com/san-lab/go4337/ui/rpcui"
	"github.com/san-lab/go4337/ui/signui"
	"github.com/san-lab/go4337/zksyncera"
)

var ZkSelectTxItem = &Item{Label: "Select zkSync 712 Transaction"}
var ZkWorkTxItem = &Item{Label: "Work on zkSync 712 Transaction"}
var ZkRemoveTxItem = &Item{Label: "Remove zkSync 712 Transaction"}
var ZkAddTxItem = &Item{Label: "Add zkSync 712 Transaction"}
var ZkDecodeItem = &Item{Label: "Decode zk Sync 712 Transaction from RLP"}

func ZkSyncEraUI() {

	for {
		items := []*Item{}
		if ZkWorkTxItem.Value != nil {
			items = append(items, ZkWorkTxItem)
		}
		items = append(items, ZkSelectTxItem, ZkAddTxItem, ZkDecodeItem, ZkRemoveTxItem, Back)
		spr := promptui.Select{Label: "ZkSync Era", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case ZkAddTxItem.Label:
			NewERAUI()
		case ZkWorkTxItem.Label:
			WorkWithERAUI(ZkWorkTxItem.Value.(*zksyncera.ZkSyncTxRLP))
		case ZkRemoveTxItem.Label:
			RemoveERAUI()
		case ZkDecodeItem.Label:
			ERAFromRLPUI()
		case ZkSelectTxItem.Label:
			SelectERAUI()
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}
}

func SelectERAUI() {
	items := []*Item{}
	for _, name := range state.ListZkERA712s() {
		items = append(items, &Item{Label: name})
	}
	items = append(items, Back)
	for {
		spr := promptui.Select{Label: "Select ZkSync Era", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		default:
			if tx, ok := state.GetZkERA712(sel); ok {
				ZkWorkTxItem.Value = tx
				ZkWorkTxItem.DisplayValueString = sel
				return
			}
		}
	}
}

func ERAFromRLPUI() {
	nit := &Item{Label: "Name of the new ZkSync Era"}
	err := InputNewStringUI(nit)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := nit.Value.(string)
	if _, ok := state.GetZkERA712(name); ok || name == "" {
		fmt.Println("ZkSync Era with name", name, "already exists")
		return
	}
	enctx, _ := hex.DecodeString(zksyncera.EncodedTransaction1[2:])
	rit := &Item{Label: "RLP bytes in hex", Value: enctx}
	err = InputBytes(rit, -1)
	if err != nil {
		fmt.Println(err)
	} else if rlpb, ok := rit.Value.([]byte); ok && len(rlpb) > 0 {
		zktx := new(zksyncera.ZkSyncTxRLP)
		err := zktx.Decode(rlpb)
		if err != nil {
			fmt.Println(err)
			return
		}
		state.AddZkERA712(name, zktx)
		ZkWorkTxItem.Value = zktx
		ZkWorkTxItem.DisplayValueString = name
	} else {
		fmt.Println("Invalid RLP bytes")
	}
}

func RemoveERAUI() {
	items := []*Item{}
	for _, name := range state.ListZkERA712s() {
		items = append(items, &Item{Label: name})
	}
	items = append(items, Back)
	for {
		spr := promptui.Select{Label: "Select ZkSync Era to remove", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		default:
			state.RemoveZkERA712(sel)
		}
	}
}

func NewERAUI() {
	nit := &Item{Label: "Name of the new ZkSync Era"}
	err := InputNewStringUI(nit)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := nit.Value.(string)
	if _, ok := state.GetZkERA712(name); ok || name == "" {
		fmt.Println("ZkSync Era with name", name, "already exists")
		return
	}
	zktx := new(zksyncera.ZkSyncTxRLP)
	state.AddZkERA712(name, zktx)
	ZkWorkTxItem.Value = zktx
}

var ZkEIP712TxItem = &Item{Label: "ZkTx Content"}
var PrintItem = &Item{Label: "Print"}
var EncodeItem = &Item{Label: "RLP Encode"}

var CurlItem = &Item{Label: "CURL call"}
var SignEraItem = &Item{Label: "Sign"}

func WorkWithERAUI(lera *zksyncera.ZkSyncTxRLP) {
	defer state.Save()
	items := []*Item{
		ZkEIP712TxItem,
		PrintItem,
		EncodeItem,
		CurlItem,
		EraCallItem,
		Back,
	}
	for {
		spr := promptui.Select{Label: "ZkSync Era", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case PrintItem.Label:
			bt, err := json.MarshalIndent(lera, "", "  ")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(bt))
			}
		case ZkEIP712TxItem.Label:
			ZkEIP712TxUI(lera)
		case EncodeItem.Label:
			bt, err := lera.Encode()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("0x%x\n", bt)
			}
		case CurlItem.Label:
			enctx, err := lera.Encode()
			if err != nil {
				fmt.Println(err)
				continue
			}
			url := "http://localhost:3000"
			call := fmt.Sprintf(`curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0x%x"],"id":1}' %s`, enctx, url)
			fmt.Println(call)
		case EraCallItem.Label:
			ERACallUI(lera)
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

var ZkChainIDItem = &Item{Label: "Chain ID"}
var ZkNonceItem = &Item{Label: "Nonce"}
var ZkFromItem = &Item{Label: "From"}
var ZkToItem = &Item{Label: "To"}
var ZkValueItem = &Item{Label: "Value"}
var ZkDataItem = &Item{Label: "Data"}
var ZkGasLimitItem = &Item{Label: "Gas Limit"}
var ZkGasTipCapItem = &Item{Label: "Gas Tip Cap"}
var ZkGasFeeCapItem = &Item{Label: "Gas Fee Cap"}
var ZkGasPerPubdataItem = &Item{Label: "Gas Per Pubdata"}
var ZkCustomSignatureItem = &Item{Label: "Custom Signature"}
var ZkFactoryDepsItem = &Item{Label: "Factory Dependencies"}
var ZkPaymasterParamsItem = &Item{Label: "Paymaster Parameters"}

func ZkEIP712TxUI(zktx *zksyncera.ZkSyncTxRLP) {
	ZkChainIDItem.Value = zktx.ChainId2
	ZkNonceItem.Value = zktx.Nonce
	ZkFromItem.Value = zktx.From
	ZkToItem.Value = zktx.To
	ZkDataItem.Value = zktx.Data
	ZkGasLimitItem.Value = zktx.GasLimit
	ZkGasTipCapItem.Value = zktx.MaxPriorityFeePerGas
	ZkGasFeeCapItem.Value = zktx.MaxFeePerGas
	ZkGasPerPubdataItem.Value = zktx.GasPerPubdata
	ZkFactoryDepsItem.Value = zktx.FactoryDeps
	ZkPaymasterParamsItem.Value = zktx.PaymasterParams
	ZkCustomSignatureItem.Value = zktx.CustomSignature
	items := []*Item{
		ZkChainIDItem,
		ZkNonceItem,
		ZkFromItem,
		ZkToItem,
		ZkValueItem,
		ZkDataItem,
		ZkGasLimitItem,
		ZkGasTipCapItem,
		ZkGasFeeCapItem,
		ZkGasPerPubdataItem,
		ZkFactoryDepsItem,
		ZkPaymasterParamsItem,
		ZkCustomSignatureItem,
		Back,
	}
	for {
		ZkFactoryDepsItem.DisplayValueString = FactoryDeptsDetails(zktx.FactoryDeps)
		spr := promptui.Select{Label: "Zk EIP712 Transaction", Items: items, Templates: ItemTemplate, Size: 18}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case ZkChainIDItem.Label:
			err := InputBigInt(ZkChainIDItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.ChainId2 = ZkChainIDItem.Value.(*big.Int)
				zktx.ChainId1 = ZkChainIDItem.Value.(*big.Int)
			}
		case ZkNonceItem.Label:
			n, err := InputUint(ZkNonceItem, 64)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.Nonce = n
			}
		case ZkFromItem.Label:
			_, addr, ok := AddressFromBookUI("zkFrom")
			if ok {
				zktx.From = addr
				ZkFromItem.Value = addr
			}
		case ZkToItem.Label:
			_, addr, ok := AddressFromBookUI("zkTo")
			if ok {
				zktx.To = addr
				ZkToItem.Value = addr
			}
		case ZkValueItem.Label:
			err = InputBigInt(ZkValueItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.Value = ZkValueItem.Value.(*big.Int)
			}
		case ZkDataItem.Label:
			calldat, err := abicalldata.PotentiallyRecursiveCallDataUI()
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.Data = calldat
				ZkDataItem.Value = calldat
			}
		case ZkGasLimitItem.Label:
			err := InputBigInt(ZkGasLimitItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.GasLimit = ZkGasLimitItem.Value.(*big.Int)
			}

		case ZkGasTipCapItem.Label:
			err := InputBigInt(ZkGasTipCapItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.MaxPriorityFeePerGas = ZkGasTipCapItem.Value.(*big.Int)
			}

		case ZkGasFeeCapItem.Label:
			err := InputBigInt(ZkGasFeeCapItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.MaxFeePerGas = ZkGasFeeCapItem.Value.(*big.Int)
			}
		case ZkGasPerPubdataItem.Label:
			err := InputBigInt(ZkGasPerPubdataItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.GasPerPubdata = ZkGasPerPubdataItem.Value.(*big.Int)
			}
		case ZkCustomSignatureItem.Label:
			SignEraUI(zktx)
			ZkCustomSignatureItem.Value = zktx.CustomSignature
		case ZkFactoryDepsItem.Label:
			ZkFactoryDepsUI(zktx)
		case ZkPaymasterParamsItem.Label:
			ZkPaymasterParamsUI()
		case PrintItem.Label:
			bt, err := json.MarshalIndent(zktx, "", "  ")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(bt))
			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func ZkFactoryDepsUI(zktx *zksyncera.ZkSyncTxRLP) {
	AppendItem := &Item{Label: "Append Factory Dependency"}
	RemoveItem := &Item{Label: "Remove Factory Dependency"}
	removing := false
	for {

		items := []*Item{}
		for i, dep := range zktx.FactoryDeps {
			depItem := &Item{Label: fmt.Sprintf("Dep#%v", i)}
			depItem.Value = dep
			depItem.DisplayValueString = FactoryDeptsDetails([]hexutil.Bytes{dep})
			depItem.Details = ShortString(fmt.Sprintf("0x%x", dep), 255)

			items = append(items, depItem)
		}
		if !removing {
			items = append(items, AppendItem, RemoveItem)
		}
		items = append(items, Back)
		label := "Factory Dependencies"
		if removing {
			label = "Select Factory Dependency to remove"
		}
		spr := promptui.Select{Label: label, Items: items, Templates: ItemTemplate, Size: 10}
		i, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case AppendItem.Label:
			nit := &Item{Label: "New factory dependency"}
			InputXexUI(nit, zktx)

		case RemoveItem.Label:
			removing = true
		default:
			if removing {
				fmt.Println("Removing Factory Dependency", i)
				if i >= 0 && i < len(zktx.FactoryDeps) {
					zktx.FactoryDeps = append(zktx.FactoryDeps[:i], zktx.FactoryDeps[i+1:]...)

				}
				removing = false
			}
		}
	}

}

func InputXexUI(nit *Item, era *zksyncera.ZkSyncTxRLP) {
	HexItem := &Item{Label: "Direct Hex Data"}
	FileItem := &Item{Label: "Hex Data from file"}
	sel := promptui.Select{Label: nit.Label, Items: []*Item{HexItem, FileItem, Back}, Templates: ItemTemplate, Size: 10}

	_, s, err := sel.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch s {
	case Back.Label:
		return
	case HexItem.Label:
		err := InputBytes(nit, -1)
		if err != nil {
			fmt.Println(err)
		} else {
			fd := nit.Value.([]byte)
			era.FactoryDeps = append(era.FactoryDeps, fd[:])
			return
		}
	case FileItem.Label:
		bt, err := InputHexFileUI("Select file with hex data")
		if err != nil {
			fmt.Println(err)
			return
		}
		era.FactoryDeps = append(era.FactoryDeps, bt)
		return
	}
	return
}

func ZkPaymasterParamsUI() {
	fmt.Println("ZkPaymasterParamsUI")
}

var ERASignerItem = &Item{Label: "Signer"}
var ERASignItem = &Item{Label: "Sign"}

func SignEraUI(era *zksyncera.ZkSyncTxRLP) {
	if era.ChainId2 == nil {
		fmt.Println("WARNING: Chain ID not set")

	} else if era.ChainId1.Cmp(era.ChainId2) != 0 {
		fmt.Println("WARNING: Chain ID mismatch, setting ChainId1 to ChainId2")
		era.ChainId1 = era.ChainId2
	}

	items := []*Item{ZkChainIDItem, ERASignerItem, ERASignItem, Back}
	for {
		spr := promptui.Select{Label: "Sign ZkSync Era EIP712 message", Items: items, Templates: ItemTemplate, Size: 10}
		_, sel, err := spr.Run()
		if err != nil {
			fmt.Println("could not run prompt:", err)
			return
		}
		switch sel {
		case Back.Label:
			return
		case ZkChainIDItem.Label:

		case ERASignerItem.Label:
			signui.SignerUI(ERASignerItem)
		case ERASignItem.Label:
			sig, err := era.Sign(ERASignerItem.Value)
			if err != nil {
				fmt.Println(err)
				continue
			}
			sig[64] += 27
			era.CustomSignature = sig
			ERASignItem.Value = sig
			fmt.Println("Signature:", hex.EncodeToString(sig[:]))
		}
	}

}

func FactoryDeptsDetails(deps []hexutil.Bytes) string {
	if len(deps) == 0 {
		return "No Factory Dependencies"
	}
	det := ""
	for i, dep := range deps {
		h, err := zksyncera.ZKBytecodeHash(dep, 1, 0)
		if err != nil {
			det = fmt.Sprintf("Error hashing at pos %v: %w", i, err)
			return det
		}
		det += fmt.Sprintf("0x%x  ", h)

	}
	return det
}

var EraCallItem = &Item{Label: "Call RPC"}

func ERACallUI(era *zksyncera.ZkSyncTxRLP) {
	address := era.To

	if address == nil {
		fmt.Println("No target address")
		return
	}
	enctx, err := era.Encode()
	if err != nil {
		fmt.Println(err)
		return
	}

	rpc, rpcOk := rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
	for {
		items := []*Item{rpcui.SendEndpointItem}
		if rpcOk {
			items = append(items, EraCallItem)

		}
		items = append(items, Back)
		prompt := promptui.Select{
			Label:     "Set the Call parameters",
			Items:     items,
			Templates: ItemTemplate,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch sel {
		case Back.Label:
			return

		case rpcui.SendEndpointItem.Label:
			rpcui.RPCEndpointsUI(rpcui.SendEndpointItem)
			rpc, rpcOk = rpcui.SendEndpointItem.Value.(*state.RPCEndpoint)
		case EraCallItem.Label:
			//ERACall(rpc, address)
			var res string
			err := rpccalls.SendMethodCall(rpc, &res, "eth_sendRawTransaction", fmt.Sprintf("0x%x", enctx))
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Transaction sent:", res)
			}

		}
	}
}
