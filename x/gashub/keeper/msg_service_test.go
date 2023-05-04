package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/x/gashub/types"
)

func (suite *KeeperTestSuite) TestMsgUpdateParams() {
	// default params
	params := types.DefaultParams()
	invalid1 := params
	invalid2 := params
	invalid1.MaxTxSize = 0
	invalid2.MinGasPerByte = 0

	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority: "invalid",
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "zero max tx size",
			input: &types.MsgUpdateParams{
				Authority: suite.gashubKeeper.GetAuthority(),
				Params:    invalid1,
			},
			expErr:    true,
			expErrMsg: "invalid max tx size",
		},
		{
			name: "zero min gas per byte",
			input: &types.MsgUpdateParams{
				Authority: suite.gashubKeeper.GetAuthority(),
				Params:    invalid2,
			},
			expErr:    true,
			expErrMsg: "invalid min gas per byte",
		},
		{
			name: "all good",
			input: &types.MsgUpdateParams{
				Authority: suite.gashubKeeper.GetAuthority(),
				Params:    params,
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.UpdateParams(suite.ctx, tc.input)

			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgSetMsgGasParams() {
	fixed := types.MsgGasParams{
		MsgTypeUrl: "fixed",
		GasParams:  &types.MsgGasParams_FixedType{FixedType: &types.MsgGasParams_FixedGasParams{FixedGas: 1200}},
	}
	dynamic := types.MsgGasParams{
		MsgTypeUrl: "dynamic",
		GasParams:  &types.MsgGasParams_MultiSendType{MultiSendType: &types.MsgGasParams_DynamicGasParams{FixedGas: 800, GasPerItem: 800}},
	}

	suite.gashubKeeper.SetMsgGasParams(suite.ctx, fixed)

	testCases := []struct {
		name      string
		input     *types.MsgSetMsgGasParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgSetMsgGasParams{
				Authority: "invalid",
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "add new msg gas params",
			input: &types.MsgSetMsgGasParams{
				Authority: suite.gashubKeeper.GetAuthority(),
				UpdateSet: []*types.MsgGasParams{&dynamic},
			},
			expErr: false,
		},
		{
			name: "delete msg gas params",
			input: &types.MsgSetMsgGasParams{
				Authority: suite.gashubKeeper.GetAuthority(),
				DeleteSet: []string{"fixed"},
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.SetMsgGasParams(suite.ctx, tc.input)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
