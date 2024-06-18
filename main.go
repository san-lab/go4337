package main

import (
	"fmt"

	"github.com/san-lab/go4337/ecsigner"
	"github.com/san-lab/go4337/state"
	"github.com/san-lab/go4337/ui"
)

func main() {
	fmt.Println("Welcome to go4337!")
	ecsigner.Init()
	state.InitState()
	ui.RootUI()
	state.State.Save()
}
