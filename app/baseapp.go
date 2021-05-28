package app

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"errors"

	"github.com/gogo/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/orientwalt/htdf/app/protocol"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	sdkerrors "github.com/orientwalt/htdf/types/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	// junying-todo,2020-01-17
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info" //trace/debug/info/warn/error/parse/fatal/panic
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.FatalLevel //TraceLevel/DebugLevel/InfoLevel/WarnLevel/ErrorLevel/ParseLevel/FatalLevel/PanicLevel
	}
	// set global log level
	logrus.SetLevel(ll)
	logrus.SetFormatter(&logrus.TextFormatter{}) //&log.JSONFormatter{})
}

// Key to store the consensus params in the main store.
var mainConsensusParamsKey = []byte("consensus_params")

// Enum mode for app.runTx
type runTxMode uint8

const (
	// Check a transaction
	runTxModeCheck runTxMode = iota
	// Simulate a transaction
	runTxModeSimulate runTxMode = iota
	// Deliver a transaction
	runTxModeDeliver runTxMode = iota
	// Recheck a (pending) transaction after a commit
	runTxModeReCheck runTxMode = iota
	// MainStoreKey is the string representation of the main store
	MainStoreKey = "main"
)

type state struct {
	ms  sdk.CacheMultiStore
	ctx sdk.Context
}

// BaseApp reflects the ABCI application implementation.
type BaseApp struct {
	// initialized on creation
	logger      log.Logger
	name        string               // application name from abci.Info
	db          dbm.DB               // common DB backend
	cms         sdk.CommitMultiStore // Main (uncached) state
	router      sdk.Router           // handle any kind of message
	queryRouter sdk.QueryRouter      // router for redirecting query calls
	txDecoder   sdk.TxDecoder        // unmarshal []byte into sdk.Tx

	// set upon LoadVersion or LoadLatestVersion.
	baseKey *sdk.KVStoreKey // Main KVStore in cms

	anteHandler    sdk.AnteHandler  // ante handler for fee and auth
	initChainer    sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker   sdk.BeginBlocker // logic to run before any txs
	endBlocker     sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	addrPeerFilter sdk.PeerFilter   // filter peers by address and port
	idPeerFilter   sdk.PeerFilter   // filter peers by node ID
	fauxMerkleMode bool             // if true, IAVL MountStores uses MountStoresDB for simulation speed.

	// --------------------
	// Volatile state
	// checkState is set on initialization and reset on Commit.
	// deliverState is set in InitChain and BeginBlock and cleared on Commit.
	// See methods setCheckState and setDeliverState.
	checkState   *state          // for CheckTx
	deliverState *state          // for DeliverTx
	voteInfos    []abci.VoteInfo // absent validators from begin block

	// paramStore is used to query for ABCI consensus parameters from an
	// application parameter store.
	paramStore ParamStore

	// consensus params
	// TODO: Move this in the future to baseapp param store on main store.
	consensusParams *abci.ConsensusParams

	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. This is mainly used for DoS and spam prevention.
	minGasPrices sdk.Coins

	// flag for sealing options and parameters to a BaseApp
	sealed bool

	Engine *protocol.ProtocolEngine

	// block height at which to halt the chain and gracefully shutdown
	haltHeight uint64

	// minimum block time (in Unix seconds) at which to halt the chain and gracefully shutdown
	haltTime uint64

	// application's version string
	appVersion string
	// genesis block's initial height
	initialHeight int64
	//
	lastblkheight int64
}

// var _ abci.Application = (*BaseApp)(nil)

// NewBaseApp returns a reference to an initialized BaseApp. It accepts a
// variadic number of option functions, which act on the BaseApp to set
// configuration choices.
//
// NOTE: The db is used to store the version number for now.
func NewBaseApp(
	name string, logger log.Logger, db dbm.DB, txDecoder sdk.TxDecoder, options ...func(*BaseApp),
) *BaseApp {

	app := &BaseApp{
		logger: logger,
		name:   name,
		db:     db,
		cms:    store.NewCommitMultiStore(db),
		// router:         NewRouter(),
		// queryRouter:    NewQueryRouter(),
		txDecoder:      txDecoder,
		fauxMerkleMode: false,
	}
	for _, option := range options {
		option(app)
	}

	return app
}

// Name returns the name of the BaseApp.
func (app *BaseApp) Name() string {
	return app.name
}

// Logger returns the logger of the BaseApp.
func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

// SetCommitMultiStoreTracer sets the store tracer on the BaseApp's underlying
// CommitMultiStore.
func (app *BaseApp) SetCommitMultiStoreTracer(w io.Writer) {
	app.cms.SetTracer(w)
}

// Mount IAVL stores to the provided keys in the BaseApp multistore
func (app *BaseApp) MountStoresIAVL(keys []*sdk.KVStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeIAVL)
	}
}

// Mount stores to the provided keys in the BaseApp multistore
func (app *BaseApp) MountStoresTransient(keys []*sdk.TransientStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeTransient)
	}
}

