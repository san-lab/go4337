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
	bytelab := fmt.Sprint(bytecount)
	if bytecount < 0 {
		bytelab = ""
	}
	label := fmt.Sprintf("%s (bytes%s)", item.Label, bytelab)
	defval := ""
	if item.Value != nil {
		defbytes, ok := item.Value.([]byte)
		if ok {
			defval = hex.EncodeToString(defbytes)
		}
	}

	prompt := promptui.Prompt{
		Label:   label,
		Default: defval,
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
		//item.DisplayValue = "0x" + hex.EncodeToString(ibt)
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
		//item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	case 8:
		fbt := [8]byte{}
		copy(fbt[8-len(ibt):], ibt)
		item.Value = fbt
		//item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	case 16:
		fbt := [16]byte{}
		copy(fbt[16-len(ibt):], ibt)
		item.Value = fbt
		//item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	case 32:
		fbt := [32]byte{}
		copy(fbt[32-len(ibt):], ibt)
		item.Value = fbt
		//item.DisplayValue = "0x" + hex.EncodeToString(fbt[:])
	default:
		return fmt.Errorf("unsupported byte count: %d", bytecount)
	}

	return nil
}

func InputBigInt(item *Item) error {
	res := new(big.Int)
	prompt := promptui.Prompt{
		Label: item.Label,
	}
	a, err := prompt.Run()
	if err != nil {
		return err
	}
	base := 10
	if strings.HasPrefix(a, "0x") {
		base = 16
		a = a[2:]
	}
	res, ok := res.SetString(a, base)
	if !ok {
		return fmt.Errorf("error parsing big int: %s", a)
	}
	item.Value = res
	//item.DisplayValue = res.String()

	return nil
}

func InputUint(item *Item, size int) (uint64, error) {
	prompt := promptui.Prompt{
		Label:   item.Label,
		Default: fmt.Sprint(item.Value),
	}
	s, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	base := 10
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
		base = 16
	}
	u, err := strconv.ParseUint(s, base, size)
	if err != nil {
		return 0, err
	}
	switch size {
	case 8:
		item.Value = uint8(u)
	case 16:
		item.Value = uint16(u)
	case 32:
		item.Value = uint32(u)
	case 48, 64:
		item.Value = u
	default:
		return 0, fmt.Errorf("unsupported size: %d", size)

	}
	//item.DisplayValue = fmt.Sprint(u)
	return u, nil
}

func InputNewStringUI(item *Item) error {
	prompt := promptui.Prompt{
		Label: item.Label,
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	item.Value = s
	//item.DisplayValue = s
	return nil
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

func InputBool(item *Item) error {
	prompt := promptui.Select{
		Label: item.Label,
		Items: []string{"true", "false"},
	}
	_, sel, err := prompt.Run()
	if err != nil {
		return err
	}
	item.Value = sel == "true"
	//item.DisplayValue() = sel
	return nil
}

func InputString(item *Item) error {
	prompt := promptui.Prompt{
		Label:   item.Label,
		Default: fmt.Sprint(item.Value),
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	item.Value = s
	//item.DisplayValue = s
	return nil
}
