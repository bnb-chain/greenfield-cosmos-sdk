package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

func (suite *TestSuite) TestQueryParams() {
	res, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(suite.app.OracleKeeper.GetParams(suite.ctx), res.GetParams())
}
