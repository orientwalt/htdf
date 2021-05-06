package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	"github.com/orientwalt/htdf/app/protocol"
	"github.com/orientwalt/htdf/params"
	"github.com/orientwalt/htdf/x/evm/core/vm"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	appParams "github.com/orientwalt/htdf/params"
	apptypes "github.com/orientwalt/htdf/types"
	sdk "github.com/orientwalt/htdf/types"
	sdkerrors "github.com/orientwalt/htdf/types/errors"
	"github.com/orientwalt/htdf/x/auth"
	vmcore "github.com/orientwalt/htdf/x/evm/core"
	tmtypes "github.com/tendermint/tendermint/types"

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

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation, homestead bool, isEIP2028 bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation && homestead {
		gas = params.DefaultMsgGasContractCreation // 53000 -> 60000
	} else {
		gas = params.DefaultMsgGas // 21000 -> 30000
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		// if (math.MaxUint64-gas)/ethparams.TxDataNonZeroGas < nz {
		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := ethparams.TxDataNonZeroGasFrontier
		if isEIP2028 {
			nonZeroGas = ethparams.TxDataNonZeroGasEIP2028
		}

		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * nonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/ethparams.TxDataZeroGas < z {
			return 0, vm.ErrOutOfGas
		}
		gas += z * ethparams.TxDataZeroGas
	}
	return gas, nil
}

//
type StateTransition struct {
	gpGasWanted *ethcore.GasPool
	initialGas  uint64
	StateDB     *CommitStateDB //vm.StateDB
	evm         *vm.EVM
	////
	TxHash   *common.Hash
	simulate bool // i.e CheckTx execution
	////
	msg              MsgEthereumTx
	sender           common.Address
	recipient        *common.Address
	ContractCreation bool
	ContractAddress  *common.Address
	payload          []byte
	amount           *big.Int
	gasLimit         uint64   //unit: gallon
	gasPrice         *big.Int //unit: satoshi/gallon
	GasUsed          uint64
}

///
func NewStateTransition(ctx sdk.Context, msg MsgEthereumTx) (*StateTransition, error) {

	st := &StateTransition{
		gpGasWanted:      new(ethcore.GasPool).AddGas(msg.GasWanted),
		msg:              msg,
		ContractCreation: true,
		ContractAddress:  nil,
		sender:           sdk.ToEthAddress(msg.From),
		amount:           msg.Amount.AmountOf(sdk.DefaultDenom).BigInt(),
		gasLimit:         msg.GasWanted,
		gasPrice:         big.NewInt(int64(msg.GasPrice)),
		GasUsed:          0,
		simulate:         ctx.IsCheckTx(),
	}
	var to common.Address
	if msg.To != nil {
		to = common.BytesToAddress(msg.To.Bytes())
		st.recipient = &to
		st.ContractCreation = false
	}

	payload, err := hex.DecodeString(msg.Data)
	if err != nil {
		return nil, err
	}

	st.payload = payload

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)
	st.TxHash = &ethHash

	return st, nil

}

//
func (st *StateTransition) GetSender() common.Address {
	return st.sender
}

//
func (st *StateTransition) GetRecipient() common.Address {
	return *st.recipient
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gasLimit < amount {
		return vm.ErrOutOfGas
	}
	st.gasLimit -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	st.gasLimit = st.msg.GasWanted
	st.initialGas = st.gasLimit
	fmt.Printf("msgGas=%d\n", st.initialGas)

	eaSender := apptypes.ToEthAddress(st.msg.From)

	msgGasVal := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.GasWanted), st.gasPrice)
	// fmt.Printf("msgGasVal=%s\n", msgGasVal.String())
	// fmt.Printf("msgGasVal=%s\n", eaSender.String())
	fmt.Printf("st.StateDB.GetBalance(Before)=%v\n", st.StateDB.GetBalance(eaSender))

	if st.StateDB.GetBalance(eaSender).Cmp(msgGasVal) < 0 {
		return errors.New("insufficient balance for gas")
	}

	// try call subGas method, to check gas limit
	if err := st.gpGasWanted.SubGas(st.msg.GasWanted); err != nil {
		fmt.Printf("SubGas error|err=%s\n", err)
		return err
	}

	st.StateDB.SubBalance(eaSender, msgGasVal)
	fmt.Printf("st.StateDB.GetBalance(After)=%v\n", st.StateDB.GetBalance(eaSender))
	return nil
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.StateDB.GetRefund() {
		refund = st.StateDB.GetRefund()
	}

	st.gasLimit += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	eaSender := apptypes.ToEthAddress(st.msg.From)

	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gasLimit), st.gasPrice)
	st.StateDB.AddBalance(eaSender, remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gpGasWanted.AddGas(st.gasLimit)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gasLimit
}

