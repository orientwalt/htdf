package slashing

import (
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/orientwalt/htdf/types"
)

// slashing begin block functionality
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, sk Keeper) sdk.Tags {

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		sk.handleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
	}

	// Iterate through any newly discovered evidence of infraction
	// Slash any validators (and since-unbonded stake within the unbonding period)
	// who contributed to valid infractions
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			sk.handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
		default:
			ctx.Logger().With("module", "app/v1/slashing").Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}

	return sdk.EmptyTags()
}

// slashing end block functionality
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, sk Keeper) (tags sdk.Tags) {
	ctx = ctx.WithCoinFlowTrigger(sdk.SlashEndBlocker)
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "app/v1/slashing"))
	// Tag the height
	tags = sdk.NewTags("height", []byte(strconv.FormatInt(req.Height, 10)))

	if int64(ctx.CheckValidNum()) < ctx.BlockHeader().NumTxs {
		proposalCensorshipTag := sk.handleProposerCensorship(ctx,
			ctx.BlockHeader().ProposerAddress,
			ctx.BlockHeight())
		tags = tags.AppendTags(proposalCensorshipTag)
	}
	return
}
