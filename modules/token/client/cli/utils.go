package cli

import (
	"context"

	"github.com/orientwalt/htdf/client"

	"github.com/orientwalt/htdf/modules/token/types"
)

// queryTokenFees retrieves the fees of issuance and minting for the specified symbol
func queryTokenFees(cliCtx client.Context, symbol string) (types.QueryFeesResponse, error) {
	queryClient := types.NewQueryClient(cliCtx)

	resp, err := queryClient.Fees(context.Background(), &types.QueryFeesRequest{Symbol: symbol})
	if err != nil {
		return types.QueryFeesResponse{}, err
	}

	return *resp, err
}
