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
func TestRoundInfo_ClearSignature(t *testing.T) {
	// Create an RoundInfo and set it's signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}

	testRoundInfo := &RoundInfo{
		RsaSignature: sig,
	}

	// Clear the signature
	testRoundInfo.ClearSignature()

	// Check that the signature's values are nil after clearing
	if testRoundInfo.GetNonce() != nil && testRoundInfo.GetSignature() != nil {
		t.Errorf("Signature's values should be nil after a ClearSignature() call!"+
			"\n\tSignature is: %+v", testRoundInfo.RsaSignature)
	}
}

// ------------------------------- Nonce tests -----------------------------------------------------

// Happy path
func TestRoundInfo_GetNonce(t *testing.T) {
	expectedNonce := []byte{1, 2, 45, 67, 42}

	// Create message with nonce value
	sig := &RSASignature{Nonce: expectedNonce}
	testRoundInfo := &RoundInfo{
		RsaSignature: sig,
	}

	// Retrieve the nonce
	receivedNonce := testRoundInfo.GetNonce()

	// Compare to the value originally set
	if bytes.Compare(expectedNonce, receivedNonce) != 0 {
		t.Errorf("Nonce does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, receivedNonce)

	}
}

// Error path (nil object)
func TestRoundInfo_GetNonce_NilObject(t *testing.T) {
	// Create RoundInfo w/o signature object
	testRoundInfo := &RoundInfo{}

	// Attempt to get nonce
	receivedSig := testRoundInfo.GetNonce()

	// Received nonce should be nil
	if receivedSig != nil {
		t.Errorf("Nonce should default to nil if not set!")
	}

}

//
func TestRoundInfo_SetNonce(t *testing.T) {
	// Create RoundInfo message
	tempVal := []byte("fail Fail fail")
	tempSig := &RSASignature{Nonce: tempVal}
	testRoundInfo := &RoundInfo{RsaSignature: tempSig}

	// Set the sig
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testRoundInfo.SetNonce(expectedNonce)

	// Check that the RoundInfo's signature is identical to the one set
	if bytes.Compare(testRoundInfo.RsaSignature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testRoundInfo.RsaSignature.Nonce)
	}
}

// Happy path
func TestRoundInfo_SetNonce_NilObject(t *testing.T) {
	testRoundInfo := &RoundInfo{}

	// Set the sig w/o signature being initialized
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testRoundInfo.SetNonce(expectedNonce)

	// Sig should be set to expected value
	if bytes.Compare(testRoundInfo.RsaSignature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testRoundInfo.RsaSignature.Nonce)
	}
}

// Error path
func TestRoundInfo_SetNonce_SetNil(t *testing.T) {
	// Create signature object
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{
		Signature: expectedSig,
		Nonce:     expectedSig,
	}

	// Create RoundInfo message
	testRoundInfo := &RoundInfo{RsaSignature: sig}

	// Set the sig to nil (error case)
	err := testRoundInfo.SetNonce(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// -------------------------- Signature tests --------------------------------------

// Happy path
func TestRoundInfo_SetSignature(t *testing.T) {
	testSign := []byte{1, 2, 45, 67, 42}

	testRoundInfo := &RoundInfo{}

	// Set the sig
	testRoundInfo.SetSignature(testSign)

	// Check that the RoundInfo's signature is identical to the one set
	if bytes.Compare(testRoundInfo.GetSignature(), testSign) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", testSign, testRoundInfo.GetSignature())
	}
}

// Error path
func TestRoundInfo_SetSignature_Error(t *testing.T) {
	testRoundInfo := &RoundInfo{}

	// Set the sig
	err := testRoundInfo.SetSignature(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// Happy path
func TestRoundInfo_GetSignature(t *testing.T) {
	// Create roundInfo and set signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}

	testRoundInfo := &RoundInfo{
		RsaSignature: sig,
	}

	// Fetch signature
	receivedSig := testRoundInfo.GetSignature()

	// Compare fetched value to expected value
	if bytes.Compare(expectedSig, receivedSig) != 0 {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, receivedSig)
	}
}

// Error path (nil signature)
func TestRoundInfo_GetSignature_NilObject(t *testing.T) {
	// Create RoundInfo w/o signature object
	testRoundInfo := &RoundInfo{}

	// Attempt to get signature
	receivedSig := testRoundInfo.GetSignature()

	// Received sig should be nil
	if receivedSig != nil {
		t.Errorf("Signature should default to nil if not set!")
	}

}

// ------------------------------ Sign/Verify tests -----------------------------------

// Happy path
func TestRoundInfo_SignVerify(t *testing.T) {
	// Create roundInfo object (to be used for RoundInfo object)
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

	// Ensure message type conforms to genericSignable interface
	signature.Sign(testRoundInfo)

	// Verify signature
	if !signature.Verify(testRoundInfo) {
		t.Error("Expected happy path: Failed to verify!")
	}
}

// Error path
func TestRoundInfo_SignVerify_Error(t *testing.T) {
	// Create roundInfo object
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

	// Ensure message type conforms to genericSignable interface
	signature.Sign(testRoundInfo)

	// Reset Topology value so verify()'s signature won't match
	testRoundInfo.Topology = []string{"fail", "fa", "il", "failfail"}
	// Verify signature
	if !signature.Verify(testRoundInfo) {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}
