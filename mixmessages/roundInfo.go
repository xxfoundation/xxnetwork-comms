////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundInfo message type conform to a generic
// signing interface

package mixmessages

import "github.com/pkg/errors"

// SetSignature sets RoundInfo's signature to the newSig argument
func (m *RoundInfo) SetSig(newSig []byte) error {
	// Cannot set signature to nil
	if newSig == nil {
		return errors.New("Cannot set signature to nil value")
	}

	// If the signature object is nil, create it and set value
	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Signature: newSig}
		return nil
	}

	// Set value as normal otherwise
	m.RsaSignature.Signature = newSig
	return nil
}

// ClearSignature clears out roundInfo's signature
func (m *RoundInfo) ClearSig() {
	m.RsaSignature = &RSASignature{}
}

// GetNonce gets the value of the nonce
func (m *RoundInfo) GetNonce() []byte {
	// If the signature object is nil, then value is nil
	if m.RsaSignature == nil {
		return nil
	}

	return m.RsaSignature.GetNonce()
}

// SetSignature sets RoundInfo's nonce to the newNonce argument
func (m *RoundInfo) SetNonce(newNonce []byte) error {
	// Cannot set nonce to nil
	if newNonce == nil {
		return errors.New("Cannot set nonce to nil")
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
func (m *RoundInfo) GetSig() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.GetRsaSignature().GetSignature()
}
