package testutil

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	authztestutil "github.com/cosmos/cosmos-sdk/x/authz/client/testutil"

	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
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
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	ethHd "github.com/evmos/ethermint/crypto/hd"
)

type testParams struct {
	description       string
	proposalID        string
	enableGrant       bool
	from              string
	relayerAddress    string
	challengerAddress string
	blsKey            string
}

var testcases = []testParams{
	{
		description:       "crate success",
		proposalID:        "1",
		enableGrant:       true,
		from:              "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
		relayerAddress:    "0x8b70dC9B691fCeB4e1c69dF8cbF8c077AD4b5853",
		challengerAddress: "0xD1d6bF74282782B0b3eb1413c901D6eCF02e8e28",
		blsKey:            "926f1853304b482634d1ca8ef9652524228202d2a3ca376f3e5ab040430c42703f9de038a41917f944cf4e653cbed45a",
	},
	{
		description:       "no grant delegate authorization to gov module account",
		proposalID:        "2",
		enableGrant:       false,
		from:              "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
		relayerAddress:    "0x36B810C7E246b1c042335bAF1933Da71C5B5AeC5",
		challengerAddress: "0x95CD3fB2037847780fA454DA719e6a01F43D4E2f",
		blsKey:            "a4c392f233e140b124c35a509d8880c52bd9bc3d13042fec3d84c1d55cc4d5be30fe3003d096bc33bc924af563a8de9f",
	},
	{
		description:       "duplicated relayer address",
		proposalID:        "3",
		enableGrant:       true,
		from:              "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
		relayerAddress:    "0x8b70dC9B691fCeB4e1c69dF8cbF8c077AD4b5853",
		challengerAddress: "0x95CD3fB2037847780fA454DA719e6a01F43D4E2f",
		blsKey:            "819a55435aed37ff4f917d3644c74da1f6354c0d4fbc7ca425ebe6b98a0416fc4ef5eb13fceffb3a109eec4af12bba0a",
	},
	{
		description:       "duplicated bls pub key",
		proposalID:        "4",
		enableGrant:       true,
		from:              "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
		relayerAddress:    "0x9E6de7BF11C459E8dA3a5f36c1A87A6FfaDCbA9d",
		challengerAddress: "0x95CD3fB2037847780fA454DA719e6a01F43D4E2f",
		blsKey:            "926f1853304b482634d1ca8ef9652524228202d2a3ca376f3e5ab040430c42703f9de038a41917f944cf4e653cbed45a",
	},
	{
		description:       "duplicated challenger address",
		proposalID:        "5",
		enableGrant:       true,
		from:              "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
		relayerAddress:    "0x8b70dC9B691fCeB4e1c69dF8cbF8c077AD4b5853",
		challengerAddress: "0xD1d6bF74282782B0b3eb1413c901D6eCF02e8e28",
		blsKey:            "819a55435aed37ff4f917d3644c74da1f6354c0d4fbc7ca425ebe6b98a0416fc4ef5eb13fceffb3a109eec4af12bba0a",
	},
	{
		description:       "invalid from",
		proposalID:        "6",
		enableGrant:       true,
		from:              "0x8b70dC9B691fCeB4e1c69dF8cbF8c077AD4b5853",
		relayerAddress:    "0x73Fd0b049bBF30b3A165e5b2eAf9895B015A2515",
		challengerAddress: "0x95CD3fB2037847780fA454DA719e6a01F43D4E2f",
		blsKey:            "b700697ccca38c0d56cde325571dc5412de99a0fa2c28eeb1078cb4eeaa31c238720bee9cce47422e986a7af4f3b19c8",
	},
}

type CreateValidatorTestSuite struct {
	suite.Suite

	cfg         network.Config
	network     *network.Network
	proposalIDs []string
	validators  []sdk.AccAddress
}

func NewCreateValidatorTestSuite(cfg network.Config) *CreateValidatorTestSuite {
	return &CreateValidatorTestSuite{cfg: cfg}
}

func (s *CreateValidatorTestSuite) SetupSuite() {
	s.T().Log("setting up test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	for _, testcase := range testcases {
		newVal := s.submitProposal(testcase)
		s.voteProposal(testcase.proposalID)
		s.validators = append(s.validators, newVal)
		s.proposalIDs = append(s.proposalIDs, testcase.proposalID)
	}
}

func (s *CreateValidatorTestSuite) SetupNewSuite() {
	s.T().Log("setting up new test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *CreateValidatorTestSuite) TearDownSuite() {
	s.T().Log("tearing down test suite")
	s.network.Cleanup()
}

func (s *CreateValidatorTestSuite) TestQuerySuccessfulCreatedValidator() {
	clientCtx := s.network.Validators[0].ClientCtx
	proposalID := s.proposalIDs[0]
	newVal := s.validators[0]

	// waiting for voting period to end
	time.Sleep(10 * time.Second)

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
		[]string{newVal.String(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
	)
	s.Require().NoError(err)
	var result types.Validator
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
	s.Require().Equal(result.GetStatus(), types.Bonded, fmt.Sprintf("validator %s not in bonded status", newVal.String()))
}

func (s *CreateValidatorTestSuite) TestQueryNoGrantAuthorization() {
	clientCtx := s.network.Validators[0].ClientCtx
	proposalID := s.proposalIDs[1]

	// query proposal
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().NoError(err)
	var proposal v1.Proposal
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &proposal), out.String())
	s.Require().Equal(v1.ProposalStatus_PROPOSAL_STATUS_FAILED, proposal.Status, out.String())
}

