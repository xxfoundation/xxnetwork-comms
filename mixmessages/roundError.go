////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundError message type conform to a generic
// signing interface

package mixmessages

import "github.com/pkg/errors"

// SetSignature sets RoundError's signature to the newSig argument
func (m *RoundError) SetSignature(newSig []byte) error {
	if newSig == nil {
		return errors.Errorf("Cannot set signature to nil value")
	}

	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Signature: newSig}
	}

	m.RsaSignature.Signature = newSig
	return nil
}

// ClearSignature clears out RoundError's signature by setting it to nil
func (m *RoundError) ClearSignature() {
	m.RsaSignature = nil
}

// GetNonce gets the value of the nonce
func (m *RoundError) GetNonce() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.RsaSignature.GetNonce()
}

// SetSignature sets RoundError's nonce to the newNonce argument
func (m *RoundError) SetNonce(newNonce []byte) error {
	if newNonce == nil {
		return errors.Errorf("Cannot set nonce to nil")
	}
	if m.RsaSignature == nil {
		return nil
	}

	m.RsaSignature.Nonce = newNonce

	return nil
}

//GetSignature
func (m *RoundError) GetSignature() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.GetRsaSignature().GetSignature()
}