func (app *BaseApp) SetProtocolEngine(pe *protocol.ProtocolEngine) {
	if app.sealed {
		panic("SetProtocolEngine() on sealed BaseApp")
	}
	app.Engine = pe
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountStores(keys ...sdk.StoreKey) {
	for _, key := range keys {
		switch key.(type) {
		case *sdk.KVStoreKey:
			if !app.fauxMerkleMode {
				app.MountStore(key, sdk.StoreTypeIAVL)
			} else {
				// StoreTypeDB doesn't do anything upon commit, and it doesn't
				// retain history, but it's useful for faster simulation.
				app.MountStore(key, sdk.StoreTypeDB)
			}
		case *sdk.TransientStoreKey:
			app.MountStore(key, sdk.StoreTypeTransient)
		default:
			panic("Unrecognized store key type " + reflect.TypeOf(key).Name())
		}
	}
}

// MountStoreWithDB mounts a store to the provided key in the BaseApp
// multistore, using a specified DB.
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

func (app *BaseApp) GetKVStore(key sdk.StoreKey) sdk.KVStore {
	return app.cms.GetKVStore(key)
}

// LoadLatestVersion loads the latest application version. It will panic if
// called more than once on a running BaseApp.
func (app *BaseApp) LoadLatestVersion(baseKey *sdk.KVStoreKey) error {
	err := app.cms.LoadLatestVersion()
	if err != nil {
		return err
	}
	return app.initFromMainStore(baseKey)
}

// LoadVersion loads the BaseApp application version. It will panic if called
// more than once on a running baseapp.
func (app *BaseApp) LoadVersion(version int64, baseKey *sdk.KVStoreKey, overwrite bool) error {
	err := app.cms.LoadVersion(version) //, overwrite)
	if err != nil {
		return fmt.Errorf("failed to load version %d: %w", version, err)
	}
	return app.initFromMainStore(baseKey)
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	lastcommit := app.cms.LastCommitID().Version
	logger().Debugln("LastBlockHeight()-lastcommit: ", lastcommit)
	// if lastcommit < app.initialHeight {
	// 	lastcommit = app.initialHeight - 1
	// }
	return lastcommit
}

// initializes the remaining logic from app.cms
func (app *BaseApp) initFromMainStore(baseKey *sdk.KVStoreKey) error {
	mainStore := app.cms.GetKVStore(baseKey)
	if mainStore == nil {
		return errors.New("baseapp expects MultiStore with 'main' KVStore")
	}

	// memoize baseKey
	// if app.baseKey != nil {
	// 	panic("app.baseKey expected to be nil; duplicate init?")
	// }
	app.baseKey = baseKey

	// Load the consensus params from the main store. If the consensus params are
	// nil, it will be saved later during InitChain.
	//
	// TODO: assert that InitChain hasn't yet been called.
	consensusParamsBz := mainStore.Get(mainConsensusParamsKey)
	if consensusParamsBz != nil {
		var consensusParams = &abci.ConsensusParams{}

		err := proto.Unmarshal(consensusParamsBz, consensusParams)
		if err != nil {
			panic(err)
		}

		app.setConsensusParams(consensusParams)
	} else {
		// It will get saved later during InitChain.
		if app.LastBlockHeight() != 0 {
			panic(errors.New("consensus params is empty"))
		}
	}

	// needed for `gaiad export`, which inits from store but never calls initchain
	app.setCheckState(abci.Header{})
	app.Seal()

	return nil
}

func (app *BaseApp) setMinGasPrices(gasPrices sdk.Coins) {
	app.minGasPrices = gasPrices
}

// // Router returns the router of the BaseApp.
// func (app *BaseApp) Router() Router {
// 	if app.sealed {
// 		// We cannot return a router when the app is sealed because we can't have
// 		// any routes modified which would cause unexpected routing behavior.
// 		panic("Router() on sealed BaseApp")
// 	}
// 	return app.router
// }

// QueryRouter returns the QueryRouter of a BaseApp.
func (app *BaseApp) QueryRouter() sdk.QueryRouter { return app.queryRouter }

// Seal seals a BaseApp. It prohibits any further modifications to a BaseApp.
func (app *BaseApp) Seal() { app.sealed = true }

// IsSealed returns true if the BaseApp is sealed and false otherwise.
func (app *BaseApp) IsSealed() bool { return app.sealed }

// setCheckState sets checkState with the cached multistore and
// the context wrapping it.
// It is called by InitChain() and Commit()
func (app *BaseApp) setCheckState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	app.checkState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, true, app.logger).WithMinGasPrices(app.minGasPrices),
	}
}

// setCheckState sets checkState with the cached multistore and
// the context wrapping it.
// It is called by InitChain() and BeginBlock(),
// and deliverState is set nil on Commit().
func (app *BaseApp) setDeliverState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	app.deliverState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, app.logger),
	}
}

