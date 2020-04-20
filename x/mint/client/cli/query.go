package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/mint"
)

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryParameters)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params mint.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryInflation(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "inflation",
		Short: "Query the current minting inflation value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryInflation)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}

// GetCmdQueryAnnualProvisions implements a command to return the current minting
// annual provisions value.
func GetCmdQueryAnnualProvisions(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "annual-provisions",
		Short: "Query the current minting annual provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryAnnualProvisions)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}

// GetCmdQueryTotalProvisions implements a command to return the current minting
// total provisions value.
func GetCmdQueryTotalProvisions(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "total-provisions",
		Short: "Query the current minting total provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryTotalProvisions)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			// var provisions sdk.Int
			var provisions mint.TotalProvisions
			if err := cdc.UnmarshalJSON(res, &provisions); err != nil {
				return err
			}

			return cliCtx.PrintOutput(provisions)
		},
	}
}

// GetCmdQueryBlockRewards implements a command to return a block reward
// at given height.
func GetCmdQueryBlockRewards(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "rewards [block height]",
		Short: "Get verified data for a block reward at given height",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			bz, err := cliCtx.Codec.MarshalJSON(mint.NewQueryBlockRewardParams(height))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryBlockRewards)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var reward mint.BlockReward
			// var reward sdk.Int
			// var reward int64
			if err := cliCtx.Codec.UnmarshalJSON(res, &reward); err != nil {
				return err
			}

			return cliCtx.PrintOutput(reward)
			// return cliCtx.PrintOutput(sdk.NewDec(reward))
		},
	}
}
