package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

// GetQueryCmd returns the transaction commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryParamsCmd(),
		QueryInturnRelayerCmd(),
	)

	return cmd
}

// QueryParamsCmd returns the command handler for evidence parameter querying.
func QueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current oracle parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query the current oracle parameters:

$ <appd> query oracle params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryParamsCmd returns the command handler for evidence parameter querying.
func QueryInturnRelayerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inturn-relayer",
		Short: "Query the inturn relayer",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query the inturn relayer:

$ <appd> query oracle inturn-relayer
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.InturnRelayer(cmd.Context(), &types.QueryInturnRelayerRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
