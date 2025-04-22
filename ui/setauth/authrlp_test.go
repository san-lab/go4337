package setauth

import (
	"bytes"
	"fmt"
	"testing"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const encodedtx = `b8d204f8cf83aa36a781ad8502540be4008506fc23ac00830186a094aab05558448c8a9597287db9f61e2d751645b12a8080c0f85cf85a8094ab37a30d244e2f8a2868ffc1b164e08b5b46bfc58080a0e0ea724978ce1bb9ca82ba616234d59a88f51273237a288004386846a30c9203a0270db12fe3b8950981ff55f63860530e33f437f11f016f5c9c329748d82cdc1b80a031ae14a4d3216ec2998b1206fda190b9913af96f0b4aadf0268127a1eb128220a03d852e45f06845693229022a0886c28fe10223e31b8f4882d3f4d58ad877b367`

func TestDecoding(t *testing.T) {
	tx := &types.Transaction{}
	tx.DecodeRLP(rlp.NewStream(bytes.NewReader(common.FromHex(encodedtx)), 0))
	jb, _ := json.MarshalIndent(tx, "", "  ")
	fmt.Println(string(jb))
	w := new(bytes.Buffer)
	tx.EncodeRLP(w)
	fmt.Printf("\n%x\n", w.Bytes())
}
