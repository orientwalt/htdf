package keeper_test

import (
	gocontext "context"

	"github.com/orientwalt/htdf/baseapp"

	"github.com/orientwalt/htdf/modules/mint/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryPoolParameters() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.MintKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	// Query Params
	resp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(app.MintKeeper.GetParamSet(ctx), resp.Params)
}
