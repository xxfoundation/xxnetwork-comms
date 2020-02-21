////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package mixmessages

import (
	"bytes"
	"encoding/binary"
	"gitlab.com/elixxir/crypto/signature"
	"strconv"
	"testing"
)

// Happy path
func TestRoundInfo_ClearSignature(t *testing.T) {
	// Create an ndf and set it's signature
	testSign := []byte{1, 2, 45, 67, 42}
	testRoundInfo := &RoundInfo{
		Signature: testSign,
	}

	// Clear the signature
	testRoundInfo.ClearSignature()

	// Check that the signature is indeed nil after clearing
	if testRoundInfo.Signature != nil {
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
	if bytes.Compare(testRoundInfo.Signature, testSign) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", testSign, testRoundInfo.Signature)
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
func TestRoundInfo_Marshal(t *testing.T) {
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

	serializedData := testRoundInfo.Marshal()

	// ------------- Replicate marshal logic ------------------

	// Create the byte array
	testData := make([]byte, 0)

	// Serialize the id into a temp buffer of uint64 size (ie 8 bytes)
	tmp := make([]byte, 8)
	binary.PutUvarint(tmp, testId)

	// Append that temp buffer into the return buffer
	testData = append(testData, tmp...)

	// Serialize the boolean value
	testData = strconv.AppendBool(testData, testRealtime)

	// Serialize the batchSize into a temp buffer of uint32 size (ie 4 bytes)
	tmp = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, testBatch)

	// Append that temp buffer into the return buffer
	testData = append(testData, tmp...)

	// Serialize the entire topology
	for _, val := range testTopology {
		testData = append(testData, []byte(val)...)
	}

	// Serialize the StartTime into a temp buffer of uint64 size (ie 8 bytes)
	tmp = make([]byte, 8)
	binary.PutUvarint(tmp, testTime)

	// Append that temp buffer into the return buffer
	testData = append(testData, tmp...)

	// Compare the replicated data and the actual data
	if bytes.Compare(serializedData, testData) != 0 {
		t.Errorf("Marshalled data does not match contents!"+
			"Expected: %+v \n\t"+
			"Recieved: %+v", testData, serializedData)
	}
}

// Error path
func TestRoundInfo_Marshal_Error(t *testing.T) {
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

	serializedData := testRoundInfo.Marshal()

	// ------------- Replicate marshal logic ------------------

	// Create the byte array
	badSerializedData := make([]byte, 0)

	// Serialize the id into a temp buffer of uint64 size (ie 8 bytes)
	tmp := make([]byte, 8)
	binary.PutUvarint(tmp, testId)

	// Append that temp buffer into the return buffer
	badSerializedData = append(badSerializedData, tmp...)

	/* Omit serializing arbitrary fields

	// Serialize the boolean value
	b = strconv.AppendBool(b, testRealtime)

	// Serialize the batchSize into a temp buffer of uint32 size (ie 4 bytes)
	tmp = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, testBatch)

	// Append that temp buffer into the return buffer
	b = append(b, tmp...)

	// Serialize the entire topology
	for _, val := range testTopology {
		b = append(b, []byte(val)...)
	}
	*/

	// Serialize the StartTime into a temp buffer of uint64 size (ie 8 bytes)
	tmp = make([]byte, 8)
	binary.PutUvarint(tmp, testTime)

	// Append that temp buffer into the return buffer
	badSerializedData = append(badSerializedData, tmp...)

	// Compare an incomplete serialization to the marsalled data
	if bytes.Compare(serializedData, badSerializedData) != 0 {
		return
	}

	t.Errorf("Expected error path: Marshaled data should not match " +
		"manually locally serialized data")
}

// Happy path
func TestRoundInfo_GetSignature(t *testing.T) {
	// Create roundInfo and set signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	testRoundInfo := &RoundInfo{
		Signature: expectedSig,
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
