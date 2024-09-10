package rpccalls

import (
	"fmt"
	"testing"
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
