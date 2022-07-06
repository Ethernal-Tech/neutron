package cli

import (
	"context"
	"fmt"
	"strconv"

	// "strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/neutron-org/neutron/x/interchainqueries/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group interchainqueries queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryRegisteredQueries())
	cmd.AddCommand(CmdQueryRegisteredQuery())
	cmd.AddCommand(CmdQueryRegisteredQueryResult())
	cmd.AddCommand(CmdQueryTransactionsSearchResult())
	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdQueryRegisteredQuery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registered-query [id]",
		Short: "queries all the interchain queries in the module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			queryID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse query id: %w", err)
			}

			res, err := queryClient.RegisteredQuery(context.Background(), &types.QueryRegisteredQueryRequest{QueryId: queryID})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryRegisteredQueries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registered-queries",
		Short: "queries all the interchain queries in the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.RegisteredQueries(context.Background(), &types.QueryRegisteredQueriesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryRegisteredQueryResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-result [query-id]",
		Short: "queries result for registered query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			queryID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse query id: %w", err)
			}

			res, err := queryClient.QueryResult(context.Background(), &types.QueryRegisteredQueryResultRequest{QueryId: queryID})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryTransactionsSearchResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-transactions-search-result [query-id] [start] [end]",
		Short: "queries result for transactions search",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			queryID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse query id: %w", err)
			}

			start, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse start: %w", err)
			}

			end, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse end: %w", err)
			}

			res, err := queryClient.QueryTransactions(context.Background(), &types.QuerySubmittedTransactionsRequest{
				QueryId: queryID,
				Start:   start,
				End:     end,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
