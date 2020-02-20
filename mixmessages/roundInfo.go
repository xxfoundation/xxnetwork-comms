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

	// Serialize the batchSize into a temp buffer of uint32 size (ie 4 bytes)
	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, m.BatchSize)

	// Append that temp buffer into the return buffer
	b = append(b, tmp...)

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
