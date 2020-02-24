////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundInfo message type conform to a generic
// signing interface

package mixmessages

import (
	"github.com/pkg/errors"
)

// SetSignature sets RoundInfo's signature to the newSig argument
func (m *RoundInfo) SetSignature(newSig []byte) error {
	if newSig == nil {
		return errors.Errorf("Cannot set signature to nil value")
	}

	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Signature: newSig}
	}

	m.RsaSignature.Signature = newSig
	return nil
}

// ClearSignature clears out roundInfo's signature by
// setting it to nil
func (m *RoundInfo) ClearSignature() {

	m.RsaSignature = nil
}

// GetNonce gets the value of the nonce
func (m *RoundInfo) GetNonce() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.RsaSignature.GetNonce()
}

// SetSignature sets RoundError's nonce to the newNonce argument
func (m *RoundInfo) SetNonce(newNonce []byte) error {
	if newNonce == nil {
		return errors.Errorf("Cannot set nonce to nil")
	}
	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Nonce: newNonce}
	}

	m.RsaSignature.Nonce = newNonce

	return nil
}

//GetSignature
func (m *RoundInfo) GetSignature() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.GetRsaSignature().GetSignature()
}
