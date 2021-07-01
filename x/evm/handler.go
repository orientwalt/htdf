package evm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/orientwalt/htdf/app/protocol"
	"github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/evm/types"
	log "github.com/sirupsen/logrus"
)

const ZeroBlockHash = "0000000000000000000000000000000000000000000000000000000000000000"

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

// BlockchainContext for evm opcodes, such as BLOCKTIME and BLOCKHASH
type BlockchainContext struct {
	blockTime     time.Time
	parentHash    common.Hash
	keeper        Keeper // evm.Keeper to get blockhash by block number
	ctx           sdk.Context
	blockGasLimit uint64 // block
}

func NewBlockchainContext(ctx sdk.Context, keeper Keeper, blockGasLimit uint64) *BlockchainContext {
	return &BlockchainContext{
		ctx:           ctx,
		blockTime:     ctx.BlockHeader().Time,
		parentHash:    common.BytesToHash(ctx.BlockHeader().LastBlockId.Hash),
		keeper:        keeper,
		blockGasLimit: blockGasLimit,
	}
}

// GetHeader to compatible with ethereum.
// In acutually, we could get blockhash by blocknumber directly, but we don't do that.
func (self BlockchainContext) GetHeader(_ common.Hash, number uint64) *ethtypes.Header {

	preBlockNumber := int64(number) - 1
	bzParentHash, ok := self.keeper.GetBlockHashByNumber(self.ctx, preBlockNumber)
	if !ok {
		logger().Errorf("GetBlockHashByNumber(%d) error:", preBlockNumber)
		bzParentHash = []byte{}
	}

	return &ethtypes.Header{
		ParentHash: common.BytesToHash(bzParentHash),
		Difficulty: big.NewInt(1),
		Number:     big.NewInt(int64(number)),
		GasLimit:   self.blockGasLimit,
		GasUsed:    0,
		Extra:      nil,

		// NOTE: Time should be deterministic, and should only used by
		//       opCode 'BLOCKTIME' which could only get current block time.
		Time: uint64(self.blockTime.Unix()),
	}
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

	// st.StateDB = k.CommitStateDB.WithContext(ctx)

	logger().Debugf("==========HandleMsgEthereumTx: %s\n", (*st.TxHash).String())
	logger().Debugf("handler:*st.TxHash[%s]\n", (*st.TxHash).String())
	// Prepare db for logs
	// TODO: block hash,
	// ? yqq,  BlockHash is belong to current block or previous block ?
	// Current blockhash hasn't generated until all transactions be packaged into block.
	k.CommitStateDB.Prepare(*st.TxHash, common.Hash{}, k.TxCount)
	k.TxCount++

	// yqq, 2021-05-17, fix issue #70
	if st.StateDB == nil {
		logger().Debugln("Creating CommitStateDB")
		st.StateDB, err = types.NewCommitStateDB(ctx, &k.AccountKeeper, protocol.KeyStorage, protocol.KeyCode)
		if err != nil {
			panic(err)
		}
		st.StateDB.Prepare(*st.TxHash, common.Hash{}, k.TxCount)
	}

	chainCtx := NewBlockchainContext(ctx, k, params.TxGasLimit)
	evmResult, err := st.TransitionDb(ctx, chainCtx, k.AccountKeeper, k.FeeCollectionKeeper)

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

// do switch
func BeginBlocker(ctx sdk.Context, evmk Keeper) {

	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "htdf/x/evm"))
	logger := ctx.Logger()

	logger.Info("=========evm.BeginBlocker ==========")

	var preBlockNumber int64
	var preBlockHash []byte
	preBlockNumber = ctx.BlockHeight() - 1
	logger.Info(fmt.Sprintf("preBlockNumber: %v", preBlockNumber))
	// Block 1 is genesis block. Becasue we start from block 1, instead of block 0.
	logger.Info(fmt.Sprintf("initialHeight: %v", ctx.InitialHeight()))
	logger.Info(fmt.Sprintf("blockHeight: %v", ctx.BlockHeight()))
	if ctx.BlockHeight() == ctx.InitialHeight() {
		// compatible with core.vm.instructions.opBlockhash
		preBlockHash, _ = hex.DecodeString(ZeroBlockHash)
	} else if ctx.BlockHeight() > ctx.InitialHeight() {
		preBlockHash = ctx.BlockHeader().LastBlockId.Hash
	}
	logger.Info(fmt.Sprintf("preBlockHash %v", hex.EncodeToString(preBlockHash)))
	evmk.SetBlockNumberToHash(ctx, preBlockNumber, preBlockHash)
	evmk.SetBlockHashToNumber(ctx, preBlockHash, preBlockNumber)
}
