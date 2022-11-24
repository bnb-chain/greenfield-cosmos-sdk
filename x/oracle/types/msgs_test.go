package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestBlsClaim(t *testing.T) {
	claim := &BlsClaim{
		ChainId:   1,
		Sequence:  1,
		Timestamp: 1000,
		Payload:   []byte("test payload"),
	}

	signBytes := claim.GetSignBytes()

	require.Equal(t, "954d4fe4c768c275f14ef32929ab83e182a4de3c0aef38964efdf0bc8f76eaff",
		hex.EncodeToString(signBytes[:]))
}

func TestValidateBasic(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Codec
	addr, _, err := testutil.GenerateCoinKey(hd.Secp256k1, cdc)
	require.NoError(t, err)

	tests := []struct {
		claimMsg     MsgClaim
		expectedPass bool
		errorMsg     string
	}{
		{
			MsgClaim{
				FromAddress:    "random string",
				ChainId:        1,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"invalid from address",
		},
		{
			MsgClaim{
				FromAddress:    addr.String(),
				ChainId:        math.MaxUint16 + 1,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"chain id should not be larger than",
		},
		{
			MsgClaim{
				FromAddress:    addr.String(),
				ChainId:        100,
				Sequence:       1,
				Payload:        []byte{},
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"payload should not be empty",
		},
		{
			MsgClaim{
				FromAddress:    addr.String(),
				ChainId:        100,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			fmt.Sprintf("length of vote addresse set should be %d", ValidatorBitSetLength),
		},
		{
			MsgClaim{
				FromAddress:    addr.String(),
				ChainId:        100,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1, 2, 3},
				AggSignature:   []byte("test sig"),
			},
			false,
			fmt.Sprintf("length of signature should be %d", BLSSignatureLength),
		},
		{
			MsgClaim{
				FromAddress:    addr.String(),
				ChainId:        100,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1, 2, 3},
				AggSignature:   bytes.Repeat([]byte{0}, BLSSignatureLength),
			},
			false,
			"timestamp should not be 0",
		},
		{
			MsgClaim{
				FromAddress:    addr.String(),
				ChainId:        100,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1, 2, 3},
				AggSignature:   bytes.Repeat([]byte{0}, BLSSignatureLength),
				Timestamp:      uint64(time.Now().Unix()),
			},
			true,
			"",
		},
	}

	for i, test := range tests {
		if test.expectedPass {
			require.Nil(t, test.claimMsg.ValidateBasic(), "test: %v", i)
		} else {
			err := test.claimMsg.ValidateBasic()
			require.ErrorContains(t, err, test.errorMsg)
		}
	}
}
