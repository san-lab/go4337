package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/zksyncera"
)

var ZkEIP712TxItem = &Item{Label: "ZkTx Content"}
var PrintItem = &Item{Label: "Print"}
var EncodeItem = &Item{Label: "RLP Encode"}
var DecodeItem = &Item{Label: "Decode from RLP"}
var CurlItem = &Item{Label: "CURL call"}
var SignEraItem = &Item{Label: "Sign"}

func ZkSyncEraUI() {
	items := []*Item{
		ZkEIP712TxItem,
		PrintItem,
		EncodeItem,
		DecodeItem,
		CurlItem,
		SignEraItem,
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
			bt, err := json.MarshalIndent(gzktx, "", "  ")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(bt))
			}
		case ZkEIP712TxItem.Label:
			ZkEIP712TxUI(gzktx)
		case DecodeItem.Label:
			enctx, _ := hex.DecodeString(EncodedTransaction1[2:])
			rit := &Item{Label: "RLP bytes in hex", Value: enctx}
			err := InputBytes(rit, -1)
			if err != nil {
				fmt.Println(err)
			} else if rlpb, ok := rit.Value.([]byte); ok && len(rlpb) > 0 {

				err := gzktx.Decode(rlpb)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Invalid RLP bytes")
			}
		case EncodeItem.Label:
			bt, err := gzktx.Encode()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("0x%x\n", bt)
			}
		case CurlItem.Label:
			enctx, err := gzktx.Encode()
			if err != nil {
				fmt.Println(err)
				continue
			}
			url := "http://localhost:3000"
			call := fmt.Sprintf(`curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[0x%x],"id":1}' %s`, enctx, url)
			fmt.Println(call)
		case SignEraItem.Label:
			SignEraUI(gzktx)
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

var gzktx = new(zksyncera.ZkSyncTxRLP)

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
				gzktx.ChainId2 = ZkChainIDItem.Value.(*big.Int)
				gzktx.ChainId1 = ZkChainIDItem.Value.(*big.Int)
			}
		case ZkNonceItem.Label:
			n, err := InputUint(ZkNonceItem, 64)
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.Nonce = n
			}
		case ZkFromItem.Label:
			_, addr, ok := AddressFromBookUI("zkFrom")
			if ok {
				gzktx.From = addr
				ZkFromItem.Value = addr
			}
		case ZkToItem.Label:
			_, addr, ok := AddressFromBookUI("zkTo")
			if ok {
				gzktx.To = addr
				ZkToItem.Value = addr
			}
		case ZkValueItem.Label:
			err = InputBigInt(ZkValueItem)
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.Value = ZkValueItem.Value.(*big.Int)
			}
		case ZkDataItem.Label:
			calldat, err := PotentiallyRecursiveCallDataUI()
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.Data = calldat
				ZkDataItem.Value = calldat
			}
		case ZkGasLimitItem.Label:
			err := InputBigInt(ZkGasLimitItem)
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.GasLimit = ZkGasLimitItem.Value.(*big.Int)
			}

		case ZkGasTipCapItem.Label:
			err := InputBigInt(ZkGasTipCapItem)
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.MaxPriorityFeePerGas = ZkGasTipCapItem.Value.(*big.Int)
			}

		case ZkGasFeeCapItem.Label:
			err := InputBigInt(ZkGasFeeCapItem)
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.MaxFeePerGas = ZkGasFeeCapItem.Value.(*big.Int)
			}
		case ZkGasPerPubdataItem.Label:
			err := InputBigInt(ZkGasPerPubdataItem)
			if err != nil {
				fmt.Println(err)
			} else {
				gzktx.GasPerPubdata = ZkGasPerPubdataItem.Value.(*big.Int)
			}
		case ZkCustomSignatureItem.Label:
			SignEraUI(gzktx)
		case ZkFactoryDepsItem.Label:
			ZkFactoryDepsUI(zktx)
		case ZkPaymasterParamsItem.Label:
			ZkPaymasterParamsUI()
		case PrintItem.Label:
			bt, err := json.MarshalIndent(gzktx, "", "  ")
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
	for {
		items := []*Item{}
		for i, dep := range zktx.FactoryDeps {
			depItem := &Item{Label: fmt.Sprintf("Dep#%v", i)}
			depItem.Value = dep
			items = append(items, depItem)
		}
		items = append(items, AppendItem, RemoveItem, Back)

		spr := promptui.Select{Label: "Factory Dependencies", Items: items, Templates: ItemTemplate, Size: 10}
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
			err := InputBytes(nit, -1)
			if err != nil {
				fmt.Println(err)
			} else {
				fd := nit.Value.([]byte)
				zktx.FactoryDeps = append(zktx.FactoryDeps, fd[:])
				continue
			}

		case RemoveItem.Label:
			if i > 0 && i < len(zktx.FactoryDeps) {
				zktx.FactoryDeps = append(zktx.FactoryDeps[:i], zktx.FactoryDeps[i+1:]...)
				continue
			}
		default:
			fmt.Println("Not implemented yet:", sel)
		}
	}

}

func ZkPaymasterParamsUI() {
	fmt.Println("ZkPaymasterParamsUI")
}

