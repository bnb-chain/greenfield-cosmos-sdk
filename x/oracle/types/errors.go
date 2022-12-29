package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrChannelNotRegistered   = sdkerrors.Register(ModuleName, 1, "channel is not registered")
	ErrInvalidReceiveSequence = sdkerrors.Register(ModuleName, 2, "receive sequence is invalid")
	ErrInvalidPayloadHeader   = sdkerrors.Register(ModuleName, 3, "payload header is invalid")
	ErrInvalidPackageType     = sdkerrors.Register(ModuleName, 4, "package type is invalid")
	ErrInvalidPackage         = sdkerrors.Register(ModuleName, 6, "package is invalid")
	ErrInvalidPayload         = sdkerrors.Register(ModuleName, 7, "payload is invalid")
	ErrValidatorSet           = sdkerrors.Register(ModuleName, 8, "validator set is invalid")
	ErrBlsPubKey              = sdkerrors.Register(ModuleName, 9, "public key is invalid")
	ErrBlsVotesNotEnough      = sdkerrors.Register(ModuleName, 10, "bls votes is not enough")
	ErrInvalidBlsSignature    = sdkerrors.Register(ModuleName, 11, "bls signature is invalid")
	ErrNotValidator           = sdkerrors.Register(ModuleName, 12, "sender is not validator")
	ErrValidatorNotInTurn     = sdkerrors.Register(ModuleName, 13, "validator is not in turn")
	ErrInvalidDestChainId     = sdkerrors.Register(ModuleName, 14, "dest chain id is invalid")
	ErrInvalidSrcChainId      = sdkerrors.Register(ModuleName, 15, "src chain id is invalid")
)
