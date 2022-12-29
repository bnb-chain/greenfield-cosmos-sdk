package types_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/x/oracle/types"
)

func TestBlsClaim(t *testing.T) {
	claim := &types.BlsClaim{
		SrcChainId:  1,
		DestChainId: 2,
		Sequence:    1,
		Timestamp:   1000,
		Payload:     []byte("test payload"),
	}

	signBytes := claim.GetSignBytes()

	require.Equal(t, "0a0b49ef40324d4c511d7a81e1edeeccaa10b768e55cece473b5cd99137f05f6",
		hex.EncodeToString(signBytes[:]))
}

func TestValidateBasic(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig().Codec
	addr, _, err := testutil.GenerateCoinKey(hd.Secp256k1, cdc)
	require.NoError(t, err)

	tests := []struct {
		claimMsg     types.MsgClaim
		expectedPass bool
		errorMsg     string
	}{
		{
			types.MsgClaim{
				FromAddress:    "random string",
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"invalid from address",
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     math.MaxUint16 + 1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"chain id should not be larger than",
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     1,
				DestChainId:    math.MaxUint16 + 1,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"chain id should not be larger than",
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte{},
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			"payload should not be empty",
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1},
				AggSignature:   []byte("test sig"),
			},
			false,
			fmt.Sprintf("length of vote addresse set should be %d", types.ValidatorBitSetLength),
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1, 2, 3},
				AggSignature:   []byte("test sig"),
			},
			false,
			fmt.Sprintf("length of signature should be %d", types.BLSSignatureLength),
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1, 2, 3},
				AggSignature:   bytes.Repeat([]byte{0}, types.BLSSignatureLength),
			},
			false,
			"timestamp should not be 0",
		},
		{
			types.MsgClaim{
				FromAddress:    addr.String(),
				SrcChainId:     1,
				DestChainId:    2,
				Sequence:       1,
				Payload:        []byte("test payload"),
				VoteAddressSet: []uint64{0, 1, 2, 3},
				AggSignature:   bytes.Repeat([]byte{0}, types.BLSSignatureLength),
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
