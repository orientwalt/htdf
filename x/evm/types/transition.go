package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	"github.com/orientwalt/htdf/params"
	"github.com/orientwalt/htdf/x/evm/core/vm"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/orientwalt/htdf/app/protocol"
	appParams "github.com/orientwalt/htdf/params"
	apptypes "github.com/orientwalt/htdf/types"
	sdk "github.com/orientwalt/htdf/types"
	sdkerrors "github.com/orientwalt/htdf/types/errors"
	"github.com/orientwalt/htdf/x/auth"
	vmcore "github.com/orientwalt/htdf/x/evm/core"
	tmtypes "github.com/tendermint/tendermint/types"
)

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation, homestead bool) (uint64, error) {
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
		if (math.MaxUint64-gas)/ethparams.TxDataNonZeroGas < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * ethparams.TxDataNonZeroGas

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
	stateDB     *CommitStateDB //vm.StateDB
	evm         *vm.EVM
	////
	txhash   *common.Hash
	simulate bool // i.e CheckTx execution
	////
	msg              MsgSend
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
func NewStateTransition(ctx sdk.Context, msg MsgSend) *StateTransition {

	st := &StateTransition{
		gpGasWanted:      new(ethcore.GasPool).AddGas(msg.GasWanted),
		msg:              msg,
		ContractCreation: true,
		ContractAddress:  nil,
		sender:           sdk.ToEthAddress(msg.From),
		amount:           big.NewInt(0),
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

	if st.ContractCreation {
		st.amount = msg.Amount.AmountOf(sdk.DefaultDenom).BigInt()
	}

	payload, err := hex.DecodeString(msg.Data)
	if err != nil {
		return nil
	}
	st.payload = payload
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

func (st *StateTransition) newEVM(ctx sdk.Context, stateDB vm.StateDB) *vm.EVM {
	// Create context for evm
	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	evmCtx := vmcore.NewEVMContext(st.msg, &st.sender, uint64(ctx.BlockHeight()))
	evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	st.evm = evm
	return evm
}

// TransitionDb will transition the state by applying the current transaction and
// returning the evm execution result.
// NOTE: State transition checks are run during AnteHandler execution.
func (st *StateTransition) TransitionDb(ctx sdk.Context, accountKeeper auth.AccountKeeper, feeCollectionKeeper auth.FeeCollectionKeeper) (*ExecutionResult, error) {
	cost, err := IntrinsicGas(st.payload, st.ContractCreation, true)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "invalid intrinsic gas for transaction")
	}
	err = st.useGas(cost)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "useGas error")
	}
	// This gas limit the the transaction gas limit with intrinsic gas subtracted
	gasLimit := st.gasLimit - ctx.GasMeter().GasConsumed()

	stateDB, err := NewCommitStateDB(ctx, &accountKeeper, protocol.KeyStorage, protocol.KeyCode)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "newStateDB error")
	}

	st.stateDB = stateDB
	evm := st.newEVM(ctx, stateDB)

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
	switch st.ContractCreation {
	case true:
		ret, contractAddress, leftOverGas, err = evm.Create(senderRef, st.payload, gasLimit, st.amount)
		recipientLog = fmt.Sprintf("contract address %s", contractAddress.String())
	default:
		// Increment the nonce for the next transaction	(just for evm state transition)
		stateDB.SetNonce(st.sender, stateDB.GetNonce(st.sender)+1)
		ret, leftOverGas, err = evm.Call(senderRef, *st.recipient, st.payload, gasLimit, st.amount)
		recipientLog = fmt.Sprintf("recipient address %s", st.recipient.String())
	}

	st.GasUsed = gasLimit - leftOverGas

	if err != nil {
		st.GasUsed = st.gasLimit //? this waste-all part is still necessary
		// evmOutput = fmt.Sprintf("evm Create error|err=%s\n", err)
	} else {
		st.refundGas()
		st.GasUsed = st.gasUsed()
		st.ContractAddress = &contractAddress
		// evmOutput = sdk.ToAppAddress(contractAddress).String() //for contract creation
		// evmOutput = hex.EncodeToString(ret)//for contract call
	}

	if err != nil {
		// Consume gas before returning
		// ctx.GasMeter().ConsumeGas(st.GasUsed, "evm execution consumption")
		return nil, err
	}

	// Resets nonce to value pre state transition
	stateDB.SetNonce(st.sender, currentNonce)

	// Generate bloom filter to be saved in tx receipt data
	bloomInt := big.NewInt(0)

	var (
		bloomFilter ethtypes.Bloom
		logs        []*ethtypes.Log
	)

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)
	st.txhash = &ethHash

	if st.txhash != nil && !st.simulate {
		logs, err := stateDB.GetLogs(*st.txhash)
		if err != nil {
			return nil, err
		}
		bloomInt = ethtypes.LogsBloom(logs)
		bloomFilter = ethtypes.BytesToBloom(bloomInt.Bytes())
	}

	if !st.simulate {
		gasUsed := new(big.Int).Mul(new(big.Int).SetUint64(st.GasUsed), st.gasPrice)
		feeCollectionKeeper.AddCollectedFees(ctx, sdk.Coins{sdk.NewCoin(sdk.DefaultDenom, sdk.NewIntFromBigInt(gasUsed))})
		// Finalise state if not a simulated transaction
		// TODO: change to depend on config
		if _, err := stateDB.Commit(false); err != nil {
			return nil, err
		}
	}

	// Encode all necessary data into slice of bytes to return in sdk result
	resultData := ResultData{
		Bloom:  bloomFilter,
		Logs:   logs,
		Ret:    ret,
		TxHash: *st.txhash,
	}

	if st.ContractCreation {
		resultData.ContractAddress = contractAddress
	}

	resBz, err := EncodeResultData(resultData)
	if err != nil {
		return nil, err
	}

	resultLog := fmt.Sprintf(
		"executed EVM state transition; sender address %s; %s", st.sender.String(), recipientLog,
	)

	executionResult := &ExecutionResult{
		Logs:  logs,
		Bloom: bloomInt,
		Result: &sdk.Result{
			Data: resBz,
			Log:  resultLog,
		},
		GasInfo: GasInfo{
			GasConsumed: st.GasUsed,
			GasLimit:    gasLimit,
			GasRefunded: leftOverGas,
		},
	}

	// TODO: Refund unused gas here, if intended in future

	// Consume gas from evm execution
	// Out of gas check does not need to be done here since it is done within the EVM execution
	// ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(gasConsumed, "EVM execution consumption")

	return executionResult, nil
}
