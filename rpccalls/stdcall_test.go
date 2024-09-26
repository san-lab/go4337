package rpccalls

import (
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestEntryPoints(t *testing.T) {
	// Test for SupportedEntryPoints
	key := "i8rGhGi6snXwgEX0oqXLD1wWY3WM2KbI"
	url := "https://eth-sepolia.g.alchemy.com/v2"
	res, err := SupportedEntryPoints(url, key)
	if err != nil {
		t.Errorf("Error making API call: %v", err)
	}
	fmt.Println("Supported Entry Points:")
	for _, ep := range *res {
		fmt.Println(ep)
	}

}

func TestTemplate(t *testing.T) {
	// Test for Template
	key := "i8rGhGi6snXwgEX0oqXLD1wWY3WM2KbI"
	url := "https://eth-sepolia.g.alchemy.com/v2/{{.}}"
	turl, err := template.New("apiurl").Parse(url)
	if err != nil {
		t.Errorf("Error parsing template: %v", err)
	}
	wr := &strings.Builder{}
	turl.Execute(wr, key)
	fmt.Println(wr.String())
}
