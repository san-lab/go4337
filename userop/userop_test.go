package userop

import (
	"fmt"
	"math/big"
	"testing"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

func TestPackUints(t *testing.T) {
	var a uint64 = 0xaaaa
	var b uint64 = 0xbbbb
	bt := PackUints(a, b)
	fmt.Printf("bt: %x | %x\n", bt[0:16], bt[16:])
	fmt.Println(UnpackUints(bt))
}

/*
sender: '0x63D993703e0514bE8C693D6CC91362A79B92Ce82',

	nonce: 0,
	initCode: '0x',
	callData: '0xb61d27f6000000000000000000000000754925c070c5368ae28d1d552dedae302280336d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000084cadcff3b0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000e6c8fe3ad93f726543689c94864ea0a45e3a7b680000000000000000000000000000000000000000000000000000000000000002613300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
	callGasLimit: 4000000,
	verificationGasLimit: 200000,
	preVerificationGas: 20000,
	maxFeePerGas: 1000,
	maxPriorityFeePerGas: 1000000,
	paymaster: '0x0000000000000000000000000000000000000000',
	paymasterData: '0x',
	paymasterVerificationGasLimit: 0,
	paymasterPostOpGasLimit: 0,
	signature: '0x7b7b8b69b6c2746006c84e6afe5a0470e2c31e2344721a91e8a8483725ab8efe5c10e7e46b56fb5092001e9e5b7070112e6bfa64b87d4ca994af5da3c71f39a71c'
*/
var TestUserOp UserOperation

func init() {
	sender := common.HexToAddress("0x63D993703e0514bE8C693D6CC91362A79B92Ce82")
	factory := common.HexToAddress("0x123456789a")
	TestUserOp = UserOperation{
		Sender:                        &sender,
		Nonce:                         (*U256)(big.NewInt(0)),
		Factory:                       &factory,
		FactoryData:                   []byte{0xff, 0xfe, 0xfd},
		CallData:                      common.Hex2Bytes("b61d27f6000000000000000000000000754925c070c5368ae28d1d552dedae302280336d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000084cadcff3b0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000e6c8fe3ad93f726543689c94864ea0a45e3a7b680000000000000000000000000000000000000000000000000000000000000002613300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		CallGasLimit:                  4000000,
		VerificationGasLimit:          200000,
		PreVerificationGas:            20000,
		MaxFeePerGas:                  1000,
		MaxPriorityFeePerGas:          1000000,
		Paymaster:                     nil,
		PaymasterData:                 []byte{},
		PaymasterVerificationGasLimit: 0,
		PaymasterPostOpGasLimit:       0,
		Signature:                     common.Hex2Bytes("7b7b8b69b6c2746006c84e6afe5a0470e2c31e2344721a91e8a8483725ab8efe5c10e7e46b56fb5092001e9e5b7070112e6bfa64b87d4ca994af5da3c71f39a71c"),
	}
}

func TestHash(t *testing.T) {
	uop := TestUserOp
	hash, err := uop.EncodeToHash()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(hash)
}

func TestFactoryCode(t *testing.T) {
	uop := &TestUserOp
	rem := uop.Pack().MarshalRemix()

	fmt.Println(rem)

	bt, err := json.Marshal(uop)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bt))
	uop2 := new(UserOperation)
	err = json.Unmarshal(bt, uop2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(uop2.Pack().MarshalRemix())
}

func TestPaymasterHash(t *testing.T) {
	fmt.Println("Obsoleted")
	/*
		uop := TestUserOp
		paymaster := common.HexToAddress("0xd4662c4530c9cB1d194Cc2e8c11A13413148Fc6F")
		uop.Paymaster = &paymaster
		bts := PackUints(uop.PaymasterVerificationGasLimit, uop.PaymasterPostOpGasLimit)
		uop.PaymasterData = bts[0:32]

		hashbytes, hash, err := GetPaymasterV7Hash(uop.Pack(), 1, 4, 7)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%x\n", hash)

		fmt.Println(hex.EncodeToString(hashbytes))

		fmt.Println(uop.Pack().MarshalRemix())
	*/
}

func TestUnmarshalWithAuth(t *testing.T) {
	uopJson := `{
      "sender": "0x06E5CaFAdF44ef2085C6cd9AEbB547f8b6d04563",
      "nonce": 12,
      "factory": "0x0000000000000000000000000000000000000000",
      "factoryData": "0x",
      "callData": "0x",
      "callGasLimit": 2000000,
      "verificationGasLimit": 200000,
      "preVerificationGas": 20000,
      "maxFeePerGas": 1000,
      "maxPriorityFeePerGas": 1000000,
      "paymaster": "0x0000000000000000000000000000000000000000",
      "paymasterVerificationGasLimit": 1000,
      "paymasterPostOpGasLimit": 100,
      "paymasterData": "0x",
      "signature": "0x214844ba54356b42cca3e3bd89422d73c4a40e4a07d61c5408b5275db1252ad26dfebfc28b45c0d2b27d1c65fff25afd471746d02b4864d97ce98f1f9cd5d1e11b",
      "eip7702Auth": {
        "chainId": "0x0",
        "address": "0x0000000000000000000000000000000000000000",
        "nonce": "0x11",
        "yParity": "0x0",
        "r": "0x0",
        "s": "0x0"
      }
    }`

	usop := new(UserOperation)
	err := json.Unmarshal([]byte(uopJson), usop)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(usop.Nonce)

	b, _ := usop.Pack().MarshalJSON()
	fmt.Println(string(b))
}
