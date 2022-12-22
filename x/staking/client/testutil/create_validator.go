package testutil

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authztestutil "github.com/cosmos/cosmos-sdk/x/authz/client/testutil"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

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

	newVal := s.submitProposal()
	s.validators = append(s.validators, newVal)
	proposalID := fmt.Sprintf("%d", 1)
	s.proposalIDs = append(s.proposalIDs, proposalID)
	s.voteProposal(proposalID)
}

func (s *CreateValidatorTestSuite) SetupNewSuite() {
	s.T().Log("setting up new test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *CreateValidatorTestSuite) submitProposal() sdk.AccAddress {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// Get coin from current validator
	k, _, err := val.ClientCtx.Keyring.NewMnemonic("NewAccount", keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)

	pub, err := k.GetPubKey()
	s.Require().NoError(err)

	newVal := sdk.AccAddress(pub.Address())
	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newVal,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200000000))),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)

	// Grant delegate authorization
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

	args := append([]string{
		s.createValidatorProposal(newVal).Name(),
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

func (s *CreateValidatorTestSuite) TearDownSuite() {
	s.T().Log("tearing down test suite")
	s.network.Cleanup()
}

func (s *CreateValidatorTestSuite) TestQuerySuccessfulCreatedValidator() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
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
		val.ClientCtx, queryCmd,
		[]string{newVal.String(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
	)
	s.Require().NoError(err)
	var result types.Validator
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(res.Bytes(), &result))
	s.Require().Equal(result.GetStatus(), types.Bonded, fmt.Sprintf("validator %s not in bonded status", newVal.String()))
}

func (s *CreateValidatorTestSuite) createValidatorProposal(valAddr sdk.AccAddress) *os.File {
	pubKey := base64.StdEncoding.EncodeToString(ed25519.GenPrivKey().PubKey().Bytes())

	blsSecretKey, _ := bls.RandKey()
	blsPk := hex.EncodeToString(blsSecretKey.PublicKey().Marshal())

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
			"relayer_blskey":"%s"
		}
	],
	"metadata": "%s",
	"deposit": "%s"
}`,
		valAddr.String(),
		valAddr.String(),
		pubKey,
		authtypes.NewModuleAddress(gov.ModuleName),
		valAddr.String(),
		blsPk,
		base64.StdEncoding.EncodeToString(propMetadata),
		sdk.NewCoin(s.cfg.BondDenom, v1.DefaultMinDepositTokens),
	)

	proposalFile := testutil.WriteToNewTempFile(s.T(), proposal)
	return proposalFile
}
