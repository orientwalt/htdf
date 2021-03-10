package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/orientwalt/htdf/app/protocol"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/x/auth"

	sdk "github.com/orientwalt/htdf/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"

	v0 "github.com/orientwalt/htdf/app/v0"
	v1 "github.com/orientwalt/htdf/app/v1"
	v2 "github.com/orientwalt/htdf/app/v2"
	"github.com/orientwalt/htdf/server"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
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

const (
	appName = "HtdfServiceApp"

	appPrometheusNamespace = "htdf"
	//
	RouterKey = "htdfservice"
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"

	DefaultCacheSize = 100 // Multistore saves last 100 blocks

	DefaultSyncableHeight = 10000 // Multistore saves a snapshot every 10000 blocks
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.hscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.hsd")
)

// Extended ABCI application
type HtdfServiceApp struct {
	*BaseApp
	// cdc *codec.Codec

	invCheckPeriod uint
}

// NewHtdfServiceApp is a constructor function for htdfServiceApp
func NewHtdfServiceApp(logger log.Logger, config *cfg.InstrumentationConfig, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*BaseApp)) *HtdfServiceApp {

	cdc := MakeLatestCodec()

	bApp := NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	var app = &HtdfServiceApp{
		BaseApp:        bApp,
		invCheckPeriod: invCheckPeriod,
	}
	protocolKeeper := sdk.NewProtocolKeeper(protocol.KeyMain)
	logrus.Traceln("/---------protocolKeeper----------/", protocolKeeper)
	engine := protocol.NewProtocolEngine(protocolKeeper)
	app.SetProtocolEngine(&engine)
	app.MountStoresIAVL(engine.GetKVStoreKeys())
	app.MountStoresTransient(engine.GetTransientStoreKeys())

	var err error
	if viper.GetBool(server.FlagReplay) {
		lastHeight := Replay(app.logger)
		err = app.LoadVersion(lastHeight, protocol.KeyMain, true)
	} else {
		err = app.LoadLatestVersion(protocol.KeyMain)
	} // app is now sealed
	if err != nil {
		cmn.Exit(err.Error())
	}
	//Duplicate prometheus config
	appPrometheusConfig := *config
	//Change namespace to appName
	appPrometheusConfig.Namespace = appPrometheusNamespace
	engine.Add(v0.NewProtocolV0(0, logger, protocolKeeper, app.invCheckPeriod, &appPrometheusConfig))
	engine.Add(v1.NewProtocolV1(1, logger, protocolKeeper, app.invCheckPeriod, &appPrometheusConfig))
	engine.Add(v2.NewProtocolV2(2, logger, protocolKeeper, app.invCheckPeriod, &appPrometheusConfig))
	//engine.Add(v2.NewProtocolV1(2, ...))
	logrus.Traceln("KeyMain----->	", app.GetKVStore(protocol.KeyMain))
	loaded, current := engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))

	logrus.Infof("currVersion=%v\n", engine.GetCurrentProtocol().GetVersion())
	logrus.Infof("LastBlockHeight=%v\n", app.BaseApp.LastBlockHeight())

	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	app.BaseApp.txDecoder = auth.DefaultTxDecoder(engine.GetCurrentProtocol().GetCodec())
	engine.GetCurrentProtocol().InitMetrics(app.cms)
	logrus.Traceln("keystorage----->	", app.GetKVStore(protocol.KeyStorage))
	return app
}

func (app *HtdfServiceApp) ExportOrReplay(replayHeight int64) (replay bool, height int64) {
	lastBlockHeight := app.BaseApp.LastBlockHeight()
	if replayHeight > lastBlockHeight {
		replayHeight = lastBlockHeight
	}

	if lastBlockHeight-replayHeight <= DefaultCacheSize {
		err := app.LoadVersion(replayHeight, protocol.KeyMain, false)
		if err != nil {
			cmn.Exit(err.Error())
		}
		return false, replayHeight
	}

	loadHeight := app.replayToHeight(replayHeight, app.logger)
	err := app.LoadVersion(loadHeight, protocol.KeyMain, true)
	if err != nil {
		cmn.Exit(err.Error())
	}
	app.logger.Info(fmt.Sprintf("Load store at %d, start to replay to %d", loadHeight, replayHeight))
	return true, replayHeight

}

