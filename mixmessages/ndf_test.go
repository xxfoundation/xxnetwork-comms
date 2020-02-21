////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package mixmessages

import (
	"bytes"
	"gitlab.com/elixxir/crypto/signature"
	"testing"
)

// Happy path
func TestNDF_ClearSignature(t *testing.T) {
	// Create an ndf and set it's signature
	testSign := []byte{1, 2, 45, 67, 42}
	testNdf := &NDF{
		Signature: testSign,
	}

	// Clear the signature
	testNdf.ClearSignature()

	// Check that the signature is indeed nil after clearing
	if testNdf.Signature != nil {
		t.Errorf("Signature should be nil after a clear signature call")
	}
}

// Happy path
func TestNDF_SetSignature(t *testing.T) {
	testSign := []byte{1, 2, 45, 67, 42}
	testNdf := &NDF{}

	// Set the sig
	testNdf.SetSignature(testSign)

	// Check that the ndf's signature is identical to the one set
	if bytes.Compare(testNdf.Signature, testSign) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", testSign, testNdf.Signature)
	}
}

// Error path
func TestNDF_SetSignature_Error(t *testing.T) {
	testNdf := &NDF{}

	// Set the sig to nil (error case)
	err := testNdf.SetSignature(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// Happy path
func TestNDF_Marshal(t *testing.T) {
	// Create ndf object
	ourNdf := []byte{25, 254, 123, 42}
	testNdf := &NDF{
		Ndf: ourNdf,
	}

	// Marshal and compare to the original ndf bytes
	serializedData := testNdf.Marshal()

	// This test assumes serialized ndf message is just an ndf
	// If the Marshal() logic ever changes, this test may fail
	if bytes.Compare(serializedData, ourNdf) != 0 {
		t.Errorf("Marshalled data does not match contents!"+
			"Expected: %+v \n\t"+
			"Recieved: %+v", ourNdf, serializedData)
	}
}

// Happy path
func TestNDF_GetSignature(t *testing.T) {
	// Create ndf and set signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	testNdf := &NDF{
		Signature: expectedSig,
	}

	// Fetch signature
	receivedSig := testNdf.GetSignature()

	// Compare fetched value to expected value
	if bytes.Compare(expectedSig, receivedSig) != 0 {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, receivedSig)
	}
}

// Happy path
func TestNdf_Sign(t *testing.T) {
	// Create ndf object
	ourNdf := []byte{25, 254, 123, 42}
	testNdf := &NDF{
		Ndf: ourNdf,
	}

	// Ensure message type conforms to genericSignable interface
	signature.Sign(testNdf)

	// Verify signature
	if !signature.Verify(testNdf) {
		t.Error("Expected happy path: Failed to verify!")
	}
}

// Error path
func TestNdf_Sign_Error(t *testing.T) {
	// Create ndf object
	ourNdf := []byte{25, 254, 123, 42}
	testNdf := &NDF{
		Ndf: ourNdf,
	}

	// Ensure message type conforms to genericSignable interface
	signature.Sign(testNdf)

	// Reset ndf value so verify()'s signature won't match
	testNdf.Ndf = []byte{1}

	// Verify signature
	if !signature.Verify(testNdf) {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}
