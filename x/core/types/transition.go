package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	"github.com/orientwalt/htdf/evm/state"
	evmstate "github.com/orientwalt/htdf/evm/state"
	"github.com/orientwalt/htdf/evm/vm"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/orientwalt/htdf/app/protocol"
	vmcore "github.com/orientwalt/htdf/evm/core"
	appParams "github.com/orientwalt/htdf/params"
	apptypes "github.com/orientwalt/htdf/types"
	sdk "github.com/orientwalt/htdf/types"
	sdkerrors "github.com/orientwalt/htdf/types/errors"
	"github.com/orientwalt/htdf/x/auth"
)

//
type StateTransition struct {
	gpGasWanted *ethcore.GasPool
	initialGas  uint64
	stateDB     *evmstate.CommitStateDB //vm.StateDB
	evm         *vm.EVM
	////
	txhash   *common.Hash
	simulate bool // i.e CheckTx execution
	////
	msg              MsgSend
	sender           common.Address
	recipient        *common.Address
	contractCreation bool
	payload          []byte
	amount           *big.Int
	gasLimit         uint64   //unit: gallon
	gasPrice         *big.Int //unit: satoshi/gallon
}

///
func NewStateTransition(ctx sdk.Context, msg MsgSend) *StateTransition {

	st := &StateTransition{
		gpGasWanted:      new(ethcore.GasPool).AddGas(msg.GasWanted),
		msg:              msg,
		contractCreation: true,
		sender:           sdk.ToEthAddress(msg.From),
		amount:           big.NewInt(0),
		gasLimit:         msg.GasWanted,
		gasPrice:         big.NewInt(int64(msg.GasPrice)),
		simulate:         ctx.IsCheckTx(),
	}
	var to common.Address
	if msg.To != nil {
		to = common.BytesToAddress(msg.To.Bytes())
		st.recipient = &to
		st.contractCreation = false
	}

	if st.contractCreation {
		st.amount = msg.Amount.AmountOf(sdk.DefaultDenom).BigInt()
	}

	payload, err := hex.DecodeString(msg.Data)
	if err != nil {
		return nil
	}
	return st

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
	fmt.Printf("msgGasVal=%s\n", msgGasVal.String())

	if st.stateDB.GetBalance(eaSender).Cmp(msgGasVal) < 0 {
		return errors.New("insufficient balance for gas")
	}

	// try call subGas method, to check gas limit
	if err := st.gpGasWanted.SubGas(st.msg.GasWanted); err != nil {
		fmt.Printf("SubGas error|err=%s\n", err)
		return err
	}

	st.stateDB.SubBalance(eaSender, msgGasVal)
	return nil
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.stateDB.GetRefund() {
		refund = st.stateDB.GetRefund()
	}

	st.gasLimit += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	eaSender := apptypes.ToEthAddress(st.msg.From)

	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gasLimit), st.gasPrice)
	st.stateDB.AddBalance(eaSender, remaining)

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

// GasInfo returns the gas limit, gas consumed and gas refunded from the EVM transition
// execution
type GasInfo struct {
	GasLimit    uint64
	GasConsumed uint64
	GasRefunded uint64
}

// ExecutionResult represents what's returned from a transition
type ExecutionResult struct {
	Logs    []*ethtypes.Log
	Bloom   *big.Int
	Result  *sdk.Result
	GasInfo GasInfo
}

func (st StateTransition) newEVM(ctx sdk.Context, stateDB *vm.StateDB) *vm.EVM {
	// Create context for evm
	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	evmCtx := vmcore.NewEVMContext(st.msg, &st.sender, uint64(ctx.BlockHeight()))
	return vm.NewEVM(evmCtx, stateDB, config, vmConfig)
}

