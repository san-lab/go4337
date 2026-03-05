package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/rpccalls"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui/abicalldata"
	. "github.com/san-lab/go4337/ui/common"
	rpcui "github.com/san-lab/go4337/ui/rpcui"
	"github.com/san-lab/go4337/ui/signui"
	"github.com/san-lab/go4337/zksyncera"
)

var ZkSelectTxItem = &Item[struct{}]{Label: "Select zkSync 712 Transaction"}
var ZkWorkTxItem = &Item[*zksyncera.ZkSyncTxRLP]{Label: "Work on zkSync 712 Transaction"}
var ZkRemoveTxItem = &Item[struct{}]{Label: "Remove zkSync 712 Transaction"}
var ZkAddTxItem = &Item[struct{}]{Label: "Add zkSync 712 Transaction"}
var ZkDecodeItem = &Item[struct{}]{Label: "Decode zk Sync 712 Transaction from RLP"}

func ZkSyncEraUI() {

	for {
		items := []MenuItem{}
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
			WorkWithERAUI(ZkWorkTxItem.Value)
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
	items := []MenuItem{}
	for _, name := range state.ListZkERA712s() {
		nameCopy := name
		items = append(items, &Item[string]{Label: nameCopy})
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
				ZkWorkTxItem.Details = sel
				return
			}
		}
	}
}

func ERAFromRLPUI() {
	nit := &Item[string]{Label: "Name of the new ZkSync Era"}
	err := InputNewStringUI(nit)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := nit.Value
	if _, ok := state.GetZkERA712(name); ok || name == "" {
		fmt.Println("ZkSync Era with name", name, "already exists")
		return
	}
	enctx, _ := hex.DecodeString(zksyncera.EncodedTransaction1[2:])
	rit := &Item[[]byte]{Label: "RLP bytes in hex", Value: enctx}
	err = InputBytes(rit, -1)
	if err != nil {
		fmt.Println(err)
	} else if len(rit.Value) > 0 {
		zktx := new(zksyncera.ZkSyncTxRLP)
		err := zktx.Decode(rit.Value)
		if err != nil {
			fmt.Println(err)
			return
		}
		state.AddZkERA712(name, zktx)
		ZkWorkTxItem.Value = zktx
		ZkWorkTxItem.Details = name
	} else {
		fmt.Println("Invalid RLP bytes")
	}
}

