package htdfservice

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	vmcore "github.com/orientwalt/htdf/evm/core"
	"github.com/orientwalt/htdf/evm/state"
	"github.com/orientwalt/htdf/evm/vm"
	appParams "github.com/orientwalt/htdf/params"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
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

//
type SendTxResp struct {
	ErrCode         sdk.CodeType `json:"code"`
	ErrMsg          string       `json:"message"`
	ContractAddress string       `json:"contract_address"`
	EvmOutput       string       `json:"evm_output"`
}

//
func (rsp SendTxResp) String() string {
	rsp.ErrMsg = sdk.GetErrMsg(rsp.ErrCode)
	data, _ := json.Marshal(&rsp)
	return string(data)
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
		case MsgSend:
			return HandleMsgSend(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		default:
			return HandleUnknownMsg(msg)
		}
	}

}

// junying-todo, 2019-08-26
func HandleUnknownMsg(msg sdk.Msg) sdk.Result {
	var sendTxResp SendTxResp
	log.Debugf("msgType error|mstType=%v\n", msg.Type())
	sendTxResp.ErrCode = sdk.ErrCode_Param
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
}

// junying-todo, 2019-08-26
func HandleMsgSend(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg MsgSend) sdk.Result {
	// initialize
	var sendTxResp SendTxResp
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
		evmOutput, gasUsed, err = HandleCreateContract(ctx, accountKeeper, feeCollectionKeeper, keyStorage, keyCode, msg)
		if err != nil {
			sendTxResp.ErrCode = sdk.ErrCode_CreateContract
		}
		sendTxResp.ContractAddress = evmOutput
	}
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String(), GasUsed: gasUsed}
}

//
func FeeCollecting(ctx sdk.Context,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	stateDB *state.CommitStateDB,
	gasused uint64,
	gasprice *big.Int) {
	gasUsed := new(big.Int).Mul(new(big.Int).SetUint64(gasused), gasprice)
	log.Debugf("FeeCollecting:gasUsed=%s\n", gasUsed.String())
	feeCollectionKeeper.AddCollectedFees(ctx, sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewIntFromBigInt(gasUsed))})
	stateDB.Commit(false)
	log.Debugln("FeeCollecting:stateDB commited!")
}

