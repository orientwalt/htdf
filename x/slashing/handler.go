package slashing

import (
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/slashing/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking module").Result()
		}
	}
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func handleMsgUnjail(ctx sdk.Context, msg types.MsgUnjail, k Keeper) sdk.Result {
	validator := k.validatorSet.Validator(ctx, msg.ValidatorAddr)
	if validator == nil {
		return types.ErrNoValidatorForAddress(k.codespace).Result()
	}

	// cannot be unjailed if no self-delegation exists
	selfDel := k.validatorSet.Delegation(ctx, sdk.AccAddress(msg.ValidatorAddr), msg.ValidatorAddr)
	if selfDel == nil {
		return types.ErrMissingSelfDelegation(k.codespace).Result()
	}

	if validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
		return types.ErrSelfDelegationTooLowToUnjail(k.codespace).Result()
	}

	// cannot be unjailed if not jailed
	if !validator.IsJailed() {
		return types.ErrValidatorNotJailed(k.codespace).Result()
	}

	consAddr := sdk.ConsAddress(validator.GetConsPubKey().Address())

	info, found := k.getValidatorSigningInfo(ctx, consAddr)
	if !found {
		return types.ErrNoValidatorForAddress(k.codespace).Result()
	}

	// cannot be unjailed if tombstoned
	if info.Tombstoned {
		return types.ErrValidatorJailed(k.codespace).Result()
	}

	// cannot be unjailed until out of jail
	if ctx.BlockHeader().Time.Before(info.JailedUntil) {
		return types.ErrValidatorJailed(k.codespace).Result()
	}

	// unjail the validator
	k.validatorSet.Unjail(ctx, consAddr)

	// tags := sdk.NewTags(
	// 	tags.Action, tags.ActionValidatorUnjailed,
	// 	tags.Validator, msg.ValidatorAddr.String(),
	// )

	// return sdk.Result{
	// 	Tags: tags,
	// }

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().ABCIEvents()}
}
