package zksyncera

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/san-lab/go4337/signer"
)

func (zktx *ZkSyncTxRLP) Decode(input []byte) error {
	return rlp.DecodeBytes(input[1:], zktx)
}

/*
	tx.Nonce = new(big.Int).SetUint64(decodedTx.Nonce)
	tx.GasTipCap = decodedTx.MaxPriorityFeePerGas
	tx.GasFeeCap = decodedTx.MaxFeePerGas
	tx.Gas = decodedTx.GasLimit
	tx.To = decodedTx.To
	tx.Value = decodedTx.Value
	tx.Data = decodedTx.Data
	tx.ChainID = decodedTx.ChainId2
	tx.From = decodedTx.From
	tx.GasPerPubdata = decodedTx.GasPerPubdata
	tx.CustomSignature = decodedTx.CustomSignature
	tx.FactoryDeps = decodedTx.FactoryDeps
	tx.PaymasterParams = decodedTx.PaymasterParams
	return nil
*/

func (zktx *ZkSyncTxRLP) Encode() ([]byte, error) {
	res, err := rlp.EncodeToBytes(zktx)
	if err != nil {
		return nil, fmt.Errorf("failed to encode RLP bytes: %w", err)
	}
	return append([]byte{0x71}, res...), nil
}

/*
	{
		Nonce:                tx.Nonce.Uint64(),
		MaxPriorityFeePerGas: tx.GasTipCap,
		MaxFeePerGas:         tx.GasFeeCap,
		GasLimit:             tx.Gas,
		To:                   tx.To,
		Value:                tx.Value,
		Data:                 tx.Data,
		ChainId1:             tx.ChainID,
		ChainId2:             tx.ChainID,
		From:                 tx.From,
		GasPerPubdata:        tx.GasPerPubdata,
		FactoryDeps:          tx.FactoryDeps,
		CustomSignature:      tx.CustomSignature,
		PaymasterParams:      tx.PaymasterParams,
	}
	if len(zkSyncTxRLP.CustomSignature) == 0 {
		zkSyncTxRLP.CustomSignature = sig
	}
*/

func (zktx *ZkSyncTxRLP) ChainedHash() []byte {
	return nil
}

func (zktx *ZkSyncTxRLP) TypedDataMessage() (apitypes.TypedDataMessage, error) {

	paymaster := big.NewInt(0)
	paymasterInput := hexutil.Bytes{}
	if zktx.PaymasterParams != nil {
		paymaster = big.NewInt(0).SetBytes(zktx.PaymasterParams.Paymaster.Bytes())
		paymasterInput = zktx.PaymasterParams.PaymasterInput
	}
	value := `0x0`
	if zktx.Value != nil {
		value = zktx.Value.String()
	}
	//Here the zkSync version bytes are hardcoded!!!
	factoryDepsHashes, err := hashArray(zktx.FactoryDeps, 1, 0)
	if err != nil {
		return nil, err
	}
	return apitypes.TypedDataMessage{
		"txType":                 `0x71`, // TxType,
		"from":                   big.NewInt(0).SetBytes(zktx.From.Bytes()).String(),
		"to":                     big.NewInt(0).SetBytes(zktx.To.Bytes()).String(),
		"gasLimit":               zktx.GasLimit.String(),
		"gasPerPubdataByteLimit": zktx.GasPerPubdata.String(),
		"maxFeePerGas":           zktx.MaxFeePerGas.String(),
		"maxPriorityFeePerGas":   zktx.MaxPriorityFeePerGas.String(),
		"paymaster":              paymaster.String(),
		"nonce":                  fmt.Sprint(zktx.Nonce),
		"value":                  value,
		"data":                   zktx.Data,
		"factoryDeps":            factoryDepsHashes,
		"paymasterInput":         paymasterInput,
	}, nil
}

func (zktx *ZkSyncTxRLP) TypedData() (*apitypes.TypedData, error) {
	tmsg, err := zktx.TypedDataMessage()
	if err != nil {
		return nil, err
	}
	domain := GetZksyncDomain(*zktx.ChainId2, nil, 0)
	return &apitypes.TypedData{
		Types: apitypes.Types{
			"Transaction":  ZKTypes(),
			"EIP712Domain": Types(domain),
		},
		PrimaryType: "Transaction",
		Domain:      domain, //domain.TypedData(),
		Message:     tmsg,
	}, nil
}

