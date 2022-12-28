package testutil

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/prysmaticlabs/prysm/crypto/bls/blst"
	blscmn "github.com/prysmaticlabs/prysm/crypto/bls/common"
)

type Vote struct {
	PubKey    [48]byte
	Signature [96]byte
}

func (vote *Vote) Verify(eventHash []byte) error {
	blsPubKey, err := bls.PublicKeyFromBytes(vote.PubKey[:])

	if err != nil {
		return err
	}
	sig, err := bls.SignatureFromBytes(vote.Signature[:])

	if err != nil {
		return err
	}
	if !sig.Verify(blsPubKey, eventHash[:]) {
		return fmt.Errorf("verify sig error")
	}
	return nil
}

func AggregatedSignature(votes []*Vote) ([]byte, error) {
	// Prepare aggregated vote signature
	signatures := make([][]byte, 0, len(votes))
	for _, v := range votes {
		signatures = append(signatures, v.Signature[:])
	}
	sigs, err := bls.MultipleSignaturesFromBytes(signatures)
	if err != nil {
		return nil, err
	}
	return bls.AggregateSignatures(sigs).Marshal(), nil
}

type VoteSigner struct {
	privkey blscmn.SecretKey
	pubKey  blscmn.PublicKey
}

func NewVoteSignerV2(privkey []byte) (*VoteSigner, error) {
	privKey, err := blst.SecretKeyFromBytes(privkey)
	if err != nil {
		return nil, err
	}
	pubKey := privKey.PublicKey()
	return &VoteSigner{
		privkey: privKey,
		pubKey:  pubKey,
	}, nil
}

// SignVote sign a vote, data is used to signed to generate the signature
func (signer *VoteSigner) SignVote(vote *Vote, data []byte) error {
	signature := signer.privkey.Sign(data[:])
	copy(vote.PubKey[:], signer.pubKey.Marshal()[:])
	copy(vote.Signature[:], signature.Marshal()[:])
	return nil
}

func GenerateBlsSig(privKeys []bls.SecretKey, data []byte) []byte {
	privateKey1 := hex.EncodeToString(privKeys[0].Marshal())
	privateKey2 := hex.EncodeToString(privKeys[1].Marshal())
	privateKey3 := hex.EncodeToString(privKeys[2].Marshal())

	validatorSigner1, _ := NewVoteSignerV2(common.Hex2Bytes(privateKey1))
	validatorSigner2, _ := NewVoteSignerV2(common.Hex2Bytes(privateKey2))
	validatorSigner3, _ := NewVoteSignerV2(common.Hex2Bytes(privateKey3))

	var vote1 Vote
	validatorSigner1.SignVote(&vote1, data)
	err := vote1.Verify(data)
	if err != nil {
		panic("verify sig error")
	}

	var vote2 Vote
	validatorSigner2.SignVote(&vote2, data)
	err = vote2.Verify(data)
	if err != nil {
		panic("verify sig error")
	}

	var vote3 Vote
	validatorSigner3.SignVote(&vote3, data)
	err = vote2.Verify(data)
	if err != nil {
		panic("verify sig error")
	}

	var votes []*Vote
	votes = append(votes, &vote1)
	votes = append(votes, &vote2)
	votes = append(votes, &vote3)

	aggreatedSigature, _ := AggregatedSignature(votes)
	return aggreatedSigature
}
