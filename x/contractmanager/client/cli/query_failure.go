package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/neutron-org/neutron/x/contractmanager/types"
	"github.com/spf13/cobra"
)

func CmdAllFailures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-failures",
		Short: "shows all failures",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllFailureRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AllFailures(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdContractFailures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-failures [address]",
		Short: "shows a failures for specific contract address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			params := &types.QueryGetFailureRequest{
				Address: argIndex,
			}

			res, err := queryClient.Failure(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}