func RemoveERAUI() {
	items := []MenuItem{}
	for _, name := range state.ListZkERA712s() {
		nameCopy := name
		items = append(items, &Item[string]{Label: nameCopy})
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
	nit := &Item[string]{Label: "Name of the new ZkSync Era"}
	err := InputNewStringUI(nit)
	if err != nil {
		fmt.Println(err)
		return
	}
	name := nit.Value
	if _, ok := state.GetZkERA712(name); ok || name == "" {
		fmt.Println("ZkSync Era with name", name, "already exists")
		return
	}
	zktx := new(zksyncera.ZkSyncTxRLP)
	state.AddZkERA712(name, zktx)
	ZkWorkTxItem.Value = zktx
}

var ZkEIP712TxItem = &Item[struct{}]{Label: "ZkTx Content"}
var PrintItem = &Item[struct{}]{Label: "Print"}
var EncodeItem = &Item[struct{}]{Label: "RLP Encode"}

var CurlItem = &Item[struct{}]{Label: "CURL call"}
var SignEraItem = &Item[struct{}]{Label: "Sign"}

func WorkWithERAUI(lera *zksyncera.ZkSyncTxRLP) {
	defer state.Save()
	items := []MenuItem{
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

var ZkChainIDItem = &Item[*big.Int]{Label: "Chain ID"}
var ZkNonceItem = &Item[uint64]{Label: "Nonce"}
var ZkFromItem = &Item[any]{Label: "From"}
var ZkToItem = &Item[any]{Label: "To"}
var ZkValueItem = &Item[*big.Int]{Label: "Value"}
var ZkDataItem = &Item[[]byte]{Label: "Data"}
var ZkGasLimitItem = &Item[*big.Int]{Label: "Gas Limit"}
var ZkGasTipCapItem = &Item[*big.Int]{Label: "Gas Tip Cap"}
var ZkGasFeeCapItem = &Item[*big.Int]{Label: "Gas Fee Cap"}
var ZkGasPerPubdataItem = &Item[*big.Int]{Label: "Gas Per Pubdata"}
var ZkCustomSignatureItem = &Item[[]byte]{Label: "Custom Signature"}
var ZkFactoryDepsItem = &Item[any]{Label: "Factory Dependencies", Display: func(v any) string {
	if deps, ok := v.([]hexutil.Bytes); ok {
		return FactoryDeptsDetails(deps)
	}
	return ""
}}
var ZkPaymasterParamsItem = &Item[any]{Label: "Paymaster Parameters"}

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
	items := []MenuItem{
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
		ZkFactoryDepsItem.Value = zktx.FactoryDeps
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
			_, err := InputBigInt(ZkChainIDItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.ChainId2 = ZkChainIDItem.Value
				zktx.ChainId1 = ZkChainIDItem.Value
			}
		case ZkNonceItem.Label:
			InputUint64(ZkNonceItem)
			zktx.Nonce = ZkNonceItem.Value
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
			_, err = InputBigInt(ZkValueItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.Value = ZkValueItem.Value
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
			_, err := InputBigInt(ZkGasLimitItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.GasLimit = ZkGasLimitItem.Value
			}

		case ZkGasTipCapItem.Label:
			_, err := InputBigInt(ZkGasTipCapItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.MaxPriorityFeePerGas = ZkGasTipCapItem.Value
			}

		case ZkGasFeeCapItem.Label:
			_, err := InputBigInt(ZkGasFeeCapItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.MaxFeePerGas = ZkGasFeeCapItem.Value
			}
		case ZkGasPerPubdataItem.Label:
			_, err := InputBigInt(ZkGasPerPubdataItem)
			if err != nil {
				fmt.Println(err)
			} else {
				zktx.GasPerPubdata = ZkGasPerPubdataItem.Value
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
	AppendItem := &Item[struct{}]{Label: "Append Factory Dependency"}
	RemoveItem := &Item[struct{}]{Label: "Remove Factory Dependency"}
	removing := false
	for {

		items := []MenuItem{}
		for i, dep := range zktx.FactoryDeps {
			depCopy := dep
			depItem := &Item[[]byte]{
				Label:   fmt.Sprintf("Dep#%v", i),
				Value:   depCopy,
				Details: ShortString(fmt.Sprintf("0x%x", dep), 255),
				Display: func(v []byte) string { return FactoryDeptsDetails([]hexutil.Bytes{v}) },
			}
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
			nit := &Item[[]byte]{Label: "New factory dependency"}
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

func InputXexUI(nit *Item[[]byte], era *zksyncera.ZkSyncTxRLP) {
	HexItem := &Item[struct{}]{Label: "Direct Hex Data"}
	FileItem := &Item[struct{}]{Label: "Hex Data from file"}
	sel := promptui.Select{Label: nit.Label, Items: []MenuItem{HexItem, FileItem, Back}, Templates: ItemTemplate, Size: 10}

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
			fd := nit.Value
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

var ERASignerItem = &Item[signer.Signer]{Label: "Signer"}
var ERASignItem = &Item[[]byte]{Label: "Sign"}

func SignEraUI(era *zksyncera.ZkSyncTxRLP) {
	if era.ChainId2 == nil {
		fmt.Println("WARNING: Chain ID not set")

	} else if era.ChainId1.Cmp(era.ChainId2) != 0 {
		fmt.Println("WARNING: Chain ID mismatch, setting ChainId1 to ChainId2")
		era.ChainId1 = era.ChainId2
	}

	items := []MenuItem{ZkChainIDItem, ERASignerItem, ERASignItem, Back}
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
			// chain ID display only

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
			det = fmt.Sprintf("Error hashing at pos %v: %v", i, err)
			return det
		}
		det += fmt.Sprintf("0x%x  ", h)

	}
	return det
}

var EraCallItem = &Item[struct{}]{Label: "Call RPC"}

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

	rpc := rpcui.SendEndpointItem.Value
	rpcOk := rpc != nil
	for {
		items := []MenuItem{rpcui.SendEndpointItem}
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
			rpc = rpcui.SendEndpointItem.Value
			rpcOk = rpc != nil
		case EraCallItem.Label:
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
