package evm

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/orientwalt/htdf/types"
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
// HandleMsgEthereumTx, HandleMsgAdd upgraded to EVM version
// commented by junying, 2019-08-21
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgEthereumTx:
			return HandleMsgEthereumTx(ctx, k, msg)
		default:
			return HandleUnknownMsg(msg)
		}
	}

}

// junying-todo, 2019-08-26
func HandleUnknownMsg(msg sdk.Msg) sdk.Result {
	var sendTxResp types.SendTxResp
	logger().Debugf("msgType error|mstType=%v\n", msg.Type())
	sendTxResp.ErrCode = sdk.ErrCodeParam
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
}

// junying-todo, 2019-08-26
func HandleMsgEthereumTx(ctx sdk.Context,
	k Keeper,
	msg types.MsgEthereumTx) sdk.Result {
	// initialize
	// var sendTxResp types.SendTxResp

	st, err := types.NewStateTransition(ctx, msg)
	if st == nil {
		return sdk.Result{Code: sdk.ErrCodeParsing, Log: fmt.Sprintf("%s\n", err)}
	}

	st.StateDB = k.CommitStateDB.WithContext(ctx)

	logger().Infof("==========HandleMsgEthereumTx: %s\n", (*st.TxHash).String())
	logger().Debugf("handler:*st.TxHash[%s]\n", (*st.TxHash).String())
	// Prepare db for logs
	// TODO: block hash
	k.CommitStateDB.Prepare(*st.TxHash, common.Hash{}, k.TxCount)
	k.TxCount++

	evmResult, err := st.TransitionDb(ctx, k.AccountKeeper, k.FeeCollectionKeeper)

	//
	if evmResult == nil {
		return sdk.Result{Code: sdk.ErrCodeIntrinsicGas, Log: fmt.Sprintf("err: %s\n", err), GasUsed: st.GasUsed}
	}
	//
	if err != nil {
		if st.ContractCreation {
			evmResult.Result.Code = sdk.ErrCodeCreateContract
		} else {
			evmResult.Result.Code = sdk.ErrCodeOpenContract
		}
		return *evmResult.Result
	}

	logger().Debugf("handler:evmResult.Log[%v]\n", evmResult.Logs)

	// update block bloom filter
	k.Bloom.Or(k.Bloom, evmResult.Bloom)

	// update transaction logs in KVStore
	err = k.SetLogs(ctx, *st.TxHash, evmResult.Logs)
	if err != nil {
		panic(err)
	}

	// log successful execution
	k.Logger(ctx).Info(evmResult.Result.Log)

	// evmResult.Result.Code = sendTxResp.ErrCode

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMsgEthereumTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, st.GetSender().String()),
		),
	})

	if msg.To != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMsgEthereumTx,
				sdk.NewAttribute(types.AttributeKeyRecipient, st.GetRecipient().String()),
			),
		)
	}
	// evmResult.Result.Log = sendTxResp.String()
	evmResult.Result.Events = ctx.EventManager().Events().ToABCIEvents()

	// set the events to the result
	return *evmResult.Result //sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String(), GasUsed: st.GasUsed, Events: ctx.EventManager().Events().ToABCIEvents(), Data: evmResult.}
}
