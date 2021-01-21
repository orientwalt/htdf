package cli

import (
	"fmt"
	"strings"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
	"github.com/orientwalt/htdf/x/slashing"
	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unjail [address]",
		Short: "unjail validator previously jailed for downtime",
		Example: `
  hscli tx slashing unjail htdfvaloper1z7lmlh2qjnsjkefpz38krkmzsgl23zwcvd3tjh

  or

  hscli tx slashing unjail htdf1z7lmlh2qjnsjkefpz38krkmzsgl23zwcx5fj9u
		`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			strAddress := args[0]
			var validatorAddr sdk.AccAddress

			if strings.HasPrefix(strAddress, sdk.GetConfig().GetBech32ValidatorAddrPrefix()) {
				bzAddr, err := sdk.ValAddressFromBech32(strAddress)
				if err != nil {
					return err
				}
				validatorAddr = sdk.AccAddress(bzAddr)
			} else if strings.HasPrefix(strAddress, sdk.GetConfig().GetBech32AccountAddrPrefix()) {
				validatorAddr, err = sdk.AccAddressFromBech32(strAddress)
				if err != nil {
					return err
				}
			} else {
				err = fmt.Errorf("Invalid address prefix, it must be one of htdfvaloper1 or htdf1 . Please use --help for more details.")
				return
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(validatorAddr))
			err = hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, validatorAddr)
			return
		},
	}
}
