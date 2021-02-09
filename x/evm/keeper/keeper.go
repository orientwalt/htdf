package keeper

import (
	"math/big"

	"github.com/orientwalt/htdf/codec"
	evmstate "github.com/orientwalt/htdf/evm/state"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/distribution/types"
	"github.com/orientwalt/htdf/x/params"
)

// keeper of the htdfservice store
type Keeper struct {
	ctx      sdk.Context
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	paramSpace          params.Subspace
	bankKeeper          types.BankKeeper
	stakingKeeper       types.StakingKeeper
	feeCollectionKeeper types.FeeCollectionKeeper

	// codespace
	codespace sdk.CodespaceType
	//
	CommitStateDB *evmstate.CommitStateDB
	TxCount       int
	Bloom         *big.Int
}
