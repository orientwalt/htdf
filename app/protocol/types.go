package protocol

import (
	"encoding/json"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Protocol interface {
	GetVersion() uint64
	GetRouter() Router
	GetQueryRouter() QueryRouter
	GetAnteHandler() sdk.AnteHandler                   // ante handler for fee and auth
	GetFeeRefundHandler() sdk.FeeRefundHandler         // fee handler for fee refund
	GetFeePreprocessHandler() sdk.FeePreprocessHandler // fee handler for fee preprocessor
	ExportAppStateAndValidators(ctx sdk.Context, forZeroHeight bool, jailWhiteList []string) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error)
	ValidateTx(ctx sdk.Context, txBytes []byte, msgs []sdk.Msg) sdk.Error

	// may be nil
	GetInitialHeight(ctx sdk.Context) int64 // get initial height

	GetInitChainer() sdk.InitChainer1  // initialize state with validators and state blob
	GetBeginBlocker() sdk.BeginBlocker // logic to run before any txs
	GetEndBlocker() sdk.EndBlocker     // logic to run after all txs, and to determine valset changes

	GetKVStoreKeyList() []*sdk.KVStoreKey
	Load()
	Init()
	GetCodec() *codec.Codec
	InitMetrics(store sdk.CommitMultiStore) // init metrics

	SetInitialHeight(ctx sdk.Context, initialheight uint64)
}
