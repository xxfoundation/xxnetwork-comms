///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundInfo message type conform to a generic
// signing interface

package mixmessages

import (
	"crypto"
	"encoding/binary"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"hash"
)

// GetSig returns the RSA signature.
// IF none exists, it creates it, adds it to the object, then returns it.
func (m *RoundInfo) GetSig() *messages.RSASignature {
	if m.Signature != nil {
		return m.Signature
	}

	m.Signature = new(messages.RSASignature)

	return m.Signature
}

// SetSignature sets RoundError's signature to the newSig argument
func (m *RoundInfo) SetSignature(newSig, nonce []byte) error {
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
func (m *RoundInfo) Digest(nonce []byte, h hash.Hash) []byte {
	h.Reset()

	// Serialize  and hash RoundId
	h.Write(serializeUin64(m.ID))

	// Serialize and hash UpdateId
	h.Write(serializeUin64(m.UpdateID))

	// Serialize and hash state
	h.Write(serializeUin32(m.State))

	// Serialize and hash batch size
	h.Write(serializeUin32(m.BatchSize))

	// Hash the topology
	for _, node := range m.Topology {
		h.Write(node)
	}

	// Serialize and hash the timestamps
	for _, timeStamp := range m.Timestamps {
		h.Write(serializeUin64(timeStamp))
	}

	// Hash the
	for _, roundErr := range m.Errors {
		sha := crypto.SHA256.New()
		data := roundErr.Digest(nonce, sha)
		h.Write(data)
	}

	// Hash nonce
	h.Write(nonce)

	// Return the hash
	return h.Sum(nil)
}

// serializeUin64 is a helper function which serializes
// any uint64 data into a byte slice for hashing purposes
func serializeUin64(data uint64) []byte {
	serializedData := make([]byte, 8)
	binary.LittleEndian.PutUint64(serializedData, data)
	return serializedData
}

// serializeUin32 is a helper function which serializes
// any uint32 data into a byte slice for hashing purposes
func serializeUin32(data uint32) []byte {
	serializedData := make([]byte, 4)
	binary.LittleEndian.PutUint32(serializedData, data)
	return serializedData
}

// GetActivity gets the state of the node
func (m *RoundInfo) GetRoundState() states.Round {
	return states.Round(m.State)
}

func (m *RoundInfo) GetRoundId() id.Round {
	return id.Round(m.ID)
}
