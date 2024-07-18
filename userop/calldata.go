package userop

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func EncodeWithParams(abis *abi.ABI, method string, params ...interface{}) (data []byte, err error) {
	data, err = abis.Pack(method, params...)
	if err != nil {
		err = fmt.Errorf("Error packing params: %w", err)
	}
	return
}