func (st *StateTransition) tokenUsed() uint64 {
	return new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice).Uint64()
}

func (st *StateTransition) newEVM(ctx sdk.Context, chainCtx vmcore.ChainContext, stateDB vm.StateDB) *vm.EVM {
	// Create context for evm
	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	blockCtx := vmcore.NewEVMBlockContext(ctx.BlockHeader(), chainCtx, &st.sender)
	txCtx := vmcore.NewEVMTxContext(st.msg)

	// evmCtx := vmcore.NewEVMContext(st.msg, &st.sender, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time)
	evm := vm.NewEVM(blockCtx, txCtx, stateDB, config, vmConfig)
	st.evm = evm
	return evm
}

// TransitionDb will transition the state by applying the current transaction and
// returning the evm execution result.
// NOTE: State transition checks are run during AnteHandler execution.
func (st *StateTransition) TransitionDb(ctx sdk.Context, chainCtx vmcore.ChainContext, ak auth.AccountKeeper, fck FeeCollectionKeeper) (*ExecutionResult, error) {
	if st.StateDB == nil {
		stateDB, err := NewCommitStateDB(ctx, &ak, protocol.KeyStorage, protocol.KeyCode)
		if err != nil {
			panic(err)
			// return nil, sdkerrors.Wrapf(err, "newStateDB error")
		}
		st.StateDB = stateDB
	}

	stateDB := st.StateDB

	evm := st.newEVM(ctx, chainCtx, stateDB)

	// commented by junying, 2019-08-22
	// subtract GasWanted*gasprice from sender
	err := st.buyGas()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "buyGas error")
	}
	// Intrinsic gas calc
	// commented by junying, 2019-08-22
	// default non-contract tx gas: 21000
	// default contract tx gas: 53000 + f(tx.data)
	logging().Debugf("in TransitionDb\n")
	cost, err := IntrinsicGas(st.payload, st.ContractCreation, true, true)
	logging().Debugf("in TransitionDb:cost[%d]\n", cost)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "invalid intrinsic gas for transaction")
	}
	// commented by junying, 2019-08-22
	// check if tx.gas >= calculated gas
	err = st.useGas(cost)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "useGas error")
	}
	// This gas limit the the transaction gas limit with intrinsic gas subtracted
	gasLimit := st.gasLimit - ctx.GasMeter().GasConsumed()
	logging().Debugf("in TransitionDb:gasLimit[%d]\n", gasLimit)

	var (
		ret             []byte
		leftOverGas     uint64
		contractAddress common.Address
		recipientLog    string
		senderRef       = vm.AccountRef(st.sender)
	)

	// Get nonce of account outside of the EVM
	currentNonce := st.StateDB.GetNonce(st.sender)
	// Set nonce of sender account before evm state transition for usage in generating Create address
	// st.StateDB.SetNonce(st.sender, st.AccountNonce)

	logging().Infof("in TransitionDb:currentNonce[%d]\n", currentNonce)

	// logger().Debugf("in TransitionDb:st[%v]\n", st)
	// logging().Debugln(st.ContractCreation)

	// stateDB.SetNonce(st.sender, currentNonce+1)
	// create contract or execute call
	switch st.ContractCreation {
	case true:
		ret, contractAddress, leftOverGas, err = evm.Create(senderRef, st.payload, gasLimit, st.amount)
		recipientLog = fmt.Sprintf("contract address: %s", contractAddress.String())
		logger().Infof("NEW CREATED CONTRACT ADDRESS:%s", contractAddress.String())
	default:
		if code := evm.StateDB.GetCode(st.GetRecipient()); len(st.payload) > 0 && len(code) == 0 {
			// copy from v1.3.0 2021-04-07
			// added by yqq 2020-12-07
			// To fix issue #14, we disable transaction which has a not empty `MsgSend.Data`
			// and `MsgSend.To` is not contract address.
			ret, leftOverGas, err = nil, 0, fmt.Errorf("invalid contract address")
			recipientLog = fmt.Sprintf("contract address: %s", st.GetRecipient().String())
		} else {
			// Increment the nonce for the next transaction	(just for evm state transition)
			ret, leftOverGas, err = evm.Call(senderRef, *st.recipient, st.payload, gasLimit, st.amount)
			recipientLog = fmt.Sprintf("recipient address: %s", st.recipient.String())
		}
	}

	recipientLog = fmt.Sprintf("%s, output: %s", recipientLog, hex.EncodeToString(ret))
	logger().Debugf("in TransitionDb:recipientLog[%s]\n", recipientLog)
	logger().Debugf("in TransitionDb:leftOverGas[%d]\n", leftOverGas)
	// st.GasUsed = gasLimit - leftOverGas
	// logging().Debugf("in TransitionDb:st.gasLimit[%d]\n", gasLimit)
	// logging().Debugf("in TransitionDb:st.GasUsed[%d]\n", st.GasUsed)

	txReceiptStatus := uint(0)
	if err != nil {
		st.GasUsed = st.initialGas
		// st.GasUsed = st.gasLimit //? this waste-all part is still necessary
		reason, _ := abi.UnpackRevert(ret)
		recipientLog = fmt.Sprintf("%s, err: %s, reason:%s", recipientLog, err, reason)
		logger().Warnf("evm revert reason: %s", reason)
		txReceiptStatus = 0

		// Consume gas before returning
		// ctx.GasMeter().ConsumeGas(st.GasUsed, "evm execution consumption")
		// return nil, err
	} else {
		st.gasLimit = leftOverGas
		st.refundGas()
		st.GasUsed = st.gasUsed()
		st.ContractAddress = &contractAddress
		logger().Debugf("in TransitionDb:contractAddress[%s]\n", contractAddress.String())
		txReceiptStatus = 1
	}

	// Resets nonce to value pre state transition
	stateDB.SetNonce(st.sender, currentNonce+1)

	if !st.simulate {
		// logging().Debugf("in TransitionDb:st.gasPrice[%d]\n", st.gasPrice)
		gasUsed := new(big.Int).Mul(new(big.Int).SetUint64(st.GasUsed), st.gasPrice)
		fck.AddCollectedFees(ctx, sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewIntFromBigInt(gasUsed))})
		logging().Debugf("in TransitionDb:feeCollectionKeeper.gasUsed[%v]\n", gasUsed)
		// Finalise state if not a simulated transaction
		// TODO: change to depend on config
		if _, _err := stateDB.Commit(false); _err != nil {
			panic(_err)
		}
	}

	// Generate bloom filter to be saved in tx receipt data
	bloomInt := big.NewInt(0)

	var (
		bloomFilter ethtypes.Bloom
		logs        []*ethtypes.Log
	)

	if st.TxHash != nil && !st.simulate {
		var _err error
		logs, _err = stateDB.GetLogs(*st.TxHash)
		if _err != nil {
			return nil, _err
		}
		// bloomInt = ethtypes.LogsBloom(logs)
		var bloom ethtypes.Bloom
		bzBloom := ethtypes.LogsBloom(logs)
		bloom.SetBytes(bzBloom)
		bloomInt = bloom.Big()

		bloomFilter = ethtypes.BytesToBloom(bzBloom)
	}

	resultData := sdk.ResultData{
		Bloom:  bloomFilter,
		Logs:   logs,
		Ret:    ret,
		Status: txReceiptStatus,
		TxHash: *st.TxHash,
	}

	if st.ContractCreation {
		resultData.ContractAddress = contractAddress
	}

	resBz, _err := sdk.EncodeResultData(resultData)
	if _err != nil {
		panic(_err)
	}

	resultLog := fmt.Sprintf(
		"executed EVM state transition; sender address: %s, %s", st.sender.String(), recipientLog,
	)

	executionResult := &ExecutionResult{
		Logs:  logs,
		Bloom: bloomInt,
		Result: &sdk.Result{
			Data:      resBz,
			Log:       resultLog,
			GasUsed:   st.GasUsed,
			GasWanted: st.msg.GasWanted,
		},
		GasInfo: GasInfo{
			GasConsumed: st.GasUsed,
			GasLimit:    gasLimit,
			GasRefunded: leftOverGas,
		},
	}
	logging().Debugf("in TransitionDb:st.GasUsed[%d]\n", st.GasUsed)
	logging().Debugf("in TransitionDb:st.GasRefunded[%d]\n", leftOverGas)
	// TODO: Refund unused gas here, if intended in future

	// Consume gas from evm execution
	// Out of gas check does not need to be done here since it is done within the EVM execution
	// ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(gasConsumed, "EVM execution consumption")

	return executionResult, err
}
