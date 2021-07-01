package server

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"net/http"
	_ "net/http/pprof"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	pvm "github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

// Tendermint full-node start flags
const (
	flagWithTendermint = "with-tendermint"
	flagAddress        = "address"
	flagTraceStore     = "trace-store"
	flagPruning        = "pruning"
	FlagMinGasPrices   = "minimum-gas-prices"

	FlagReplay        = "replay-last-block"
	FlagInitialHeight = "initial-height"
)

// StartCmd runs the service passed in, either stand-alone or in-process with
// Tendermint.
func StartCmd(ctx *Context, appCreator AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if !viper.GetBool(flagWithTendermint) {
			// 	ctx.Logger.Info("Starting ABCI without Tendermint")
			// 	return startStandAlone(ctx, appCreator)
			// }

			ctx.Logger.Info("Starting ABCI with Tendermint")

			// yqq, 2021-05-26
			// we start a pprof at test environment for profiling
			if _, ok := os.LookupEnv("HTDF_TEST_ENV"); ok {
				go func() {
					if err := http.ListenAndServe(":9999", nil); err != nil {
						ctx.Logger.Info("pprof=====>" + err.Error())
					}
				}()
			}
			if err := startInProcess(ctx, appCreator); err != nil {
				ctx.Logger.Error(err.Error())
			}
			return nil
		},
	}

	cmd.Flags().Bool(FlagReplay, false, "Replay the last block")
	cmd.Flags().Int64(FlagInitialHeight, 1, "genesis block's initial height")
	// core flags for the ABCI application
	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	// cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().String(flagPruning, "nothing", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().String(
		FlagMinGasPrices, "",
		"Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photino;0.0001stake)",
	)

	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}

func startInProcess(ctx *Context, appCreator AppCreator) error {
	cfg := ctx.Config
	home := cfg.RootDir
	traceWriterFile := viper.GetString(flagTraceStore)

	db, err := openDB(home)
	if err != nil {
		return err
	}
	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	app := appCreator(ctx.Logger, db, traceWriter, ctx.Config.Instrumentation)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return err
	}

	UpgradeOldPrivValFile(cfg)
	// create & start tendermint node
	tmNode, err := node.NewNode(
		cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	if err != nil {
		return err
	}

	err = tmNode.Start()
	if err != nil {
		return err
	}

	// FIXBUG: yqq 2021-07-01
	// TrapSignal(ctx.Logger, func() {
	// 	if tmNode.IsRunning() {
	// 		_ = tmNode.Stop()
	// 	}
	// })

	defer func() {
		if tmNode.IsRunning() {
			_ = tmNode.Stop()
		}
		ctx.Logger.Info("exiting...")
	}()

	return WaitForQuitSignals()
}
