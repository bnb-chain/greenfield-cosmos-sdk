package signing_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestVerifySignature(t *testing.T) {
	priv, pubKey, addr := testdata.KeyTestPubAddrEthSecp256k1(require.New(t))

	const (
		memo    = "testmemo"
		chainId = "test-chain"
	)

	encCfg := moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{})
	key := sdk.NewKVStoreKey(types.StoreKey)

	maccPerms := map[string][]string{
		"fee_collector":          nil,
		"mint":                   {"minter"},
		"bonded_tokens_pool":     {"burner", "staking"},
		"not_bonded_tokens_pool": {"burner", "staking"},
		"multiPerm":              {"burner", "minter", "staking"},
		"random":                 {"random"},
	}

	accountKeeper := keeper.NewAccountKeeper(
		encCfg.Codec,
		key,
		types.ProtoBaseAccount,
		maccPerms,
		types.NewModuleAddress("gov").String(),
	)

	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{})
	encCfg.Amino.RegisterConcrete(testdata.TestMsg{}, "cosmos-sdk/Test", nil)

	acc1 := accountKeeper.NewAccountWithAddress(ctx, addr)
	accountKeeper.SetAccount(ctx, acc1)
	acc, err := ante.GetSignerAcc(ctx, accountKeeper, addr)
	require.NoError(t, err)

	msgs := []sdk.Msg{testdata.NewTestMsg(addr)}
	fee := legacytx.NewStdFee(50000, sdk.Coins{sdk.NewInt64Coin("atom", 150)})
	signerData := signing.SignerData{
		Address:       addr.String(),
		ChainID:       chainId,
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
		PubKey:        pubKey,
	}
	signBytes := legacytx.StdSignBytes(signerData.ChainID, signerData.AccountNumber, signerData.Sequence, 10, fee, msgs, memo, nil)
	signature, err := priv.Sign(signBytes)
	require.NoError(t, err)

	stdSig := legacytx.StdSignature{PubKey: pubKey, Signature: signature}
	sigV2, err := legacytx.StdSignatureToSignatureV2(encCfg.Amino, stdSig)
	require.NoError(t, err)

	handler := MakeTestHandlerMap()
	stdTx := legacytx.NewStdTx(msgs, fee, []legacytx.StdSignature{stdSig}, memo)
	stdTx.TimeoutHeight = 10
	err = signing.VerifySignature(sdk.Context{}, pubKey, signerData, sigV2.Data, handler, stdTx)
	require.NoError(t, err)
}
