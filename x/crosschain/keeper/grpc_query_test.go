package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
)

func (s *TestSuite) TestQueryParams() {
	res, err := s.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(s.crossChainKeeper.GetParams(s.ctx), res.GetParams())
}