// GetConsensusParams returns the current consensus parameters from the BaseApp's
// ParamStore. If the BaseApp has no ParamStore defined, nil is returned.
func (app *BaseApp) GetConsensusParams(ctx sdk.Context) *abci.ConsensusParams {
	if app.paramStore == nil {
		return nil
	}

	cp := new(abci.ConsensusParams)

	if app.paramStore.Has(ctx, ParamStoreKeyBlockParams) {
		var bp abci.BlockParams
		app.paramStore.Get(ctx, ParamStoreKeyBlockParams, &bp)
		cp.Block = &bp
	}

	if app.paramStore.Has(ctx, ParamStoreKeyEvidenceParams) {
		var ep abci.EvidenceParams
		app.paramStore.Get(ctx, ParamStoreKeyEvidenceParams, &ep)
		cp.Evidence = &ep
	}

	if app.paramStore.Has(ctx, ParamStoreKeyValidatorParams) {
		var vp abci.ValidatorParams
		app.paramStore.Get(ctx, ParamStoreKeyValidatorParams, &vp)
		cp.Validator = &vp
	}

	return cp
}

// setConsensusParams memoizes the consensus params.
func (app *BaseApp) setConsensusParams(consensusParams *abci.ConsensusParams) {
	app.consensusParams = consensusParams
}

// setConsensusParams stores the consensus params to the main store.
func (app *BaseApp) StoreConsensusParams(consensusParams *abci.ConsensusParams) {
	consensusParamsBz, err := proto.Marshal(consensusParams)
	if err != nil {
		panic(err)
	}
	mainStore := app.cms.GetKVStore(app.baseKey)
	mainStore.Set(mainConsensusParamsKey, consensusParamsBz)
}

// // StoreConsensusParams sets the consensus parameters to the baseapp's param store.
// func (app *BaseApp) StoreConsensusParams(ctx sdk.Context, cp *abci.ConsensusParams) {
// 	if app.paramStore == nil {
// 		panic("cannot store consensus params with no params store set")
// 	}
// 	if cp == nil {
// 		return
// 	}

// 	app.paramStore.Set(ctx, ParamStoreKeyBlockParams, cp.Block)
// 	app.paramStore.Set(ctx, ParamStoreKeyEvidenceParams, cp.Evidence)
// 	app.paramStore.Set(ctx, ParamStoreKeyValidatorParams, cp.Validator)
// }

// // getMaximumBlockGas gets the maximum gas from the consensus params. It panics
// // if maximum block gas is less than negative one and returns zero if negative
// // one.
// func (app *BaseApp) getMaximumBlockGas() uint64 {
// 	if app.consensusParams == nil || app.consensusParams.Block == nil {
// 		return 0
// 	}

// 	maxGas := app.consensusParams.Block.MaxGas
// 	switch {
// 	case maxGas < -1:
// 		panic(fmt.Sprintf("invalid maximum block gas: %d", maxGas))

// 	case maxGas == -1:
// 		return 0

// 	default:
// 		return uint64(maxGas)
// 	}
// }

// getMaximumBlockGas gets the maximum gas from the consensus params. It panics
// if maximum block gas is less than negative one and returns zero if negative
// one.
func (app *BaseApp) getMaximumBlockGas(ctx sdk.Context) uint64 {
	cp := app.GetConsensusParams(ctx)
	if cp == nil || cp.Block == nil {
		return 0
	}

	maxGas := cp.Block.MaxGas
	switch {
	case maxGas < -1:
		panic(fmt.Sprintf("invalid maximum block gas: %d", maxGas))

	case maxGas == -1:
		return 0

	default:
		return uint64(maxGas)
	}
}

// ----------------------------------------------------------------------------
// ABCI

// // Info implements the ABCI interface.
// func (app *BaseApp) Info(req abci.RequestInfo) abci.ResponseInfo {
// 	lastCommitID := app.cms.LastCommitID()

// 	return abci.ResponseInfo{
// 		AppVersion:       version.ProtocolVersion,
// 		Data:             app.name,
// 		LastBlockHeight:  lastCommitID.Version,
// 		LastBlockAppHash: lastCommitID.Hash,
// 	}
// }

// // SetOption implements the ABCI interface.
// func (app *BaseApp) SetOption(req abci.RequestSetOption) (res abci.ResponseSetOption) {
// 	// TODO: Implement!
// 	return
// }

// // FilterPeerByAddrPort filters peers by address/port.
// func (app *BaseApp) FilterPeerByAddrPort(info string) abci.ResponseQuery {
// 	if app.addrPeerFilter != nil {
// 		return app.addrPeerFilter(info)
// 	}
// 	return abci.ResponseQuery{}
// }

// // FilterPeerByIDfilters peers by node ID.
// func (app *BaseApp) FilterPeerByID(info string) abci.ResponseQuery {
// 	if app.idPeerFilter != nil {
// 		return app.idPeerFilter(info)
// 	}
// 	return abci.ResponseQuery{}
// }

// // Splits a string path using the delimiter '/'.
// // e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
// func splitPath(requestPath string) (path []string) {
// 	path = strings.Split(requestPath, "/")
// 	// first element is empty string
// 	if len(path) > 0 && path[0] == "" {
// 		path = path[1:]
// 	}
// 	return path
// }

// // Query implements the ABCI interface. It delegates to CommitMultiStore if it
// // implements Queryable.
// func (app *BaseApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
// 	path := splitPath(req.Path)
// 	if len(path) == 0 {
// 		msg := "no query path provided"
// 		return sdk.ErrUnknownRequest(msg).QueryResult()
// 	}