func HashEra(era *zksyncera.ZkSyncTxRLP) ([]byte, error) {
	td, err := era.TypedData()
	if err != nil {
		return nil, err
	}
	h, _, err := apitypes.TypedDataAndHash(*td)
	if err != nil {
		return nil, err
	}
	return h, nil
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
	h, err := HashEra(era)
	if err != nil {
		fmt.Println(err)
		return
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
			SignerUI(ERASignerItem)
		case ERASignItem.Label:
			signer, ok := ERASignerItem.Value.(signer.Signer)
			if !ok {
				fmt.Println("Invalid Signer")
				return
			}
			ps := signer.PlainSigner()
			sig, err := ps.Sign(h[:])
			if err != nil {
				fmt.Println("error signing:", err)
				return
			}
			if len(sig) != 65 {
				fmt.Println("Invalid signature length:", len(sig))
				return
			}
			if sig[64] < 27 {
				sig[64] += 27
			}
			era.CustomSignature = sig
			ERASignItem.Value = sig
			fmt.Println("Signature:", hex.EncodeToString(sig[:]))
		}
	}

}

var EncodedTransaction1 = "0x71f909750f84017d784084017d78408310c34094000000000000000000000000000000000000800680b8849c4d535b000000000000000000000000000000000000000000000000000000000000000001000043f2ab5ca62ea1739a8279e29bff4c37bb9f850c5ca08fc984a65201240000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000082012c808082012c94aab05558448c8a9597287db9f61e2d751645b12a82c350f90863b908600000008003000039000000400030043f0000000100200190000000440000c13d000000000201001900000060022002700000003405200197000000040050008c0000005f0000413d000000000201043b000000e002200270000000360020009c0000004c0000613d000000370020009c0000005f0000c13d000000240050008c0000005f0000413d0000000002000416000000000002004b0000005f0000c13d0000000402100370000000000302043b0000003c0030009c0000005f0000213d0000002302300039000000000052004b0000005f0000813d0000000404300039000000000241034f000000000202043b0000003c0020009c0000005f0000213d00000024033000390000000006320019000000000056004b0000005f0000213d000000000500041a000000010750019000000001065002700000007f0660618f0000001f0060008c00000000050000390000000105002039000000000057004b000000590000c13d000000200060008c0000003d0000413d0000001f0520003900000005055002700000003d0550009a000000200020008c00000038050040410000001f0660003900000005066002700000003d0660009a000000000065004b0000003d0000813d000000000005041b0000000105500039000000000065004b000000390000413d0000001f0020008c000000a40000a13d0000003e05200198000000b00000c13d00000038040000410000000006000019000000ba0000013d0000000001000416000000000001004b0000005f0000c13d0000002001000039000001000010044300000120000004430000003501000041000000cc0001042e0000000001000416000000000001004b0000005f0000c13d000000000200041a000000010320019000000001012002700000007f0110618f0000001f0010008c00000000040000390000000104002039000000000442013f0000000100400190000000610000613d0000003a01000041000000000010043f0000002201000039000000040010043f0000003b01000041000000cd000104300000000001000019000000cd00010430000000800010043f000000000003004b0000006a0000613d000000000000043f000000000001004b000000700000c13d0000002003000039000000a001000039000000860000013d0000003f02200197000000a00020043f000000000001004b000000c001000039000000a0010060390000007b0000013d000000380200004100000000040000190000000003040019000000000402041a000000a005300039000000000045043500000001022000390000002004300039000000000014004b000000720000413d000000c001300039000000610110008a0000003e03100197000000390030009c000000850000a13d0000003a01000041000000000010043f0000004101000039000000040010043f0000003b01000041000000cd000104300000008001300039000000400010043f00000020020000390000000000210435000000a004300039000000800200043d0000000000240435000000c003300039000000000002004b000000970000613d00000000040000190000000005340019000000a006400039000000000606043300000000006504350000002004400039000000000024004b000000900000413d0000001f042000390000003e04400197000000000232001900000000000204350000004002400039000000340020009c00000034020080410000006002200210000000340010009c00000034010080410000004001100210000000000112019f000000cc0001042e000000000002004b0000000003000019000000aa0000613d0000002003400039000000000131034f000000000301043b0000000301200210000000400110027f0000004001100167000000000313016f0000000101200210000000c70000013d000000380400004100000000060000190000000007360019000000000771034f000000000707043b000000000074041b00000001044000390000002006600039000000000056004b000000b20000413d000000000025004b000000c50000813d0000000305200210000000f80550018f000000400550027f00000040055001670000000003360019000000000131034f000000000101043b000000000151016f000000000014041b00000001010000390000000103200210000000000113019f000000000010041b0000000001000019000000cc0001042e000000cb00000432000000cc0001042e000000cd000104300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff000000020000000000000000000000000000004000000100000000000000000000000000000000000000000000000000000000000000000000000000febb0f7e0000000000000000000000000000000000000000000000000000000097bc14aa290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563000000000000000000000000000000000000000000000000ffffffffffffff7f4e487b71000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffd6f21326ab749d5729fcba5677c79037b459436ab7bff709c9d06ce9f10c1a9dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000000000009de69d7a19df2d7504993086cb7920feda759ff4e09aef6677d85d0ddea87c9ab841efe2bd74d8be4f196eb01202fff708b937319c7b0ff3857eed6712eb155d3c4659a4a5ff310fcbb63931e587047ab826b9a21d94f80b90ca99a5f5449b40c3cc1cc0"
