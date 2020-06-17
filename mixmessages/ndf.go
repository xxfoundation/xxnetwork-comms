///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functions to make the ndf message type conform to a generic
// signing interface

package mixmessages

import "github.com/pkg/errors"

// SetSignature sets NDF's signature to the newSig argument
func (m *NDF) SetSig(newSig []byte) error {
	// Cannot set signature to nil
	if newSig == nil {
		return errors.New("Cannot set signature to nil")
	}

	// If the signature object is nil, create it and set value
	if m.Signature == nil {
		m.Signature = &RSASignature{Signature: newSig}
		return nil
	}

	// Set value otherwise
	m.Signature.Signature = newSig
	return nil
}

// ClearSignature clears out NDF's signature
func (m *NDF) ClearSig() {
	if m.Signature != nil {
		m.Signature.Signature = nil
	}
}

// GetNonce gets the value of the nonce
func (m *NDF) GetNonce() []byte {
	// If the signature object is nil, then value is nil
	if m.Signature == nil {
		return nil
	}

	return m.Signature.GetNonce()
}

// SetSignature sets NDF's nonce to the newNonce argument
func (m *NDF) SetNonce(newNonce []byte) error {
	// Cannot set nonce to nil
	if newNonce == nil {
		return errors.New("Cannot set nonce to nil")
	}

	// If the signature object is nil, create it and set value
	if m.Signature == nil {
		m.Signature = &RSASignature{Nonce: newNonce}
		return nil
	}

	// Set value otherwise
	m.Signature.Nonce = newNonce

	return nil
}

// GetSignature gets the value of the signature in RSASignature
func (m *NDF) GetSig() []byte {
	if m.Signature == nil {
		return nil
	}

	return m.GetSignature().Signature
}
