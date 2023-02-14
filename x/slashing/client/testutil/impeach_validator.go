package testutil

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type testParams struct {
	discription string
	proposalID  string
	from        string
}

var testcases = []testParams{
	{
		discription: "impeach the last validator",
		proposalID:  "1",
		from:        "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
	},
	{
		discription: "impeach the last validator again",
		proposalID:  "2",
		from:        "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
	},
	{
		discription: "impeach the last validator, but invalid from",
		proposalID:  "3",
		from:        "0x8b70dC9B691fCeB4e1c69dF8cbF8c077AD4b5853",
	},
}

type ImpeachValidatorTestSuite struct {
	suite.Suite

	cfg         network.Config
	network     *network.Network
	proposalIDs []string
}

func NewImpeachValidatorTestSuite(cfg network.Config) *ImpeachValidatorTestSuite {
	return &ImpeachValidatorTestSuite{cfg: cfg}
}

func (s *ImpeachValidatorTestSuite) SetupSuite() {
	s.T().Log("setting up test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	for _, testcase := range testcases {
		s.submitProposal(testcase)
		s.voteProposal(testcase.proposalID)
		s.proposalIDs = append(s.proposalIDs, testcase.proposalID)
	}
}

func (s *ImpeachValidatorTestSuite) SetupNewSuite() {
	s.T().Log("setting up new test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *ImpeachValidatorTestSuite) TearDownSuite() {
	s.T().Log("tearing down test suite")
	s.network.Cleanup()
}

func (s *ImpeachValidatorTestSuite) TestQuerySuccessfulImpeachedValidator() {
	proposalID := s.proposalIDs[0]
	targetVal := s.network.Validators[len(s.network.Validators)-1]
	clientCtx := s.network.Validators[0].ClientCtx

	// query proposal
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().NoError(err)
	var proposal v1.Proposal
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &proposal), out.String())
	s.Require().Equal(v1.ProposalStatus_PROPOSAL_STATUS_PASSED, proposal.Status, out.String())

	// query validator
	queryCmd := cli.GetCmdQueryValidator()
	res, err := clitestutil.ExecTestCLICmd(
		clientCtx, queryCmd,
		[]string{targetVal.Address.String(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
	)
	s.Require().NoError(err)
	var result types.Validator
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
	s.Require().NotEqual(result.GetStatus(), types.Bonded, fmt.Sprintf("validator %s not in bonded status", targetVal.Address.String()))
	s.Require().Equal(result.Jailed, true)
}

func (s *ImpeachValidatorTestSuite) TestQueryDoubleImpeachedValidator() {
	proposalID := s.proposalIDs[1]
	targetVal := s.network.Validators[len(s.network.Validators)-1]
	clientCtx := s.network.Validators[0].ClientCtx

	// query proposal
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().NoError(err)
	var proposal v1.Proposal
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &proposal), out.String())
	s.Require().Equal(v1.ProposalStatus_PROPOSAL_STATUS_PASSED, proposal.Status, out.String())

	// query validator
	queryCmd := cli.GetCmdQueryValidator()
	res, err := clitestutil.ExecTestCLICmd(
		clientCtx, queryCmd,
		[]string{targetVal.Address.String(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
	)
	s.Require().NoError(err)
	var result types.Validator
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
	s.Require().NotEqual(result.GetStatus(), types.Bonded, fmt.Sprintf("validator %s not in bonded status", targetVal.Address.String()))
	s.Require().Equal(result.Jailed, true)
}

func (s *ImpeachValidatorTestSuite) TestQueryInvalidFromAddress() {
	proposalID := s.proposalIDs[2]
	clientCtx := s.network.Validators[0].ClientCtx

	// query proposal, should not be found because of invalid from the proposal will be rejected.
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	_, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().Error(err)
}

func (s *ImpeachValidatorTestSuite) submitProposal(params testParams) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// Always impeach the last validator.
	targetVal := s.network.Validators[len(s.network.Validators)-1]
	args := []string{
		s.impeachValidatorProposal(targetVal.Address, params.from).Name(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.ValAddress.String()),
		fmt.Sprintf("--gas=%s", fmt.Sprintf("%d", flags.DefaultGasLimit+100000)),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
	}

	_, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdSubmitProposal(), args)
	s.Require().NoError(err)
}

func (s *ImpeachValidatorTestSuite) impeachValidatorProposal(valAddr sdk.AccAddress, from string) *os.File {
	propMetadata := []byte{42}
	proposal := fmt.Sprintf(`
{
	"messages": [
		{
			"@type": "/cosmos.slashing.v1beta1.MsgImpeach",
			"from":"%s",
			"validator_address":"%s"
		}
	],
	"metadata": "%s",
	"deposit": "%s"
}`,
		from,
		valAddr.String(),
		base64.StdEncoding.EncodeToString(propMetadata),
		sdk.NewCoin(s.cfg.BondDenom, v1.DefaultMinDepositTokens),
	)

	proposalFile := testutil.WriteToNewTempFile(s.T(), proposal)
	return proposalFile
}

func (s *ImpeachValidatorTestSuite) voteProposal(proposalID string) {
	for i := 0; i < len(s.network.Validators); i++ {
		clientCtx := s.network.Validators[i].ClientCtx
		clientCtx.Client = s.network.Validators[0].RPCClient

		// The last validator vote no, others vote yes.
		voteOption := "yes"
		if i == len(s.network.Validators)-1 {
			voteOption = "no"
		}

		args := []string{
			proposalID,
			voteOption,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, s.network.Validators[i].Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
		}

		_, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdVote(), args)
		s.Require().NoError(err)
	}
}
