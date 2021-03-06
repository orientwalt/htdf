package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/params"

	"github.com/orientwalt/htdf/server"
	"github.com/orientwalt/htdf/store"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/orientwalt/htdf/app"
	hsinit "github.com/orientwalt/htdf/init"
	lite "github.com/orientwalt/htdf/lite/cmd"
	guardian "github.com/orientwalt/htdf/x/guardian/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagOverwrite = "overwrite"
)

var (
	invCheckPeriod uint
	GitCommit      = ""
	GitBranch      = ""
)

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
	rootCmd.AddCommand(hsinit.RealNetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(guardian.AddGuardianAccountCmd(ctx, cdc))
	rootCmd.AddCommand(hsinit.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(lite.Commands())
	rootCmd.AddCommand(versionCmd(ctx, cdc))

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

func versionCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cbCmd := &cobra.Command{
		Use:   "version",
		Short: "print version, api security level",
		Run: func(cmd *cobra.Command, args []string) {
			md5Sum, _ := getCurrentExeMd5Sum()
			fmt.Printf("GoVersion=%s|GitCommit=%s|version=%s|GitBranch=%s|md5sum=%s\n",
				runtime.Version(), GitCommit, params.Version, GitBranch, md5Sum)
		},
	}

	return cbCmd
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

func getCurrentExeMd5Sum() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	filePath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	var md5Sum string
	fp, err := os.Open(filePath)
	if err != nil {
		return md5Sum, err
	}
	defer fp.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, fp); err != nil {
		return md5Sum, err
	}
	hashInBytes := hash.Sum(nil)[:4] // only show 4 bytes
	// hashInBytes := hash.Sum(nil)
	md5Sum = hex.EncodeToString(hashInBytes)
	return md5Sum, nil
}
