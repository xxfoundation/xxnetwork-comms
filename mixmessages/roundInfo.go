////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundInfo message type conform to a generic
// signing interface

package mixmessages

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"strconv"
)

// Marshal serializes all the data needed for a signature
func (m *RoundInfo) Marshal() []byte {
	// Create the byte array
	b := make([]byte, 0)

	// Serialize the id into a temp buffer of uint64 size (ie 8 bytes)
	tmp := make([]byte, 8)
	binary.PutUvarint(tmp, m.ID)

	// Append that temp buffer into the return buffer
	b = append(b, tmp...)

	// Serialize the boolean value
	b = strconv.AppendBool(b, m.Realtime)

	// Serialize the batchSize into a temp buffer of uint32 size (ie 4 bytes)
	tmp = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, m.BatchSize)

	// Append that temp buffer into the return buffer
	b = append(b, tmp...)

	// Serialize the entire topology
	for _, val := range m.Topology {
		b = append(b, []byte(val)...)
	}

	// Serialize the StartTime into a temp buffer of uint64 size (ie 8 bytes)
	tmp = make([]byte, 8)
	binary.PutUvarint(tmp, m.StartTime)

	// Append that temp buffer into the return buffer
	b = append(b, tmp...)

	return b
}

// SetSignature sets RoundInfo's signature to the newSig argument
func (m *RoundInfo) SetSignature(newSig []byte) error {
	if newSig == nil {
		return errors.Errorf("Cannot set signature to nil value")
	}

	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Signature: newSig}
	}

	m.RsaSignature.Signature = newSig
	return nil
}

// ClearSignature clears out roundInfo's signature by
// setting it to nil
func (m *RoundInfo) ClearSignature() {

	m.RsaSignature = nil
}

// GetNonce gets the value of the nonce
func (m *RoundInfo) GetNonce() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.RsaSignature.GetNonce()
}

// SetSignature sets RoundError's nonce to the newNonce argument
func (m *RoundInfo) SetNonce(newNonce []byte) error {
	if newNonce == nil {
		return errors.Errorf("Cannot set nonce to nil")
	}
	if m.RsaSignature == nil {
		m.RsaSignature = &RSASignature{Nonce: newNonce}
	}

	m.RsaSignature.Nonce = newNonce

	return nil
}

//GetSignature
func (m *RoundInfo) GetSignature() []byte {
	if m.RsaSignature == nil {
		return nil
	}

	return m.GetRsaSignature().GetSignature()
}
