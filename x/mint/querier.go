package mint

import (
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
)

// Query endpoints supported by the minting querier
const (
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
	QueryBlockRewards     = "rewards"
	QueryTotalProvisions  = "total_provisions"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, k)

		case QueryInflation:
			return queryInflation(ctx, k)

		case QueryAnnualProvisions:
			return queryAnnualProvisions(ctx, k)

		case QueryTotalProvisions:
			return queryTotalProvisions(ctx, k)

		case QueryBlockRewards:
			return queryBlockRewards(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

func queryInflation(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, minter.Inflation)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

func queryAnnualProvisions(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, minter.AnnualProvisions)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

// This is for block-reward query inferencing(cmd & rest)
type TotalProvisions struct {
	Provision sdk.Int
}

func NewTotalProvisions(p sdk.Int) TotalProvisions {
	return TotalProvisions{
		Provision: p,
	}
}

func (tp TotalProvisions) String() string {
	return tp.Provision.String()
}

// junying-todo, 2020-03-09
func queryTotalProvisions(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {

	res, err := codec.MarshalJSONIndent(k.cdc, NewTotalProvisions(k.sk.TotalTokens(ctx)))
	// res, err := codec.MarshalJSONIndent(k.cdc, k.sk.TotalTokens(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

// defines the params for query: "custom/mint/rewards"
type QueryBlockRewardParams struct {
	Height int64
}

func NewQueryBlockRewardParams(h int64) QueryBlockRewardParams {
	return QueryBlockRewardParams{
		Height: h,
	}
}

// This is for block-reward query inferencing(cmd & rest)
type BlockReward struct {
	Reward int64
}

func NewBlockReward(r int64) BlockReward {
	return BlockReward{
		Reward: r,
	}
}

func (br BlockReward) String() string {
	return strconv.FormatInt(int64(br.Reward), 10)
}

// querydata part that queries a block reward
// it can be called by route "/custom/mint/rewards" from cmd & rest
func queryBlockRewards(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryBlockRewardParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	reward := keeper.GetReward(ctx, params.Height)
	if reward < 0 {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("height %s does not exist", strconv.FormatInt(params.Height, 10)))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, NewBlockReward(reward))
	// res, err := codec.MarshalJSONIndent(keeper.cdc, sdk.NewInt(reward))
	// res, err := codec.MarshalJSONIndent(keeper.cdc, reward)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}