// TransitionDb will transition the state by applying the current transaction and
// returning the evm execution result.
// NOTE: State transition checks are run during AnteHandler execution.
func (st StateTransition) TransitionDb(ctx sdk.Context, accountKeeper auth.AccountKeeper, feeCollectionKeeper auth.FeeCollectionKeeper) (*ExecutionResult, error) {
	cost, err := IntrinsicGas(st.payload, st.contractCreation, true)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "invalid intrinsic gas for transaction")
	}

	// This gas limit the the transaction gas limit with intrinsic gas subtracted
	gasLimit := st.gasLimit - ctx.GasMeter().GasConsumed()

	stateDB, err := state.NewCommitStateDB(ctx, &accountKeeper, protocol.KeyStorage, protocol.KeyCode)
	if err != nil {
		evmOutput := fmt.Sprintf("newStateDB error\n")
		return nil, sdkerrors.Wrapf(err, "newStateDB error")
	}

	st.stateDB = stateDB
	st.evm = newEVM(ctx, stateDB)

	var (
		ret             []byte
		leftOverGas     uint64
		contractAddress common.Address
		recipientLog    string
		senderRef       = vm.AccountRef(st.sender)
	)

	// Get nonce of account outside of the EVM
	currentNonce := st.stateDB.GetNonce(st.sender)
	// Set nonce of sender account before evm state transition for usage in generating Create address
	// st.stateDB.SetNonce(st.sender, st.AccountNonce)

	// create contract or execute call
	switch st.contractCreation {
	case true:
		ret, contractAddress, leftOverGas, err = evm.Create(senderRef, st.payload, gasLimit, st.amount)
		recipientLog = fmt.Sprintf("contract address %s", contractAddress.String())
	default:
		// Increment the nonce for the next transaction	(just for evm state transition)
		stateDB.SetNonce(st.sender, stateDB.GetNonce(st.sender)+1)
		ret, leftOverGas, err = evm.Call(senderRef, *st.recipient, st.payload, gasLimit, st.amount)
		recipientLog = fmt.Sprintf("recipient address %s", st.recipient.String())
	}

	gasConsumed := gasLimit - leftOverGas

	if err != nil {
		gasUsed = msg.GasWanted
		evmOutput = fmt.Sprintf("evm Create error|err=%s\n", err)
	} else {
		st.gas = gasLeftover
		st.refundGas()
		gasUsed = st.gasUsed()
		evmOutput = sdk.ToAppAddress(contractAddress).String() //for contract creation
		// evmOutput = hex.EncodeToString(ret)//for contract call
	}

	if err != nil {
		// Consume gas before returning
		ctx.GasMeter().ConsumeGas(gasConsumed, "evm execution consumption")
		return nil, err
	}

	// Resets nonce to value pre state transition
	st.stateDB.SetNonce(st.sender, currentNonce)

	// Generate bloom filter to be saved in tx receipt data
	bloomInt := big.NewInt(0)

	var (
		bloomFilter ethtypes.Bloom
		logs        []*ethtypes.Log
	)

	if st.txhash != nil && !st.simulate {
		logs = stateDB.GetLogs(*st.txhash)
		bloomInt = ethtypes.LogsBloom(logs)
		bloomFilter = ethtypes.BytesToBloom(bloomInt.Bytes())
	}

	if !st.simulate {
		gasUsed := new(big.Int).Mul(new(big.Int).SetUint64(gasused), gasprice)
		feeCollectionKeeper.AddCollectedFees(ctx, sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewIntFromBigInt(gasUsed))})
		// Finalise state if not a simulated transaction
		// TODO: change to depend on config
		if _, err := st.stateDB.Commit(false); err != nil {
			return nil, err
		}
	}

	// Encode all necessary data into slice of bytes to return in sdk result
	resultData := ResultData{
		Bloom:  bloomFilter,
		Logs:   logs,
		Ret:    ret,
		TxHash: *st.TxHash,
	}

	if st.contractCreation {
		resultData.ContractAddress = contractAddress
	}

	resBz, err := EncodeResultData(resultData)
	if err != nil {
		return nil, err
	}

	resultLog := fmt.Sprintf(
		"executed EVM state transition; sender address %s; %s", st.Sender.String(), recipientLog,
	)

	executionResult := &ExecutionResult{
		Logs:  logs,
		Bloom: bloomInt,
		Result: &sdk.Result{
			Data: resBz,
			Log:  resultLog,
		},
		GasInfo: GasInfo{
			GasConsumed: gasConsumed,
			GasLimit:    gasLimit,
			GasRefunded: leftOverGas,
		},
	}

	// TODO: Refund unused gas here, if intended in future

	// Consume gas from evm execution
	// Out of gas check does not need to be done here since it is done within the EVM execution
	ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(gasConsumed, "EVM execution consumption")

	return executionResult, nil
}
