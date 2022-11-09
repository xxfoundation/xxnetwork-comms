////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package mixmessages

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"gitlab.com/elixxir/primitives/states"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/ec"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"reflect"
	"testing"
)

// Ensure message type conforms to genericSignable interface
// If this ever fails, check for modifications in the source library
//  as well as for this message type
var _ = signature.GenericRsaSignable(&RoundInfo{})
var _ = signature.GenericEccSignable(&RoundInfo{})

// -------------------------- Get tests --------------------------------------

// Happy path
func TestRoundInfo_GetSignature(t *testing.T) {
	// Create roundErr and set signature (without using setSignature)
	expectedSig := []byte("expectedSig")
	expectedNonce := []byte("expectedNonce")
	expectedRsaSig := &messages.RSASignature{
		Signature: expectedSig,
		Nonce:     expectedNonce,
	}

	testRoundError := &RoundInfo{Signature: expectedRsaSig}

	// Fetch signature
	receivedSig := testRoundError.GetSig()

	// Compare fetched value to expected value
	if !reflect.DeepEqual(expectedRsaSig, receivedSig) {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedRsaSig, receivedSig)
	}

}

// -------------------- Digest tests -------------------------------

func TestRoundInfo_DigestTestHelper(t *testing.T) {
	testRoundInfo := &RoundInfo{}
	checkdigest(t, testRoundInfo)
}

// Consistency test
func TestRoundInfo_Digest_Consistency(t *testing.T) {
	// Generate a message
	testId := uint64(25)
	testUpdateId := uint64(26)
	testState := uint32(42)
	testBatch := uint32(23)
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testRoundInfo := &RoundInfo{
		ID:        testId,
		UpdateID:  testUpdateId,
		State:     testState,
		BatchSize: testBatch,
		Topology:  testTopology,
	}
	// Hardcoded digest output. Any changes are a smoke test of changing of
	// crypto libraries
	expectedDigestEncoded := "5JubUJL0XrIbx4G69V+SNHzY0w9Bs3LcZQnNd+ZJMZQ="

	// Generate a digest
	sha := crypto.SHA256.New()
	testNonce := []byte("expectedNonce")
	digest := testRoundInfo.Digest(testNonce, sha)

	// Encode outputted digest to base64 encoded string
	receivedDigestEncoded := base64.StdEncoding.EncodeToString(digest)

	// Check the consistency of generated digest and hard-coded digest
	if expectedDigestEncoded != receivedDigestEncoded {
		t.Errorf("Consistency test failed for testNDF."+
			"\n\tExpected: %v"+
			"\n\tRecieved: %v", expectedDigestEncoded, receivedDigestEncoded)
	}
}

// Test that digest output matches manual digest creation
func TestRoundInfo_Digest(t *testing.T) {
	// Generate a message
	testId := uint64(25)
	testUpdateId := uint64(26)
	testState := uint32(42)
	testBatch := uint32(23)
	testResourceQueueTimeout := uint32(1000)
	testAddressSpaceSize := uint32(10)
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testRoundInfo := &RoundInfo{
		ID:                         testId,
		UpdateID:                   testUpdateId,
		State:                      testState,
		BatchSize:                  testBatch,
		Topology:                   testTopology,
		ResourceQueueTimeoutMillis: testResourceQueueTimeout,
		AddressSpaceSize:           testAddressSpaceSize,
	}

	// Generate a digest
	sha := crypto.SHA256.New()
	testNonce := []byte("expectedNonce")
	receivedDigest := testRoundInfo.Digest(testNonce, sha)

	// Manually generate the digest
	sha.Reset()
	sha.Write(serializeUin64(testId))
	sha.Write(serializeUin64(testUpdateId))
	sha.Write(serializeUin32(testState))
	sha.Write(serializeUin32(testBatch))
	sha.Write(serializeUin32(testResourceQueueTimeout))
	sha.Write(serializeUin32(testAddressSpaceSize))
	for _, node := range testTopology {
		sha.Write(node)
	}
	sha.Write(testNonce)
	expectedDigest := sha.Sum(nil)
	// Check that manual digest matches expected digest
	if !bytes.Equal(receivedDigest, expectedDigest) {
		t.Errorf("Digest did not output expected result."+
			"\n\tExpected: %v"+
			"\n\tRecieved: %v", expectedDigest, receivedDigest)
	}

}

// ------------------------------ Sign/Verify tests -----------------------------------

// Happy path
func TestRoundInfo_SignVerify(t *testing.T) {
	// Create roundInfo object
	testId := uint64(25)
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testBatch := uint32(23)
	testRoundInfo := &RoundInfo{
		ID:        testId,
		Topology:  testTopology,
		BatchSize: testBatch,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Ensure message type conforms to genericSignable interface
	err = signature.SignRsa(testRoundInfo, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Verify signature
	err = signature.VerifyRsa(testRoundInfo, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Failed to verify: %+v", err)
	}

}

// Error path
func TestRoundInfo_SignVerify_Error(t *testing.T) {
	// Create roundInfo object
	testId := uint64(25)
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testBatch := uint32(23)
	testRoundInfo := &RoundInfo{
		ID:        testId,
		Topology:  testTopology,
		BatchSize: testBatch,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Ensure message type conforms to genericSignable interface
	err = signature.SignRsa(testRoundInfo, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Reset Topology value so verify()'s signature won't match
	testRoundInfo.Topology = [][]byte{[]byte("I"), []byte("am"), []byte("totally"), []byte("failing right now")}
	// Verify signature
	err = signature.VerifyRsa(testRoundInfo, pubKey)
	if err != nil {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}

// Happy path
func TestNDF_SignVerifyEddsa(t *testing.T) {
	// Create roundInfo object
	testId := uint64(25)
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testBatch := uint32(23)
	testRoundInfo := &RoundInfo{
		ID:        testId,
		Topology:  testTopology,
		BatchSize: testBatch,
	}
	// Generate keys
	privateKey, err := ec.NewKeyPair(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Sign message
	err = signature.SignEddsa(testRoundInfo, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Verify signature
	err = signature.VerifyEddsa(testRoundInfo, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Failed to verify: %+v", err)
	}

}

// Error path: Change internals of message between signing and verifying
func TestNdf_SignVerifyEddsa_Error(t *testing.T) {
	// Create roundInfo object
	testId := uint64(25)
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testBatch := uint32(23)
	testRoundInfo := &RoundInfo{
		ID:        testId,
		Topology:  testTopology,
		BatchSize: testBatch,
	}
	// Generate keys
	privateKey, err := ec.NewKeyPair(rand.Reader)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Sign message
	err = signature.SignEddsa(testRoundInfo, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Reset Topology value so verify()'s signature won't match
	testRoundInfo.Topology = [][]byte{[]byte("I"), []byte("am"), []byte("totally"), []byte("failing right now")}

	// Verify signature
	err = signature.VerifyEddsa(testRoundInfo, pubKey)
	// Verify signature
	if err == nil {
		t.Error("Expected error path: Should not have verified!")
	}

}

// ---------------------- Non-genericSignable tests -----------------------------------

func TestRoundInfo_GetActivity(t *testing.T) {
	expected := uint32(45)
	testRoundInfo := &RoundInfo{
		State: expected,
	}

	received := testRoundInfo.GetRoundState()

	if received != states.Round(expected) {
		t.Errorf("Received does not match expected for getter function! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expected, received)
	}
}
