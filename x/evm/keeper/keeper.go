package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/libs/log"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/orientwalt/htdf/app/protocol"
	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/store/prefix"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	evmtypes "github.com/orientwalt/htdf/x/evm/types"
)

// keeper of the htdfservice store
type Keeper struct {
	// ctx      sdk.Context
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	// paramSpace          params.Subspace
	AccountKeeper       auth.AccountKeeper
	FeeCollectionKeeper evmtypes.FeeCollectionKeeper

	// codespace
	codespace sdk.CodespaceType
	//
	CommitStateDB *evmtypes.CommitStateDB
	TxCount       int
	Bloom         *big.Int
}

// NewKeeper generates new evm module keeper
func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, ak auth.AccountKeeper, bk evmtypes.BankKeeper, fck evmtypes.FeeCollectionKeeper,
) Keeper {
	csdb, err := evmtypes.NewCommitStateDB(sdk.Context{}, &ak, protocol.KeyStorage, protocol.KeyCode)
	if err != nil {
		csdb = nil
	}
	return Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		CommitStateDB:       csdb,
		AccountKeeper:       ak,
		FeeCollectionKeeper: fck,
		TxCount:             0,
		Bloom:               big.NewInt(0),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", evmtypes.ModuleName))
}

// ----------------------------------------------------------------------------
// Block hash mapping functions
// Required by Web3 API.
//  TODO: remove once tendermint support block queries by hash.
// ----------------------------------------------------------------------------

// GetBlockHash gets block height from block consensus hash
func (k Keeper) GetBlockNumberByHash(ctx sdk.Context, hash []byte) (int64, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBlockHash)
	bz := store.Get(hash)
	if len(bz) == 0 {
		return 0, false
	}

	height := binary.BigEndian.Uint64(bz)
	return int64(height), true
}

// SetBlockHash sets the mapping from block consensus hash to block height
func (k Keeper) SetBlockHashToNumber(ctx sdk.Context, hash []byte, height int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBlockHash)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	store.Set(hash, bz)
}

// SetBlockNumber sets the mapping from block height to block consensus hash
func (k Keeper) SetBlockNumberToHash(ctx sdk.Context, height int64, hash []byte) {
	db := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(db, evmtypes.KeyPrefixBlockNumber)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	store.Set(bz, hash)
}

// GetBlockNumber  gets block hash by block height
func (k Keeper) GetBlockHashByNumber(ctx sdk.Context, height int64) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBlockNumber)
	bzHeight := sdk.Uint64ToBigEndian(uint64(height))
	bzHash := store.Get(bzHeight)
	if len(bzHash) == 0 {
		return []byte{}, false
	}
	return bzHash, true
}

// ----------------------------------------------------------------------------
// Block bloom bits mapping functions
// Required by Web3 API.
// ----------------------------------------------------------------------------

// GetBlockBloom gets bloombits from block height
func (k Keeper) GetBlockBloom(ctx sdk.Context, height int64) (ethtypes.Bloom, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBloom)
	bz := store.Get(evmtypes.BloomKey(height))
	if len(bz) == 0 {
		return ethtypes.Bloom{}, false
	}

	return ethtypes.BytesToBloom(bz), true
}

// SetBlockBloom sets the mapping from block height to bloom bits
func (k Keeper) SetBlockBloom(ctx sdk.Context, height int64, bloom ethtypes.Bloom) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), evmtypes.KeyPrefixBloom)
	store.Set(evmtypes.BloomKey(height), bloom.Bytes())
}

// GetAllTxLogs return all the transaction logs from the store.
func (k Keeper) GetAllTxLogs(ctx sdk.Context) []evmtypes.TransactionLogs {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, evmtypes.KeyPrefixLogs)
	defer iterator.Close()

	txsLogs := []evmtypes.TransactionLogs{}
	for ; iterator.Valid(); iterator.Next() {
		hash := ethcmn.BytesToHash(iterator.Key())
		var logs []*ethtypes.Log
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &logs)

		// add a new entry
		txLog := evmtypes.NewTransactionLogs(hash, logs)
		txsLogs = append(txsLogs, txLog)
	}
	return txsLogs
}
