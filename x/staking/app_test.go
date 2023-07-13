package staking_test

import (
	"encoding/hex"
	"testing"

	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/stretchr/testify/require"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/ethsecp256k1"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	priv1, _ = ethsecp256k1.GenPrivKey()
	addr1    = sdk.AccAddress(priv1.PubKey().Address())
	priv2, _ = ethsecp256k1.GenPrivKey()
	addr2    = sdk.AccAddress(priv2.PubKey().Address())

	valKey          = ed25519.GenPrivKey()
	commissionRates = types.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec())
)

func TestStakingMsgs(t *testing.T) {
	genTokens := sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	bondTokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	genCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

	acc1 := &authtypes.BaseAccount{Address: addr1.String()}
	acc2 := &authtypes.BaseAccount{Address: addr2.String()}
	accs := []simtestutil.GenesisAccount{
		{GenesisAccount: acc1, Coins: sdk.Coins{genCoin}},
		{GenesisAccount: acc2, Coins: sdk.Coins{genCoin}},
	}

	var (
		bankKeeper    bankKeeper.Keeper
		stakingKeeper *stakingKeeper.Keeper
	)

	startupCfg := simtestutil.DefaultStartUpConfig()
	startupCfg.GenesisAccounts = accs

	app, err := simtestutil.SetupWithConfiguration(testutil.AppConfig, startupCfg, &bankKeeper, &stakingKeeper)
	require.NoError(t, err)
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{ChainID: sdktestutil.DefaultChainId})

	require.True(t, sdk.Coins{genCoin}.IsEqual(bankKeeper.GetAllBalances(ctxCheck, addr1)))
	require.True(t, sdk.Coins{genCoin}.IsEqual(bankKeeper.GetAllBalances(ctxCheck, addr2)))

	// create validator
	description := types.NewDescription("foo_moniker", "", "", "", "")
	blsSecretKey, _ := bls.RandKey()
	blsPubKey := hex.EncodeToString(blsSecretKey.PublicKey().Marshal())
	blsProofBuf := blsSecretKey.Sign(tmhash.Sum(blsSecretKey.PublicKey().Marshal()))
	blsProof1 := hex.EncodeToString(blsProofBuf.Marshal())
	createValidatorMsg, err := types.NewMsgCreateValidator(
		addr1, valKey.PubKey(),
		bondCoin, description, commissionRates, sdk.OneInt(),
		addr1, addr1, addr1, addr1, blsPubKey, blsProof1,
	)
	require.NoError(t, err)

	header := tmproto.Header{ChainID: sdktestutil.DefaultChainId, Height: app.LastBlockHeight() + 1}
	txConfig := moduletestutil.MakeTestEncodingConfig().TxConfig
	_, _, err = simtestutil.SignCheckDeliver(t, txConfig, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, sdktestutil.DefaultChainId, []uint64{0}, []uint64{0}, true, true, []cryptotypes.PrivKey{priv1}, simtestutil.SetMockHeight(app.BaseApp, 0))
	require.NoError(t, err)
	ctxCheck = app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, sdk.Coins{genCoin.Sub(bondCoin)}.IsEqual(bankKeeper.GetAllBalances(ctxCheck, addr1)))

	header = tmproto.Header{ChainID: sdktestutil.DefaultChainId, Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctxCheck = app.BaseApp.NewContext(true, tmproto.Header{})
	validator, found := stakingKeeper.GetValidator(ctxCheck, addr1)
	require.True(t, found)
	require.Equal(t, addr1.String(), validator.OperatorAddress)
	require.Equal(t, types.Bonded, validator.Status)
	require.True(math.IntEq(t, bondTokens, validator.BondedTokens()))

	header = tmproto.Header{ChainID: sdktestutil.DefaultChainId, Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// edit the validator
	description = types.NewDescription("bar_moniker", "", "", "", "")
	editValidatorMsg := types.NewMsgEditValidator(
		addr1, description, nil, nil,
		sdk.AccAddress(""), sdk.AccAddress(""), "", "",
	)
	header = tmproto.Header{ChainID: sdktestutil.DefaultChainId, Height: app.LastBlockHeight() + 1}
	_, _, err = simtestutil.SignCheckDeliver(t, txConfig, app.BaseApp, header, []sdk.Msg{editValidatorMsg}, sdktestutil.DefaultChainId, []uint64{0}, []uint64{1}, true, true, []cryptotypes.PrivKey{priv1})
	require.NoError(t, err)

	ctxCheck = app.BaseApp.NewContext(true, tmproto.Header{})
	validator, found = stakingKeeper.GetValidator(ctxCheck, addr1)
	require.True(t, found)
	require.Equal(t, description, validator.Description)

	// delegate
	require.True(t, sdk.Coins{genCoin}.IsEqual(bankKeeper.GetAllBalances(ctxCheck, addr2)))
	delegateMsg := types.NewMsgDelegate(addr2, addr1, bondCoin)

	header = tmproto.Header{ChainID: sdktestutil.DefaultChainId, Height: app.LastBlockHeight() + 1}
	_, _, err = simtestutil.SignCheckDeliver(t, txConfig, app.BaseApp, header, []sdk.Msg{delegateMsg}, sdktestutil.DefaultChainId, []uint64{1}, []uint64{0}, true, true, []cryptotypes.PrivKey{priv2})
	require.NoError(t, err)

	ctxCheck = app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, sdk.Coins{genCoin.Sub(bondCoin)}.IsEqual(bankKeeper.GetAllBalances(ctxCheck, addr2)))
	_, found = stakingKeeper.GetDelegation(ctxCheck, addr2, addr1)
	require.True(t, found)

	// begin unbonding
	beginUnbondingMsg := types.NewMsgUndelegate(addr2, addr1, bondCoin)
	header = tmproto.Header{ChainID: sdktestutil.DefaultChainId, Height: app.LastBlockHeight() + 1}
	_, _, err = simtestutil.SignCheckDeliver(t, txConfig, app.BaseApp, header, []sdk.Msg{beginUnbondingMsg}, sdktestutil.DefaultChainId, []uint64{1}, []uint64{1}, true, true, []cryptotypes.PrivKey{priv2})
	require.NoError(t, err)

	// delegation should exist anymore
	ctxCheck = app.BaseApp.NewContext(true, tmproto.Header{})
	_, found = stakingKeeper.GetDelegation(ctxCheck, addr2, addr1)
	require.False(t, found)

	// balance should be the same because bonding not yet complete
	require.True(t, sdk.Coins{genCoin.Sub(bondCoin)}.IsEqual(bankKeeper.GetAllBalances(ctxCheck, addr2)))
}
