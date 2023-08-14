package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal       = sdkerrors.Register(ModuleName, 2, "unknown proposal")
	ErrInactiveProposal      = sdkerrors.Register(ModuleName, 3, "inactive proposal")
	ErrAlreadyActiveProposal = sdkerrors.Register(ModuleName, 4, "proposal already active")
	// Errors 5 & 6 are legacy errors related to v1beta1.Proposal.
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 5, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 6, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 7, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 8, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 9, "no handler exists for proposal type")
	ErrUnroutableProposalMsg   = sdkerrors.Register(ModuleName, 10, "proposal message not recognized by router")
	ErrNoProposalMsgs          = sdkerrors.Register(ModuleName, 11, "no messages proposed")
	ErrInvalidProposalMsg      = sdkerrors.Register(ModuleName, 12, "invalid proposal message")
	ErrInvalidSigner           = sdkerrors.Register(ModuleName, 13, "expected gov account as only signer for proposal message")
	ErrInvalidSignalMsg        = sdkerrors.Register(ModuleName, 14, "signal message is invalid")
	ErrMetadataTooLong         = sdkerrors.Register(ModuleName, 15, "metadata too long")
	ErrMinDepositTooSmall      = sdkerrors.Register(ModuleName, 16, "minimum deposit is too small")
)

var (
	ErrEmptyChange = errors.Register(ModuleName, 22, "crosschain: change is empty")
	ErrEmptyValue  = errors.Register(ModuleName, 23, "crosschain: value  is empty")
	ErrEmptyTarget = errors.Register(ModuleName, 24, "crosschain: target is empty")

	ErrAddressSizeNotMatch     = errors.Register(ModuleName, 25, "number of old address not equal to new addresses")
	ErrAddressNotValid         = errors.Register(ModuleName, 26, "address format is not valid")
	ErrExceedParamsChangeLimit = errors.Register(ModuleName, 27, "exceed params change limit")
	ErrInvalidUpgradeProposal  = errors.Register(ModuleName, 28, "invalid sync params package")
	ErrInvalidValue            = errors.Register(ModuleName, 29, "decode hex value failed")

	ErrChainNotSupported = errors.Register(ModuleName, 30, "crosschain: chain is not supported")
)
