package rpccalls

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/san-lab/go4337/entrypoint/entrypointv6/baseaccv6"
	"github.com/san-lab/go4337/state"
)

func GetWalletNonce(rpc *state.RPCEndpoint, addr common.Address) (uint64, error) {
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

func GetStdNonce(rpc *state.RPCEndpoint, addr common.Address) (uint64, error) {
	client, err := ethclient.Dial(rpc.URL)
	if err != nil {
		return 0, fmt.Errorf("could not connect to rpc: %v", err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return 0, fmt.Errorf("could not get nonce: %v", err)
	}
	return nonce, nil
}

func GetBalance(rpc *state.RPCEndpoint, addr common.Address) (*big.Int, error) {
	client, err := ethclient.Dial(rpc.URL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to rpc: %v", err)
	}
	bal, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get balance: %v", err)
	}
	return bal, nil
}
