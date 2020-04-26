package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"

	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"

	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
	"github.com/orientwalt/htdf/x/distribution/client/common"
	"github.com/orientwalt/htdf/x/distribution/types"
	log "github.com/sirupsen/logrus"
)

var (
	flagOnlyFromValidator = "only-from-validator"
	flagIsValidator       = "is-validator"
	flagComission         = "commission"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *amino.Codec) *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:   "dist",
		Short: "Distribution transactions subcommands",
	}

	distTxCmd.AddCommand(client.PostCommands(
		GetCmdWithdrawRewards(cdc),
		GetCmdSetWithdrawAddr(cdc),
	)...)

	return distTxCmd
}

// command to withdraw rewards
func GetCmdWithdrawRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards [delegator-addr] [validator-addr]",
		Short: "witdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
		Long: strings.TrimSpace(`witdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator:

$ hscli tx distr withdraw-rewards htdf1keyvaa4u5rcjwq3gncvct4hrmq553fpkremp5v htdfvaloper1keyvaa4u5rcjwq3gncvct4hrmq553fpkfqrcr8
$ hscli tx distr withdraw-rewards htdf1keyvaa4u5rcjwq3gncvct4hrmq553fpkremp5v htdfvaloper1keyvaa4u5rcjwq3gncvct4hrmq553fpkfqrcr8 --commission
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{types.NewMsgWithdrawDelegatorReward(delAddr, valAddr)}
			if viper.GetBool(flagComission) {
				msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(valAddr))
			}
			str, err := cliCtx.Codec.MarshalJSON(msgs)
			log.Infoln(string(str))
			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs, delAddr)
		},
	}
	cmd.Flags().Bool(flagComission, false, "also withdraw validator's commission")
	return cmd
}

// command to withdraw all rewards
func GetCmdWithdrawAllRewards(cdc *codec.Codec, queryRoute string) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-all-rewards",
		Short: "withdraw all delegations rewards for a delegator",
		Long: strings.TrimSpace(`Withdraw all rewards for a single delegator:

$ hscli tx distr withdraw-all-rewards htdf1keyvaa4u5rcjwq3gncvct4hrmq553fpkremp5v
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgs, err := common.WithdrawAllDelegatorRewards(cliCtx, cdc, queryRoute, delAddr)
			if err != nil {
				return err
			}

			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs, delAddr)
		},
	}
}

// command to replace a delegator's withdrawal address
func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-withdraw-addr [withdraw-addr]",
		Short: "change the default withdraw address for rewards associated with an address",
		Long: strings.TrimSpace(`Set the withdraw address for rewards associated with a delegator address:

$ hscli tx set-withdraw-addr htdf1keyvaa4u5rcjwq3gncvct4hrmq553fpkremp5v htdf14psya76ttdx5qvqq5zzz2q6v63k2g3h2k599zd
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			withdrawAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)
			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, delAddr)
		},
	}
	return cmd
}
