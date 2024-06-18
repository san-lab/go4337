package signer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Signer interface {
	Sign([]byte) ([]byte, error)
	String() string
	Type() string
	Marshal() ([]byte, error)
}

/*
OpenZeppelin:

	function toEthSignedMessageHash(bytes memory message) internal pure returns (bytes32) {
	        return
	            keccak256(bytes.concat("\x19Ethereum Signed Message:\n", bytes(Strings.toString(message.length)), message));
	    }
*/
const SigMsgHeader = "\x19Ethereum Signed Message:\n"

func ToEthSignedMessageHash(message []byte) [32]byte {
	buf := []byte(SigMsgHeader)
	buf = append(buf, ([]byte(fmt.Sprint(len(message))))...)
	buf = append(buf, message...)
	return [32]byte(crypto.Keccak256(buf))
}

type AddSigner func() error
