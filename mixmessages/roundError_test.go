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
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/comms/signature"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"reflect"
	"testing"
)

// Ensure message type conforms to genericSignable interface
// If this ever fails, check for modifications in the crypto library
//  as well as for this message type
var _ = signature.GenericRsaSignable(&RoundError{})

// -------------------------------- Get tests -----------------------------------------------------------

// Happy path
func TestRoundError_GetSignature(t *testing.T) {
	// Create roundErr and set signature (without using setSignature)
	expectedSig := []byte("expectedSig")
	expectedNonce := []byte("expectedNonce")
	expectedRsaSig := &messages.RSASignature{
		Signature: expectedSig,
		Nonce:     expectedNonce,
	}

	testRoundError := &RoundError{Signature: expectedRsaSig}

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

func TestRoundError_DigestTestHelper(t *testing.T) {
	testRoundErr := &RoundError{}
	checkdigest(t, testRoundErr)
}

// Consistency test
func TestRoundError_Digest_Consistency(t *testing.T) {
	// Generate a message
	testNodeId := []byte("nodeId")
	testError := "I failed. Fix me now!"
	testID := uint64(0)
	testRoundErr := &RoundError{
		Id:     testID,
		NodeId: testNodeId,
		Error:  testError,
	}

	// Hardcoded digest output. Any changes are a smoke test of changing of
	// crypto libraries
	expectedDigestEncoded := "M6v7dGA97cSP4bAmp3wJoRkrnUse34ouEFQoMYgZG2w="

	// Generate a digest
	sha := crypto.SHA256.New()
	testNonce := []byte("expectedNonce")
	digest := testRoundErr.Digest(testNonce, sha)

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
func TestRoundError_Digest(t *testing.T) {
	// Generate a message
	testNodeId := []byte("nodeId")
	testError := "I failed. Fix me now!"
	testID := uint64(240)
	testRoundErr := &RoundError{
		Id:     testID,
		NodeId: testNodeId,
		Error:  testError,
	}
	// Generate a digest
	sha := crypto.SHA256.New()
	testNonce := []byte("expectedNonce")
	receivedDigest := testRoundErr.Digest(testNonce, sha)

	// Manually generate the digest
	sha.Reset()
	sha.Write(testNodeId)
	sha.Write([]byte(testError))
	sha.Write(serializeUin64(testID))

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
func TestRoundError_SignVerify(t *testing.T) {

	// Generate a message
	testNodeId := []byte("nodeId")
	testError := "I failed. Fix me now!"
	testID := uint64(240)
	testRoundErr := &RoundError{
		Id:     testID,
		NodeId: testNodeId,
		Error:  testError,
	}
	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Sign message
	err = signature.SignRsa(testRoundErr, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Verify signature
	err = signature.VerifyRsa(testRoundErr, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Failed to verify: %+v", err)
	}

}

// Error path: Change internals of message between signing and verifying
func TestRoundError_SignVerify_Error(t *testing.T) {
	// Generate a message
	testNodeId := []byte("nodeId")
	testError := "I failed. Fix me now!"
	testID := uint64(240)
	testRoundErr := &RoundError{
		Id:     testID,
		NodeId: testNodeId,
		Error:  testError,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Ensure message type conforms to genericSignable interface
	err = signature.SignRsa(testRoundErr, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Reset Error value so verify()'s signature won't match
	testRoundErr.Error = "invalidChange"

	// Verify signature
	err = signature.VerifyRsa(testRoundErr, pubKey)
	if err == nil {
		t.Error("Expected error path: Should not have verified!")

	}

}
