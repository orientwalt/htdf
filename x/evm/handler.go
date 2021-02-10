package htdfservice

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	"github.com/orientwalt/htdf/x/evm/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	// junying-todo,2020-01-17
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info" //trace/debug/info/warn/error/parse/fatal/panic
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.FatalLevel //TraceLevel/DebugLevel/InfoLevel/WarnLevel/ErrorLevel/ParseLevel/FatalLevel/PanicLevel
	}
	// set global log level
	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{}) //&log.JSONFormatter{})
}

func logger() *log.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return log.WithField("file", filename).WithField("function", fn)
}

// New HTDF Message Handler
// connected to handler.go
// HandleMsgSend, HandleMsgAdd upgraded to EVM version
// commented by junying, 2019-08-21
func NewHandler(accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSend:
			return HandleMsgSend(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		default:
			return HandleUnknownMsg(msg)
		}
	}

}

// junying-todo, 2019-08-26
func HandleUnknownMsg(msg sdk.Msg) sdk.Result {
	var sendTxResp types.SendTxResp
	logger().Debugf("msgType error|mstType=%v\n", msg.Type())
	sendTxResp.ErrCode = sdk.ErrCode_Param
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
}

// junying-todo, 2019-08-26
func HandleMsgSend(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg types.MsgSend) sdk.Result {
	// initialize
	var sendTxResp types.SendTxResp

	st := types.NewStateTransition(ctx, msg)

	evmResult, err := st.TransitionDb(ctx, accountKeeper, feeCollectionKeeper)
	// logger().Debugf("handler:evmResult[%v]\n", evmResult)
	if evmResult == nil {
		sendTxResp.EvmOutput = fmt.Sprintf("%s\n", err)
		if st.ContractCreation {
			sendTxResp.ErrCode = sdk.ErrCode_CreateContract
		} else {
			sendTxResp.ErrCode = sdk.ErrCode_OpenContract
		}
		return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String(), GasUsed: st.GasUsed}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMsgSend,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	if msg.To != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMsgSend,
				sdk.NewAttribute(types.AttributeKeyRecipient, msg.To.String()),
			),
		)
	}
	if st.ContractCreation {
		sendTxResp.ContractAddress = sdk.ToAppAddress(*st.ContractAddress).String()
	}
	evmResult.Result.Events = ctx.EventManager().Events().ToABCIEvents()
	// set the events to the result
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String(), GasUsed: st.GasUsed, Events: ctx.EventManager().Events().ToABCIEvents()}
}
