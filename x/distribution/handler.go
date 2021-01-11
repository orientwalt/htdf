package distribution

import (
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/distribution/keeper"
	"github.com/orientwalt/htdf/x/distribution/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, msg, k)
		case types.MsgWithdrawDelegatorReward:
			return handleMsgWithdrawDelegatorReward(ctx, msg, k)
		case types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in distribution module").Result()
		}
	}
}

// These functions assume everything has been authenticated (ValidateBasic passed, and signatures checked)

func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.MsgSetWithdrawAddress, k keeper.Keeper) sdk.Result {

	err := k.SetWithdrawAddr(ctx, msg.DelegatorAddress, msg.WithdrawAddress)
	if err != nil {
		return err.Result()
	}

	// tags := sdk.NewTags(
	// 	tags.Delegator, []byte(msg.DelegatorAddress.String()),
	// )
	// return sdk.Result{
	// 	Tags: tags,
	// }

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().ABCIEvents()}
}

func handleMsgWithdrawDelegatorReward(ctx sdk.Context, msg types.MsgWithdrawDelegatorReward, k keeper.Keeper) sdk.Result {
	rewards, err := k.WithdrawDelegationRewards(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return err.Result()
	}

	// return sdk.Result{
	// 	Tags: sdk.NewTags(
	// 		tags.Rewards, rewards.String(),
	// 		tags.Delegator, []byte(msg.DelegatorAddress.String()),
	// 		tags.Validator, []byte(msg.ValidatorAddress.String()),
	// 	),
	// }
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()), // inserted by junying
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().ABCIEvents()}
}

func handleMsgWithdrawValidatorCommission(ctx sdk.Context, msg types.MsgWithdrawValidatorCommission, k keeper.Keeper) sdk.Result {
	commission, err := k.WithdrawValidatorCommission(ctx, msg.ValidatorAddress)
	if err != nil {
		return err.Result()
	}

	// return sdk.Result{
	// 	Tags: sdk.NewTags(
	// 		tags.Commission, commission.String(),
	// 		tags.Validator, []byte(msg.ValidatorAddress.String()),
	// 	),
	// }

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()), // inserted by junying
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().ABCIEvents()}
}
