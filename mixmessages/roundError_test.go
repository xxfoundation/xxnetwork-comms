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
func TestRoundError_ClearSignature(t *testing.T) {
	// Create an ndf and set it's signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}

	testRoundError := &RoundError{
		RsaSignature: sig,
	}

	// Clear the signature
	testRoundError.ClearSignature()

	// Check that the signature is indeed nil after clearing
	if testRoundError.RsaSignature != nil {
		t.Errorf("Signature should be nil after a clear signature call")
	}
}

// Happy path
func TestRoundError_SetSignature(t *testing.T) {
	expectedSig := []byte{1, 2, 45, 67, 42}
	testRoundError := &RoundError{}

	// Set the sig
	testRoundError.SetSignature(expectedSig)

	// Check that the ndf's signature is identical to the one set
	if bytes.Compare(testRoundError.GetSignature(), expectedSig) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, testRoundError.GetSignature())
	}
}

// Error path
func TestRoundError_SetSignature_Error(t *testing.T) {
	testRoundError := &RoundError{}

	// Set the sig to nil (error case)
	err := testRoundError.SetSignature(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// Happy path
func TestRoundError_GetSignature(t *testing.T) {
	// Create roundErr and set signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}

	testRoundError := &RoundError{
		RsaSignature: sig,
	}

	// Fetch signature
	receivedSig := testRoundError.GetSignature()

	// Compare fetched value to expected value
	if bytes.Compare(expectedSig, receivedSig) != 0 {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, receivedSig)
	}
}

// Happy path
func TestRoundError_SignVerify(t *testing.T) {
	// Create roundInfo object (to be used for roundError object)
	testId := uint64(25)
	testTopology := []string{"test", "te", "st", "testtest"}
	testRealtime := false
	testTime := uint64(49)
	testBatch := uint32(23)
	testRoundInfo := &RoundInfo{
		ID:        testId,
		Realtime:  testRealtime,
		Topology:  testTopology,
		StartTime: testTime,
		BatchSize: testBatch,
	}

	// Create RoundError object
	testError := "I failed. Fix me now!"
	testRoundError := &RoundError{
		Info:  testRoundInfo,
		Error: testError,
	}

	// Ensure message type conforms to genericSignable interface
	signature.Sign(testRoundError)

	// Verify signature
	if !signature.Verify(testRoundError) {
		t.Error("Expected happy path: Failed to verify!")
	}

}

// Error path
func TestRoundError_SignVerify_Error(t *testing.T) {
	// Create roundInfo object (to be used for roundError object)
	testId := uint64(25)
	testTopology := []string{"test", "te", "st", "testtest"}
	testRealtime := false
	testTime := uint64(49)
	testBatch := uint32(23)
	testRoundInfo := &RoundInfo{
		ID:        testId,
		Realtime:  testRealtime,
		Topology:  testTopology,
		StartTime: testTime,
		BatchSize: testBatch,
	}

	// Create RoundError object
	testError := "I failed. Fix me now!"
	testRoundError := &RoundError{
		Info:  testRoundInfo,
		Error: testError,
	}

	// Ensure message type conforms to genericSignable interface
	signature.Sign(testRoundError)

	// Reset Error value so verify()'s signature won't match
	testRoundError.Error = "Not an expected error message"

	// Verify signature
	if !signature.Verify(testRoundError) {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}
