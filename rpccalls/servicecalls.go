package rpccalls

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/san-lab/go4337/entrypoint/entrypointv6/baseaccv6"
	"github.com/san-lab/go4337/state"
)

func GetNonce(rpc *state.RPCEndpoint, addr common.Address) (uint64, error) {
	client, err := ethclient.Dial(rpc.URL)
	if err != nil {
		return 0, fmt.Errorf("could not connect to rpc: %v", err)
	}
	bac, err := baseaccv6.NewBaseaccv6Caller(addr, client)
	if err != nil {
		return 0, fmt.Errorf("could not create baseacc contract caller: %v", err)
	}
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}
	nonce, err := bac.GetNonce(callOpts)
	if err != nil {
		return 0, fmt.Errorf("could not get nonce: %v", err)
	}
	return nonce.Uint64(), nil

}
