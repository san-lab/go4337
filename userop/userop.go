package userop

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/san-lab/go4337/entrypoint/entrypointv6"
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

type UserOperation struct {
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

func NewUserOperationWithDefaults() *UserOperation {
	uop := &UserOperation{
		CallGasLimit:                  DefaultCallGasLimit,
		VerificationGasLimit:          DefaultVerificationGasLimit,
		PreVerificationGas:            DefaultPreVerificationGas,
		MaxFeePerGas:                  DefaultMaxFeePerGas,
		MaxPriorityFeePerGas:          DefaultMaxPriorityFeePerGas,
		PaymasterVerificationGasLimit: DefaultPaymasterVerificationGasLimit,
		PaymasterPostOpGasLimit:       DefaultPaymasterPostOpGasLimit,
	}
	return uop
}

type UsOpJsonAdapter struct {
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
}

func (u UserOperation) MarshalJSON() ([]byte, error) {
	var senderstring = ""
	if u.Sender != nil {
		senderstring = u.Sender.Hex()
	}
	return json.Marshal(
		&UsOpJsonAdapter{
			Sender:                        senderstring,
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

func (u *UserOperation) UnmarshalJSON(data []byte) error {
	adapter := &UsOpJsonAdapter{}
	err := json.Unmarshal(data, adapter)
	if err != nil {
		return err
	}
	if adapter.Sender != "" {
		addr := common.HexToAddress(adapter.Sender)
		u.Sender = &addr
	}
	u.Nonce = adapter.Nonce
	if adapter.Factory != "" {
		u.Factory = new(common.Address)
		*u.Factory = common.HexToAddress(adapter.Factory)
	}
	u.FactoryData = common.FromHex(adapter.FactoryData)
	u.CallData = common.FromHex(adapter.CallData)
	u.CallGasLimit = adapter.CallGasLimit
	u.VerificationGasLimit = adapter.VerificationGasLimit
	u.PreVerificationGas = adapter.PreVerificationGas
	u.MaxFeePerGas = adapter.MaxFeePerGas
	u.MaxPriorityFeePerGas = adapter.MaxPriorityFeePerGas
	if adapter.Paymaster != "" {
		u.Paymaster = new(common.Address)
		*u.Paymaster = common.HexToAddress(adapter.Paymaster)
	}
	u.PaymasterVerificationGasLimit = adapter.PaymasterVerificationGasLimit
	u.PaymasterPostOpGasLimit = adapter.PaymasterPostOpGasLimit
	u.PaymasterData = common.FromHex(adapter.PaymasterData)
	u.Signature = common.FromHex(adapter.Signature)
	return nil
}

func (u *UserOperation) InitData() []byte {
	if u.Factory != nil && *(u.Factory) != (common.Address{}) {
		return append(u.Factory.Bytes(), u.FactoryData...)
	}
	return []byte{}
}

func (u *UserOperation) PaymasterAndData() []byte {
	if u.Paymaster != nil && *(u.Paymaster) != (common.Address{}) {
		return append(u.Paymaster.Bytes(), u.PaymasterData...)
	}
	return []byte{}
}

func (u *UserOperation) MarshalRemixV6() string {
	uop6 := new(entrypointv6.UserOperation)
	uop6.Sender = *u.Sender
	uop6.Nonce = big.NewInt(int64(u.Nonce))
	//if u.Factory != nil && *(u.Factory) != (common.Address{}) {
	//	uop6.InitCode = append(u.Factory.Bytes(), u.FactoryData...)
	//}
	uop6.InitCode = u.InitData()
	uop6.CallData = u.CallData
	uop6.CallGasLimit = big.NewInt(int64(u.CallGasLimit))
	uop6.VerificationGasLimit = big.NewInt(int64(u.VerificationGasLimit))
	uop6.PreVerificationGas = big.NewInt(int64(u.PreVerificationGas))
	uop6.MaxFeePerGas = big.NewInt(int64(u.MaxFeePerGas))
	uop6.MaxPriorityFeePerGas = big.NewInt(int64(u.MaxPriorityFeePerGas))
	//if u.Paymaster != nil && *(u.Paymaster) != (common.Address{}) {
	//	uop6.PaymasterAndData = append(u.Paymaster.Bytes(), u.PaymasterData...)
	//}
	uop6.PaymasterAndData = u.PaymasterAndData()
	uop6.Signature = u.Signature

	return fmt.Sprintf("[%s, %d, %s, %s, %d, %d, %d, %d, %d, %s, %s]", ToRemixHex(uop6.Sender.Bytes()), uop6.Nonce, ToRemixHex(uop6.InitCode), ToRemixHex(uop6.CallData),
		uop6.CallGasLimit, uop6.VerificationGasLimit, uop6.PreVerificationGas, uop6.MaxFeePerGas, uop6.MaxPriorityFeePerGas,
		ToRemixHex(uop6.PaymasterAndData), ToRemixHex(uop6.Signature))
}

func (u *UserOperation) MarshalValuesV6() []interface{} {
	values := make([]interface{}, 11)
	source := u
	values[0] = *source.Sender
	values[1] = big.NewInt(int64(source.Nonce))
	values[2] = source.InitData()
	values[3] = source.CallData
	values[4] = big.NewInt(int64(source.CallGasLimit))
	values[5] = big.NewInt(int64(source.VerificationGasLimit))
	values[6] = big.NewInt(int64(source.PreVerificationGas))
	values[7] = big.NewInt(int64(source.MaxFeePerGas))
	values[8] = big.NewInt(int64(source.MaxPriorityFeePerGas))
	values[9] = source.PaymasterAndData()
	values[10] = source.Signature
	return values
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
	DefaultCallGasLimit                  = uint64(2000000)
	DefaultVerificationGasLimit          = uint64(200000)
	DefaultPreVerificationGas            = uint64(20000)
	DefaultMaxFeePerGas                  = uint64(1000)
	DefaultMaxPriorityFeePerGas          = uint64(1000000)
	DefaultPaymasterVerificationGasLimit = uint64(1000)
	DefaultPaymasterPostOpGasLimit       = uint64(100)
)

func (u *UserOperation) Pack() *PackedUserOp {
	if u == nil {
		fmt.Println("UserOp is nil")
		return nil
	}
	initCode := []byte{}
	if u.Factory != nil && *(u.Factory) != (common.Address{}) {
		initCode = append(u.Factory.Bytes(), u.FactoryData...)
	} else {
		fmt.Println("Factory is nil")
	}
	paymasterData := []byte{}
	if u.Paymaster != nil && *(u.Paymaster) != (common.Address{}) {
		gasbytes := PackUints(u.PaymasterVerificationGasLimit, u.PaymasterPostOpGasLimit)
		paymasterData = append(u.Paymaster.Bytes(), gasbytes[:]...)
		paymasterData = append(paymasterData, u.PaymasterData...)
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
	binary.BigEndian.PutUint64(buf[24:], b)
	binary.BigEndian.PutUint64(buf[8:], a)
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
