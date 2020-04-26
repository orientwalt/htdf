package htdfservice

import (
	"fmt"
	"os"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/orientwalt/htdf/codec"
	vmcore "github.com/orientwalt/htdf/evm/core"
	"github.com/orientwalt/htdf/evm/state"
	"github.com/orientwalt/htdf/evm/vm"
	appParams "github.com/orientwalt/htdf/params"
	"github.com/orientwalt/htdf/types"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	log "github.com/sirupsen/logrus"
	abci "github.com/tendermint/tendermint/abci/types"
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

// Query endpoints supported by the core querier
const (
	QueryContract = "contract"
	//
	ZeroAddress = "0000000000000000000000000000000000000000"
	//
	TxGasLimit = 100000
)

// NewQuerier returns a htdfservice Querier handler.
func NewQuerier(accountKeeper auth.AccountKeeper, keyStorage *sdk.KVStoreKey, keyCode *sdk.KVStoreKey) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryContract:
			return queryContract(ctx, req, accountKeeper, keyStorage, keyCode)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

// defines the params for query: "custom/hs/contract"
// junying-todo, 2020-03-30
type QueryContractParams struct {
	Address sdk.AccAddress
	Code    string
}

func NewQueryContractParams(addr sdk.AccAddress, code string) QueryContractParams {
	return QueryContractParams{
		Address: addr,
		Code:    code,
	}
}

//
type MsgTest struct {
	From sdk.AccAddress
}

func NewMsgTest(addr sdk.AccAddress) MsgTest {
	return MsgTest{
		From: addr,
	}
}
func (msg MsgTest) FromAddress() common.Address {
	return types.ToEthAddress(msg.From)
}

// junying-todo, 2020-03-30
func queryContract(ctx sdk.Context, req abci.RequestQuery, accountKeeper auth.AccountKeeper, keyStorage *sdk.KVStoreKey, keyCode *sdk.KVStoreKey) ([]byte, sdk.Error) {
	var params QueryContractParams
	if err := codec.New().UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	//
	stateDB, err := state.NewCommitStateDB(ctx, &accountKeeper, keyStorage, keyCode)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("newStateDB error: %s", err))
	}
	//
	contractAddress := sdk.ToEthAddress(params.Address)

	inputCode, err := hex.DecodeString(params.Code)
	log.Debugf("inputCode=%s\n", hex.EncodeToString(inputCode))
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("DecodeString error: %s", err))
	}
	//

	msg := NewMsgTest(params.Address)
	fromAddress := msg.FromAddress()
	//
	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	evmCtx := vmcore.NewEVMContext(msg, &fromAddress, uint64(ctx.BlockHeight()))
	evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)
	contractRef := vm.AccountRef(fromAddress)
	//
	outputs, _, err := evm.StaticCall(contractRef, contractAddress, inputCode, TxGasLimit)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("evm call error|err=: %s", err))
	}
	log.Debugf("itrsGas|gas=%d\n", outputs)
	//
	bz, err := codec.MarshalJSONIndent(codec.New(), hex.EncodeToString(outputs))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
