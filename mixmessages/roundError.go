////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the RoundError message type conform to a generic
// signing interface

package mixmessages

// Marshal serializes all the data needed for a signature
func (m *RoundError) Marshal() []byte {
	// Create the byte array
	b := make([]byte, 0)

	// Marshall the roundInfo data and append to byte array
	b = append(b, m.Info.Marshal()...)

	// Serialize the error message
	b = append(b, []byte(m.Error)...)

	return b
}

// SetSignature sets RoundError's signature to the newSig argument
func (m *RoundError) SetSignature(newSig []byte) {
	m.Signature = newSig
}

// ClearSignature clears out RoundError's signature by setting it to nil
func (m *RoundError) ClearSignature() {
	m.Signature = nil
}
