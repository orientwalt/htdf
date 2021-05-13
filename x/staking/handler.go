package staking

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	abci "github.com/tendermint/tendermint/abci/types"
	common "github.com/tendermint/tendermint/libs/strings"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/staking/keeper"
	"github.com/orientwalt/htdf/x/staking/types"
)

func init() {
	// junying-todo,2020-01-17
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info" //trace/debug/info/warn/error/parse/fatal/panic
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.FatalLevel //TraceLevel/DebugLevel/InfoLevel/WarnLevel/ErrorLevel/ParseLevel/FatalLevel/PanicLevel
	}
	// set global log level
	logrus.SetLevel(ll)
	logrus.SetFormatter(&logrus.TextFormatter{}) //&log.JSONFormatter{})
}

func logger() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return logrus.WithField("file", filename).WithField("function", fn)
}

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, msg, k)

		case types.MsgEditValidator:
			return handleMsgEditValidator(ctx, msg, k)

		case types.MsgDelegate:
			return handleMsgDelegate(ctx, msg, k)

		case types.MsgBeginRedelegate:
			return handleMsgBeginRedelegate(ctx, msg, k)

		case types.MsgUndelegate:
			return handleMsgUndelegate(ctx, msg, k)
		case types.MsgSetUndelegateStatus:
			return handleMsgSetDelegatorStatus(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking module").Result()
		}
	}
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	// move to ProtocolV0.EndBlocker , yqq, 2021-05-12,
	// ctx = ctx.WithEventManager(sdk.NewEventManager())

	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates := k.ApplyAndReturnValidatorSetUpdates(ctx)

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		balances, err := k.CompleteUnbonding(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			continue
		}

		// 2021-05-12, yqq
		// we emit events to log unbonding records for application layer, such as block explorer.
		logger().Infof("===================Unbonding Compiled ================")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				// yqq, 2021-05-12,
				// whther balances contains the rewards of delegation after begin-unbonding?
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress.String()),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		balances, err := k.CompleteRedelegation(ctx, dvvTriplet.DelegatorAddress,
			dvvTriplet.ValidatorSrcAddress, dvvTriplet.ValidatorDstAddress)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress.String()),
			),
		)
	}

	return validatorUpdates
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgCreateValidator(ctx sdk.Context, msg types.MsgCreateValidator, k keeper.Keeper) sdk.Result {
	logger().Traceln()
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, msg.ValidatorAddress); found {
		return ErrValidatorOwnerExists(k.Codespace()).Result()
	}
	logger().Traceln()
	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); found {
		return ErrValidatorPubKeyExists(k.Codespace()).Result()
	}
	logger().Traceln()
	if msg.Value.Denom != k.GetParams(ctx).BondDenom {
		return ErrBadDenom(k.Codespace()).Result()
	}
	logger().Traceln()
	if _, err := msg.Description.EnsureLength(); err != nil {
		return err.Result()
	}
	logger().Traceln()
	if ctx.ConsensusParams() != nil {
		logger().Traceln()
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PubKey)
		logger().Traceln()
		if !common.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			logger().Traceln()
			return ErrValidatorPubKeyTypeUnsupported(k.Codespace(),
				tmPubKey.Type,
				ctx.ConsensusParams().Validator.PubKeyTypes).Result()
		}
	}
	logger().Traceln()
	validator := NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)
	commission := NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)
	logger().Traceln()
	validator, err := validator.SetInitialCommission(commission)
	if err != nil {
		return err.Result()
	}
	logger().Traceln()
	validator.MinSelfDelegation = msg.MinSelfDelegation

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)
	logger().Traceln()
	// call the after-creation hook
	k.AfterValidatorCreated(ctx, validator.OperatorAddress)
	logger().Traceln()
	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	_, err = k.Delegate(ctx, msg.DelegatorAddress, msg.Value.Amount, validator, true)
	if err != nil {
		return err.Result()
	}

	logger().Traceln()
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateValidator,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})
	logger().Traceln()
	return sdk.Result{Events: ctx.EventManager().ABCIEvents()} //, nil
}

func handleMsgEditValidator(ctx sdk.Context, msg types.MsgEditValidator, k keeper.Keeper) sdk.Result {
	// validator must already be registered
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return ErrNoValidatorFound(k.Codespace()).Result()
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return err.Result()
	}

	validator.Description = description

	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return err.Result()
		}

		// call the before-modification hook since we're about to update the commission
		k.BeforeValidatorModified(ctx, msg.ValidatorAddress)

		validator.Commission = commission
	}

	if msg.MinSelfDelegation != nil {
		if !(*msg.MinSelfDelegation).GT(validator.MinSelfDelegation) {
			return ErrMinSelfDelegationDecreased(k.Codespace()).Result()
		}
		if (*msg.MinSelfDelegation).GT(validator.Tokens) {
			return ErrSelfDelegationBelowMinimum(k.Codespace()).Result()
		}
		validator.MinSelfDelegation = (*msg.MinSelfDelegation)
	}

	k.SetValidator(ctx, validator)

	// tags := sdk.NewTags(
	// 	tags.DstValidator, msg.ValidatorAddress.String(),
	// 	tags.Moniker, description.Moniker,
	// 	tags.Identity, description.Identity,
	// )

	// return sdk.Result{
	// 	Tags: tags,
	// }

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditValidator,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyCommissionRate, validator.Commission.String()),
			sdk.NewAttribute(types.AttributeKeyMinSelfDelegation, validator.MinSelfDelegation.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().ABCIEvents()}
}

