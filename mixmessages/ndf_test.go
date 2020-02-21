////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package mixmessages

import (
	"bytes"
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

func TestNDF_SetSignature_Error(t *testing.T) {
	testNdf := &NDF{}

	// Set the sig
	err := testNdf.SetSignature(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

func TestNDF_Marshal(t *testing.T) {
	ourNdf := []byte{25, 254, 123, 42}
	testNdf := &NDF{
		Ndf: ourNdf,
	}

	serializedData := testNdf.Marshal()

	// This test assumes serialized ndf message is just an ndf
	// If the Marshal() logic ever changes, this test may fail
	if bytes.Compare(serializedData, ourNdf) != 0 {
		t.Errorf("Marshalled data does not match contents!"+
			"Expected: %+v \n\t"+
			"Recieved: %+v", ourNdf, serializedData)
	}
}

func TestNDF_GetSignature(t *testing.T) {
	testSign := []byte{1, 2, 45, 67, 42}
	testNdf := &NDF{
		Signature: testSign,
	}

	ourSig := testNdf.GetSignature()

	if bytes.Compare(testSign, ourSig) != 0 {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", testSign, ourSig)
	}
}
