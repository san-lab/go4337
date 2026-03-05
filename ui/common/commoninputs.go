package common

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
	"github.com/holiman/uint256"
	"github.com/manifoldco/promptui"
	"github.com/san-lab/go4337/state"
)

func InputBytes(item *Item[[]byte], bytecount int) error {
	bytelab := fmt.Sprint(bytecount)
	if bytecount < 0 {
		bytelab = ""
	}
	label := fmt.Sprintf("%s (bytes%s)", item.Label, bytelab)
	defval := ""
	if item.Value != nil {
		defval = hex.EncodeToString(item.Value)
	}

	prompt := promptui.Prompt{
		Label:   label,
		Default: defval,
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	ibt := common.FromHex(s)
	if bytecount > 0 && len(ibt) > bytecount {
		return fmt.Errorf("expected %d bytes, got %d", bytecount, len(ibt))
	}
	item.Value = ibt
	return nil
}

func InputBigInt(item *Item[*big.Int]) (*big.Int, error) {

	prompt := promptui.Prompt{
		Label: item.Label,
	}
	if item.Value != nil {
		prompt.Default = fmt.Sprint(item.Value)
	}
	a, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	res, err := ParseBigInt(a)
	if err != nil {
		return nil, fmt.Errorf("error parsing big int: %w", err)
	}
	item.Value = res

	return res, nil
}

func InputUint256(item *Item[*uint256.Int]) (*uint256.Int, error) {
	bigItem := &Item[*big.Int]{Label: item.Label}
	_, err := InputBigInt(bigItem)
	if err != nil {
		return nil, err
	}
	if bigItem.Value == nil {
		return nil, fmt.Errorf("error parsing nil big int")
	}
	uint256int, _ := uint256.FromBig(bigItem.Value)
	item.Value = uint256int
	return uint256int, nil
}

// InputUint64 is the primary helper for uint64 items.
func InputUint64(item *Item[uint64]) (uint64, error) {
	return InputUint(item, 64)
}

func InputUint(item *Item[uint64], size int) (uint64, error) {
	prompt := promptui.Prompt{
		Label:   item.Label,
		Default: fmt.Sprint(item.Value),
	}
	s, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	big, err := ParseBigInt(s)
	if err != nil {
		return 0, err
	}
	u := big.Uint64()
	item.Value = u
	return u, nil
}

func InputNewStringUI(item *Item[string]) error {
	prompt := promptui.Prompt{
		Label: item.Label,
	}
	s, err := prompt.Run()
	if err != nil {
		return err
	}
	item.Value = s
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

func InputBool(item *Item[bool]) error {
	prompt := promptui.Select{
		Label: item.Label,
		Items: []string{"true", "false"},
	}
	_, sel, err := prompt.Run()
	if err != nil {
		return err
	}
	item.Value = sel == "true"
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

func InputFloatUI(it *Item[float64]) error {
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

func InputHexFileUI(label string) ([]byte, error) {
	fmt.Println(label)
	hbytes, err := SelectFileFromFS("")
	if err != nil {
		return nil, err
	}
	sbytes := string(hbytes)
	state.Log("Selected file:", sbytes)
	//trim witespaces
	sbytes = strings.TrimSpace(sbytes)
	sbytes = strings.TrimPrefix(sbytes, "0x")
	return hex.DecodeString(sbytes)

}

