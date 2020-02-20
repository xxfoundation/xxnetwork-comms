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
	"strconv"
)

// Marshal serializes all the data needed for a signature
func (m *RoundInfo) Marshal() []byte {
	// Create the byte array
	b := make([]byte, 0)

	// Serialize the id into the array
	binary.PutUvarint(b, m.ID)

	// Serialize the boolean value
	b = strconv.AppendBool(b, m.Realtime)

	// Serialize the batchsize
	binary.LittleEndian.PutUint32(b, m.BatchSize)

	// Serialize the entire topology
	for _, val := range m.Topology {
		b = append(b, []byte(val)...)
	}

	// Serialize the start time
	binary.PutUvarint(b, m.StartTime)

	return b
}

// SetSignature sets RoundInfo's signature to the newSig argument
func (m *RoundInfo) SetSignature(newSig []byte) {
	m.Signature = newSig
}

// ClearSignature clears out roundInfo's signature by
// setting it to nil
func (m *RoundInfo) ClearSignature() {
	m.Signature = nil
}
