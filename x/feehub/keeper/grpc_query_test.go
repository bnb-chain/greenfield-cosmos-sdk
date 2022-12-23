package keeper_test

import (
	gocontext "context"
	"github.com/cosmos/cosmos-sdk/x/feehub/types"
)

func (suite *IntegrationTestSuite) TestQueryParams() {
	res, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(suite.app.FeehubKeeper.GetParams(suite.ctx), res.GetParams())
}
