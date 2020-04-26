package cli

import (
	"fmt"
	"strconv"

	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/gov"
	"github.com/pkg/errors"

	"strings"

	"github.com/spf13/cobra"

	govClientUtils "github.com/orientwalt/htdf/x/gov/client/utils"

	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
	"github.com/spf13/viper"
)

type proposal struct {
	Title       string
	Description string
	Type        string
	Deposit     string
}

var proposalFlags = []string{
	flagTitle,
	flagDescription,
	flagProposalType,
	flagDeposit,
}

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit a proposal along with an initial deposit",
		Long: strings.TrimSpace(`
Submit a proposal along with an initial deposit. Proposal title, description, type and deposit can be given directly or through a proposal JSON file. For example:

$ hscli gov submit-proposal htdf1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 --proposal="path/to/proposal.json" --from mykey

where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10test"
}

is equivalent to

$ hscli gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10test"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			proposal, err := parseSubmitProposalFlags()
			if err != nil {
				return err
			}

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			// Get from address
			from, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// Find deposit amount
			amount, err := sdk.ParseCoins(proposal.Deposit)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			// if !account.GetCoins().IsAllGTE(amount) {
			// 	return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			// }

			proposalType, err := gov.ProposalTypeFromString(proposal.Type)
			if err != nil {
				return err
			}

			msg := gov.NewMsgSubmitProposal(proposal.Title, proposal.Description, proposalType, from, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			if proposalType == gov.ProposalTypeSoftwareUpgrade {

				version_ := viper.GetInt64(flagVersion)
				if version_ < 0 {
					return errors.Errorf("Version must greater than or equal to zero")
				}

				version := uint64(version_)
				software := viper.GetString(flagSoftware)

				switchHeight_ := viper.GetInt64(flagSwitchHeight)
				if switchHeight_ < 0 {
					return errors.Errorf("SwitchHeight must greater than or equal to zero")
				}
				switchHeight := uint64(switchHeight_)

				thresholdStr := viper.GetString(flagThreshold)
				threshold, err := sdk.NewDecFromStr(thresholdStr)
				if err != nil {
					return err
				}
				msg := gov.NewMsgSubmitSoftwareUpgradeProposal(msg, version, software, switchHeight, threshold)

				err = msg.ValidateBasic()
				if err != nil {
					return err
				}
				return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, from)
			}
			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, from)
		},
	}

	cmd.Flags().String(flagTitle, "", "title of proposal")
	cmd.Flags().String(flagDescription, "", "description of proposal")
	cmd.Flags().String(flagProposalType, "", "proposalType of proposal, types: text/parameter_change/software_upgrade")
	cmd.Flags().String(flagDeposit, "", "deposit of proposal")
	cmd.Flags().String(flagProposal, "", "proposal file path (if this path is given, other proposal flags are ignored)")

	cmd.Flags().String(flagVersion, "0", "the version of the new protocol")
	cmd.Flags().String(flagSoftware, " ", "the software of the new protocol")
	cmd.Flags().String(flagSwitchHeight, "0", "the switchheight of the new protocol")
	cmd.Flags().String(flagThreshold, "0.8", "the upgrade signal threshold of the software upgrade")

	cmd.MarkFlagRequired(flagTitle)
	cmd.MarkFlagRequired(flagDescription)
	cmd.MarkFlagRequired(flagProposalType)

	return cmd
}

// GetCmdDeposit implements depositing tokens for an active proposal.
func GetCmdDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [proposal-id] [deposit]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens for activing proposal",
		Long: strings.TrimSpace(`
Submit a deposit for an acive proposal. You can find the proposal-id by running hscli query gov proposals:

$ hscli tx gov deposit htdf1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 1 10stake
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = govClientUtils.QueryProposalByID(proposalID, cliCtx, cdc, queryRoute)
			if err != nil {
				return fmt.Errorf("Failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// Get from address
			from, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// Fetch associated account
			// account, err := cliCtx.GetAccount(from)
			// if err != nil {
			// 	return err
			// }

			// Get amount of coins
			amount, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			// ensure account has enough coins
			// if !account.GetCoins().IsAllGTE(amount) {
			// 	return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			// }

			msg := gov.NewMsgDeposit(from, proposalID, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, from)
		},
	}
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(3),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(`
Submit a vote for an acive proposal. You can find the proposal-id by running hscli query gov proposals:

$ hscli tx gov vote htdf1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 1 yes --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			// Get voting address
			from, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[1])
			}

			// check to see if the proposal is in the store
			_, err = govClientUtils.QueryProposalByID(proposalID, cliCtx, cdc, queryRoute)
			if err != nil {
				return fmt.Errorf("Failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// Find out which vote option user chose
			byteVoteOption, err := gov.VoteOptionFromString(govClientUtils.NormalizeVoteOption(args[2]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := gov.NewMsgVote(from, proposalID, byteVoteOption)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, from)
		},
	}
}

// DONTCOVER
