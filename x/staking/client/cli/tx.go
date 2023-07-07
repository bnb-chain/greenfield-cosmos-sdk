package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"cosmossdk.io/math"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// default values
var (
	DefaultTokens                  = sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	defaultAmount                  = DefaultTokens.String() + sdk.DefaultBondDenom
	defaultCommissionRate          = "0.1"
	defaultCommissionMaxRate       = "0.2"
	defaultCommissionMaxChangeRate = "0.01"
	defaultMinSelfDelegation       = "1"
)

// NewTxCmd returns a root CLI command handler for all x/staking transaction commands.
func NewTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		NewCreateValidatorCmd(),
		NewEditValidatorCmd(),
		NewDelegateCmd(),

		NewUnbondCmd(),
		NewCancelUnbondingDelegation(),
	)

	return stakingTxCmd
}

// NewCreateValidatorCmd returns a CLI command handler for creating a MsgCreateValidator transaction.
func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator [path/to/create_validator_proposal.json]",
		Short: "submit a create new validator proposal",
		Args:  cobra.ExactArgs(1),
		Long:  `Submit a create new validator proposal by submitting a JSON file with the new validator details, once the proposal has been passed, create a new validator initialized with a self-delegation.`,
		Example: strings.TrimSpace(
			fmt.Sprintf(`
$ %s tx staking create-validator path/to/create_validator_proposal.json --from keyname

Where create_validator_proposal.json contains:

{
	"messages": [
		{
			"@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
			"description": {
				"moniker": "${NODE_NAME}",
				"identity": "",
				"website": "",
				"security_contact": "",
				"details": ""
			},
			"commission": {
				"rate": "0.070000000000000000",
				"max_rate": "1.000000000000000000",
				"max_change_rate": "0.010000000000000000"
			},
			"min_self_delegation": "1000000000000000000000",
			"delegator_address": "${VALIDATOR_ADDR}",
			"validator_address": "${VALIDATOR_ADDR}",
			"pubkey": {
				"@type": "/cosmos.crypto.ed25519.PubKey",
				"key": "${VALIDATOR_NODE_PUB_KEY}"
			},
			"value": {
				"denom": "BNB",
				"amount": "1000000000000000000000"
			},
			"from": "0x7b5Fe22B5446f7C62Ea27B8BD71CeF94e03f3dF2",
			"relayer_address": "${RELAYER_ADDR}",
			"challenger_address": "${CHALLENGER_ADDR}",
			"bls_key": "${VALIDATOR_BLS}"
		}
	],
	"metadata": "",
	"title": "Create ${NODE_NAME} Validator",
	"summary": "create ${NODE_NAME} validator",
	"deposit": "1000000000000000000BNB"
}

modify the related configrations as you need, where you can get the pubkey using "%s tendermint show-validator"
`, version.AppName, version.AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msgs, metadata, title, summary, deposit, err := govcli.ParseSubmitProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			govMsg, err := v1.NewMsgSubmitProposal(msgs, deposit, clientCtx.GetFromAddress().String(), metadata, title, summary)
			if err != nil {
				return fmt.Errorf("invalid message: %w", err)
			}

			if len(msgs) != 1 {
				return fmt.Errorf("invalid message length: %d", len(msgs))
			}

			valMsg, ok := msgs[0].(*types.MsgCreateValidator)
			if !ok || valMsg.ValidateBasic() != nil {
				return fmt.Errorf("invalid create validator message")
			}

			delAddr, err := sdk.AccAddressFromHexUnsafe(valMsg.DelegatorAddress)
			if err != nil {
				return err
			}
			if !delAddr.Equals(clientCtx.GetFromAddress()) {
				return fmt.Errorf("the from address should be the self delegator address: %s", delAddr.String())
			}

			valAddr, err := sdk.AccAddressFromHexUnsafe(valMsg.ValidatorAddress)
			if err != nil {
				return err
			}

			grantee := authtypes.NewModuleAddress(govtypes.ModuleName)
			authorization, err := types.NewStakeAuthorization([]sdk.AccAddress{valAddr}, nil, types.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE, &valMsg.Value)
			if err != nil {
				return err
			}
			authzMsg, err := authz.NewMsgGrant(clientCtx.GetFromAddress(), grantee, authorization, nil)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), authzMsg, govMsg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewEditValidatorCmd returns a CLI command handler for creating a MsgEditValidator transaction.
func NewEditValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator",
		Short: "edit an existing validator account",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			valAddr := clientCtx.GetFromAddress()
			moniker, _ := cmd.Flags().GetString(FlagEditMoniker)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			security, _ := cmd.Flags().GetString(FlagSecurityContact)
			details, _ := cmd.Flags().GetString(FlagDetails)
			description := types.NewDescription(moniker, identity, website, security, details)
			relayer := sdk.AccAddress("")
			challenger := sdk.AccAddress("")

			var newRate *sdk.Dec

			commissionRate, _ := cmd.Flags().GetString(FlagCommissionRate)
			if commissionRate != "" {
				rate, err := sdk.NewDecFromStr(commissionRate)
				if err != nil {
					return fmt.Errorf("invalid new commission rate: %v", err)
				}

				newRate = &rate
			}

			var newMinSelfDelegation *math.Int

			minSelfDelegationString, _ := cmd.Flags().GetString(FlagMinSelfDelegation)
			if minSelfDelegationString != "" {
				msb, ok := sdk.NewIntFromString(minSelfDelegationString)
				if !ok {
					return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
				}

				newMinSelfDelegation = &msb
			}

			relayerAddr, _ := cmd.Flags().GetString(FlagAddressRelayer)
			if relayerAddr != "" {
				relayer, err = sdk.AccAddressFromHexUnsafe(relayerAddr)
				if err != nil {
					return fmt.Errorf("invalid relayer address: %v", err)
				}
			}

			challengerAddr, _ := cmd.Flags().GetString(FlagAddressChallenger)
			if challengerAddr != "" {
				challenger, err = sdk.AccAddressFromHexUnsafe(challengerAddr)
				if err != nil {
					return fmt.Errorf("invalid challenger address: %v", err)
				}
			}

			blsPk, _ := cmd.Flags().GetString(FlagBlsKey)
			blsProof, _ := cmd.Flags().GetString(FlagBlsProof)

			msg := types.NewMsgEditValidator(
				valAddr, description, newRate, newMinSelfDelegation,
				relayer, challenger, blsPk, blsProof,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetDescriptionEdit())
	cmd.Flags().AddFlagSet(flagSetCommissionUpdate())
	cmd.Flags().AddFlagSet(FlagSetMinSelfDelegation())
	cmd.Flags().AddFlagSet(FlagSetRelayerAddress())
	cmd.Flags().AddFlagSet(FlagSetChallengerAddress())
	cmd.Flags().AddFlagSet(FlagSetBlsKey())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewDelegateCmd returns a CLI command handler for creating a MsgDelegate transaction.
func NewDelegateCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "delegate [validator-addr] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Delegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate an amount of liquid coins to a validator from your wallet.

Example:
$ %s tx staking delegate %s1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 1000stake --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.AccAddressFromHexUnsafe(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegate(delAddr, valAddr, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRedelegateCmd returns a CLI command handler for creating a MsgBeginRedelegate transaction.
func NewRedelegateCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "redelegate [src-validator-addr] [dst-validator-addr] [amount]",
		Short: "Redelegate illiquid tokens from one validator to another",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Redelegate an amount of illiquid staking tokens from one validator to another.

Example:
$ %s tx staking redelegate %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj %s1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 100stake --from mykey
`,
				version.AppName, bech32PrefixValAddr, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr := clientCtx.GetFromAddress()
			valSrcAddr, err := sdk.AccAddressFromHexUnsafe(args[0])
			if err != nil {
				return err
			}

			valDstAddr, err := sdk.AccAddressFromHexUnsafe(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnbondCmd returns a CLI command handler for creating a MsgUndelegate transaction.
func NewUnbondCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "unbond [validator-addr] [amount]",
		Short: "Unbond shares from a validator",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbond an amount of bonded shares from a validator.

Example:
$ %s tx staking unbond %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.AccAddressFromHexUnsafe(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUndelegate(delAddr, valAddr, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCancelUnbondingDelegation returns a CLI command handler for creating a MsgCancelUnbondingDelegation transaction.
func NewCancelUnbondingDelegation() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "cancel-unbond [validator-addr] [amount] [creation-height]",
		Short: "Cancel unbonding delegation and delegate back to the validator",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel Unbonding Delegation and delegate back to the validator.

Example:
$ %s tx staking cancel-unbond %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake 2 --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		Example: fmt.Sprintf(`$ %s tx staking cancel-unbond %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake 2 --from mykey`,
			version.AppName, bech32PrefixValAddr),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr := clientCtx.GetFromAddress()
			valAddr, err := sdk.AccAddressFromHexUnsafe(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			creationHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return sdkerrors.Wrap(fmt.Errorf("invalid height: %d", creationHeight), "invalid height")
			}

			msg := types.NewMsgCancelUnbondingDelegation(delAddr, valAddr, creationHeight, amount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Return the flagset, particular flags, and a description of defaults
// this is anticipated to be used with the gen-tx
func CreateValidatorMsgFlagSet(ipDefault string) (fs *flag.FlagSet, defaultsDesc string) {
	fsCreateValidator := flag.NewFlagSet("", flag.ContinueOnError)
	fsCreateValidator.String(FlagIP, ipDefault, "The node's public P2P IP")
	fsCreateValidator.Uint(FlagP2PPort, 26656, "The node's public P2P port")
	fsCreateValidator.String(FlagNodeID, "", "The node's NodeID")
	fsCreateValidator.String(FlagMoniker, "", "The validator's (optional) moniker")
	fsCreateValidator.String(FlagWebsite, "", "The validator's (optional) website")
	fsCreateValidator.String(FlagSecurityContact, "", "The validator's (optional) security contact email")
	fsCreateValidator.String(FlagDetails, "", "The validator's (optional) details")
	fsCreateValidator.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	fsCreateValidator.AddFlagSet(FlagSetCommissionCreate())
	fsCreateValidator.AddFlagSet(FlagSetMinSelfDelegation())
	fsCreateValidator.AddFlagSet(FlagSetAmount())
	fsCreateValidator.AddFlagSet(FlagSetPublicKey())

	defaultsDesc = fmt.Sprintf(`
	delegation amount:           %s
	commission rate:             %s
	commission max rate:         %s
	commission max change rate:  %s
	minimum self delegation:     %s
`, defaultAmount, defaultCommissionRate,
		defaultCommissionMaxRate, defaultCommissionMaxChangeRate,
		defaultMinSelfDelegation)

	return fsCreateValidator, defaultsDesc
}

type TxCreateValidatorConfig struct {
	ChainID string
	NodeID  string
	Moniker string

	Amount string

	CommissionRate          string
	CommissionMaxRate       string
	CommissionMaxChangeRate string
	MinSelfDelegation       string

	PubKey cryptotypes.PubKey

	IP              string
	P2PPort         uint
	Website         string
	SecurityContact string
	Details         string
	Identity        string

	Validator  sdk.AccAddress
	Delegator  sdk.AccAddress
	Relayer    sdk.AccAddress
	Challenger sdk.AccAddress
	BlsKey     string
	BLSProof   string
}

func PrepareConfigForTxCreateValidator(flagSet *flag.FlagSet, moniker, nodeID, chainID string, valPubKey cryptotypes.PubKey) (TxCreateValidatorConfig, error) {
	c := TxCreateValidatorConfig{}

	ip, err := flagSet.GetString(FlagIP)
	if err != nil {
		return c, err
	}

	if ip == "" {
		_, _ = fmt.Fprintf(os.Stderr, "failed to retrieve an external IP; the tx's memo field will be unset")
	}

	p2pPort, err := flagSet.GetUint(FlagP2PPort)
	if err != nil {
		return c, err
	}

	website, err := flagSet.GetString(FlagWebsite)
	if err != nil {
		return c, err
	}

	securityContact, err := flagSet.GetString(FlagSecurityContact)
	if err != nil {
		return c, err
	}

	details, err := flagSet.GetString(FlagDetails)
	if err != nil {
		return c, err
	}

	identity, err := flagSet.GetString(FlagIdentity)
	if err != nil {
		return c, err
	}

	c.Amount, err = flagSet.GetString(FlagAmount)
	if err != nil {
		return c, err
	}

	c.CommissionRate, err = flagSet.GetString(FlagCommissionRate)
	if err != nil {
		return c, err
	}

	c.CommissionMaxRate, err = flagSet.GetString(FlagCommissionMaxRate)
	if err != nil {
		return c, err
	}

	c.CommissionMaxChangeRate, err = flagSet.GetString(FlagCommissionMaxChangeRate)
	if err != nil {
		return c, err
	}

	c.MinSelfDelegation, err = flagSet.GetString(FlagMinSelfDelegation)
	if err != nil {
		return c, err
	}

	c.IP = ip
	c.P2PPort = p2pPort
	c.Website = website
	c.SecurityContact = securityContact
	c.Identity = identity
	c.NodeID = nodeID
	c.PubKey = valPubKey
	c.Website = website
	c.SecurityContact = securityContact
	c.Details = details
	c.Identity = identity
	c.ChainID = chainID
	c.Moniker = moniker

	if c.Amount == "" {
		c.Amount = defaultAmount
	}

	if c.CommissionRate == "" {
		c.CommissionRate = defaultCommissionRate
	}

	if c.CommissionMaxRate == "" {
		c.CommissionMaxRate = defaultCommissionMaxRate
	}

	if c.CommissionMaxChangeRate == "" {
		c.CommissionMaxChangeRate = defaultCommissionMaxChangeRate
	}

	if c.MinSelfDelegation == "" {
		c.MinSelfDelegation = defaultMinSelfDelegation
	}

	return c, nil
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(clientCtx client.Context, config TxCreateValidatorConfig, txBldr tx.Factory, generateOnly bool) (tx.Factory, sdk.Msg, error) {
	amounstStr := config.Amount
	amount, err := sdk.ParseCoinNormalized(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	from := clientCtx.GetFromAddress()

	description := types.NewDescription(
		config.Moniker,
		config.Identity,
		config.Website,
		config.SecurityContact,
		config.Details,
	)

	// get the initial validator commission parameters
	rateStr := config.CommissionRate
	maxRateStr := config.CommissionMaxRate
	maxChangeRateStr := config.CommissionMaxChangeRate
	commissionRates, err := buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txBldr, nil, err
	}

	// get the initial validator min self delegation
	msbStr := config.MinSelfDelegation
	minSelfDelegation, ok := sdk.NewIntFromString(msbStr)

	if !ok {
		return txBldr, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	msg, err := types.NewMsgCreateValidator(
		config.Validator, config.PubKey,
		amount, description, commissionRates, minSelfDelegation,
		from, config.Delegator, config.Relayer, config.Challenger, config.BlsKey, config.BLSProof)
	if err != nil {
		return txBldr, msg, err
	}

	if generateOnly {
		ip := config.IP
		p2pPort := config.P2PPort
		nodeID := config.NodeID

		if nodeID != "" && ip != "" && p2pPort > 0 {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:%d", nodeID, ip, p2pPort))
		}
	}

	return txBldr, msg, nil
}
