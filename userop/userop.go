package userop

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	Nonce  uint64          `json:"nonce"` // internaly, for the ease of increment, etc
	// Factory address, only for new accounts
	Factory *common.Address `json:"factory,omitempty"`
	// Data for account factory (only if account factory exists)
	FactoryData []byte `json:"factoryData,omitempty"`
	// Data to pass to the sender during the main execution call
	CallData []byte `json:"callData"`
	// The amount of gas to allocate the main execution call
	CallGasLimit uint64 `json:"callGasLimit"`
	// The amount of gas to allocate for the verification step
	VerificationGasLimit uint64 `json:"verificationGasLimit"`
	// Extra gas to pay the bunder
	PreVerificationGas uint64 `json:"preVerificationGas"`
	// Maximum fee per gas (similar to EIP-1559 max_fee_per_gas)
	MaxFeePerGas uint64 `json:"maxFeePerGas"`
	// Maximum priority fee per gas (similar to EIP-1559 max_priority_fee_per_gas)
	MaxPriorityFeePerGas uint64 `json:"maxPriorityFeePerGas"`
	// Address of paymaster contract, (or empty, if account pays for itself)
	Paymaster *common.Address `json:"paymaster,omitempty"`
	// The amount of gas to allocate for the paymaster validation code
	PaymasterVerificationGasLimit uint64 `json:"paymasterVerificationGasLimit"`
	// The amount of gas to allocate for the paymaster post-operation code
	PaymasterPostOpGasLimit uint64 `json:"paymasterPostOpGasLimit"`
	// Data for paymaster (only if paymaster exists)
	PaymasterData []byte `json:"paymasterData"`
	// Data passed into the account to verify authorization
	Signature   []byte                      `json:"signature"`
	EIP7702Auth *types.SetCodeAuthorization `json:"eip7702Auth,omitempty"`
}

/*
"eip7702Auth": {
    "chainId": "0xaa36a7",       // Chain ID (e.g., 11155111)
    "address": "0x4e59b44847b379578588920ca78fbf26c0b4956c", // The Smart Account Implementation to delegate to
    "nonce": "0x0",               // The EOA's nonce for this specific authorization
    "yParity": "0x1",             // v (0 or 1)
    "r": "0x4f3b...",             // r component of the EOA's signature
    "s": "0x2a1c..."              // s component of the EOA's signature
  }
*/

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
	Sender                        string                      `json:"sender"`
	Nonce                         uint64                      `json:"nonce"`
	Factory                       any                         `json:"factory,omitempty"`
	FactoryData                   string                      `json:"factoryData,omitempty"`
	CallData                      string                      `json:"callData"`
	CallGasLimit                  uint64                      `json:"callGasLimit"`
	VerificationGasLimit          uint64                      `json:"verificationGasLimit"`
	PreVerificationGas            uint64                      `json:"preVerificationGas"`
	MaxFeePerGas                  uint64                      `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          uint64                      `json:"maxPriorityFeePerGas"`
	Paymaster                     any                         `json:"paymaster,omitempty"`
	PaymasterVerificationGasLimit uint64                      `json:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       uint64                      `json:"paymasterPostOpGasLimit"`
	PaymasterData                 string                      `json:"paymasterData"`
	Signature                     string                      `json:"signature"`
	EIP7702Auth                   *types.SetCodeAuthorization `json:"eip7702Auth,omitempty"`
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
			Factory:                       JSAddressHex(u.Factory, ""),
			FactoryData:                   BytesToString(u.FactoryData),
			CallData:                      BytesToString(u.CallData),
			CallGasLimit:                  u.CallGasLimit,
			VerificationGasLimit:          u.VerificationGasLimit,
			PreVerificationGas:            u.PreVerificationGas,
			MaxFeePerGas:                  u.MaxFeePerGas,
			MaxPriorityFeePerGas:          u.MaxPriorityFeePerGas,
			Paymaster:                     JSAddressHex(u.Paymaster, ""),
			PaymasterVerificationGasLimit: u.PaymasterVerificationGasLimit,
			PaymasterPostOpGasLimit:       u.PaymasterPostOpGasLimit,
			PaymasterData:                 BytesToString(u.PaymasterData),
			Signature:                     BytesToString(u.Signature),
			EIP7702Auth:                   u.EIP7702Auth,
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
		*u.Factory = common.HexToAddress(fmt.Sprintf("%s", adapter.Factory))
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
		*u.Paymaster = common.HexToAddress(fmt.Sprintf("%s", adapter.Paymaster))
	}
	u.PaymasterVerificationGasLimit = adapter.PaymasterVerificationGasLimit
	u.PaymasterPostOpGasLimit = adapter.PaymasterPostOpGasLimit
	u.PaymasterData = common.FromHex(adapter.PaymasterData)
	u.Signature = common.FromHex(adapter.Signature)
	u.EIP7702Auth = adapter.EIP7702Auth
	return nil
}

