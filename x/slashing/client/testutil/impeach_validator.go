package testutil

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/client/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

type ImpeachValidatorTestSuite struct {
	suite.Suite

	cfg         network.Config
	network     *network.Network
	proposalIDs []string
	validators  []sdk.AccAddress
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

	// Invalid proposal (--from address not gov)
	s.submitProposal(s.network.Validators[0].ValAddress.String(), false)

	// Valid proposal
	s.submitProposal(gov.ModuleName, true)
	proposalID := fmt.Sprintf("%d", 1)
	s.proposalIDs = append(s.proposalIDs, proposalID)

	// Valid proposal to impeach the same validator that was already impeached
	s.submitProposal(gov.ModuleName, true)
	proposalIDInvalidImpeached := fmt.Sprintf("%d", 2)
	s.proposalIDs = append(s.proposalIDs, proposalIDInvalidImpeached)

	for _, proposal := range s.proposalIDs {
		s.voteProposal(proposal, s.network.Validators[0], "yes")
		s.voteProposal(proposal, s.network.Validators[1], "yes")
		s.voteProposal(proposal, s.network.Validators[2], "no")
	}

	time.Sleep(10 * time.Second)

	s.TestQuerySuccessfulImpeachedValidator()
	s.TestQueryInvalidAlreadyImpeached()
}

func (s *ImpeachValidatorTestSuite) SetupNewSuite() {
	s.T().Log("setting up new test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *ImpeachValidatorTestSuite) submitProposal(from string, isValidTestCase bool) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	kickedVal := s.network.Validators[2]

	//nolint:staticcheck
	args := append([]string{
		s.ImpeachValidatorProposal(kickedVal.Address, from).Name(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.ValAddress.String()),
		fmt.Sprintf("--gas=%s", fmt.Sprintf("%d", flags.DefaultGasLimit+100000)),
	}, stakingtestutil.CommonArgs...)

	response, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdSubmitProposal(), args)

	if err == nil {
		splits := strings.Split(response.String(), ",")
		splits = strings.Split(splits[2], ":")
		codespace := strings.Trim(splits[1], "\"")
		if string(codespace) != "" {
			err = gov.ErrInvalidProposalMsg
		}
	}

	if isValidTestCase {
		s.Require().NoError(err)
	} else {
		s.Require().Error(err)
	}
}

func (s *ImpeachValidatorTestSuite) getCoin(recipientVal sdk.Address) {
	val := s.network.Validators[0]

	// Get coin from current validator
	_, err := banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		recipientVal.Bytes(),
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200000000))),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)
}

func (s *ImpeachValidatorTestSuite) voteProposal(proposalID string, val *network.Validator, voteOption string) {
	clientCtx := val.ClientCtx
	clientCtx.Client = s.network.Validators[0].RPCClient

	//nolint:staticcheck
	args := append([]string{
		proposalID,
		voteOption,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
	}, stakingtestutil.CommonArgs...)
	_, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdVote(), args)
	s.Require().NoError(err)
}

func (s *ImpeachValidatorTestSuite) TearDownSuite() {
	s.T().Log("tearing down test suite")
	s.network.Cleanup()
}

func (s *ImpeachValidatorTestSuite) TestQuerySuccessfulImpeachedValidator() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	proposalID := s.proposalIDs[0]
	kickedVal := s.network.Validators[2]

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
		val.ClientCtx, queryCmd,
		[]string{kickedVal.Address.String(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
	)
	s.Require().NoError(err)
	var result types.Validator
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
	s.Require().NotEqual(result.GetStatus(), types.Bonded, fmt.Sprintf("validator %s not in bonded status", kickedVal.Address.String()))
	s.Require().Equal(result.Jailed, true)
}

func (s *ImpeachValidatorTestSuite) TestQueryInvalidAlreadyImpeached() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	proposalID := s.proposalIDs[1]
	kickedVal := s.network.Validators[2]

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
		val.ClientCtx, queryCmd,
		[]string{kickedVal.Address.String(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
	)
	s.Require().NoError(err)
	var result types.Validator
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
	s.Require().NotEqual(result.GetStatus(), types.Bonded, fmt.Sprintf("validator %s not in bonded status", kickedVal.Address.String()))
	s.Require().Equal(result.Jailed, true)
}

func (s *ImpeachValidatorTestSuite) ImpeachValidatorProposal(valAddr sdk.AccAddress, from string) *os.File {
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
		authtypes.NewModuleAddress(from),
		valAddr.String(),
		base64.StdEncoding.EncodeToString(propMetadata),
		sdk.NewCoin(s.cfg.BondDenom, v1.DefaultMinDepositTokens),
	)

	proposalFile := testutil.WriteToNewTempFile(s.T(), proposal)
	return proposalFile
}
