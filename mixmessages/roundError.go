///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundError message type conform to a generic
// signing interface

package mixmessages

import (
	"github.com/pkg/errors"
	"gitlab.com/xx_network/comms/messages"
	"hash"
)

// GetSig returns the RSA signature.
// IF none exists, it creates it, adds it to the object, then returns it.
func (m *RoundError) GetSig() *messages.RSASignature {
	if m.Signature != nil {
		return m.Signature
	}

	m.Signature = new(messages.RSASignature)

	return m.Signature
}

// SetSignature sets RoundError's signature to the newSig argument
func (m *RoundError) SetSignature(newSig, nonce []byte) error {
	// Cannot set signature to nil
	if newSig == nil || nonce == nil {
		return errors.New("Cannot set signature to nil")
	}

	// Set the signature value
	m.Signature = &messages.RSASignature{
		Signature: newSig,
		Nonce:     nonce,
	}

	return nil
}

// Digest hashes the contents of the message in a repeatable manner
// using the provided cryptographic hash. It includes the nonce in the hash
func (m *RoundError) Digest(nonce []byte, h hash.Hash) []byte {
	h.Reset()

	// Hash the nodeId
	h.Write(m.NodeId)
	h.Write([]byte(m.Error))

	// Serialize and hash the round ID
	h.Write(serializeUin64(m.Id))

	// Hash the nonce
	h.Write(nonce)

	// Return the hash
	return h.Sum(nil)
}