func (u *UserOperation) TotalGasLimit() uint64 {
	tot := u.CallGasLimit + u.VerificationGasLimit + u.PreVerificationGas +
		u.PaymasterVerificationGasLimit + u.PaymasterPostOpGasLimit
	return tot
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

func (u *UserOperation) MarshalV6UserOp() entrypointv6.UserOperation {
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

	return *uop6
}

func (u *UserOperation) MarshalRemixV6() string {
	uop6 := u.MarshalV6UserOp()

	return fmt.Sprintf("[%s, %d, %s, %s, %d, %d, %d, %d, %d, %s, %s]", ToRemixHex(uop6.Sender.Bytes()), uop6.Nonce, ToRemixHex(uop6.InitCode), ToRemixHex(uop6.CallData),
		uop6.CallGasLimit, uop6.VerificationGasLimit, uop6.PreVerificationGas, uop6.MaxFeePerGas, uop6.MaxPriorityFeePerGas,
		ToRemixHex(uop6.PaymasterAndData), ToRemixHex(uop6.Signature))
}

func (u *UserOperation) MarshalValuesV6() []interface{} {
	values := make([]interface{}, 11)
	source := u
	values[0] = *source.Sender
	values[1] = source.Nonce
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

func (u *UserOperation) MarshalAlchemy() string {
	// TODO use omitempty instead!
	json := fmt.Sprintf("{\"sender\": \"%s\", ", u.Sender.Hex())
	json += fmt.Sprintf("\"nonce\": \"0x%x\", ", u.Nonce)
	if u.Factory != nil && u.Factory.Cmp(common.BytesToAddress([]byte{0})) != 0 {
		// Non zero address, otherwise skip fields
		json += fmt.Sprintf("\"factory\": \"%s\", ", u.Factory.Hex())
		json += fmt.Sprintf("\"factoryData\": \"0x%x\", ", u.FactoryData)
	}
	json += fmt.Sprintf("\"callData\": \"0x%x\", ", u.CallData)
	json += fmt.Sprintf("\"callGasLimit\": \"0x%x\", ", u.CallGasLimit)
	json += fmt.Sprintf("\"maxPriorityFeePerGas\": \"0x%x\", ", u.MaxPriorityFeePerGas)
	json += fmt.Sprintf("\"maxFeePerGas\": \"0x%x\", ", u.MaxFeePerGas)
	json += fmt.Sprintf("\"preVerificationGas\": \"0x%x\", ", u.PreVerificationGas)
	json += fmt.Sprintf("\"verificationGasLimit\": \"0x%x\", ", u.VerificationGasLimit)
	json += fmt.Sprintf("\"paymasterVerificationGasLimit\": \"0x%x\", ", u.PaymasterVerificationGasLimit)
	json += fmt.Sprintf("\"paymasterPostOpGasLimit\": \"0x%x\", ", u.PaymasterPostOpGasLimit)
	if u.Paymaster != nil && u.Paymaster.Cmp(common.BytesToAddress([]byte{0})) != 0 {
		// Bundler answers with paymaster signature failure...
		// json += fmt.Sprintf("paymaster: %s,", u.Paymaster.Hex())
		// json += fmt.Sprintf("paymasterData: 0x%x,", u.PaymasterData)
	}
	json += fmt.Sprintf("\"signature\": \"0x%x\"}", u.Signature)
	return json
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
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

func (puop *PackedUserOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Sender             string   `json:"sender"`
			Nonce              *big.Int `json:"nonce"`
			InitCode           string   `json:"initCode"`
			CallData           string   `json:"callData"`
			AccountGasLimits   string   `json:"accountGasLimits"`
			PreVerificationGas *big.Int `json:"preVerificationGas"`
			GasFees            string   `json:"gasFees"`
			PaymasterAndData   string   `json:"paymasterAndData"`
			Signature          string   `json:"signature"`
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
		//fmt.Println("Factory is nil")
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
		Nonce:              big.NewInt(int64(u.Nonce)),
		InitCode:           initCode,
		CallData:           u.CallData,
		AccountGasLimits:   PackUints(u.VerificationGasLimit, u.CallGasLimit),
		PreVerificationGas: big.NewInt(int64(u.PreVerificationGas)),
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

func JSAddressHex(addr *common.Address, provider string) any {
	var zeroval any
	if provider == "pimlico" || provider == "alchemy" {
		zeroval = nil
	} else {
		zeroval = "0x"
	}
	if addr == nil || (*addr).Hex() == (common.Address{}).Hex() {
		return zeroval
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
