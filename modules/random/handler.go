package random

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
	sdkerrors "github.com/orientwalt/htdf/types/errors"

	"github.com/orientwalt/htdf/modules/random/keeper"
	"github.com/orientwalt/htdf/modules/random/types"
)

// NewHandler returns a handler for all random msgs
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRequestRandom:
			return handleMsgRequestRandom(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

// handleMsgRequestRandom handles MsgRequestRandom
func handleMsgRequestRandom(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRequestRandom) (*sdk.Result, error) {
	request, err := k.RequestRandom(ctx, msg.Consumer, msg.BlockInterval, msg.Oracle, msg.ServiceFeeCap)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
			),
			sdk.NewEvent(
				types.EventTypeRequestRandom,
				sdk.NewAttribute(types.AttributeKeyRequestID, hex.EncodeToString(types.GenerateRequestID(request))),
				sdk.NewAttribute(types.AttributeKeyGenHeight, fmt.Sprintf("%d", request.Height+int64(msg.BlockInterval))),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
