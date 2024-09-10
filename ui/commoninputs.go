package ui

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
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

	prompt := promptui.Prompt{
		Label: item.Label,
	}
	a, err := prompt.Run()
	if err != nil {
		return err
	}
	res, err := ParseBigInt(a)
	if err != nil {
		return fmt.Errorf("error parsing big int: %w", err)
	}
	item.Value = res
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

func YesNoPromptUI(label string) bool {
	prpt := promptui.Prompt{Label: label + "(yes/no)", Default: "no"}
	y, err := prpt.Run()
	return err == nil && y == "yes"
}

func InputNewAddressUI(label string) (string, *common.Address, error) {
	prompt := promptui.Prompt{
		Label: "Name for " + label,
	}
	name, err := prompt.Run()
	if err != nil {
		return "", nil, err
	}

	prompt = promptui.Prompt{
		Label: "Input address of " + name + " (0x...)",
	}
	s, err := prompt.Run()
	if err != nil {
		return "", nil, err
	}
	s = strings.Trim(s, " ")
	common.BytesToAddress(common.FromHex(strings.TrimSpace(s)))
	addr := common.HexToAddress(s)
	return name, &addr, nil
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

var bigIntParserDec = regexp.MustCompile(`(\d*)(\.\d+)?([KMGTPE]?)$`)
var two = big.NewInt(2)
var kilo = big.NewInt(1000)
var mega = big.NewInt(1000000)
var giga = big.NewInt(1000000000)
var tera = big.NewInt(1000000000000)
var peta = big.NewInt(1000000000000000)
var exa = big.NewInt(1000000000000000000)
var ten = big.NewInt(10)
var sixteen = big.NewInt(16)

// Decimal Parsing
func ParseBigInt(input string) (*big.Int, error) {
	input = strings.TrimSpace(input)
	res := new(big.Int)
	if strings.HasPrefix(input, "0x") {
		input = input[2:]
		_, ok := res.SetString(input, 16)
		if !ok {
			return nil, errors.New("invalid hex number")
		}
		return res, nil
	}
	match := bigIntParserDec.FindStringSubmatch(input)
	if match == nil {
		return nil, errors.New("invalid number")
	}
	multiplier := big.NewInt(1)
	divisor := big.NewInt(1)
	alldigits := match[1]
	divpower := 0
	if len(match[2]) > 0 {
		alldigits += match[2][1:]
		divpower = len(match[2]) - 1
	}
	_, ok := res.SetString(alldigits, 10)
	if !ok {
		return nil, errors.New("invalid decimal number")
	}
	switch match[3] {
	case "K":
		multiplier = kilo
	case "M":
		multiplier = mega
	case "G":
		multiplier = giga
	case "T":
		multiplier = tera
	case "P":
		multiplier = peta
	case "E":
		multiplier = exa
	}
	if divpower > 0 {
		divisor.Exp(ten, big.NewInt(int64(divpower)), nil)
	}
	if divisor.Cmp(multiplier) > 0 {
		return nil, errors.New("invalid number")
	}
	res.Mul(res, multiplier)
	res.Div(res, divisor)
	return res, nil
}

func InputFloatUI(it *Item) error {
	prompt := promptui.Prompt{
		Label: it.Label,
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	it.Value = f
	return nil
}
