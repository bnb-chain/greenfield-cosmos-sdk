package keeper_test

import (
	gocontext "context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

func (suite *KeeperTestSuite) TestQueryParams() {
	// It should be empty params
	res, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(suite.gashubKeeper.GetParams(suite.ctx), res.GetParams())
}

func (suite *KeeperTestSuite) TestQueryMsgGasParams() {
	ctx, gashubKeeper := suite.ctx, suite.gashubKeeper

	fixed := types.MsgGasParams{
		MsgTypeUrl: "fixed",
		GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
	}
	dynamic := types.MsgGasParams{
		MsgTypeUrl: "dynamic",
		GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
	}

	gashubKeeper.SetMsgGasParams(ctx, fixed)
	gashubKeeper.SetMsgGasParams(ctx, dynamic)

	tests := []struct {
		name string
		req  *types.QueryMsgGasParamsRequest
		exp  *types.QueryMsgGasParamsResponse
	}{
		{
			name: "empty urls list",
			req:  &types.QueryMsgGasParamsRequest{MsgTypeUrls: []string{}},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: []*types.MsgGasParams{
					{
						MsgTypeUrl: "fixed",
						GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
					},
					{
						MsgTypeUrl: "dynamic",
						GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: nil,
					Total:   2,
				},
			},
		},
		{
			name: "limit 1",
			req: &types.QueryMsgGasParamsRequest{
				Pagination: &query.PageRequest{
					Limit:      1,
					CountTotal: true,
				},
			},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: []*types.MsgGasParams{
					{
						MsgTypeUrl: "fixed",
						GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
					},
				},
				Pagination: &query.PageResponse{
					NextKey: types.LengthPrefix([]byte("dynamic")),
					Total:   2,
				},
			},
		},
		{
			name: "just fixed",
			req:  &types.QueryMsgGasParamsRequest{MsgTypeUrls: []string{"fixed"}},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: []*types.MsgGasParams{
					{
						MsgTypeUrl: "fixed",
						GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
					},
				},
				Pagination: nil,
			},
		},
		{
			name: "just dynamic",
			req:  &types.QueryMsgGasParamsRequest{MsgTypeUrls: []string{"dynamic"}},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: []*types.MsgGasParams{
					{
						MsgTypeUrl: "dynamic",
						GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
					},
				},
				Pagination: nil,
			},
		},
		{
			name: "just an unknown type",
			req:  &types.QueryMsgGasParamsRequest{MsgTypeUrls: []string{"unknown"}},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: nil,
				Pagination:   nil,
			},
		},
		{
			name: "both fixed dynamic",
			req:  &types.QueryMsgGasParamsRequest{MsgTypeUrls: []string{"fixed", "dynamic"}},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: []*types.MsgGasParams{
					{
						MsgTypeUrl: "fixed",
						GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
					},
					{
						MsgTypeUrl: "dynamic",
						GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
					},
				},
				Pagination: nil,
			},
		},
		{
			name: "both fixed dynamic and an unknown",
			req:  &types.QueryMsgGasParamsRequest{MsgTypeUrls: []string{"fixed", "dynamic", "unknown"}},
			exp: &types.QueryMsgGasParamsResponse{
				MsgGasParams: []*types.MsgGasParams{
					{
						MsgTypeUrl: "fixed",
						GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
					},
					{
						MsgTypeUrl: "dynamic",
						GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
					},
				},
				Pagination: nil,
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			resp, err := suite.queryClient.MsgGasParams(gocontext.Background(), tc.req)
			suite.Require().NoError(err)
			if !suite.Assert().Equal(tc.exp, resp) {
				if !suite.Assert().Len(resp.MsgGasParams, len(tc.exp.MsgGasParams)) {
					for i := range tc.exp.MsgGasParams {
						suite.Assert().Equal(tc.exp.MsgGasParams[i].MsgTypeUrl, resp.MsgGasParams[i].MsgTypeUrl, fmt.Sprintf("MsgGasParams[%d].MsgTypeUrl", i))
						suite.Assert().Equal(tc.exp.MsgGasParams[i].GasParams, resp.MsgGasParams[i].GasParams, fmt.Sprintf("MsgGasParams[%d].GasParams", i))
					}
				}
				if !suite.Assert().Equal(tc.exp.Pagination, resp.Pagination, "Pagination") && tc.exp.Pagination != nil && resp.Pagination != nil {
					suite.Assert().Equal(tc.exp.Pagination.NextKey, resp.Pagination.NextKey, "Pagination.NextKey")
					suite.Assert().Equal(tc.exp.Pagination.Total, resp.Pagination.Total, "Pagination.Total")
				}
			}
			suite.Require().Equal(tc.exp, resp)
		})
	}
}
