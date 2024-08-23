package rpccalls

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/san-lab/go4337/state"
)

type KeyContainer interface {
	GetKey() *ecdsa.PrivateKey
}

func CreateSignedTransaction(rpc *state.RPCEndpoint, from, to common.Address, value *big.Int, calldata []byte, gasLimit uint64, key *ecdsa.PrivateKey) (*types.Transaction, error) {
	client, err := ethclient.Dial(rpc.URL)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not get gas price: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, fmt.Errorf("could not get nonce: %v", err)
	}

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, calldata)

	return types.SignTx(tx, types.NewEIP155Signer(rpc.ChainId), key)

}

func SendTransaction(rpc *state.RPCEndpoint, signedTx *types.Transaction) (*common.Hash, error) {
	client, err := ethclient.Dial(rpc.URL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to rpc: %v", err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("could not send tx: %v", err)
	}
	h := signedTx.Hash()
	return &h, nil
}

func CreateAndSendTransaction(rpc *state.RPCEndpoint, from, to common.Address, value *big.Int, calldata []byte, gasLimit uint64, key KeyContainer) (*common.Hash, error) {
	signedTx, err := CreateSignedTransaction(rpc, from, to, value, calldata, gasLimit, key.GetKey())
	if err != nil {
		return nil, fmt.Errorf("could not create signed tx: %v", err)
	}
	return SendTransaction(rpc, signedTx)
}

// CallContract calls a contract with the given calldata
// Attempts to parse the response with the given return type
// Returns the result of parsing and the raw response
// or error, if any
func CallContract(rpc *state.RPCEndpoint, from, to *common.Address, calldata []byte, retType abi.Arguments) ([]interface{}, []byte, error) {
	client, err := ethclient.Dial(rpc.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to rpc: %v", err)
	}
	msg := ethereum.CallMsg{
		From: *from,
		To:   to,
		Data: calldata,
	}
	res, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not call contract: %v", err)
	}
	var ret []interface{}
	if retType != nil {
		ret, err = retType.Unpack(res)
		if err != nil {
			return nil, res, ErrRetParse
		}
	}
	return ret, res, nil
}

var ErrRetParse = fmt.Errorf("could not unpack return values")
