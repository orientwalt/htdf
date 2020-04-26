package cli

import (
	"fmt"

	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	htdfservice "github.com/orientwalt/htdf/x/core"
	"github.com/spf13/cobra"
)

// junying-todo, 2020-04-01
// GetCmdCall is the CLI command for call contract.
func GetCmdCall(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract [contract-address] [callcode]",
		Short: "query contract data",
		Long:  "hscli query contract htdf...  7839124400000000...",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// load sign tx from string
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			contractaddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			callcode := args[1]
			//
			bz, err := cliCtx.Codec.MarshalJSON(htdfservice.NewQueryContractParams(contractaddr, callcode))
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", htdfservice.QuerierRoute, htdfservice.QueryContract)
			res, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			var answer string
			if err := cliCtx.Codec.UnmarshalJSON(res, &answer); err != nil {
				return err
			}
			//
			// cliCtx.PrintOutput(res)
			fmt.Println(answer)
			return nil

		},
	}
	return client.PostCommands(cmd)[0]
}
