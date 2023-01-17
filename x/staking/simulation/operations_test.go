package simulation_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

// TestWeightedOperations tests the weights of the operations.
func TestWeightedOperations(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)
	app, ctx, accs := createTestApp(t, false, r, 3)

	ctx = ctx.WithChainID(simapp.DefaultChainId)

	cdc := app.AppCodec()
	appParams := make(simtypes.AppParams)

	weightesOps := simulation.WeightedOperations(appParams, cdc, app.AccountKeeper,
		app.BankKeeper, app.StakingKeeper,
	)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{simappparams.DefaultWeightMsgEditValidator, types.ModuleName, types.TypeMsgEditValidator},
		{simappparams.DefaultWeightMsgDelegate, types.ModuleName, types.TypeMsgDelegate},
		{simappparams.DefaultWeightMsgUndelegate, types.ModuleName, types.TypeMsgUndelegate},
		{simappparams.DefaultWeightMsgBeginRedelegate, types.ModuleName, types.TypeMsgBeginRedelegate},
		{simappparams.DefaultWeightMsgCancelUnbondingDelegation, types.ModuleName, types.TypeMsgCancelUnbondingDelegation},
	}

	for i, w := range weightesOps {
		operationMsg, _, _ := w.Op()(r, app.BaseApp, ctx, accs, ctx.ChainID())
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		require.Equal(t, expected[i].weight, w.Weight(), "weight should be the same")
		require.Equal(t, expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		require.Equal(t, expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

// TestSimulateMsgCreateValidator tests the normal scenario of a valid message of type TypeMsgCreateValidator.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateValidator(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)
	app, ctx, accounts := createTestApp(t, false, r, 3)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{ChainID: simapp.DefaultChainId, Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := simulation.SimulateMsgCreateValidator(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	_, _, err := op(r, app.BaseApp, ctx, accounts, simapp.DefaultChainId)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid signer")
}

// TestSimulateMsgCancelUnbondingDelegation tests the normal scenario of a valid message of type TypeMsgCancelUnbondingDelegation.
// Abonormal scenarios, where the message is
func TestSimulateMsgCancelUnbondingDelegation(t *testing.T) {
	s := rand.NewSource(5)
	r := rand.New(s)
	app, ctx, accounts := createTestApp(t, false, r, 3)

	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// remove genesis validator account
	accounts = accounts[1:]

	// setup accounts[0] as validator
	validator0 := getTestingValidator0(t, app, ctx, accounts)

	// setup delegation
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)
	delegator := accounts[1]
	delegation := types.NewDelegation(delegator.Address, validator0.GetOperator(), issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)
	app.DistrKeeper.SetDelegatorStartingInfo(ctx, validator0.GetOperator(), delegator.Address, distrtypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	setupValidatorRewards(app, ctx, validator0.GetOperator())

	// unbonding delegation
	udb := types.NewUnbondingDelegation(delegator.Address, validator0.GetOperator(), app.LastBlockHeight(), blockTime.Add(2*time.Minute), delTokens)
	app.StakingKeeper.SetUnbondingDelegation(ctx, udb)
	setupValidatorRewards(app, ctx, validator0.GetOperator())

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := simulation.SimulateMsgCancelUnbondingDelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	accounts = []simtypes.Account{accounts[1]}
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgCancelUnbondingDelegation
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, types.TypeMsgCancelUnbondingDelegation, msg.Type())
	require.Equal(t, delegator.Address.String(), msg.DelegatorAddress)
	require.Equal(t, validator0.GetOperator().String(), msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgEditValidator tests the normal scenario of a valid message of type TypeMsgEditValidator.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgEditValidator(t *testing.T) {
	s := rand.NewSource(3)
	r := rand.New(s)
	app, ctx, accounts := createTestApp(t, false, r, 3)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// remove genesis validator account
	accounts = accounts[1:]

	// setup accounts[0] as validator
	_ = getTestingValidator0(t, app, ctx, accounts)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := simulation.SimulateMsgEditValidator(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgEditValidator
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "0.697435912730803287", msg.CommissionRate.String())
	require.Equal(t, "UAkxytCUyv", msg.Description.Moniker)
	require.Equal(t, "HlkvYDcFSP", msg.Description.Identity)
	require.Equal(t, "PEHhWhpbMX", msg.Description.Website)
	require.Equal(t, "DYWWAOJwiv", msg.Description.SecurityContact)
	require.Equal(t, types.TypeMsgEditValidator, msg.Type())
	require.Equal(t, "0x87C4f0688CB0C0650d819a43F27d9934bd51a2b5", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgDelegate tests the normal scenario of a valid message of type TypeMsgDelegate.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgDelegate(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)
	app, ctx, accounts := createTestApp(t, false, r, 3)

	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// execute operation
	op := simulation.SimulateMsgDelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgDelegate
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "0xd4BFb1CB895840ca474b0D15abb11Cf0f26bc88a", msg.DelegatorAddress)
	require.Equal(t, "98100858108421259236", msg.Amount.Amount.String())
	require.Equal(t, "stake", msg.Amount.Denom)
	require.Equal(t, types.TypeMsgDelegate, msg.Type())
	require.Equal(t, "0x5cEEa0528c3b88442d6580c548753DD89b99a213", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgUndelegate tests the normal scenario of a valid message of type TypeMsgUndelegate.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgUndelegate(t *testing.T) {
	s := rand.NewSource(3)
	r := rand.New(s)
	app, ctx, accounts := createTestApp(t, false, r, 3)

	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// remove genesis validator account
	accounts = accounts[1:]

	// setup accounts[0] as validator
	validator0 := getTestingValidator0(t, app, ctx, accounts)

	// setup delegation
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)
	delegator := accounts[1]
	delegation := types.NewDelegation(delegator.Address, validator0.GetOperator(), issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)
	app.DistrKeeper.SetDelegatorStartingInfo(ctx, validator0.GetOperator(), delegator.Address, distrtypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	setupValidatorRewards(app, ctx, validator0.GetOperator())

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := simulation.SimulateMsgUndelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg types.MsgUndelegate
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "0xf05B84B69DB4490B1e395f698F20092e39472017", msg.DelegatorAddress)
	require.Equal(t, "1850357417337650264", msg.Amount.Amount.String())
	require.Equal(t, "stake", msg.Amount.Denom)
	require.Equal(t, types.TypeMsgUndelegate, msg.Type())
	require.Equal(t, "0x87C4f0688CB0C0650d819a43F27d9934bd51a2b5", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgBeginRedelegate tests the normal scenario of a valid message of type TypeMsgBeginRedelegate.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func TestSimulateMsgBeginRedelegate(t *testing.T) {
	s := rand.NewSource(5)
	r := rand.New(s)
	app, ctx, accounts := createTestApp(t, false, r, 4)

	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// remove genesis validator account
	accounts = accounts[1:]

	// setup accounts[0] as validator0 and accounts[1] as validator1
	validator0 := getTestingValidator0(t, app, ctx, accounts)
	validator1 := getTestingValidator1(t, app, ctx, accounts)

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)

	// setup accounts[2] as delegator
	delegator := accounts[2]
	delegation := types.NewDelegation(delegator.Address, validator1.GetOperator(), issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)
	app.DistrKeeper.SetDelegatorStartingInfo(ctx, validator1.GetOperator(), delegator.Address, distrtypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	setupValidatorRewards(app, ctx, validator0.GetOperator())
	setupValidatorRewards(app, ctx, validator1.GetOperator())

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := simulation.SimulateMsgBeginRedelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, simapp.DefaultChainId)
	require.NoError(t, err)

	var msg types.MsgBeginRedelegate
	types.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "0x1c52A0a7e7c01F242541F639f6f5D14a1200e567", msg.DelegatorAddress)
	require.Equal(t, "98319947001459696", msg.Amount.Amount.String())
	require.Equal(t, "stake", msg.Amount.Denom)
	require.Equal(t, types.TypeMsgBeginRedelegate, msg.Type())
	require.Equal(t, "0x9534A0660d545f5D6875d1f63e8bC18d0Ef605dd", msg.ValidatorDstAddress)
	require.Equal(t, "0xfC9D85d26FdA017d9BAe5cA16581d43bC9cB1658", msg.ValidatorSrcAddress)
	require.Len(t, futureOperations, 0)
}

// returns context and an app with updated mint keeper
func createTestApp(t *testing.T, isCheckTx bool, r *rand.Rand, n int) (*simapp.SimApp, sdk.Context, []simtypes.Account) {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

	accounts := simtypes.RandomAccounts(r, n)
	// create validator set with single validator
	account := accounts[0]
	tmPk, err := cryptocodec.ToTmPubKeyInterface(account.PubKey)
	require.NoError(t, err)
	validator := tmtypes.NewValidator(tmPk, 1)

	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey, _ := ethsecp256k1.GenerateKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	app := simapp.SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, true, balance)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{ChainID: simapp.DefaultChainId})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))

	// remove genesis validator account
	accs := accounts[1:]

	// add coins to the accounts
	for _, account := range accs {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		require.NoError(t, testutil.FundAccount(app.BankKeeper, ctx, account.Address, initCoins))
	}

	return app, ctx, accounts
}

func getTestingValidator0(t *testing.T, app *simapp.SimApp, ctx sdk.Context, accounts []simtypes.Account) types.Validator {
	commission0 := types.NewCommission(sdk.ZeroDec(), sdk.OneDec(), sdk.OneDec())
	return getTestingValidator(t, app, ctx, accounts, commission0, 0)
}

func getTestingValidator1(t *testing.T, app *simapp.SimApp, ctx sdk.Context, accounts []simtypes.Account) types.Validator {
	commission1 := types.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	return getTestingValidator(t, app, ctx, accounts, commission1, 1)
}

func getTestingValidator(t *testing.T, app *simapp.SimApp, ctx sdk.Context, accounts []simtypes.Account, commission types.Commission, n int) types.Validator {
	account := accounts[n]
	valPubKey := account.PubKey
	valAddr := sdk.AccAddress(account.PubKey.Address().Bytes())
	validator := teststaking.NewValidator(t, valAddr, valPubKey)
	validator, err := validator.SetInitialCommission(commission)
	require.NoError(t, err)

	validator.DelegatorShares = sdk.NewDec(100)
	validator.Tokens = app.StakingKeeper.TokensFromConsensusPower(ctx, 100)

	app.StakingKeeper.SetValidator(ctx, validator)

	return validator
}

func setupValidatorRewards(app *simapp.SimApp, ctx sdk.Context, valAddress sdk.AccAddress) {
	decCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.OneDec())}
	historicalRewards := distrtypes.NewValidatorHistoricalRewards(decCoins, 2)
	app.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddress, 2, historicalRewards)
	// setup current revards
	currentRewards := distrtypes.NewValidatorCurrentRewards(decCoins, 3)
	app.DistrKeeper.SetValidatorCurrentRewards(ctx, valAddress, currentRewards)
}
