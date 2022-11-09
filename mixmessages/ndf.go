////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the ndf message type conform to a generic
// signing interface

package mixmessages

import (
	"gitlab.com/xx_network/comms/messages"
	"hash"
)

// GetSig returns the RSA signature.
// IF none exists, it creates it, adds it to the object, then returns it.
func (m *NDF) GetSig() *messages.RSASignature {
	if m.Signature != nil {
		return m.Signature
	}

	m.Signature = new(messages.RSASignature)

	return m.Signature
}

// Digest hashes the contents of the message in a repeatable manner
// using the provided cryptographic hash. It includes the nonce in the hash
func (m *NDF) Digest(nonce []byte, h hash.Hash) []byte {
	h.Reset()

	// Hash the ndf and the nonce
	h.Write(m.Ndf)
	h.Write(nonce)

	// Return the hash
	return h.Sum(nil)
}
