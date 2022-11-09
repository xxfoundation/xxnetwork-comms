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
	"crypto"
	"encoding/binary"
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

func (m *RoundInfo) GetEccSig() *messages.ECCSignature {
	if m.EccSignature != nil {
		return m.EccSignature
	}

	m.EccSignature = new(messages.ECCSignature)

	return m.EccSignature
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

	// Serialize and hash resource queue timeout
	h.Write(serializeUin32(m.ResourceQueueTimeoutMillis))

	// Serialize and hash address space size
	h.Write(serializeUin32(m.AddressSpaceSize))

	// Hash the topology
	for _, node := range m.Topology {
		h.Write(node)
	}

	// Serialize and hash the timestamps
	for _, timeStamp := range m.Timestamps {
		h.Write(serializeUin64(timeStamp))
	}

	// Hash ClientErrors
	for _, clientError := range m.ClientErrors {
		sha := crypto.SHA256.New()
		data := clientError.Digest(nonce, sha)
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

func (m *RoundInfo) DeepCopy() *RoundInfo {
	return &RoundInfo{
		ID:                         m.ID,
		UpdateID:                   m.UpdateID,
		State:                      m.State,
		BatchSize:                  m.BatchSize,
		Topology:                   m.Topology,
		Timestamps:                 m.Timestamps,
		Errors:                     m.Errors,
		ClientErrors:               m.ClientErrors,
		ResourceQueueTimeoutMillis: m.ResourceQueueTimeoutMillis,
		Signature:                  m.Signature,
		AddressSpaceSize:           m.AddressSpaceSize,
		EccSignature:               m.EccSignature,
	}
}
