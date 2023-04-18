package main

import (
	"github.com/cosmos/cosmos-sdk/tools/hubl/internal"
)

func main() {
	cmd, err := internal.RootCommand()
	if err != nil {
		panic(err)
	}

	if err = cmd.Execute(); err != nil {
		panic(err)
	}
}
