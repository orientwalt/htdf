package htdfservice

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	vmcore "github.com/orientwalt/htdf/evm/core"
	"github.com/orientwalt/htdf/evm/state"
	"github.com/orientwalt/htdf/evm/vm"
	appParams "github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	xcore "github.com/orientwalt/htdf/x/core"
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

// New HTDF Message Handler
// connected to handler.go
// HandleMsgSend, HandleMsgAdd upgraded to EVM version
// commented by junying, 2019-08-21
func NewHandler(accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case xcore.MsgSend:
			return HandleMsgSend(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		default:
			return xcore.HandleUnknownMsg(msg)
		}
	}

}

// copy from x/core/handler.go
func HandleMsgSend(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg xcore.MsgSend) sdk.Result {
	// initialize
	var sendTxResp xcore.SendTxResp
	var gasUsed uint64
	var evmOutput string
	var err error

	if !msg.To.Empty() {
		// open smart contract
		evmOutput, gasUsed, err = HandleOpenContract(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		if err != nil {
			sendTxResp.ErrCode = sdk.ErrCode_OpenContract
		}
		sendTxResp.EvmOutput = evmOutput
	} else {
		// create smart contract
		evmOutput, gasUsed, err = xcore.HandleCreateContract(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		if err != nil {
			sendTxResp.ErrCode = sdk.ErrCode_CreateContract
		}
		sendTxResp.ContractAddress = evmOutput
	}
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String(), GasUsed: gasUsed}
}

// HandleOpenContract copy from x/core/handler, to implement v2 handler  2020-12-07
func HandleOpenContract(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg xcore.MsgSend) (evmOutput string, gasUsed uint64, err error) {

	log.Debugf("Handling MsgSend with No Contract.\n")
	log.Debugln(" HandleOpenContract0:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	stateDB, err := state.NewCommitStateDB(ctx, &accountKeeper, keyStorage, keyCode)
	if err != nil {
		evmOutput = fmt.Sprintf("newStateDB error\n")
		return
	}

	fromAddress := sdk.ToEthAddress(msg.From)
	toAddress := sdk.ToEthAddress(msg.To)

	log.Debugf("fromAddr|appFormat=%s|ethFormat=%s|\n", msg.From.String(), fromAddress.String())
	log.Debugf("toAddress|appFormat=%s|ethFormat=%s|\n", msg.To.String(), toAddress.String())

	log.Debugf("fromAddress|testBalance=%v\n", stateDB.GetBalance(fromAddress))
	log.Debugf("fromAddress|nonce=%d\n", stateDB.GetNonce(fromAddress))

	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	blockTime := ctx.BlockHeader().Time
	log.Info("=== v2.HandleOpenContract  ===")
	log.Infof("blockHeaderTime: %s", blockTime.Format(time.RFC3339Nano))

	evmCtx := vmcore.NewEVMContext(msg, &fromAddress, uint64(ctx.BlockHeight()), blockTime)
	evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	contractRef := vm.AccountRef(fromAddress)

	inputCode, err := hex.DecodeString(msg.Data)
	if err != nil {
		evmOutput = fmt.Sprintf("DecodeString error\n")
		return
	}

	log.Debugf("inputCode=%s\n", hex.EncodeToString(inputCode))

	transferAmount := msg.Amount.AmountOf(sdk.DefaultDenom).BigInt()

	log.Debugf("transferAmount: %d\n", transferAmount)
	st := xcore.NewStateTransition(evm, msg, stateDB)

	log.Debugf("gasPrice=%d|gasWanted=%d\n", msg.GasPrice, msg.GasWanted)

	// commented by junying, 2019-08-22
	// subtract GasWanted*gasprice from sender
	err = st.BuyGas()
	if err != nil {
		evmOutput = fmt.Sprintf("buyGas error|err=%s\n", err)
		return
	}

	// Intrinsic gas calc
	// commented by junying, 2019-08-22
	// default non-contract tx gas: 21000
	// default contract tx gas: 53000 + f(tx.data)
	itrsGas, err := xcore.IntrinsicGas(inputCode, true)
	log.Debugf("itrsGas|gas=%d\n", itrsGas)
	// commented by junying, 2019-08-22
	// check if tx.gas >= calculated gas
	err = st.UseGas(itrsGas)
	if err != nil {
		evmOutput = fmt.Sprintf("useGas error|err=%s\n", err)
		return
	}

	// commented by junying, 2019-08-22
	// 1. cantransfer check
	// 2. create receiver account if no exists
	// 3. execute contract & calculate gas
	log.Debugln(" HandleOpenContract1:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())

	var outputs []byte
	var gasLeftover uint64
	if code := evm.StateDB.GetCode(toAddress); len(inputCode) > 0 && len(code) == 0 {
		// added by yqq 2020-12-07
		// To fix issue #14, we disable transaction which has a not empty `MsgSend.Data`
		// and `MsgSend.To` is not contract address.
		outputs, gasLeftover, err = nil, 0, fmt.Errorf("invalid contract address")
	} else {
		outputs, gasLeftover, err = evm.Call(contractRef, toAddress, inputCode, st.GetGas(), transferAmount)
	}

	log.Debugln(" HandleOpenContract2:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	if err != nil {
		log.Debugf("evm call error|err=%s\n", err)
		// junying-todo, 2019-11-05
		gasUsed = msg.GasWanted
		evmOutput = fmt.Sprintf("evm call error|err=%s\n", err)
	} else {
		st.SetGas(gasLeftover)
		// junying-todo, 2019-08-22
		// refund(add) remaining to sender
		st.RefundGas()
		log.Debugf("gasUsed=%d\n", st.GasUsed())
		gasUsed = st.GasUsed()
		evmOutput = hex.EncodeToString(outputs)
	}
	xcore.FeeCollecting(ctx, feeCollectionKeeper, stateDB, gasUsed, st.GetGasPrice())
	return
}
