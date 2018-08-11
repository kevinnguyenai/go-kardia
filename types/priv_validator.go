package types

import (
	"crypto/ecdsa"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/crypto"
)

// PrivValidator defines the functionality of a local Kardia validator
// that signs votes, proposals, and heartbeats, and never double signs.
type PrivValidator struct {
	privKey *ecdsa.PrivateKey
}

func NewPrivValidator(privKey *ecdsa.PrivateKey) *PrivValidator {
	return &PrivValidator{
		privKey: privKey,
	}
}

func (privVal *PrivValidator) GetAddress() common.Address {
	return crypto.PubkeyToAddress(privVal.GetPubKey())
}

func (privVal *PrivValidator) GetPubKey() ecdsa.PublicKey {
	return privVal.privKey.PublicKey
}

func (privVal *PrivValidator) GetPrivKey() *ecdsa.PrivateKey {
	return privVal.privKey
}

func (privVal *PrivValidator) SignVote(chainID string, vote *Vote) error {
	panic("SignVote - not yet implemented")
}

func (privVal *PrivValidator) SignProposal(chainID string, proposal *Proposal) error {
	panic("SignProposal - not yet implemented")
}

//func (privVal *PrivValidator) SignHeartbeat(chainID string, heartbeat *Heartbeat) error {
//	panic("SignHeartbeat - not yet implemented")
//}