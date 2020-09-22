package debug

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/orientwalt/htdf/store"
	"github.com/tendermint/tendermint/libs/log"
	cmn "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/orientwalt/htdf/app"
	"github.com/orientwalt/htdf/app/protocol"
	v0 "github.com/orientwalt/htdf/app/v0"

	"encoding/json"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	cfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"
)

func runHackCmd(cmd *cobra.Command, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("Expected 1 arg")
	}

	// ".htdf"
	dataDir := args[0]
	dataDir = path.Join(dataDir, "data")

	// load the app
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	db, err := dbm.NewGoLevelDB("htdf", dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app := NewHtdfApp(logger, db, bam.SetPruning(store.PruneNothing))

	// print some info
	id := app.LastCommitID()
	lastBlockHeight := app.LastBlockHeight()
	fmt.Println("ID", id)
	fmt.Println("LastBlockHeight", lastBlockHeight)

	//----------------------------------------------------
	// XXX: start hacking!
	//----------------------------------------------------
	// eg. fuxi-2000 testnet bug
	// We paniced when iterating through the "bypower" keys.
	// The following powerKey was there, but the corresponding "trouble" validator did not exist.
	// So here we do a binary search on the past states to find when the powerKey first showed up ...

	// owner of the validator the bonds, gets revoked, later unbonds, and then later is still found in the bypower store
	trouble := hexToBytes("880497F5AA9210987CAA945C588AF9E13A69E6F0")
	// this is his "bypower" key
	powerKey := hexToBytes("05303030303030303030303033FFFFFFFFFFFF4C0C0000FFFED3DC0FF59F7C3B548B7AFA365561B87FD0208AF8")

	topHeight := lastBlockHeight
	bottomHeight := int64(0)
	checkHeight := topHeight
	for {
		// load the given version of the state
		err = app.LoadVersion(checkHeight, protocol.KeyMain, false)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ctx := app.NewContext(true, abci.Header{})

		// check for the powerkey and the validator from the store
		store := ctx.KVStore(protocol.KeyStake)
		res := store.Get(powerKey)
		val, _ := app.Engine.GetCurrentProtocol().(*v0.ProtocolV0).StakeKeeper.GetValidator(ctx, trouble)
		fmt.Println("checking height", checkHeight, res, val)
		if res == nil {
			bottomHeight = checkHeight
		} else {
			topHeight = checkHeight
		}
		checkHeight = (topHeight + bottomHeight) / 2
	}
}

func base64ToPub(b64 string) ed25519.PubKeyEd25519 {
	data, _ := base64.StdEncoding.DecodeString(b64)
	var pubKey ed25519.PubKeyEd25519
	copy(pubKey[:], data)
	return pubKey

}

func hexToBytes(h string) []byte {
	trouble, _ := hex.DecodeString(h)
	return trouble

}

//--------------------------------------------------------------------------------
// NOTE: This is all copied from app/app.go
// so we can access internal fields!

const (
	appName = "HtdfApp"
)

// Extended ABCI application
type HtdfApp struct {
	*bam.BaseApp
}

func NewHtdfApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *HtdfApp {
	cdc := bam.MakeLatestCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	// create your application object
	var app = &HtdfApp{
		BaseApp: bApp,
	}
	protocolKeeper := sdk.NewProtocolKeeper(protocol.KeyMain)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	app.SetProtocolEngine(&engine)
	app.MountStoresIAVL(engine.GetKVStoreKeys())
	app.MountStoresTransient(engine.GetTransientStoreKeys())
	err := app.LoadLatestVersion(protocol.KeyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, 0, cfg.DefaultInstrumentationConfig()))
	// engine.Add(v1.NewProtocolV1(1, ...))

	engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))

	return app
}

// export the state of htdf for a genesis file
func (app *HtdfApp) ExportAppStateAndValidators(forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, nil)
}

func (app *HtdfApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.KeyMain, false)
}
