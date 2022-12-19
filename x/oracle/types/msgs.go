package types

import (
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ValidatorBitSetLength = 4 // 256 bits
	BLSPublicKeyLength    = 48
	BLSSignatureLength    = 96
)

type BLSPublicKey [BLSPublicKeyLength]byte
type BLSSignature [BLSSignatureLength]byte

func NewMsgClaim(fromAddr string, chainId uint32, sequence uint64, timestamp uint64, payload []byte, voteAddrSet []uint64, aggSignature []byte) *MsgClaim {
	return &MsgClaim{
		FromAddress:    fromAddr,
		ChainId:        chainId,
		Sequence:       sequence,
		Timestamp:      timestamp,
		Payload:        payload,
		VoteAddressSet: voteAddrSet,
		AggSignature:   aggSignature,
	}
}

// Route implements the LegacyMsg interface.
func (m MsgClaim) Route() string { return sdk.MsgTypeURL(&m) }

// Type implements the LegacyMsg interface.
func (m MsgClaim) Type() string { return sdk.MsgTypeURL(&m) }

// GetSignBytes implements the LegacyMsg interface.
func (m MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromHexUnsafe(m.FromAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if m.ChainId > math.MaxUint16 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("chain id should not be larger than %d", math.MaxUint16))
	}

	if len(m.Payload) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("payload should not be empty"))
	}

	if len(m.VoteAddressSet) != ValidatorBitSetLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("length of vote addresse set should be %d", ValidatorBitSetLength))
	}

	if len(m.AggSignature) != BLSSignatureLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("length of signature should be %d", BLSSignatureLength),
		)
	}

	if m.Timestamp == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("timestamp should not be 0"),
		)
	}

	return nil
}

// GetSigners returns the expected signers for MsgCancelUpgrade.
func (m *MsgClaim) GetSigners() []sdk.AccAddress {
	fromAddress := sdk.MustAccAddressFromHex(m.FromAddress)
	return []sdk.AccAddress{fromAddress}
}

// GetBlsSignBytes returns the sign bytes of bls signature
func (m *MsgClaim) GetBlsSignBytes() [32]byte {
	blsClaim := &BlsClaim{
		ChainId:   m.ChainId,
		Timestamp: m.Timestamp,
		Sequence:  m.Sequence,
		Payload:   m.Payload,
	}
	return blsClaim.GetSignBytes()
}

type BlsClaim struct {
	ChainId   uint32
	Timestamp uint64
	Sequence  uint64
	Payload   []byte
}

func (c *BlsClaim) GetSignBytes() [32]byte {
	bts, err := rlp.EncodeToBytes(c)
	if err != nil {
		panic("encode bls claim error")
	}

	btsHash := sdk.Keccak256Hash(bts)
	return btsHash
}

type Packages []Package

type Package struct {
	ChannelId sdk.ChannelID
	Sequence  uint64
	Payload   []byte
}