// junying-todo, 2019-08-26
func HandleOpenContract(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg MsgSend) (evmOutput string, gasUsed uint64, err error) {

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
	log.Infof("===========> blktime: %s", blockTime.Format(time.RFC3339Nano))

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
	st := NewStateTransition(evm, msg, stateDB)

	log.Debugf("gasPrice=%d|gasWanted=%d\n", msg.GasPrice, msg.GasWanted)

	// commented by junying, 2019-08-22
	// subtract GasWanted*gasprice from sender
	err = st.buyGas()
	if err != nil {
		evmOutput = fmt.Sprintf("buyGas error|err=%s\n", err)
		return
	}

	// Intrinsic gas calc
	// commented by junying, 2019-08-22
	// default non-contract tx gas: 21000
	// default contract tx gas: 53000 + f(tx.data)
	itrsGas, err := IntrinsicGas(inputCode, true)
	log.Debugf("itrsGas|gas=%d\n", itrsGas)
	// commented by junying, 2019-08-22
	// check if tx.gas >= calculated gas
	err = st.useGas(itrsGas)
	if err != nil {
		evmOutput = fmt.Sprintf("useGas error|err=%s\n", err)
		return
	}

	// commented by junying, 2019-08-22
	// 1. cantransfer check
	// 2. create receiver account if no exists
	// 3. execute contract & calculate gas
	log.Debugln(" HandleOpenContract1:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	outputs, gasLeftover, err := evm.Call(contractRef, toAddress, inputCode, st.gas, transferAmount)
	log.Debugln(" HandleOpenContract2:ctx.GasMeter().GasConsumed()", ctx.GasMeter().GasConsumed())
	if err != nil {
		log.Debugf("evm call error|err=%s\n", err)
		// junying-todo, 2019-11-05
		gasUsed = msg.GasWanted
		evmOutput = fmt.Sprintf("evm call error|err=%s\n", err)
	} else {
		st.gas = gasLeftover
		// junying-todo, 2019-08-22
		// refund(add) remaining to sender
		st.refundGas()
		log.Debugf("gasUsed=%d\n", st.gasUsed())
		gasUsed = st.gasUsed()
		evmOutput = hex.EncodeToString(outputs)
	}
	FeeCollecting(ctx, feeCollectionKeeper, stateDB, gasUsed, st.gasPrice)
	return
}

// junying-todo, 2019-08-26
func HandleCreateContract(ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	feeCollectionKeeper auth.FeeCollectionKeeper,
	keyStorage *sdk.KVStoreKey,
	keyCode *sdk.KVStoreKey,
	msg MsgSend) (evmOutput string, gasUsed uint64, err error) {

	stateDB, err := state.NewCommitStateDB(ctx, &accountKeeper, keyStorage, keyCode)
	if err != nil {
		evmOutput = fmt.Sprintf("newStateDB error\n")
		return
	}
	fromAddress := sdk.ToEthAddress(msg.From)
	toAddress := sdk.ToEthAddress(msg.To)

	log.Debugf("fromAddr|appFormat=%s|ethFormat=%s|\n", msg.From.String(), fromAddress.String())
	log.Debugf("toAddress|appFormat=%s|ethFormat=%s|\n", msg.To.String(), toAddress.String())
	log.Debugf("fromAddress|Balance=%v\n", stateDB.GetBalance(fromAddress))

	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	log.Debugf("fromAddress|nonce=%d\n", stateDB.GetNonce(fromAddress))

	blockTime := ctx.BlockHeader().Time
	log.Infof("===========> blktime: %s", blockTime.Format(time.RFC3339Nano))

	evmCtx := vmcore.NewEVMContext(msg, &fromAddress, uint64(ctx.BlockHeight()), blockTime)
	evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	contractRef := vm.AccountRef(fromAddress)

	log.Debugf("blockHeight=%d|IsHomestead=%v|IsDAOFork=%v|IsEIP150=%v|IsEIP155=%v|IsEIP158=%v|IsByzantium=%v\n", ctx.BlockHeight(), evm.ChainConfig().IsHomestead(evm.BlockNumber),
		evm.ChainConfig().IsDAOFork(evm.BlockNumber), evm.ChainConfig().IsEIP150(evm.BlockNumber),
		evm.ChainConfig().IsEIP155(evm.BlockNumber), evm.ChainConfig().IsEIP158(evm.BlockNumber),
		evm.ChainConfig().IsByzantium(evm.BlockNumber))

	inputCode, err := hex.DecodeString(msg.Data)
	if err != nil {
		evmOutput = fmt.Sprintf("DecodeString error\n")
		return
	}

	log.Debugf("inputCode=%s\n", hex.EncodeToString(inputCode))

	st := NewStateTransition(evm, msg, stateDB)

	log.Debugf("gasPrice=%d|GasWanted=%d\n", msg.GasPrice, msg.GasWanted)

	err = st.buyGas()
	if err != nil {
		evmOutput = fmt.Sprintf("buyGas error|err=%s\n", err)
		return
	}

	//Intrinsic gas calc
	itrsGas, err := IntrinsicGas(inputCode, true)
	log.Debugf("itrsGas|gas=%d\n", itrsGas)
	err = st.useGas(itrsGas)
	if err != nil {
		evmOutput = fmt.Sprintf("useGas error|err=%s\n", err)
		return
	}

	_, contractAddr, gasLeftover, err := evm.Create(contractRef, inputCode, st.gas, big.NewInt(0))
	if err != nil {
		log.Debugf("evm Create error|err=%s\n", err)
		// junying-todo, 2019-11-05
		gasUsed = msg.GasWanted
		evmOutput = fmt.Sprintf("evm Create error|err=%s\n", err)
	} else {
		st.gas = gasLeftover
		st.refundGas()
		gasUsed = st.gasUsed()
		evmOutput = sdk.ToAppAddress(contractAddr).String()
	}

	log.Debugf("Create contract ok,contractAddr|appFormat=%s|ethFormat=%s\n", sdk.ToAppAddress(contractAddr).String(), contractAddr.String())
	FeeCollecting(ctx, feeCollectionKeeper, stateDB, gasUsed, st.gasPrice)
	return
}
