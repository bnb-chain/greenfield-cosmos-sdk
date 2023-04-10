package types

import (
	"cosmossdk.io/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal       = errors.Register(ModuleName, 2, "unknown proposal")
	ErrInactiveProposal      = errors.Register(ModuleName, 3, "inactive proposal")
	ErrAlreadyActiveProposal = errors.Register(ModuleName, 4, "proposal already active")
	// Errors 5 & 6 are legacy errors related to v1beta1.Proposal.
	ErrInvalidProposalContent  = errors.Register(ModuleName, 5, "invalid proposal content")
	ErrInvalidProposalType     = errors.Register(ModuleName, 6, "invalid proposal type")
	ErrInvalidVote             = errors.Register(ModuleName, 7, "invalid vote option")
	ErrInvalidGenesis          = errors.Register(ModuleName, 8, "invalid genesis state")
	ErrNoProposalHandlerExists = errors.Register(ModuleName, 9, "no handler exists for proposal type")
	ErrUnroutableProposalMsg   = errors.Register(ModuleName, 10, "proposal message not recognized by router")
	ErrNoProposalMsgs          = errors.Register(ModuleName, 11, "no messages proposed")
	ErrInvalidProposalMsg      = errors.Register(ModuleName, 12, "invalid proposal message")
	ErrInvalidSigner           = errors.Register(ModuleName, 13, "expected gov account as only signer for proposal message")
	ErrInvalidSignalMsg        = errors.Register(ModuleName, 14, "signal message is invalid")
	ErrMetadataTooLong         = errors.Register(ModuleName, 15, "metadata too long")
	ErrMinDepositTooSmall      = errors.Register(ModuleName, 16, "minimum deposit is too small")
	ErrProposalNotFound        = errors.Register(ModuleName, 17, "proposal is not found")
	ErrInvalidProposer         = errors.Register(ModuleName, 18, "invalid proposer")
	ErrNoDeposits              = errors.Register(ModuleName, 19, "no deposits found")
	ErrVotingPeriodEnded       = errors.Register(ModuleName, 20, "voting period already ended")
	ErrInvalidProposal         = errors.Register(ModuleName, 21, "invalid proposal")
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
)
