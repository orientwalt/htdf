package evm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/magiconair/properties/assert"
	"github.com/orientwalt/htdf/utils"
	"github.com/stretchr/testify/require"

	evmcore "github.com/orientwalt/htdf/x/evm/core"
	"github.com/orientwalt/htdf/x/evm/core/vm"

	//cosmos-sdk
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	evmtypes "github.com/orientwalt/htdf/x/evm/types"
	"github.com/orientwalt/htdf/x/params"

	//tendermint
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	//evm
	newevmtypes "github.com/orientwalt/htdf/x/evm/core/types"

	//ethereum
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	appParams "github.com/orientwalt/htdf/params"

	"testing"
)

// TODO:current test code , is base go-ethereum V1.8.0
//	when this evm package is stable ,need to update to new version, like  V1.8.23

var (
	accKey     = sdk.NewKVStoreKey("acc")
	authCapKey = sdk.NewKVStoreKey("authCapKey")
	fckCapKey  = sdk.NewKVStoreKey("fckCapKey")
	keyParams  = sdk.NewKVStoreKey("params")
	tkeyParams = sdk.NewTransientStoreKey("transient_params")

	storageKey = sdk.NewKVStoreKey("storage")
	codeKey    = sdk.NewKVStoreKey("code")

	testHash    = utils.StringToHash("zhoushx")
	fromAddress = utils.StringToAddress("UserA")
	toAddress   = utils.StringToAddress("UserB")
	amount      = big.NewInt(0)
	nonce       = uint64(0)
	gasLimit    = big.NewInt(100000)
	coinbase    = fromAddress

	nplogger = tmlog.NewNopLogger()
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Message struct {
	to                      *common.Address
	from                    common.Address
	nonce                   uint64
	amount, price, gasLimit *big.Int
	data                    []byte
	checkNonce              bool
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount, gasLimit, price *big.Int, data []byte, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		price:      price,
		gasLimit:   gasLimit,
		data:       data,
		checkNonce: checkNonce,
	}
}

func (m Message) FromAddress() common.Address { return m.from }
func (m Message) To() *common.Address         { return m.to }
func (m Message) GasPrice() *big.Int          { return m.price }
func (m Message) Value() *big.Int             { return m.amount }
func (m Message) Gas() *big.Int               { return m.gasLimit }
func (m Message) Nonce() uint64               { return m.nonce }
func (m Message) Data() []byte                { return m.data }
func (m Message) CheckNonce() bool            { return m.checkNonce }

func loadBin(filename string) []byte {
	code, err := ioutil.ReadFile(filename)
	must(err)
	return hexutil.MustDecode("0x" + string(code))
}
func loadAbi(filename string) abi.ABI {
	abiFile, err := os.Open(filename)
	must(err)
	defer abiFile.Close()
	abiObj, err := abi.JSON(abiFile)
	must(err)
	return abiObj
}

