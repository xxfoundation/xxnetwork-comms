////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package mixmessages

import "hash"

// Digest hashes the contents of the message in a repeatable manner
// using the provided cryptographic hash. It includes the nonce in the hash
func (m *ClientError) Digest(nonce []byte, h hash.Hash) []byte {
	h.Reset()

	// Hash the nodeId
	h.Write(m.ClientId)
	h.Write([]byte(m.Error))

	// Hash the nonce
	h.Write(nonce)

	// Return the hash
	return h.Sum(nil)
}