// 	switch path[0] {
// 	// "/app" prefix for special application queries
// 	case "app":
// 		return handleQueryApp(app, path, req)

// 	case "store":
// 		return handleQueryStore(app, path, req)

// 	case "p2p":
// 		return handleQueryP2P(app, path, req)

// 	case "custom":
// 		return handleQueryCustom(app, path, req)
// 	}

// 	msg := "unknown query path"
// 	return sdk.ErrUnknownRequest(msg).QueryResult()
// }

// func handleQueryApp(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
// 	if len(path) >= 2 {
// 		var result sdk.Result

// 		switch path[1] {
// 		case "simulate":
// 			txBytes := req.Data
// 			tx, err := app.txDecoder(txBytes)
// 			if err != nil {
// 				result = err.Result()
// 			} else {
// 				result = app.Simulate(txBytes, tx)
// 			}

// 		case "version":
// 			return abci.ResponseQuery{
// 				Code:      uint32(sdk.CodeOK),
// 				Codespace: string(sdk.CodespaceRoot),
// 				Value:     []byte(version.GetVersion()),
// 			}

// 		default:
// 			result = sdk.ErrUnknownRequest(fmt.Sprintf("Unknown query: %s", path)).Result()
// 		}

// 		value := codec.Cdc.MustMarshalBinaryLengthPrefixed(result)
// 		return abci.ResponseQuery{
// 			Code:      uint32(sdk.CodeOK),
// 			Codespace: string(sdk.CodespaceRoot),
// 			Value:     value,
// 		}
// 	}

// 	msg := "Expected second parameter to be either simulate or version, neither was present"
// 	return sdk.ErrUnknownRequest(msg).QueryResult()
// }

// func handleQueryStore(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
// 	// "/store" prefix for store queries
// 	queryable, ok := app.cms.(sdk.Queryable)
// 	if !ok {
// 		msg := "multistore doesn't support queries"
// 		return sdk.ErrUnknownRequest(msg).QueryResult()
// 	}

// 	req.Path = "/" + strings.Join(path[1:], "/")
// 	return queryable.Query(req)
// }

// func handleQueryP2P(app *BaseApp, path []string, _ abci.RequestQuery) (res abci.ResponseQuery) {
// 	// "/p2p" prefix for p2p queries
// 	if len(path) >= 4 {
// 		cmd, typ, arg := path[1], path[2], path[3]
// 		switch cmd {
// 		case "filter":
// 			switch typ {
// 			case "addr":
// 				return app.FilterPeerByAddrPort(arg)
// 			case "id":
// 				return app.FilterPeerByID(arg)
// 			}
// 		default:
// 			msg := "Expected second parameter to be filter"
// 			return sdk.ErrUnknownRequest(msg).QueryResult()
// 		}
// 	}

// 	msg := "Expected path is p2p filter <addr|id> <parameter>"
// 	return sdk.ErrUnknownRequest(msg).QueryResult()
// }

// func handleQueryCustom(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
// 	// path[0] should be "custom" because "/custom" prefix is required for keeper
// 	// queries.
// 	//
// 	// The queryRouter routes using path[1]. For example, in the path
// 	// "custom/gov/proposal", queryRouter routes using "gov".
// 	if len(path) < 2 || path[1] == "" {
// 		return sdk.ErrUnknownRequest("No route for custom query specified").QueryResult()
// 	}

// 	//querier := app.queryRouter.Route(path[1])
// 	querier := app.Engine.GetCurrentProtocol().GetQueryRouter().Route(path[1])
// 	if querier == nil {
// 		return sdk.ErrUnknownRequest(fmt.Sprintf("no custom querier found for route %s", path[1])).QueryResult()
// 	}

// 	// cache wrap the commit-multistore for safety
// 	ctx := sdk.NewContext(
// 		app.cms.CacheMultiStore(), app.checkState.ctx.BlockHeader(), true, app.logger,
// 	).WithMinGasPrices(app.minGasPrices)

// 	// Passes the rest of the path as an argument to the querier.
// 	//
// 	// For example, in the path "custom/gov/proposal/test", the gov querier gets
// 	// []string{"proposal", "test"} as the path.
// 	resBytes, err := querier(ctx, path[2:], req)
// 	if err != nil {
// 		return abci.ResponseQuery{
// 			Code:      uint32(err.Code()),
// 			Codespace: string(err.Codespace()),
// 			Log:       err.ABCILog(),
// 		}
// 	}

// 	return abci.ResponseQuery{
// 		Code:  uint32(sdk.CodeOK),
// 		Value: resBytes,
// 	}
// }

func (app *BaseApp) validateHeight(req abci.RequestBeginBlock) error {
	if req.Header.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Header.Height)
	}

	prevHeight := app.LastBlockHeight()
	if req.Header.Height != prevHeight+1 {
		return fmt.Errorf("invalid height: %d; expected: %d", req.Header.Height, prevHeight+1)
	}

	return nil
}

// // BeginBlock implements the ABCI application interface.
// func (app *BaseApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
// 	if app.cms.TracingEnabled() {
// 		app.cms.SetTracingContext(sdk.TraceContext(
// 			map[string]interface{}{"blockHeight": req.Header.Height},
// 		))
// 	}

