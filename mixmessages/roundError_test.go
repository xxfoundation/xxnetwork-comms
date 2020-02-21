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
func TestRoundError_ClearSignature(t *testing.T) {
	// Create an ndf and set it's signature
	testSign := []byte{1, 2, 45, 67, 42}
	testRoundError := &RoundError{
		Signature: testSign,
	}

	// Clear the signature
	testRoundError.ClearSignature()

	// Check that the signature is indeed nil after clearing
	if testRoundError.Signature != nil {
		t.Errorf("Signature should be nil after a clear signature call")
	}
}

// Happy path
func TestRoundError_SetSignature(t *testing.T) {
	testSign := []byte{1, 2, 45, 67, 42}

	testRoundError := &RoundError{}

	// Set the sig
	testRoundError.SetSignature(testSign)

	// Check that the ndf's signature is identical to the one set
	if bytes.Compare(testRoundError.Signature, testSign) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", testSign, testRoundError.Signature)
	}
}

func TestRoundError_SetSignature_Error(t *testing.T) {
	testRoundError := &RoundError{}

	// Set the sig
	err := testRoundError.SetSignature(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

func TestRoundError_Marshal(t *testing.T) {
	// ------ Set fields -----------
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

	testError := "I failed. Fix me now!"

	testRoundError := &RoundError{
		Info:  testRoundInfo,
		Error: testError,
	}

	serializedData := testRoundError.Marshal()

	// ------------- Replicate marshal logic ------------------

	// Create the byte array
	testData := make([]byte, 0)

	// Marshall the roundInfo data and append to byte array
	testData = append(testData, testRoundError.Info.Marshal()...)

	// Serialize the error message
	testData = append(testData, []byte(testError)...)

	// Compare the replicated data and the actual data
	if bytes.Compare(serializedData, testData) != 0 {
		t.Errorf("Marshalled data does not match contents!"+
			"Expected: %+v \n\t"+
			"Recieved: %+v", testData, serializedData)
	}
}

// Error path
func TestRoundError_Marshal_Error(t *testing.T) {
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

	testError := "I failed. Fix me now!"

	testRoundError := &RoundError{
		Info:  testRoundInfo,
		Error: testError,
	}

	serializedData := testRoundError.Marshal()

	// ------------- Replicate marshal logic ------------------

	// Create the byte array
	badSerializedData := make([]byte, 0)

	// Marshall the roundInfo data and append to byte array
	badSerializedData = append(badSerializedData, testRoundError.Info.Marshal()...)

	/* Omit serializing arbitrary fields

	// Serialize the error message
	testData = append(testData, []byte(testError)...)

	*/

	if bytes.Compare(serializedData, badSerializedData) != 0 {
		return
	}

	t.Errorf("Expected error path: Marshaled data should not match " +
		"manually locally serialized data")
}

func TestRoundError_GetSignature(t *testing.T) {
	testSign := []byte{1, 2, 45, 67, 42}
	testRoundError := &RoundError{
		Signature: testSign,
	}

	ourSig := testRoundError.GetSignature()

	if bytes.Compare(testSign, ourSig) != 0 {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", testSign, ourSig)
	}
}
