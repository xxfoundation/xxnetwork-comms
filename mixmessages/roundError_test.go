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
	// Create an roundError and set it's signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}

	testRoundError := &RoundError{
		RsaSignature: sig,
	}

	// Clear the signature
	testRoundError.ClearSignature()

	// Check that the signature's values are nil after clearing
	if testRoundError.GetSignature() != nil && testRoundError.GetNonce() != nil {
		t.Errorf("Signature's values should be nil after a ClearSignature() call!"+
			"\n\tSignature is: %+v", testRoundError.RsaSignature)
	}
}

// ------------------------------- Nonce tests -----------------------------------------------------

// Happy path
func TestRoundError_GetNonce(t *testing.T) {
	expectedNonce := []byte{1, 2, 45, 67, 42}

	// Create message with nonce value
	sig := &RSASignature{Nonce: expectedNonce}
	testRoundError := &RoundError{
		RsaSignature: sig,
	}

	// Retrieve the nonce
	receivedNonce := testRoundError.GetNonce()

	// Compare to the value originally set
	if bytes.Compare(expectedNonce, receivedNonce) != 0 {
		t.Errorf("Nonce does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, receivedNonce)

	}
}

// Error path (nil object)
func TestRoundError_GetNonce_NilObject(t *testing.T) {
	// Create roundError w/o signature object
	testRoundError := &RoundError{}

	// Attempt to get nonce
	receivedSig := testRoundError.GetNonce()

	// Received nonce should be nil
	if receivedSig != nil {
		t.Errorf("Nonce should default to nil if not set!")
	}

}

//
func TestRoundError_SetNonce(t *testing.T) {
	// Create roundError message
	tempVal := []byte("fail Fail fail")
	tempSig := &RSASignature{Nonce: tempVal}
	testRoundError := &RoundError{RsaSignature: tempSig}

	// Set the sig
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testRoundError.SetNonce(expectedNonce)

	// Check that the roundError's signature is identical to the one set
	if bytes.Compare(testRoundError.RsaSignature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testRoundError.RsaSignature.Nonce)
	}
}

// Happy path
func TestRoundError_SetNonce_NilObject(t *testing.T) {
	testRoundError := &RoundError{}

	// Set the sig w/o signature being initialized
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testRoundError.SetNonce(expectedNonce)

	// Sig should be set to expected value
	if bytes.Compare(testRoundError.RsaSignature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testRoundError.RsaSignature.Nonce)
	}
}

// Error path
func TestRoundError_SetNonce_SetNil(t *testing.T) {
	// Create signature object
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{
		Signature: expectedSig,
		Nonce:     expectedSig,
	}

	// Create roundError message
	testRoundError := &RoundError{RsaSignature: sig}

	// Set the sig to nil (error case)
	err := testRoundError.SetNonce(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// -------------------------------- Signature tests -----------------------------------------------------------

// Happy path
func TestRoundError_SetSignature(t *testing.T) {
	expectedSig := []byte{1, 2, 45, 67, 42}
	testRoundError := &RoundError{}

	// Set the sig
	testRoundError.SetSignature(expectedSig)

	// Check that the roundError's signature is identical to the one set
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

// Error path (nil signature)
func TestRoundError_GetSignature_NilObject(t *testing.T) {
	// Create RoundError w/o signature object
	testRoundError := &RoundError{}

	// Attempt to get signature
	receivedSig := testRoundError.GetSignature()

	// Received sig should be nil
	if receivedSig != nil {
		t.Errorf("Signature should default to nil if not set!")
	}

}

// ------------------------------ Sign/Verify tests -----------------------------------

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
