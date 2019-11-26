package main

import (
	"encoding/json"
	"io"

	"github.com/orientwalt/htdf/server"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	bam "github.com/orientwalt/htdf/app"
	hsinit "github.com/orientwalt/htdf/init"
	lite "github.com/orientwalt/htdf/lite/cmd"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagOverwrite = "overwrite"
)

var invCheckPeriod uint

func main() {
	cobra.EnableCommandSorting = false
	cdc := bam.MakeLatestCodec()
	ctx := server.NewDefaultContext()

	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:               "hsd",
		Short:             "HtdfService App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(hsinit.InitCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.LiveNetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(lite.Commands())

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "HS", bam.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagOverwrite,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, config *cfg.InstrumentationConfig) abci.Application {
	return bam.NewHtdfServiceApp(
		logger, config, db, traceStore, true, invCheckPeriod,
		bam.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		bam.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(ctx *server.Context,
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		gApp := bam.NewHtdfServiceApp(logger, ctx.Config.Instrumentation, db, traceStore, false, uint(1))
		err := gApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return gApp.ExportAppStateAndValidators(forZeroHeight)
	}
	gApp := bam.NewHtdfServiceApp(logger, ctx.Config.Instrumentation, db, traceStore, true, uint(1))
	return gApp.ExportAppStateAndValidators(forZeroHeight)
}
