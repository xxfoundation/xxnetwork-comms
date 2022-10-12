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
// If this ever fails, check for modifications in the source library
//  as well as for this message type
var _ = signature.GenericRsaSignable(&SharePiece{})

// -------------------------- Signature tests --------------------------------------

// Happy path
func TestSharePiece_GetSignature(t *testing.T) {
	// Create roundErr and set signature (without using setSignature)
	expectedSig := []byte("expectedSig")
	expectedNonce := []byte("expectedNonce")
	expectedRsaSig := &messages.RSASignature{
		Signature: expectedSig,
		Nonce:     expectedNonce,
	}

	testRoundError := &SharePiece{Signature: expectedRsaSig}

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

// Consistency test
func TestShare_Digest_Consistency(t *testing.T) {
	// Generate a message
	testRoundID := uint64(25)
	testPiece := []byte("testPiece")
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testSharePiece := &SharePiece{
		Piece:        testPiece,
		Participants: testTopology,
		RoundID:      testRoundID,
	}
	// Hardcoded digest output. Any changes are a smoke test of changing of
	// lower level crypto libraries or changes of the digest() implementation
	expectedDigestEncoded := "t42oN4TqtXvSjlv1OBYeXfsTLPoI/NyMKBk9NCXrTJ8="

	// Generate a digest
	sha := crypto.SHA256.New()
	testNonce := []byte("expectedNonce")
	digest := testSharePiece.Digest(testNonce, sha)

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
func TestSharePiece_Digest(t *testing.T) {
	// Generate a message
	testRoundID := uint64(25)
	testPiece := []byte("testPiece")
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testSharePiece := &SharePiece{
		Piece:        testPiece,
		Participants: testTopology,
		RoundID:      testRoundID,
	}

	// Generate a digest
	sha := crypto.SHA256.New()
	testNonce := []byte("expectedNonce")
	receivedDigest := testSharePiece.Digest(testNonce, sha)

	// Manually generate the digest
	sha.Reset()
	sha.Write(testPiece)
	sha.Write(serializeUin64(testRoundID))
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
func TestSharePiece_SignVerify(t *testing.T) {
	// Generate a message
	testRoundID := uint64(25)
	testPiece := []byte("testPiece")
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testSharePiece := &SharePiece{
		Piece:        testPiece,
		Participants: testTopology,
		RoundID:      testRoundID,
	}
	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Ensure message type conforms to genericSignable interface
	err = signature.SignRsa(testSharePiece, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Verify signature
	err = signature.VerifyRsa(testSharePiece, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Failed to verify: %+v", err)
	}

}

// Error path
func TestShare_SignVerify_Error(t *testing.T) {
	// Generate a message
	testRoundID := uint64(25)
	testPiece := []byte("testPiece")
	testTopology := [][]byte{[]byte("test"), []byte("te"), []byte("st"), []byte("testtest")}
	testSharePiece := &SharePiece{
		Piece:        testPiece,
		Participants: testTopology,
		RoundID:      testRoundID,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Ensure message type conforms to genericSignable interface
	err = signature.SignRsa(testSharePiece, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Reset Participants value so verify()'s signature won't match
	testSharePiece.Participants = [][]byte{[]byte("I"), []byte("am"), []byte("totally"), []byte("failing right now")}
	// Verify signature
	err = signature.VerifyRsa(testSharePiece, pubKey)
	if err != nil {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}