func (s *CreateValidatorTestSuite) TestQueryDuplicatedRelayerAddress() {
	clientCtx := s.network.Validators[0].ClientCtx
	proposalID := s.proposalIDs[2]

	// query proposal
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().NoError(err)
	var proposal v1.Proposal
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &proposal), out.String())
	s.Require().Equal(v1.ProposalStatus_PROPOSAL_STATUS_FAILED, proposal.Status, out.String())
}

func (s *CreateValidatorTestSuite) TestQueryDuplicatedBlsKey() {
	clientCtx := s.network.Validators[0].ClientCtx
	proposalID := s.proposalIDs[3]

	// query proposal
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().NoError(err)
	var proposal v1.Proposal
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &proposal), out.String())
	s.Require().Equal(v1.ProposalStatus_PROPOSAL_STATUS_FAILED, proposal.Status, out.String())
}

func (s *CreateValidatorTestSuite) TestQueryDuplicatedChallengerAddress() {
	clientCtx := s.network.Validators[0].ClientCtx
	proposalID := s.proposalIDs[4]

	// query proposal
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().NoError(err)
	var proposal v1.Proposal
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &proposal), out.String())
	s.Require().Equal(v1.ProposalStatus_PROPOSAL_STATUS_FAILED, proposal.Status, out.String())
}

func (s *CreateValidatorTestSuite) TestQueryInvalidFromAddress() {
	clientCtx := s.network.Validators[0].ClientCtx
	proposalID := s.proposalIDs[5]

	// query proposal, should not be found because of invalid from the proposal will be rejected.
	args := []string{proposalID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	_, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.GetCmdQueryProposal(), args)
	s.Require().Error(err)
}

func (s *CreateValidatorTestSuite) submitProposal(params testParams) sdk.AccAddress {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// Get coin from current validator
	k, _, err := clientCtx.Keyring.NewMnemonic("NewAccount", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, ethHd.EthSecp256k1)
	s.Require().NoError(err)

	pub, err := k.GetPubKey()
	s.Require().NoError(err)

	newVal := sdk.AccAddress(pub.Address())
	_, err = banktestutil.MsgSendExec(
		clientCtx,
		val.Address,
		newVal,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(20000000))),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// Grant delegate authorization
	if params.enableGrant {
		_, err = authztestutil.CreateGrant(val, []string{
			authtypes.NewModuleAddress(gov.ModuleName).String(),
			"delegate",
			fmt.Sprintf("--spend-limit=100000000%s", s.cfg.BondDenom),
			fmt.Sprintf("--allowed-validators=%s", newVal.String()),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, newVal.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
		})
		s.Require().NoError(err)
	}

	args := append([]string{
		s.createValidatorProposal(newVal, params.from, params.relayerAddress, params.challengerAddress, params.blsKey).Name(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, newVal.String()),
		fmt.Sprintf("--gas=%s", fmt.Sprintf("%d", flags.DefaultGasLimit+100000)),
	}, commonArgs...)

	_, err = clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdSubmitProposal(), args)
	s.Require().NoError(err)

	return newVal
}

func (s *CreateValidatorTestSuite) voteProposal(proposalID string) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	args := append([]string{
		proposalID,
		"yes",
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
	}, commonArgs...)
	_, err := clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdVote(), args)
	s.Require().NoError(err)
}

func (s *CreateValidatorTestSuite) createValidatorProposal(valAddr sdk.AccAddress, from string,
	relayerAddress string, challengerAddress string, blsKey string,
) *os.File {
	pubKey := base64.StdEncoding.EncodeToString(ed25519.GenPrivKey().PubKey().Bytes())
	propMetadata := []byte{42}
	proposal := fmt.Sprintf(`
{
	"messages": [
		{
			"@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
			"description":{
				"moniker":"testnode",
				"identity":"",
				"website":"",
				"security_contact":"",
				"details":""
			},
			"commission":{
				"rate":"0.100000000000000000",
				"max_rate":"0.200000000000000000",
				"max_change_rate":"0.010000000000000000"
			},
			"min_self_delegation":"10000000",
			"delegator_address":"%s",
			"validator_address":"%s",
			"pubkey":{
				"@type":"/cosmos.crypto.ed25519.PubKey",
				"key":"%s"
			},
			"value":{
				"denom":"stake",
				"amount":"10000000"
			},
			"from":"%s",
			"relayer_address":"%s",
			"challenger_address":"%s",
			"bls_key":"%s"
		}
	],
	"metadata": "%s",
	"deposit": "%s"
}`,
		valAddr.String(),
		valAddr.String(),
		pubKey,
		from,
		relayerAddress,
		challengerAddress,
		blsKey,
		base64.StdEncoding.EncodeToString(propMetadata),
		sdk.NewCoin(s.cfg.BondDenom, v1.DefaultMinDepositTokens),
	)

	proposalFile := testutil.WriteToNewTempFile(s.T(), proposal)
	return proposalFile
}
