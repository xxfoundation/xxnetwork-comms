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
	// Cannot set signature to nil
	if newSig == nil {
		return errors.Errorf("Cannot set signature to nil value")
	}

	// If the signature object is nil, create it and set value
	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Signature: newSig}
		return nil
	}

	// Set value otherwise
	m.RsaSignature.Signature = newSig
	return nil
}

// ClearSignature clears out RoundError's signature
func (m *RoundError) ClearSignature() {
	m.RsaSignature = &RSASignature{}
}

// GetNonce gets the value of the nonce
func (m *RoundError) GetNonce() []byte {
	// If the signature object is nil, then value is nil
	if m.RsaSignature == nil {
		return nil
	}

	return m.RsaSignature.GetNonce()
}

// SetSignature sets RoundError's nonce to the newNonce argument
func (m *RoundError) SetNonce(newNonce []byte) error {
	// Cannot set nonce to nil
	if newNonce == nil {
		return errors.Errorf("Cannot set nonce to nil")
	}

	// If the signature object is nil, create it and set value
	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Nonce: newNonce}
		return nil
	}

	// Set value otherwise
	m.RsaSignature.Nonce = newNonce

	return nil
}

// GetSignature gets the value of the signature in RSASignature
func (m *RoundError) GetSignature() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.GetRsaSignature().GetSignature()
}
