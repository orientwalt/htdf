package app

import (
	"os"
	"testing"

	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/stretchr/testify/require"

	v0 "github.com/orientwalt/htdf/app/v0"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	"github.com/orientwalt/htdf/x/crisis"
	distr "github.com/orientwalt/htdf/x/distribution"
	"github.com/orientwalt/htdf/x/gov"
	"github.com/orientwalt/htdf/x/guardian"
	"github.com/orientwalt/htdf/x/mint"
	"github.com/orientwalt/htdf/x/service"
	"github.com/orientwalt/htdf/x/slashing"
	"github.com/orientwalt/htdf/x/staking"
	"github.com/orientwalt/htdf/x/upgrade"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func setGenesis(happ *HtdfServiceApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]v0.GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = v0.NewGenesisAccount(acc)
	}

	genesisState := v0.NewGenesisState(
		genaccs,
		auth.DefaultGenesisState(),
		staking.DefaultGenesisState(),
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		upgrade.DefaultGenesisState(),
		service.DefaultGenesisState(),
		guardian.DefaultGenesisState(),
		slashing.DefaultGenesisState(),
		crisis.DefaultGenesisState(),
	)

	stateBytes, err := codec.MarshalJSONIndent(v0.MakeLatestCodec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	happ.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	happ.Commit()

	return nil
}

func TestHsdExport(t *testing.T) {
	db := db.NewMemDB()

	happ := NewHtdfServiceApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), config.TestInstrumentationConfig(), db, nil, true, 0)
	// accs added by junying, 2019-11-20
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	acc := auth.NewBaseAccountWithAddress(addr)
	setGenesis(happ, &acc)

	// Making a new app object with the db, so that initchain hasn't been called
	// panic: consensus params is empty
	_, _, err := happ.ExportAppStateAndValidators(false)
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