func newTestCodec1() *codec.Codec {
	cdc := codec.New()
	newevmtypes.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func testChainConfig(t *testing.T, evm *vm.EVM) {
	height := big.NewInt(1)

	// yqq , 2021-05-10
	// I had make Berlin fork as default chainconfig in orientwalt/go-ethereum.
	// So, all forks before Berlin always be true.
	assert.Equal(t, evm.ChainConfig().IsHomestead(height), true)
	assert.Equal(t, evm.ChainConfig().IsDAOFork(height), true)
	assert.Equal(t, evm.ChainConfig().IsEIP150(height), true)
	assert.Equal(t, evm.ChainConfig().IsEIP155(height), true)
	assert.Equal(t, evm.ChainConfig().IsEIP158(height), true)
	assert.Equal(t, evm.ChainConfig().IsByzantium(height), true)
	assert.Equal(t, evm.ChainConfig().IsConstantinople(height), true)
	assert.Equal(t, evm.ChainConfig().IsMuirGlacier(height), true)
	assert.Equal(t, evm.ChainConfig().IsPetersburg(height), true)
	assert.Equal(t, evm.ChainConfig().IsIstanbul(height), true)
	assert.Equal(t, evm.ChainConfig().IsBerlin(height), true)
	assert.Equal(t, evm.ChainConfig().IsEWASM(height), false)

}

func TestNewEvm(t *testing.T) {

	//---------------------stateDB test--------------------------------------
	dbName := "htdfnewevmtestdata3"
	dataPath, err := ioutil.TempDir("", dbName)
	require.NoError(t, err)
	db := dbm.NewDB("state", dbm.GoLevelDBBackend, dataPath)

	cdc := newTestCodec1()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(codeKey, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, nil)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, accKey, pk.Subspace(auth.DefaultParamspace), newevmtypes.ProtoBaseAccount)

	cms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(storageKey, sdk.StoreTypeIAVL, nil)

	cms.SetPruning(store.PruneNothing)

	err = cms.LoadLatestVersion()
	require.NoError(t, err)

	ms := cms.CacheMultiStore()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	stateDB, err := evmtypes.NewCommitStateDB(ctx, &ak, storageKey, codeKey)
	must(err)

	fmt.Printf("addr=%s|testBalance=%v\n", fromAddress.String(), stateDB.GetBalance(fromAddress))
	stateDB.AddBalance(fromAddress, big.NewInt(1e18))
	fmt.Printf("addr=%s|testBalance=%v\n", fromAddress.String(), stateDB.GetBalance(fromAddress))

	assert.Equal(t, stateDB.GetBalance(fromAddress).String() == "1000000000000000000", true)

	//---------------------call evm--------------------------------------
	abiFileName := "../../tests/evm/coin/coin_sol_Coin.abi"
	binFileName := "../../tests/evm/coin/coin_sol_Coin.bin"
	data := loadBin(binFileName)

	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	msg := NewMessage(fromAddress, &toAddress, nonce, amount, gasLimit, big.NewInt(0), data, false)
	// evmCtx := evmcore.NewEVMContext(msg, &fromAddress, 1000)
	// evmCtx := evmcore.NewEVMContext(msg, &fromAddress, 1000, ctx.BlockHeader().Time)
	// evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)

	blockCtx := evmcore.NewEVMBlockContext(ctx.BlockHeader(), &evmcore.FakeChainContext{}, &fromAddress)
	txCtx := evmcore.NewEVMTxContext(msg)
	evm := vm.NewEVM(blockCtx, txCtx, stateDB, config, vmConfig)

	contractRef := vm.AccountRef(fromAddress)
	contractCode, contractAddr, gasLeftover, vmerr := evm.Create(contractRef, data, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("BlockNumber=%d|IsEIP158=%v\n", evm.Context.BlockNumber.Uint64(), evm.ChainConfig().IsEIP158(evm.Context.BlockNumber))
	testChainConfig(t, evm)

	fmt.Printf("Create|str_contractAddr=%s|gasLeftOver=%d|contractCode=%x\n", contractAddr.String(), gasLeftover, contractCode)

	stateDB.SetBalance(fromAddress, big.NewInt(0).SetUint64(gasLeftover))
	testBalance := stateDB.GetBalance(fromAddress)
	fmt.Println("after create contract, testBalance =", testBalance)

	abiObj := loadAbi(abiFileName)

	input, err := abiObj.Pack("minter")
	must(err)
	outputs, gasLeftover, vmerr := evm.Call(contractRef, contractAddr, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("smartcontract func, minter|the minter addr=%s\n", common.BytesToAddress(outputs).String())

	sender := common.BytesToAddress(outputs)

	fmt.Printf("sender=%s|fromAddress=%s\n", sender.String(), fromAddress.String())

	if !bytes.Equal(sender.Bytes(), fromAddress.Bytes()) {
		fmt.Println("caller are not equal to minter!!")
		os.Exit(-1)
	}

	senderAcc := vm.AccountRef(sender)

	input, err = abiObj.Pack("mint", sender, big.NewInt(1000000))
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(senderAcc, contractAddr, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("smartcontract func, mint|senderAcc=%s\n", sender.String())

	stateDB.SetBalance(fromAddress, big.NewInt(0).SetUint64(gasLeftover))
	testBalance = evm.StateDB.GetBalance(fromAddress)

	input, err = abiObj.Pack("send", toAddress, big.NewInt(11))
	outputs, gasLeftover, vmerr = evm.Call(senderAcc, contractAddr, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("smartcontract func, send 1|senderAcc=%s|toAddress=%s\n", senderAcc.Address().String(), toAddress.String())

	//send
	input, err = abiObj.Pack("send", toAddress, big.NewInt(19))
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(senderAcc, contractAddr, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("smartcontract func, send 2|senderAcc=%s|toAddress=%s\n", senderAcc.Address().String(), toAddress.String())

	// get balance
	input, err = abiObj.Pack("balances", toAddress)
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(contractRef, contractAddr, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("smartcontract  func, balances|toAddress=%s|balance=%x\n", toAddress.String(), outputs)
	toAddressBalance := outputs

	// get balance
	input, err = abiObj.Pack("balances", sender)
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(contractRef, contractAddr, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	fmt.Printf("smartcontract  func, balances|sender=%s|balance=%x\n", sender.String(), outputs)

	// get event
	logs := stateDB.Logs()

	for _, log := range logs {
		fmt.Printf("%#v\n", log)
		for _, topic := range log.Topics {
			fmt.Printf("topic: %#v\n", topic)
		}
		fmt.Printf("data: %#v\n", log.Data)
	}

	testBalance = stateDB.GetBalance(fromAddress)
	fmt.Println("get testBalance =", testBalance)

	//commit
	stateDB.Commit(false)
	ms.Write()
	cms.Commit()
	db.Close()

	if !bytes.Equal(contractCode, stateDB.GetCode(contractAddr)) {
		fmt.Println("BUG!,the code was changed!")
		os.Exit(-1)
	}

	//reopen DB
	err = reOpenDB(t, contractCode, contractAddr.String(), toAddressBalance, dataPath)
	must(err)

	//remove DB dir
	cleanup(dataPath)
}

func cleanup(dataDir string) {
	fmt.Printf("cleaning up db dir|dataDir=%s\n", dataDir)
	os.RemoveAll(dataDir)
}

func reOpenDB(t *testing.T, lastContractCode []byte, strContractAddress string, lastBalance []byte, dataPath string) (err error) {
	fmt.Printf("strContractAddress=%s\n", strContractAddress)

	lastContractAddress := common.HexToAddress(strContractAddress)

	fmt.Printf("reOpenDB...\n")

	//---------------------stateDB test--------------------------------------
	// dataPath := "/tmp/htdfNewEvmTestData3"
	db := dbm.NewDB("state", dbm.GoLevelDBBackend, dataPath)

	cdc := newTestCodec1()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(codeKey, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, nil)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, accKey, pk.Subspace(auth.DefaultParamspace), newevmtypes.ProtoBaseAccount)

	cms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, nil)
	cms.MountStoreWithDB(storageKey, sdk.StoreTypeIAVL, nil)

	cms.SetPruning(store.PruneNothing)

	err = cms.LoadLatestVersion()
	must(err)

	ms := cms.CacheMultiStore()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	stateDB, err := evmtypes.NewCommitStateDB(ctx, &ak, storageKey, codeKey)
	must(err)

	fmt.Printf("addr=%s|testBalance=%v\n", fromAddress.String(), stateDB.GetBalance(fromAddress))

	fmt.Printf("lastContractCode=%x\n", lastContractCode)

	if !bytes.Equal(lastContractCode, stateDB.GetCode(lastContractAddress)) {
		panic("different contract code")
	}

	//---------------------call evm--------------------------------------
	abiFileName := "../../tests/evm/coin/coin_sol_Coin.abi"
	binFileName := "../../tests/evm/coin/coin_sol_Coin.bin"
	data := loadBin(binFileName)

	//	config := params.TestnetChainConfig
	config := appParams.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	msg := NewMessage(fromAddress, &toAddress, nonce, amount, gasLimit, big.NewInt(0), data, false)
	// evmCtx := evmcore.NewEVMContext(msg, &fromAddress, 1000, ctx.BlockHeader().Time)
	// evm := vm.NewEVM(evmCtx, stateDB, config, vmConfig)

	blockCtx := evmcore.NewEVMBlockContext(ctx.BlockHeader(), &evmcore.FakeChainContext{}, &fromAddress)
	txCtx := evmcore.NewEVMTxContext(msg)
	evm := vm.NewEVM(blockCtx, txCtx, stateDB, config, vmConfig)

	contractRef := vm.AccountRef(fromAddress)

	fmt.Printf("BlockNumber=%d|IsEIP158=%v\n", evm.Context.BlockNumber.Uint64(), evm.ChainConfig().IsEIP158(evm.Context.BlockNumber))
	testChainConfig(t, evm)

	abiObj := loadAbi(abiFileName)

	// get balance
	input, err := abiObj.Pack("balances", toAddress)
	must(err)
	outputs, gasLeftover, vmerr := evm.Call(contractRef, lastContractAddress, input, stateDB.GetBalance(fromAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	tmpHexString := hex.EncodeToString(input)
	rawData, _ := hex.DecodeString(tmpHexString)
	if bytes.Compare(input, rawData) != 0 {
		t.Errorf("rawData convert error|rawData=%s|tmpHexString=%s\n", hex.EncodeToString(rawData), tmpHexString)
	}

	fmt.Printf("input=%s\n", tmpHexString)

	fmt.Printf("smartcontract  func, balances|toAddress=%s|balance=%x\n", toAddress.String(), outputs)

	fmt.Printf("gasLeftover=%d\n", gasLeftover)

	if !bytes.Equal(lastBalance, outputs) {
		panic("different balance")
	}

	//commit
	stateDB.Commit(false)
	ms.Write()
	cms.Commit()
	db.Close()

	return nil
}

// type ChainContext struct{}