// 	if err := app.validateHeight(req); err != nil {
// 		panic(err)
// 	}

// 	// Initialize the DeliverTx state. If this is the first block, it should
// 	// already be initialized in InitChain. Otherwise app.deliverState will be
// 	// nil, since it is reset on Commit.
// 	if app.deliverState == nil {
// 		app.setDeliverState(req.Header)
// 	} else {
// 		// In the first block, app.deliverState.ctx will already be initialized
// 		// by InitChain. Context is now updated with Header information.
// 		app.deliverState.ctx = app.deliverState.ctx.
// 			WithBlockHeader(req.Header).
// 			WithBlockHeight(req.Header.Height).WithCheckValidNum(0)
// 	}

// 	// add block gas meter
// 	var gasMeter sdk.GasMeter
// 	if maxGas := app.getMaximumBlockGas(); maxGas > 0 {
// 		gasMeter = sdk.NewGasMeter(maxGas)
// 	} else {
// 		gasMeter = sdk.NewInfiniteGasMeter()
// 	}

// 	app.deliverState.ctx = app.deliverState.ctx.WithBlockGasMeter(gasMeter)

// 	// if app.beginBlocker != nil {
// 	// 	res = app.beginBlocker(app.deliverState.ctx, req)
// 	// }
// 	beginBlocker := app.Engine.GetCurrentProtocol().GetBeginBlocker()

// 	if beginBlocker != nil {
// 		res = beginBlocker(app.deliverState.ctx, req)
// 	}
// 	// set the signed validators for addition to context in deliverTx
// 	app.voteInfos = req.LastCommitInfo.GetVotes()
// 	return
// }

// // CheckTx implements the ABCI interface. It runs the "basic checks" to see
// // whether or not a transaction can possibly be executed, first decoding, then
// // the ante handler (which checks signatures/fees/ValidateBasic), then finally
// // the route match to see whether a handler exists.
// //
// // NOTE:CheckTx does not run the actual Msg handler function(s).
// func (app *BaseApp) CheckTx(txBytes []byte) (res abci.ResponseCheckTx) {
// 	var result sdk.Result
// 	tx, err := app.txDecoder(txBytes)
// 	logrus.Traceln("CheckTx88888888888888888888:tx", tx)
// 	if err != nil {
// 		result = err.Result()
// 	} else {
// 		result = app.runTx(runTxModeCheck, txBytes, tx)
// 	}

// 	return abci.ResponseCheckTx{
// 		Code:      uint32(result.Code),
// 		Data:      result.Data,
// 		Log:       result.Log,
// 		GasWanted: int64(result.GasWanted), // TODO: Should type accept unsigned ints?
// 		GasUsed:   int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
// 		// Tags:      result.Events,
// 		Events: result.Events,
// 	}
// }

// // DeliverTx implements the ABCI interface.
// func (app *BaseApp) DeliverTx(txBytes []byte) (res abci.ResponseDeliverTx) {
// 	var result sdk.Result

// 	tx, err := app.txDecoder(txBytes)
// 	logrus.Traceln("DeliverTx1111111111111", tx)
// 	if err != nil {
// 		result = err.Result()
// 	} else {
// 		result = app.runTx(runTxModeDeliver, txBytes, tx)
// 	}
// 	logrus.Traceln("DeliverTx1111111111111", result.Data, result.Log, result.Events)
// 	// junying-todo, 2019-10-18
// 	// this return value is written to database(blockchain)
// 	return abci.ResponseDeliverTx{
// 		Code:      uint32(result.Code),
// 		Codespace: string(result.Codespace),
// 		Data:      result.Data,
// 		Log:       result.Log,
// 		GasWanted: int64(result.GasWanted), // TODO: Should type accept unsigned ints?
// 		GasUsed:   int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
// 		Events:    result.Events,
// 	}
// }

// // junying-todo, 2019-11-13
// // ValidateBasic executes basic validator calls for all messages
// // and checking minimum for ?
// func ValidateBasic(ctx sdk.Context, tx sdk.Tx) sdk.Error {
// 	stdtx, ok := tx.(auth.StdTx)
// 	if !ok {
// 		return sdk.ErrInternal("tx must be StdTx")
// 	}
// 	// skip gentxs
// 	logrus.Traceln("Current BlockHeight:", ctx.BlockHeight())
// 	if ctx.BlockHeight() < 1 {
// 		return nil
// 	}
// 	// Validate Tx
// 	return stdtx.ValidateBasic()
// }

// retrieve the context for the tx w/ txBytes and other memoized values.
func (app *BaseApp) getContextForTx(mode runTxMode, txBytes []byte) (ctx sdk.Context) {
	if app.consensusParams != nil {
		logger().Traceln(app.consensusParams)
	}
	ctx = app.getState(mode).ctx.
		WithTxBytes(txBytes).
		WithVoteInfos(app.voteInfos).
		WithConsensusParams(app.consensusParams)

	if mode == runTxModeSimulate {
		ctx, _ = ctx.CacheContext()
	}

	return
}

