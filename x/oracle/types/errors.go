package types

import "cosmossdk.io/errors"

var (
	ErrChannelNotRegistered   = errors.Register(ModuleName, 1, "channel is not registered")
	ErrInvalidReceiveSequence = errors.Register(ModuleName, 2, "receive sequence is invalid")
	ErrInvalidPayloadHeader   = errors.Register(ModuleName, 3, "payload header is invalid")
	ErrInvalidPackageType     = errors.Register(ModuleName, 4, "package type is invalid")
	ErrInvalidPackage         = errors.Register(ModuleName, 6, "package is invalid")
	ErrInvalidPayload         = errors.Register(ModuleName, 7, "payload is invalid")
	ErrValidatorSet           = errors.Register(ModuleName, 8, "validator set is invalid")
	ErrBlsPubKey              = errors.Register(ModuleName, 9, "public key is invalid")
	ErrBlsVotesNotEnough      = errors.Register(ModuleName, 10, "bls votes is not enough")
	ErrInvalidBlsSignature    = errors.Register(ModuleName, 11, "bls signature is invalid")
	ErrNotRelayer             = errors.Register(ModuleName, 12, "sender is not a relayer")
	ErrRelayerNotInTurn       = errors.Register(ModuleName, 13, "relayer is not in turn")
	ErrInvalidDestChainId     = errors.Register(ModuleName, 14, "dest chain id is invalid")
	ErrInvalidSrcChainId      = errors.Register(ModuleName, 15, "src chain id is invalid")
	ErrInvalidAddress         = errors.Register(ModuleName, 16, "address is invalid")
	ErrInvalidMultiMessage    = errors.Register(ModuleName, 17, "multi message is invalid")
	ErrInvalidMessagesResult  = errors.Register(ModuleName, 18, "multi message result is invalid")
)