type ZkSyncTxRLP struct {
	Nonce                uint64
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
	GasLimit             *big.Int
	To                   *common.Address `rlp:"nil"` // nil means contract creation
	Value                *big.Int
	Data                 hexutil.Bytes
	// zkSync part
	ChainId1 *big.Int // legacy
	Empty1   string   // legacy
	Empty2   string   // legacy
	ChainId2 *big.Int
	From     *common.Address
	// Meta fields   *Meta
	GasPerPubdata   *big.Int
	FactoryDeps     []hexutil.Bytes
	CustomSignature hexutil.Bytes
	PaymasterParams *PaymasterParams `rlp:"nil"`
}

type PaymasterParams struct {
	Paymaster      common.Address `json:"paymaster"`      // Address of the paymaster.
	PaymasterInput []byte         `json:"paymasterInput"` // Encoded input.
}

func hashArray(harr []hexutil.Bytes, verByteOne, verByteTwo byte) ([]hexutil.Bytes, error) {
	if len(harr) == 0 {
		return []hexutil.Bytes{}, nil
	}
	res := make([]hexutil.Bytes, len(harr))
	for i, d := range harr {
		h, err := ZKBytecodeHash(d, verByteOne, verByteTwo)
		if err != nil {
			return nil, fmt.Errorf("failed to get hash of some bytecode in FactoryDeps: %w", err)
		}
		res[i] = h
	}
	return res, nil
}

func ZKBytecodeHash(bytecode []byte, one, two byte) ([]byte, error) {
	l := len(bytecode)
	if l%32 != 0 {
		return nil, fmt.Errorf("bytecode length in bytes must be divisible by 32")
	}
	l >>= 5 //divide by 32
	rawHash := sha256.Sum256(bytecode)
	if l >= 1<<16 {
		return nil, fmt.Errorf("bytecode length must be less than 2^16 bytes")
	}
	// replace first 2 bytes of hash with version
	version := make([]byte, 32)
	version[0] = one
	version[1] = two
	version[2] = byte(l / 0x100)
	version[3] = byte(l & 0xff)
	copy(version[4:], rawHash[4:])
	return version, nil
}

func GetZksyncDomain(chainId big.Int, verifyingContract *common.Address, salt uint64) apitypes.TypedDataDomain {
	vc := ""
	if verifyingContract != nil && (*verifyingContract) != (common.Address{}) {
		vc = verifyingContract.String()
	}
	ss := ""
	if salt != 0 {
		ss = fmt.Sprintf("%d", salt)
	}
	return apitypes.TypedDataDomain{
		Name:              "zkSync",
		Version:           "2",
		ChainId:           math.NewHexOrDecimal256(chainId.Int64()),
		VerifyingContract: vc,
		Salt:              ss,
	}
}

func ZKTypes() []apitypes.Type {
	return []apitypes.Type{
		{Name: "txType", Type: "uint256"},
		{Name: "from", Type: "uint256"},
		{Name: "to", Type: "uint256"},
		{Name: "gasLimit", Type: "uint256"},
		{Name: "gasPerPubdataByteLimit", Type: "uint256"},
		{Name: "maxFeePerGas", Type: "uint256"},
		{Name: "maxPriorityFeePerGas", Type: "uint256"},
		{Name: "paymaster", Type: "uint256"},
		{Name: "nonce", Type: "uint256"},
		{Name: "value", Type: "uint256"},
		{Name: "data", Type: "bytes"},
		{Name: "factoryDeps", Type: "bytes32[]"},
		{Name: "paymasterInput", Type: "bytes"},
	}
}

// Types returns the domain field types.
func Types(d apitypes.TypedDataDomain) []apitypes.Type {
	types := []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
	}
	if d.VerifyingContract != "" {
		types = append(types, apitypes.Type{Name: "verifyingContract", Type: "address"})
	}
	if d.Salt != "" {
		types = append(types, apitypes.Type{Name: "salt", Type: "uint256"})
	}
	return types
}

func (era *ZkSyncTxRLP) HashEra() ([]byte, error) {
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

// At the moment this metthod IS NOT setting the CustomSignature field
func (era *ZkSyncTxRLP) Sign(vsigner interface{}) ([]byte, error) {
	if vsigner == nil {
		return nil, fmt.Errorf("signer is nil")
	}
	signer, ok := vsigner.(signer.Signer)
	if !ok {
		return nil, fmt.Errorf("invalid signer type")
	}
	h, err := era.HashEra()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash of tx: %w", err)
	}
	return signer.SignHash(h)
}