// Check if the msg is MsgEthereumTx
func IsMsgEthereumTx(msg sdk.Msg) bool {
	if msg.Route() == "htdfservice" {
		return true
	}
	return false
}

// runMsgs iterates through all the messages and executes them.
func (app *BaseApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode runTxMode) (*sdk.Result, error) {
	msgLogs := make([]sdk.ABCIMessageLog, 0, len(msgs)) // a list of JSON-encoded logs with msg index
	events := sdk.EmptyEvents()

	var data []byte // NOTE: we just append them all (?!)
	// var tags sdk.Tags // also just append them all
	var code sdk.CodeType
	var codespace sdk.CodespaceType
	// var gasUsed uint64

	logger().Traceln("runMsgs	begin~~~~~~~~~~~~~~~~~~~~~~~~")
	start := time.Now()
	for msgIdx, msg := range msgs {
		// match message route
		msgRoute := msg.Route()
		logger().Traceln("999999999999", msgRoute)
		//handler := app.router.Route(msgRoute)
		handler := app.Engine.GetCurrentProtocol().GetRouter().Route(msgRoute)
		if handler == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s; message index: %d", msgRoute, msgIdx)
		}

		var msgResult sdk.Result
		// skip actual execution for CheckTx mode & ReCheckTx mode
		// what about simulation mode?
		if mode != runTxModeCheck && mode != runTxModeReCheck {
			logrus.Traceln(msgRoute, handler)
			msgResult = handler(ctx, msg)

		}

		logger().Traceln("runMsgs:msgResult.GasUsed=", msgResult.GasUsed)
		// NOTE: GasWanted is determined by ante handler and GasUsed by the GasMeter.

		// Result.Data must be length prefixed in order to separate each result
		data = append(data, msgResult.Data...)
		// tags = append(tags, sdk.MakeTag(sdk.TagAction, msg.Type()))
		// tags = append(tags, msgResult.Tags...)
		msgEvents := sdk.Events{
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type())),
		}
		msgEvents = msgEvents.AppendEvents(msgResult.GetEvents())
		events = events.AppendEvents(msgEvents)

		msgLog := sdk.ABCIMessageLog{MsgIndex: msgIdx, Log: msgResult.Log}

		// junying-todo, 2019-11-05
		if IsMsgEthereumTx(msg) {
			ctx.GasMeter().UseGas(sdk.Gas(msgResult.GasUsed), msgRoute)
		}

		// stop execution and return on first failed message
		if !msgResult.IsOK() {
			msgLog.Success = false
			msgLogs = append(msgLogs, msgLog)

			code = msgResult.Code
			codespace = msgResult.Codespace

			break
		}

		msgLog.Success = true
		msgLogs = append(msgLogs, msgLog)

	}
	logJSON := codec.Cdc.MustMarshalJSON(msgLogs)

	result := &sdk.Result{
		Code:      code,
		Codespace: codespace,
		Data:      data,
		Log:       strings.TrimSpace(string(logJSON)),
		GasUsed:   ctx.GasMeter().GasConsumed(),
		Events:    events.ToABCIEvents(),
	}
	logger().Traceln("runMsgs	end~~~~~~~~~~~~~~~~~~~~~~~~")
	logger().Debugf("=======>>>>> runMsgs elapsed time: %v", time.Since(start))
	return result, nil
}

// Returns the applications's deliverState if app is in runTxModeDeliver,
// otherwise it returns the application's checkstate.
func (app *BaseApp) getState(mode runTxMode) *state {
	if mode != runTxModeDeliver {
		return app.checkState
	}

	return app.deliverState
}

// cacheTxContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *BaseApp) cacheTxContext(ctx sdk.Context, txBytes []byte) (
	sdk.Context, sdk.CacheMultiStore) {

	ms := ctx.MultiStore()
	// TODO: https://github.com/cosmos/cosmos-sdk/issues/2824
	msCache := ms.CacheMultiStore()
	if msCache.TracingEnabled() {
		msCache = msCache.SetTracingContext(
			sdk.TraceContext(
				map[string]interface{}{
					"txHash": fmt.Sprintf("%X", tmhash.Sum(txBytes)),
				},
			),
		).(sdk.CacheMultiStore)
	}

	return ctx.WithMultiStore(msCache), msCache
}

// Validate Tx
func (app *BaseApp) ValidateTx(ctx sdk.Context, txBytes []byte, tx sdk.Tx) sdk.Error {
	// TxByteSize Check
	var msgs = tx.GetMsgs()
	if err := app.Engine.GetCurrentProtocol().ValidateTx(ctx, txBytes, msgs); err != nil {
		return err
	}

	// // ValidateBasic
	// if err := ValidateBasic(ctx, tx); err != nil {
	// 	logrus.Traceln("1runTx!!!!!!!!!!!!!!!!!")
	// 	return err
	// }

	// Msgs Check
	// All htdfservice Msgs: OK
	// All non-htdfservice Msgs: OK
	// htdfservice Msg(s) + non-htdfservice Msg(s): No
	// htdfservice Msg: OK, Msgs: No?
	var count = 0
	for _, msg := range msgs {
		if msg.Route() == "htdfservice" {
			count = count + 1
		}
	}
	if count > 0 && len(msgs) != count {
		return sdk.ErrInternal("mixed type of htdfservice msgs & non-htdfservice msgs can't be used")
	}
	// htdfservice Msgs: No
	// if count > 1 {
	// 	return sdk.ErrInternal("the number of htdfservice can't be more than one")
	// }
	return nil
}

