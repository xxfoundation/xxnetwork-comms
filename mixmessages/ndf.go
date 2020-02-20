////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functions to make the ndf message type conform to a generic
// signing interface

package mixmessages

// Marshal serializes all the data needed for a signature
func (m *NDF) Marshal() []byte {
	// Create the byte array
	b := make([]byte, 0)

	// Marshall the roundInfo data and append to byte array
	b = append(b, m.Ndf...)

	return b
}

// SetSignature sets NDF's signature to the newSig argument
func (m *NDF) SetSignature(newSig []byte) {
	m.Signature = newSig
}

// ClearSignature clears out NDF's signature by setting it to nil
func (m *NDF) ClearSignature() {
	m.Signature = nil
}
