package ui

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
)

func InputBytes(item *Item, bytecount int) error {
	prompt := promptui.Prompt{
		Label: item.Details,
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	thebytes := common.FromHex(s)
	if bytecount > 0 && len(thebytes) > bytecount {
		return fmt.Errorf("expected %d bytes, got %d", bytecount, len(thebytes))
	}
	ibt := common.FromHex(s)
	if bytecount < 0 {

		item.Value = ibt
		item.DisplayValue = "0x" + hex.EncodeToString(ibt)
		return nil
	}
	if len(ibt) > bytecount {
		return fmt.Errorf("expected %d bytes, got %d", bytecount, len(ibt))
	}
	switch bytecount {
	case 4:
		fbt := [4]byte{}
		copy(fbt[4-len(ibt):], ibt)
		item.Value = fbt
		item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	case 8:
		fbt := [8]byte{}
		copy(fbt[8-len(ibt):], ibt)
		item.Value = fbt
		item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	case 16:
		fbt := [16]byte{}
		copy(fbt[16-len(ibt):], ibt)
		item.Value = fbt
		item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	case 32:
		fbt := [32]byte{}
		copy(fbt[32-len(ibt):], ibt)
		item.Value = fbt
		item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	default:
		return fmt.Errorf("unsupported byte count: %d", bytecount)
	}

	return nil
}

func InputBigInt(label string) (*big.Int, error) {
	res := new(big.Int)
	prompt := promptui.Prompt{
		Label: label,
	}
	a, err := prompt.Run()
	if err != nil {
		return res, err
	}
	base := 10
	if strings.HasPrefix(a, "0x") {
		base = 16
		a = a[2:]
	}
	res, ok := res.SetString(a, base)
	if ok {
		return res, nil
	}
	return res, fmt.Errorf("error parsing uint: %s", a)
}

func InputUint(item *Item, size int) error {
	prompt := promptui.Prompt{
		Label: item.Details,
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	u, err := strconv.ParseUint(s, 10, size)
	if err != nil {
		return err
	}
	switch size {
	case 8:
		item.Value = uint8(u)
	case 16:
		item.Value = uint16(u)
	case 32:
		item.Value = uint32(u)
	case 64:
		item.Value = u
	}
	item.DisplayValue = fmt.Sprint(u)
	return nil
}

func InputNewStringUI(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	return prompt.Run()
}

func MultiLineInput(label string) (string, error) {

	fmt.Println(label)

	// Collect subsequent lines
	lines := []string{}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error reading input: %v\n", err)
			return "", err
		}
		line = strings.TrimRight(line, "\n")
		if line == "" || strings.TrimSpace(line) == "END" {
			break
		}
		lines = append(lines, line)
	}

	// Join all lines into a single multi-line string
	multiLineInput := strings.Join(lines, "\n")
	return multiLineInput, nil
}

func sanitize(abi string, l int) string {
	// Remove newlines and leading/trailing whitespace
	san := strings.TrimSpace(strings.ReplaceAll(abi, "\n", ""))
	san = strings.ReplaceAll(san, "\t", "")
	if l > 0 && len(san) > l {
		return san[:l] + "..."
	}
	return san
}

func InputNewAddressUI(label string) (*common.Address, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	s, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	s = strings.Trim(s, " ")
	common.BytesToAddress(common.FromHex(s))
	addr := common.HexToAddress(s)
	return &addr, nil
}

func ShortHex(data []byte, l int) string {
	if len(data) < 2*l+3 {
		return hex.EncodeToString(data)
	}
	return fmt.Sprintf("%s...%s", hex.EncodeToString(data[:l]), hex.EncodeToString(data[len(data)-l:]))
}
