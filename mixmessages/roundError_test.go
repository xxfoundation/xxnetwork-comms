///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package mixmessages

import (
	"bytes"
	"crypto/rand"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/xx_network/comms/messages"
	"testing"
)

// Ensure message type conforms to genericSignable interface
// If this ever fails, check for modifications in the crypto library
//  as well as for this message type
var _ = signature.GenericSignable(&RoundError{})

// Happy path
func TestRoundError_ClearSignature(t *testing.T) {
	// Create an roundError and set it's signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &messages.RSASignature{Signature: expectedSig}

	testRoundError := &RoundError{
		Signature: sig,
	}

	// Clear the signature
	testRoundError.ClearSig()

	// Check that the signature's values are nil after clearing
	if testRoundError.GetSig() != nil && testRoundError.GetNonce() != nil {
		t.Errorf("Signature's values should be nil after a ClearSignature() call!"+
			"\n\tSignature is: %+v", testRoundError.Signature)
	}
}

// ------------------------------- Nonce tests -----------------------------------------------------

// Happy path
func TestRoundError_GetNonce(t *testing.T) {
	expectedNonce := []byte{1, 2, 45, 67, 42}

	// Create message with nonce value
	sig := &messages.RSASignature{Nonce: expectedNonce}
	testRoundError := &RoundError{
		Signature: sig,
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
	tempSig := &messages.RSASignature{Nonce: tempVal}
	testRoundError := &RoundError{Signature: tempSig}

	// Set the sig
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testRoundError.SetNonce(expectedNonce)

	// Check that the roundError's signature is identical to the one set
	if bytes.Compare(testRoundError.Signature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testRoundError.Signature.Nonce)
	}
}

// Happy path
func TestRoundError_SetNonce_NilObject(t *testing.T) {
	testRoundError := &RoundError{}

	// Set the sig w/o signature being initialized
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testRoundError.SetNonce(expectedNonce)

	// Sig should be set to expected value
	if bytes.Compare(testRoundError.Signature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testRoundError.Signature.Nonce)
	}
}

// Error path
func TestRoundError_SetNonce_SetNil(t *testing.T) {
	// Create signature object
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &messages.RSASignature{
		Signature: expectedSig,
		Nonce:     expectedSig,
	}

	// Create roundError message
	testRoundError := &RoundError{Signature: sig}

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
	testRoundError.SetSig(expectedSig)

	// Check that the roundError's signature is identical to the one set
	if bytes.Compare(testRoundError.GetSig(), expectedSig) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, testRoundError.GetSig())
	}
}

// Error path
func TestRoundError_SetSignature_Error(t *testing.T) {
	testRoundError := &RoundError{}

	// Set the sig to nil (error case)
	err := testRoundError.SetSig(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// Happy path
func TestRoundError_GetSignature(t *testing.T) {
	// Create roundErr and set signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &messages.RSASignature{Signature: expectedSig}

	testRoundError := &RoundError{
		Signature: sig,
	}

	// Fetch signature
	receivedSig := testRoundError.GetSig()

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
	receivedSig := testRoundError.GetSig()

	// Received sig should be nil
	if receivedSig != nil {
		t.Errorf("Signature should default to nil if not set!")
	}

}

// ------------------------------ Sign/Verify tests -----------------------------------

// Happy path
func TestRoundError_SignVerify(t *testing.T) {

	// Create RoundError object
	testError := "I failed. Fix me now!"
	testRoundError := &RoundError{
		Error: testError,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Sign message
	err = signature.Sign(testRoundError, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Verify signature
	err = signature.Verify(testRoundError, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Failed to verify: %+v", err)
	}

}

// Error path
func TestRoundError_SignVerify_Error(t *testing.T) {
	// Create RoundError object
	testError := "I failed. Fix me now!"
	testRoundError := &RoundError{
		Error: testError,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Ensure message type conforms to genericSignable interface
	err = signature.Sign(testRoundError, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Reset Error value so verify()'s signature won't match
	testRoundError.Error = "Not an expected error message"

	// Verify signature
	err = signature.Verify(testRoundError, pubKey)
	if err != nil {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}