// export the state of htdf for a genesis file
func (app *HtdfServiceApp) ExportAppStateAndValidators(forZeroHeight bool) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	return app.Engine.GetCurrentProtocol().ExportAppStateAndValidators(ctx, forZeroHeight, []string{})
}

// load a particular height
func (app *HtdfServiceApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, protocol.KeyMain, false)
}

// MakeCodec generates the necessary codecs for Amino
func MakeLatestCodec() *codec.Codec {
	// TODO: replace v0 with v2 ??
	var cdc = v0.MakeLatestCodec() // replace with latest protocol version
	return cdc
}

func (app *HtdfServiceApp) replayToHeight(replayHeight int64, logger log.Logger) int64 {
	loadHeight := int64(0)
	if replayHeight >= DefaultSyncableHeight {
		loadHeight = replayHeight - replayHeight%DefaultSyncableHeight
	} else {
		// version 1 will always be kept for block reset
		loadHeight = 1
	}
	return loadHeight
}

// ResetOrReplay returns whether you need to reset or replay
func (app *HtdfServiceApp) ResetOrReplay(replayHeight int64) (replay bool, height int64) {
	lastBlockHeight := app.BaseApp.LastBlockHeight()
	if replayHeight > lastBlockHeight {
		replayHeight = lastBlockHeight
	}

	fmt.Println("NOTE: This Reset operation will change the application store!")
	fmt.Println("️NOTE: Backup(备份,備份,지원,Apoyo,Резервный) your node home directory before proceeding!")

	// for safety, ask user input the reset height again
	fmt.Println("Please input reset height again:")
	inputHeight, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		cmn.Exit(err.Error())
	}
	resetHeight, err := strconv.ParseInt(strings.TrimSpace(inputHeight), 10, 64)
	if resetHeight != replayHeight {
		cmn.Exit(fmt.Sprintf("The second input reset height(%v) does not match first input height(%v)!", resetHeight, replayHeight))
	}

	// for safety, check backup dir exists
	fmt.Println("Please input absolute path of your backuped node home directory for check it's exists:")
	backupPath, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		cmn.Exit(err.Error())
	}
	backupPath = strings.TrimSpace(backupPath)
	s, err := os.Stat(backupPath)
	if err != nil {
		cmn.Exit("Backup path doesn't exists: " + err.Error())
	}
	if !s.IsDir() {
		cmn.Exit(fmt.Sprintf("Backup path '%v' is not a directory!", backupPath))
	}

	// last confirm
	fmt.Printf("The last block height is %v, will reset height to %v.\n", lastBlockHeight, replayHeight)
	fmt.Println("Are you sure to proceed? (yes/n)")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		cmn.Exit(err.Error())
	}
	confirm := strings.ToLower(strings.TrimSpace(input))
	if confirm != "yes" {
		cmn.Exit("Reset operation aborted.")
	}

	if lastBlockHeight-replayHeight <= DefaultCacheSize {
		err := app.LoadVersion(replayHeight, protocol.KeyMain, true)

		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("wanted to load target %v but only found up to", replayHeight)) {
				app.logger.Info(fmt.Sprintf("Can not find the target version %d, trying to load an earlier version and replay blocks", replayHeight))
			} else {
				cmn.Exit(err.Error())
			}
		} else {
			app.logger.Info(fmt.Sprintf("The last block height is %d, loaded store at %d", lastBlockHeight, replayHeight))
			return false, replayHeight
		}
	}

	loadHeight := app.replayToHeight(replayHeight, app.logger)
	err = app.LoadVersion(loadHeight, protocol.KeyMain, true)
	if err != nil {
		cmn.Exit(err.Error())
	}

	// If reset to another protocol version, should reload Protocol and reset txDecoder
	loaded, current := app.Engine.LoadCurrentProtocol(app.GetKVStore(protocol.KeyMain))
	if !loaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	app.BaseApp.txDecoder = auth.DefaultTxDecoder(app.Engine.GetCurrentProtocol().GetCodec())

	app.logger.Info(fmt.Sprintf("The last block height is %d, want to load store at %d", lastBlockHeight, replayHeight))

	// Version 1 does not need replay
	if replayHeight == 1 {
		app.logger.Info(fmt.Sprintf("Loaded store at %d", loadHeight))
		return false, replayHeight
	}

	app.logger.Info(fmt.Sprintf("Loaded store at %d, start to replay to %d", loadHeight, replayHeight))
	return true, replayHeight

}