// runTx processes a transaction. The transactions is proccessed via an
// anteHandler. The provided txBytes may be nil in some cases, eg. in tests. For
// further details on transaction execution, reference the BaseApp SDK
// documentation.
func (app *BaseApp) runTx(mode runTxMode, txBytes []byte, tx sdk.Tx) (gInfo sdk.GasInfo, result sdk.Result, err error) {
	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
	// determined by the GasMeter. We need access to the context to get the gas
	// meter so we initialize upfront.
	result = sdk.Result{}

	var gasWanted uint64
	ctx := app.getContextForTx(mode, txBytes)
	ms := ctx.MultiStore()

	gInfo = sdk.NewGasInfo()
	// only run the tx if there is block gas remaining
	if mode == runTxModeDeliver && ctx.BlockGasMeter().IsOutOfGas() {
		return gInfo, result, sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx") //sdk.ErrOutOfGas("no block gas left to run tx")
	}

	if err := app.ValidateTx(ctx, txBytes, tx); err != nil {
		result.Code = err.Code()
		result.Codespace = err.Codespace()
		result.Log = err.ABCILog()
		return gInfo, result, err //err.Result()
	}

	var startingGas uint64
	if mode == runTxModeDeliver {
		startingGas = ctx.BlockGasMeter().GasConsumed()
	}
	logger().Traceln("runTx:startingGas", startingGas)
	if mode == runTxModeDeliver {
		app.deliverState.ctx = app.deliverState.ctx.WithCheckValidNum(app.deliverState.ctx.CheckValidNum() + 1)
	}

	defer func() {

		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				err = sdkerrors.Wrap(
					sdkerrors.ErrOutOfGas, fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gasWanted, ctx.GasMeter().GasConsumed(),
					),
				)
				logger().Traceln("")
			default:
				err = sdkerrors.Wrap(
					sdkerrors.ErrPanic, fmt.Sprintf(
						"recovered: %v\nstack:\n%v", r, string(debug.Stack()),
					),
				)
				logger().Traceln(err)
			}
			logger().Traceln("2runTx!!!!!!!!!!!!!!!!!", r)
		}

		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
	}()
	if result.Code != 0 {
		logger().Traceln("runTx:result.GasUsed", result.GasUsed)
	}
	// Add cache in fee refund. If an error is returned or panic happes during refund,
	// no value will be written into blockchain state.
	defer func() {

		// commented by junying,2019-10-30
		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}

		var refundCtx sdk.Context
		var refundCache sdk.CacheMultiStore
		refundCtx, refundCache = app.cacheTxContext(ctx, txBytes)
		feeRefundHandler := app.Engine.GetCurrentProtocol().GetFeeRefundHandler()

		// Refund unspent fee
		if (mode == runTxModeDeliver || mode == runTxModeSimulate) && feeRefundHandler != nil {
			_, err := feeRefundHandler(refundCtx, tx, result)
			if err != nil {
				return
			}
			refundCache.Write()
		}
	}()
	logger().Traceln("3runTx!!!!!!!!!!!!!!!!!")
	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	defer func() {
		if mode == runTxModeDeliver {
			ctx.BlockGasMeter().ConsumeGas(
				ctx.GasMeter().GasConsumedToLimit(), // replaced by junying,2019-10-30
				// result.GasUsed, ///////////////////////// with this
				"block gas meter",
			)

			if ctx.BlockGasMeter().GasConsumed() < startingGas {

				panic(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"})
			}
		}
	}()
	logger().Traceln("4runTx!!!!!!!!!!!!!!!!!")
	// feePreprocessHandler := app.Engine.GetCurrentProtocol().GetFeePreprocessHandler()
	// // run the fee handler
	// if feePreprocessHandler != nil && ctx.BlockHeight() != 0 {
	// 	err := feePreprocessHandler(ctx, tx)
	// 	if err != nil {

	// 		return err.Result()
	// 	}
	// }
	logger().Traceln("5runTx!!!!!!!!!!!!!!!!!")

	anteHandler := app.Engine.GetCurrentProtocol().GetAnteHandler()
	if anteHandler != nil {
		var anteCtx sdk.Context
		var msCache sdk.CacheMultiStore

		// Cache wrap context before anteHandler call in case it aborts.
		// This is required for both CheckTx and DeliverTx.
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
		//
		// NOTE: Alternatively, we could require that anteHandler ensures that
		// writes do not happen if aborted/failed.  This may have some
		// performance benefits, but it'll be more difficult to get right.
		anteCtx, msCache = app.cacheTxContext(ctx, txBytes)

		newCtx, result, abort := anteHandler(anteCtx, tx, mode == runTxModeSimulate)
		logger().Traceln("anteHandler", result.GasUsed, result.GasWanted, result.Log)
		if !newCtx.IsZero() {
			// At this point, newCtx.MultiStore() is cache-wrapped, or something else
			// replaced by the ante handler. We want the original multistore, not one
			// which was cache-wrapped for the ante handler.
			//
			// Also, in the case of the tx aborting, we need to track gas consumed via
			// the instantiated gas meter in the ante handler, so we update the context
			// prior to returning.
			ctx = newCtx.WithMultiStore(ms)
		}

		gasWanted = result.GasWanted

		if abort {
			return gInfo, result, fmt.Errorf("%s", result.String())
		}

		msCache.Write()
	}
	logger().Traceln("6runTx!!!!!!!!!!!!!!!!!", result)
	logger().Traceln("runTxMode(check0,simulate,deliver,recheck)", mode)
	if mode == runTxModeCheck || mode == runTxModeReCheck {
		return
	}
	logger().Traceln("7runTx!!!!!!!!!!!!!!!!!")
	// Create a new context based off of the existing context with a cache wrapped
	// multi-store in case message processing fails.
	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)
	logger().Traceln("8runTx!!!!!!!!!!!!!!!!!", tx.GetMsgs(), mode)
	res, err := app.runMsgs(runMsgCtx, tx.GetMsgs(), mode)
	logger().Traceln("9runTx!!!!!!!!!!!!!!!!!", tx.GetMsgs())

	// 1. simulating
	// 2. fail to find routekey, resulted in nil returned from runMsg
	if mode == runTxModeSimulate || res == nil {
		return
	}
	result = *res
	result.GasWanted = gasWanted
	logger().Traceln("10runTx!!!!!!!!!!!!!!!!!", result.IsOK(), result.Code, result.GasUsed, result.GasWanted)
	// only update state if all messages pass
	// junying-todo, 2019-11-05
	// wondering if should add some condition for evm failure
	// if result.IsOK()
	// ErrCodeOK             CodeType = 0
	// ErrCodeCreateContract CodeType = 1
	// ErrCodeOpenContract   CodeType = 2
	// ErrCodeBeZeroAmount   CodeType = 3
	// ErrCodeParam          CodeType = 4
	// ErrCodeParsing        CodeType = 5
	// ErrCodeIntrinsicGas   CodeType = 6
	if result.Code < 3 {
		logger().Traceln("11runTx!!!!!!!!!!!!!!!!!")
		msCache.Write()
	}

	return
}

