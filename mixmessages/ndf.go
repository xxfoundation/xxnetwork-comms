////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the ndf message type conform to a generic
// signing interface

package mixmessages

import "github.com/pkg/errors"

// Marshal serializes all the data needed for a signature
func (m *NDF) Marshal() []byte {
	// Create the byte array
	b := make([]byte, 0)

	// Marshall the roundInfo data and append to byte array
	b = append(b, m.Ndf...)

	return b
}

// SetSignature sets NDF's signature to the newSig argument
func (m *NDF) SetSignature(newSig []byte) error {
	if newSig == nil {
		return errors.Errorf("Cannot set signature to nil")
	}
	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Signature: newSig}
	}
	m.RsaSignature.Signature = newSig
	return nil
}

// ClearSignature clears out NDF's signature by setting it to nil
func (m *NDF) ClearSignature() {
	m.RsaSignature = nil
}

// GetNonce gets the value of the nonce
func (m *NDF) GetNonce() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.RsaSignature.GetNonce()
}

// SetSignature sets NDF's nonce to the newNonce argument
func (m *NDF) SetNonce(newNonce []byte) error {
	if newNonce == nil {
		return errors.Errorf("Cannot set nonce to nil")
	}
	if m.RsaSignature == nil {
		return nil
	}

	m.GetRsaSignature().Nonce = newNonce

	return nil
}

//GetSignature
func (m *NDF) GetSignature() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.GetRsaSignature().GetSignature()
}
