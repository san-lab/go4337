package userop

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

/* ERC 4337 UserOperation
sender	address	The account making the operation
nonce	uint256	Anti-replay parameter (see “Semi-abstracted Nonce Support” )
factory	address	account factory, only for new accounts
factoryData	bytes	data for account factory (only if account factory exists)
callData	bytes	The data to pass to the sender during the main execution call
callGasLimit	uint256	The amount of gas to allocate the main execution call
verificationGasLimit	uint256	The amount of gas to allocate for the verification step
preVerificationGas	uint256	Extra gas to pay the bunder
maxFeePerGas	uint256	Maximum fee per gas (similar to EIP-1559 max_fee_per_gas)
maxPriorityFeePerGas	uint256	Maximum priority fee per gas (similar to EIP-1559 max_priority_fee_per_gas)
paymaster	address	Address of paymaster contract, (or empty, if account pays for itself)
paymasterVerificationGasLimit	uint256	The amount of gas to allocate for the paymaster validation code
paymasterPostOpGasLimit	uint256	The amount of gas to allocate for the paymaster post-operation code
paymasterData	bytes	Data for paymaster (only if paymaster exists)
signature	bytes	Data passed into the account to verify authorization
*/

type UserOp struct {
	Sender *common.Address `json:"sender"`
	Nonce  uint64          `json:"nonce"`
	// Factory address, only for new accounts
	Factory *common.Address `json:"factory,omitempty"`
	// Data for account factory (only if account factory exists)
	FactoryData []byte `json:"factorydata,omitempty"`
	// Data to pass to the sender during the main execution call
	CallData []byte `json:"calldata"`
	// The amount of gas to allocate the main execution call
	CallGasLimit uint64 `json:"callgaslimit"`
	// The amount of gas to allocate for the verification step
	VerificationGasLimit uint64 `json:"verificationgaslimit"`
	// Extra gas to pay the bunder
	PreVerificationGas uint64 `json:"preverificationgas"`
	// Maximum fee per gas (similar to EIP-1559 max_fee_per_gas)
	MaxFeePerGas uint64 `json:"maxfeepergas"`
	// Maximum priority fee per gas (similar to EIP-1559 max_priority_fee_per_gas)
	MaxPriorityFeePerGas uint64 `json:"maxpriorityfeepergas"`
	// Address of paymaster contract, (or empty, if account pays for itself)
	Paymaster *common.Address `json:"paymaster,omitempty"`
	// The amount of gas to allocate for the paymaster validation code
	PaymasterVerificationGasLimit uint64 `json:"paymasterverificationgaslimit"`
	// The amount of gas to allocate for the paymaster post-operation code
	PaymasterPostOpGasLimit uint64 `json:"paymasterpostopgaslimit"`
	// Data for paymaster (only if paymaster exists)
	PaymasterData []byte `json:"paymasterdata"`
	// Data passed into the account to verify authorization
	Signature []byte `json:"signature"`
}

func (u *UserOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Sender                        string `json:"sender"`
			Nonce                         uint64 `json:"nonce"`
			Factory                       string `json:"factory,omitempty"`
			FactoryData                   string `json:"factorydata,omitempty"`
			CallData                      string `json:"calldata"`
			CallGasLimit                  uint64 `json:"callgaslimit"`
			VerificationGasLimit          uint64 `json:"verificationgaslimit"`
			PreVerificationGas            uint64 `json:"preverificationgas"`
			MaxFeePerGas                  uint64 `json:"maxfeepergas"`
			MaxPriorityFeePerGas          uint64 `json:"maxpriorityfeepergas"`
			Paymaster                     string `json:"paymaster,omitempty"`
			PaymasterVerificationGasLimit uint64 `json:"paymasterverificationgaslimit"`
			PaymasterPostOpGasLimit       uint64 `json:"paymasterpostopgaslimit"`
			PaymasterData                 string `json:"paymasterdata"`
			Signature                     string `json:"signature"`
		}{
			Sender:                        u.Sender.Hex(),
			Nonce:                         u.Nonce,
			Factory:                       JSAddressHex(u.Factory),
			FactoryData:                   BytesToString(u.FactoryData),
			CallData:                      BytesToString(u.CallData),
			CallGasLimit:                  u.CallGasLimit,
			VerificationGasLimit:          u.VerificationGasLimit,
			PreVerificationGas:            u.PreVerificationGas,
			MaxFeePerGas:                  u.MaxFeePerGas,
			MaxPriorityFeePerGas:          u.MaxPriorityFeePerGas,
			Paymaster:                     JSAddressHex(u.Paymaster),
			PaymasterVerificationGasLimit: u.PaymasterVerificationGasLimit,
			PaymasterPostOpGasLimit:       u.PaymasterPostOpGasLimit,
			PaymasterData:                 BytesToString(u.PaymasterData),
			Signature:                     BytesToString(u.Signature),
		},
	)
}

