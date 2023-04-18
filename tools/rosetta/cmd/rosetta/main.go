package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/log"
	rosettaCmd "github.com/cosmos/cosmos-sdk/tools/rosetta/cmd"
)

func main() {
	var (
		logger            = log.NewLoggerWithKV(os.Stdout, log.ModuleKey, "rosetta")
		interfaceRegistry = codectypes.NewInterfaceRegistry()
		cdc               = codec.NewProtoCodec(interfaceRegistry)
	)

	if err := rosettaCmd.RosettaCommand(interfaceRegistry, cdc).Execute(); err != nil {
		logger.Error("failed to run rosetta", "error", err)
		os.Exit(1)
	}
}