func handleMsgDelegate(ctx sdk.Context, msg types.MsgDelegate, k keeper.Keeper) sdk.Result {
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return ErrNoValidatorFound(k.Codespace()).Result()
	}

	if msg.Amount.Denom != k.GetParams(ctx).BondDenom {
		return ErrBadDenom(k.Codespace()).Result()
	}

	_, err := k.Delegate(ctx, msg.DelegatorAddress, msg.Amount.Amount, validator, true)
	if err != nil {
		return err.Result()
	}

	// tags := sdk.NewTags(
	// 	tags.Delegator, msg.DelegatorAddress.String(),
	// 	tags.DstValidator, msg.ValidatorAddress.String(),
	// )
	// by yqq 2020-11-10
	// log `delegation/withdrawDelegationRewards` rewards informations for application layers
	//
	// DONE, yqq, 2021-05-12,  replay tags with Events
	// var rewardsLog string
	// ts := ctx.CoinFlowTags().GetTags()
	// if ts != nil {
	// 	ctx.Logger().Debug("handleMsgDelegate ", "module", types.ModuleName, "len(ctx.CoinFlowTags().GetTags())", len(ts))
	// 	for _, tg := range ts {
	// 		tg.Key = []byte(sdk.DelegatorRewardFlow)
	// 		strTag := tg.String()
	// 		if strings.Contains(strTag, msg.DelegatorAddress.String()) &&
	// 			strings.Contains(strTag, msg.ValidatorAddress.String()) {
	// 			rewardsLog += strTag
	// 		}
	// 	}
	// }

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(sdk.AttributeKeyDelegator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		//
		// We should make event more readable and concise.
		//
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
		// Log:    rewardsLog,
	}
}

func handleMsgUndelegate(ctx sdk.Context, msg types.MsgUndelegate, k keeper.Keeper) sdk.Result {
	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount.Amount,
	)
	if err != nil {
		return err.Result()
	}

	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, shares)
	if err != nil {
		return err.Result()
	}

	finishTime := types.MsgCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	// tags := sdk.NewTags(
	// 	tags.Delegator, msg.DelegatorAddress.String(),
	// 	tags.SrcValidator, msg.ValidatorAddress.String(),
	// 	tags.EndTime, completionTime.Format(time.RFC3339),
	// )

	// return sdk.Result{Data: finishTime, Tags: tags}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return sdk.Result{Data: finishTime, Events: ctx.EventManager().ABCIEvents()}
}

func handleMsgBeginRedelegate(ctx sdk.Context, msg types.MsgBeginRedelegate, k keeper.Keeper) sdk.Result {
	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.Amount.Amount,
	)
	if err != nil {
		return err.Result()
	}

	completionTime, err := k.BeginRedelegation(
		ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress, shares,
	)
	if err != nil {
		return err.Result()
	}

	finishTime := types.MsgCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	// resTags := sdk.NewTags(
	// 	tags.Delegator, msg.DelegatorAddress.String(),
	// 	tags.SrcValidator, msg.ValidatorSrcAddress.String(),
	// 	tags.DstValidator, msg.ValidatorDstAddress.String(),
	// 	tags.EndTime, completionTime.Format(time.RFC3339),
	// )

	// return sdk.Result{Data: finishTime, Tags: resTags}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedelegate,
			sdk.NewAttribute(sdk.AttributeKeyDelegator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeySrcValidator, msg.ValidatorSrcAddress.String()),
			sdk.NewAttribute(types.AttributeKeyDstValidator, msg.ValidatorDstAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return sdk.Result{Data: finishTime, Events: ctx.EventManager().ABCIEvents()}
}

func handleMsgSetDelegatorStatus(ctx sdk.Context, msg types.MsgSetUndelegateStatus, k keeper.Keeper) sdk.Result {
	// validator must already be registered
	_, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return ErrNoValidatorFound(k.Codespace()).Result()
	}

	del, found := k.GetDelegation(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if !found {
		return types.ErrNoDelegation(k.Codespace()).Result()
	}

	// upgarede delegator status
	// Always set true?
	del.Status = true
	k.UpgradeDelegation(ctx, del)

	// tags := sdk.NewTags(
	// 	tags.Delegator, msg.DelegatorAddress.String(),
	// 	tags.SrcValidator, msg.ValidatorAddress.String(),
	// 	tags.ActionCompleteAuthorization, "true",
	// )

	// return sdk.Result{Tags: tags}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetDelegatorStatus,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyStatus, "true"), //msg.Status),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().ABCIEvents()}
}
