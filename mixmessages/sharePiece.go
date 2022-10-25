////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundInfo message type conform to a generic
// signing interface
package mixmessages

import (
	"gitlab.com/xx_network/comms/messages"
	"hash"
)

// GetSig returns the RSA signature.
// IF none exists, it creates it, adds it to the object, then returns it.
func (m *SharePiece) GetSig() *messages.RSASignature {
	if m.Signature != nil {
		return m.Signature
	}

	m.Signature = new(messages.RSASignature)

	return m.Signature
}

// Digest hashes the contents of the message in a repeatable manner
// using the provided cryptographic hash. It includes the nonce in the hash
func (m *SharePiece) Digest(nonce []byte, h hash.Hash) []byte {
	h.Reset()

	// Hash the signature piece
	h.Write(m.Piece)

	// Hash the round ID
	h.Write(serializeUin64(m.RoundID))

	// Hash the participants in the message
	for _, participant := range m.Participants {
		h.Write(participant)
	}

	// Hash the nonce
	h.Write(nonce)

	// Return the hash
	return h.Sum(nil)
}