// // EndBlock implements the ABCI interface.
// func (app *BaseApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
// 	if app.deliverState.ms.TracingEnabled() {
// 		app.deliverState.ms = app.deliverState.ms.SetTracingContext(nil).(sdk.CacheMultiStore)
// 	}

// 	// if app.endBlocker != nil {
// 	// 	res = app.endBlocker(app.deliverState.ctx, req)
// 	// }
// 	endBlocker := app.Engine.GetCurrentProtocol().GetEndBlocker()
// 	if endBlocker != nil {
// 		res = endBlocker(app.deliverState.ctx, req)
// 	}
// 	_, appVersionStr, ok := abci.GetEventByKey(res.Events, sdk.AppVersionTag)
// 	if !ok {
// 		return
// 	}

// 	appVersion, _ := strconv.ParseUint(string(appVersionStr.GetValue()), 10, 64)
// 	if appVersion <= app.Engine.GetCurrentVersion() {
// 		return
// 	}
// 	fmt.Print("111111111111	", appVersion, "	22222222222	", app.Engine.GetCurrentVersion(), "\n")
// 	success := app.Engine.Activate(appVersion)
// 	if success {
// 		app.txDecoder = auth.DefaultTxDecoder(app.Engine.GetCurrentProtocol().GetCodec())
// 		return
// 	}

// 	if upgradeConfig, ok := app.Engine.ProtocolKeeper.GetUpgradeConfigByStore(app.GetKVStore(protocol.KeyMain)); ok {
// 		res.Events = append(res.Events,
// 			sdk.MakeEvent("upgrade", tmstate.UpgradeFailureTagKey,
// 				("Please install the right application version from "+upgradeConfig.Protocol.Software)))
// 	} else {
// 		res.Events = append(res.Events,
// 			sdk.MakeEvent("upgrade", tmstate.UpgradeFailureTagKey, ("Please install the right application version")))
// 	}

// 	return
// }

// // Commit implements the ABCI interface.
// func (app *BaseApp) Commit() (res abci.ResponseCommit) {
// 	header := app.deliverState.ctx.BlockHeader()

// 	// write the Deliver state and commit the MultiStore
// 	app.deliverState.ms.Write()
// 	commitID := app.cms.Commit()
// 	app.logger.Debug("Commit synced", "commit", fmt.Sprintf("%X", commitID))

// 	// Reset the Check state to the latest committed.
// 	//
// 	// NOTE: safe because Tendermint holds a lock on the mempool for Commit.
// 	// Use the header from this latest block.
// 	app.setCheckState(header)

// 	// empty/reset the deliver state
// 	app.deliverState = nil

// 	return abci.ResponseCommit{
// 		Data: commitID.Hash,
// 	}
// }

// ----------------------------------------------------------------------------
// State

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

func (st *state) Context() sdk.Context {
	return st.ctx
}
