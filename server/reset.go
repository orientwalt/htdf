package server

import (
	"fmt"
	"path/filepath"

	"github.com/orientwalt/htdf/codec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	bc "github.com/tendermint/tendermint/blockchain"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// ResetCmd reset app state to particular height
func ResetCmd(ctx *Context, cdc *codec.Codec, appReset AppReset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset app state to the specified height",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := viper.GetString(tmcli.HomeFlag)
			traceWriterFile := viper.GetString(flagTraceStore)
			emptyState, err := isEmptyState(home)
			if err != nil {
				return err
			}

			if emptyState {
				fmt.Println("WARNING: State is not initialized.")
				return nil
			}

			db, err := openDB(home)
			if err != nil {
				return err
			}

			traceWriter, err := openTraceWriter(traceWriterFile)
			if err != nil {
				return err
			}
			height := viper.GetInt64(flagHeight)
			if height <= 0 {
				return errors.Errorf("Height must greater than zero")
			}

			if err := checkHeight(home, height); err != nil {
				return err
			}

			err = appReset(ctx, ctx.Logger, db, traceWriter, height)
			if err != nil {
				return errors.Errorf("Error reset state: %v\n", err)
			}

			fmt.Println("Reset app state successfully")
			return nil
		},
	}
	cmd.Flags().Uint64(flagHeight, 0, "Reset state from a particular height (greater than latest height means latest height)")
	cmd.MarkFlagRequired(flagHeight)
	return cmd
}

func checkHeight(home string, target int64) error {
	home = filepath.Join(home, "data")
	blockDb := loadDb("blockstore", home)
	defer func() {
		blockDb.Close()
		if r := recover(); r != nil {
			panic(errors.Errorf("height: %d not existed in block store", target))
		}

	}()
	blockStore := bc.NewBlockStore(blockDb)
	block := blockStore.LoadBlock(target)
	if block == nil {
		return errors.Errorf("height: %d not existed in block store", target)
	}
	return nil
}

func loadDb(name, path string) *dbm.GoLevelDB {
	db, err := dbm.NewGoLevelDB(name, path)
	if err != nil {
		panic(err)
	}
	return db
}
