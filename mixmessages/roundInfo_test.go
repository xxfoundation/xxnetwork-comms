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
	// Create an ndf and set it's signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}

	testRoundInfo := &RoundInfo{
		RsaSignature: sig,
	}

	// Clear the signature
	testRoundInfo.ClearSignature()

	// Check that the signature is indeed nil after clearing
	if testRoundInfo.RsaSignature != nil {
		t.Errorf("Signature should be nil after a clear signature call")
	}
}

// Happy path
func TestRoundInfo_SetSignature(t *testing.T) {
	testSign := []byte{1, 2, 45, 67, 42}

	testRoundInfo := &RoundInfo{}

	// Set the sig
	testRoundInfo.SetSignature(testSign)

	// Check that the ndf's signature is identical to the one set
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

// Happy path
func TestRoundInfo_SignVerify(t *testing.T) {
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
