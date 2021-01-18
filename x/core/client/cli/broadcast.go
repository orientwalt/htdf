package cli

import (
	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/codec"
	coretypes "github.com/orientwalt/htdf/x/core/types"
	"github.com/spf13/cobra"
)

// junying-todo-20190327
// GetCmdBroadCast is the CLI command for broadcasting a signed transaction
/*
	inspired by
	hscli tx broadcast signed.json
*/
func GetCmdBroadCast(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast [rawdata]",
		Short: "broadcast signed transaction",
		Long:  "hscli tx broadcast 72032..13123",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// load sign tx from string
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			stdTx, err := coretypes.ReadStdTxFromRawData(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}
			// convert tx to bytes
			txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(stdTx)
			if err != nil {
				return err
			}
			// broadcast
			res, err := cliCtx.BroadcastTx(txBytes)
			cliCtx.PrintOutput(res)
			return err

		},
	}
	return client.PostCommands(cmd)[0]
}
