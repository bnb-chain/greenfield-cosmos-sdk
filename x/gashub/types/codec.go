package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMsgGasParams{}, "cosmos-sdk/MsgUpdateMsgGasParams")

	cdc.RegisterInterface((*isMsgGasParams_GasParams)(nil), nil)
	cdc.RegisterConcrete(&MsgGasParams_FixedType{}, "cosmos-sdk/MsgGasParams/FixedType", nil)
	cdc.RegisterConcrete(&MsgGasParams_GrantType{}, "cosmos-sdk/MsgGasParams/GrantType", nil)
	cdc.RegisterConcrete(&MsgGasParams_MultiSendType{}, "cosmos-sdk/MsgGasParams/MultiSendType", nil)
	cdc.RegisterConcrete(&MsgGasParams_GrantAllowanceType{}, "cosmos-sdk/MsgGasParams/GrantAllowanceType", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateMsgGasParams{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}
