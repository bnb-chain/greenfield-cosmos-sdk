package crosschain

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/testutil"
	crosschaintypes "github.com/cosmos/cosmos-sdk/x/crosschain/types"

	"github.com/cosmos/gogoproto/proto"
)

func (s *E2ETestSuite) TestQueryGRPC() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress
	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"gRPC request params",
			fmt.Sprintf("%s/cosmos/crosschain/v1/params", baseURL),
			map[string]string{},
			&crosschaintypes.QueryParamsResponse{},
			&crosschaintypes.QueryParamsResponse{
				Params: crosschaintypes.DefaultParams(),
			},
		},
	}
	for _, tc := range testCases {
		resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
		s.Run(tc.name, func() {
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
