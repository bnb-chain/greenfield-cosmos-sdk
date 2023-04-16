package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "cosmos-sdk/x/gashub/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgSetMsgGasParams{}, "cosmos-sdk/MsgSetMsgGasParams")

	cdc.RegisterInterface((*isMsgGasParams_GasParams)(nil), nil)
	cdc.RegisterConcrete(&MsgGasParams_FixedType{}, "cosmos-sdk/MsgGasParams/FixedType", nil)
	cdc.RegisterConcrete(&MsgGasParams_GrantType{}, "cosmos-sdk/MsgGasParams/GrantType", nil)
	cdc.RegisterConcrete(&MsgGasParams_MultiSendType{}, "cosmos-sdk/MsgGasParams/MultiSendType", nil)
	cdc.RegisterConcrete(&MsgGasParams_GrantAllowanceType{}, "cosmos-sdk/MsgGasParams/GrantAllowanceType", nil)

	cdc.RegisterConcrete(&Params{}, "cosmos-sdk/x/gashub/Params", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgSetMsgGasParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}
