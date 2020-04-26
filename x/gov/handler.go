package gov

import (
	"fmt"
	"strconv"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/gov/tags"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgSubmitSoftwareUpgradeProposal:
			return handleMsgSubmitSoftwareUpgradeProposal(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) sdk.Result {
	var content ProposalContent
	switch msg.ProposalType {
	case ProposalTypeText:
		content = NewTextProposal(msg.Title, msg.Description)
	default:
		return ErrInvalidProposalType(keeper.codespace, msg.ProposalType).Result()
	}
	proposal, err := keeper.SubmitProposal(ctx, content)
	if err != nil {
		return err.Result()
	}
	proposalID := proposal.GetProposalID()
	proposalIDStr := fmt.Sprintf("%d", proposalID)

	err, votingStarted := keeper.AddDeposit(ctx, proposalID, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDStr,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Data: keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID),
		Tags: resTags,
	}
	return sdk.Result{}
}

//Submit upgrade software proposal
func handleMsgSubmitSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitSoftwareUpgradeProposal) sdk.Result {
	if !keeper.protocolKeeper.IsValidVersion(ctx, msg.Version) {

		return ErrCodeInvalidVersion(keeper.codespace, msg.Version).Result()
	}
	if uint64(ctx.BlockHeight()) > msg.SwitchHeight {
		return ErrCodeInvalidSwitchHeight(keeper.codespace, uint64(ctx.BlockHeight()), msg.SwitchHeight).Result()
	}
	_, found := keeper.guardianKeeper.GetProfiler(ctx, msg.Proposer)
	if !found {
		return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
	}
	if _, ok := keeper.protocolKeeper.GetUpgradeConfig(ctx); ok {
		return ErrSwitchPeriodInProcess(keeper.codespace).Result()
	}

	proposal := keeper.NewSoftwareUpgradeProposal(ctx, msg)
	err, votingStarted := keeper.AddInitialDeposit(ctx, proposal, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}
	proposalIDBytes := strconv.FormatUint(proposal.GetProposalID(), 10)

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}

	// keeper.AddProposalNum(ctx, proposal)
	return sdk.Result{
		Data: []byte(proposalIDBytes),
		Tags: resTags,
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) sdk.Result {
	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return err.Result()
	}

	proposalIDStr := fmt.Sprintf("%d", msg.ProposalID)
	resTags := sdk.NewTags(
		tags.Depositor, []byte(msg.Depositor.String()),
		tags.ProposalID, proposalIDStr,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) sdk.Result {
	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Voter, msg.Voter.String(),
			tags.ProposalID, fmt.Sprintf("%d", msg.ProposalID),
		),
	}
}
