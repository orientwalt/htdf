package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"

	"github.com/orientwalt/htdf/utils"
	"github.com/orientwalt/htdf/version"
	evmtypes "github.com/orientwalt/htdf/x/evm/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case evmtypes.QueryProtocolVersion:
			return queryProtocolVersion(keeper)
		case evmtypes.QueryBalance:
			return queryBalance(ctx, path, keeper)
		case evmtypes.QueryBlockNumber:
			return queryBlockNumber(ctx, keeper)
		case evmtypes.QueryStorage:
			return queryStorage(ctx, path, keeper)
		case evmtypes.QueryCode:
			return queryCode(ctx, path, keeper)
		case evmtypes.QueryHashToHeight:
			return queryHashToHeight(ctx, path, keeper)
		case evmtypes.QueryTransactionLogs:
			return queryTransactionLogs(ctx, path, keeper)
		case evmtypes.QueryBloom:
			return queryBlockBloom(ctx, path, keeper)
		case evmtypes.QueryLogs:
			return queryLogs(ctx, keeper)
		case evmtypes.QueryAccount:
			return queryAccount(ctx, path, keeper)
		// case evmtypes.QueryExportAccount:
		// 	return queryExportAccount(ctx, path, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown service query endpoint")
		}
	}
}

func queryProtocolVersion(keeper Keeper) ([]byte, sdk.Error) {
	vers := version.ProtocolVersion

	bz, err := codec.MarshalJSONIndent(keeper.cdc, hexutil.Uint(vers))
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryBalance(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.HexToAddress(path[1])
	balance := keeper.GetBalance(ctx, addr)
	balanceStr, err := utils.MarshalBigInt(balance)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	res := evmtypes.QueryResBalance{Balance: balanceStr}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryBlockNumber(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	num := ctx.BlockHeight()
	bnRes := evmtypes.QueryResBlockNumber{Number: num}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, bnRes)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.HexToAddress(path[1])
	key := ethcmn.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	res := evmtypes.QueryResStorage{Value: val.Bytes()}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}
	return bz, nil
}

func queryCode(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.HexToAddress(path[1])
	code := keeper.GetCode(ctx, addr)
	res := evmtypes.QueryResCode{Code: code}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryHashToHeight(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	blockHash := ethcmn.FromHex(path[1])
	blockNumber, found := keeper.GetBlockHash(ctx, blockHash)
	if !found {
		return []byte{}, sdk.ErrJsonMarshal(fmt.Errorf("block height not found for hash %s", path[1]).Error())
	}

	res := evmtypes.QueryResBlockNumber{Number: blockNumber}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryBlockBloom(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	num, err := strconv.ParseInt(path[1], 10, 64)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(fmt.Errorf("could not unmarshal block height: %w", err).Error())
	}

	bloom, found := keeper.GetBlockBloom(ctx.WithBlockHeight(num), num)
	if !found {
		return nil, sdk.ErrJsonMarshal(fmt.Errorf("block bloom not found for height %d", num).Error())
	}

	res := evmtypes.QueryBloomFilter{Bloom: bloom}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryTransactionLogs(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	txHash := ethcmn.HexToHash(path[1])

	logs, err := keeper.GetLogs(ctx, txHash)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	res := evmtypes.QueryETHLogs{Logs: logs}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}

func queryLogs(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	logs := keeper.AllLogs(ctx)

	res := evmtypes.QueryETHLogs{Logs: logs}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}
	return bz, nil
}

func queryAccount(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.HexToAddress(path[1])
	so := keeper.GetOrNewStateObject(ctx, addr)

	balance, err := utils.MarshalBigInt(so.Balance())
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	res := evmtypes.QueryResAccount{
		Balance:  balance,
		CodeHash: so.CodeHash(),
		Nonce:    so.Nonce(),
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}
	return bz, nil
}

func queryExportAccount(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := ethcmn.HexToAddress(path[1])

	var storage evmtypes.Storage
	keeper.CommitStateDB.ForEachStorage(addr, func(key, value ethcmn.Hash) bool {
		storage = append(storage, evmtypes.NewState(key, value))
		return false
	})
	// if err != nil {
	// 	return nil, err
	// }

	res := evmtypes.GenesisAccount{
		Address: addr,
		Balance: keeper.GetBalance(ctx, addr),
		Code:    keeper.GetCode(ctx, addr),
		Storage: storage,
	}

	// TODO: codec.MarshalJSONIndent doesn't call the String() method of types properly
	bz, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return nil, sdk.ErrJsonMarshal(err.Error())
	}

	return bz, nil
}