// PackedUserOp is an EntryPoint viev of UserOp
/*
sender	address
nonce	uint256
initCode	bytes	concatenation of factory address and factoryData (or empty)
callData	bytes
accountGasLimits	bytes32	concatenation of verificationGas (16 bytes) and callGas (16 bytes)
preVerificationGas	uint256
gasFees	bytes32	concatenation of maxPriorityFee (16 bytes) and maxFeePerGas (16 bytes)
paymasterAndData	bytes	concatenation of paymaster fields (or empty)
signature	bytes
*/
type PackedUserOp struct {
	Sender             *common.Address
	Nonce              uint64
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas uint64
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

func (puop *PackedUserOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Sender             string `json:"sender"`
			Nonce              uint64 `json:"nonce"`
			InitCode           string `json:"initCode"`
			CallData           string `json:"callData"`
			AccountGasLimits   string `json:"accountGasLimits"`
			PreVerificationGas uint64 `json:"preVerificationGas"`
			GasFees            string `json:"gasFees"`
			PaymasterAndData   string `json:"paymasterAndData"`
			Signature          string `json:"signature"`
		}{
			Sender:             puop.Sender.Hex(),
			Nonce:              puop.Nonce,
			InitCode:           JSBytesToString(puop.InitCode),
			CallData:           JSBytesToString(puop.CallData),
			AccountGasLimits:   fmt.Sprintf("0x%x", puop.AccountGasLimits),
			PreVerificationGas: puop.PreVerificationGas,
			GasFees:            fmt.Sprintf("0x%x", puop.GasFees),
			PaymasterAndData:   JSBytesToString(puop.PaymasterAndData),
			Signature:          JSBytesToString(puop.Signature),
		},
	)
}

func (u *PackedUserOp) MarshalRemix() string {
	return fmt.Sprintf("[%s, %d, %s, %s, %s, %d, %s, %s, %s]",
		ToRemixHex([]byte(u.Sender.Bytes())), u.Nonce, ToRemixHex(u.InitCode), ToRemixHex(u.CallData), ToRemixHex(u.AccountGasLimits[:]),
		u.PreVerificationGas, ToRemixHex(u.GasFees[:]), ToRemixHex(u.PaymasterAndData), ToRemixHex(u.Signature))
}

//Default Gas cost values
/*
 callGasLimit: 4000000,
  verificationGasLimit: 200000,
  preVerificationGas: 20000,
  maxFeePerGas: 1000,
  maxPriorityFeePerGas: 1000000,
  paymasterVerificationGasLimit: 0,
  paymasterPostOpGasLimit: 0,
*/
var (
	DefaultCallGasLimit                  = uint64(4000000)
	DefaultVerificationGasLimit          = uint64(200000)
	DefaultPreVerificationGas            = uint64(20000)
	DefaultMaxFeePerGas                  = uint64(1000)
	DefaultMaxPriorityFeePerGas          = uint64(1000000)
	DefaultPaymasterVerificationGasLimit = uint64(0)
	DefaultPaymasterPostOpGasLimit       = uint64(0)
)

func (u *UserOp) Pack() *PackedUserOp {
	initCode := []byte{}
	if u.Factory != nil && *(u.Factory) != (common.Address{}) {
		initCode = append(u.Factory.Bytes(), u.FactoryData...)
	}
	paymasterData := []byte{}
	if u.Paymaster != nil && *(u.Paymaster) != (common.Address{}) {
		paymasterData = append(u.Paymaster.Bytes(), u.PaymasterData...)
	}
	if u.Sender == nil {
		u.Sender = &common.Address{}
	}
	return &PackedUserOp{
		Sender:             u.Sender,
		Nonce:              u.Nonce,
		InitCode:           initCode,
		CallData:           u.CallData,
		AccountGasLimits:   PackUints(u.VerificationGasLimit, u.CallGasLimit),
		PreVerificationGas: u.PreVerificationGas,
		GasFees:            PackUints(u.MaxPriorityFeePerGas, u.MaxFeePerGas),
		PaymasterAndData:   paymasterData,
		Signature:          u.Signature,
	}
}

func PackUints(a, b uint64) [32]byte {
	buf := make([]byte, 32)
	binary.BigEndian.PutUint64(buf[24:], a)
	binary.BigEndian.PutUint64(buf[8:], b)
	return [32]byte(buf)
}

func UnpackUints(buf [32]byte) (uint64, uint64) {
	return binary.BigEndian.Uint64(buf[24:]), binary.BigEndian.Uint64(buf[8:16])
}

func BytesToString(bts []byte) string {
	if len(bts) == 0 {
		return "0x"
	}
	return fmt.Sprintf("0x%x", bts)
}

func JSBytesToString(bts []byte) string {

	return fmt.Sprintf("'%s'", BytesToString(bts))
}

func JSAddressHex(addr *common.Address) string {
	if addr == nil {
		return "0x"
	}
	return addr.Hex()
}

func SliceToBytes32(s []byte) ([32]byte, error) {
	if len(s) > 32 {
		return [32]byte{}, fmt.Errorf("slice is too long")
	}
	a := [32]byte{}
	copy(a[32-len(s):], s)
	return a, nil
}

func UnsafeSliceToBytes32(s []byte) [32]byte {
	a, _ := SliceToBytes32(s)
	return a
}

func ToRemixHex(data []byte) string {
	if len(data) == 0 {
		return `"0x"`
	}
	return fmt.Sprintf(`"0x%x"`, data)

}

func ParamToString(p interface{}) string {
	switch v := p.(type) {
	case *common.Address:
		return v.Hex()
	case []byte, [4]byte, [8]byte, [16]byte, [32]byte:
		return fmt.Sprintf("0x%x", v)
	case uint64, uint32, uint16, uint8, int64, int32, int16, int8:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
