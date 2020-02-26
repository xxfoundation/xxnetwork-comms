////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package mixmessages

import (
	"bytes"
	"crypto/rand"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"testing"
)

// Ensure message type conforms to genericSignable interface
// If this ever fails, check for modifications in the crypto library
//  as well as for this message type
var _ = signature.GenericSignable(&NDF{})

// Happy path
func TestNDF_ClearSignature(t *testing.T) {
	// Create signature object
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{
		Signature: expectedSig,
		Nonce:     expectedSig,
	}

	// Create ndf message
	testNdf := &NDF{
		RsaSignature: sig,
	}

	// Clear the signature
	testNdf.ClearSig()

	// Check that the signature's values are nil after clearing
	if testNdf.GetSig() != nil && testNdf.GetNonce() != nil {
		t.Errorf("Signature's values should be nil after a ClearSignature() call!"+
			"\n\tSignature is: %+v", testNdf.RsaSignature)
	}
}

// ------------------------------------ Signature tests ------------------------------------------

// Happy path
func TestNDF_SetSignature(t *testing.T) {
	// Create ndf message
	tempVal := []byte("fail Fail fail")
	tempSig := &RSASignature{Signature: tempVal}
	testNdf := &NDF{RsaSignature: tempSig}

	// Set the sig
	expectedSig := []byte{1, 2, 45, 67, 42}
	testNdf.SetSig(expectedSig)

	// Check that the ndf's signature is identical to the one set
	if bytes.Compare(testNdf.RsaSignature.Signature, expectedSig) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, testNdf.RsaSignature.Signature)
	}
}

// Happy path
func TestNDF_SetSignature_NilObject(t *testing.T) {
	testNdf := &NDF{}

	// Set the sig w/o signature being initialized
	expectedSig := []byte{1, 2, 45, 67, 42}
	testNdf.SetSig(expectedSig)

	// Sig should be set to expected value
	if bytes.Compare(testNdf.RsaSignature.Signature, expectedSig) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, testNdf.RsaSignature.Signature)
	}

}

// Error path
func TestNDF_SetSignature_SetNil(t *testing.T) {
	// Create signature object
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{
		Signature: expectedSig,
		Nonce:     expectedSig,
	}

	// Create ndf message
	testNdf := &NDF{RsaSignature: sig}

	// Set the sig to nil (error case)
	err := testNdf.SetSig(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// Happy path
func TestNDF_GetSignature(t *testing.T) {
	// Create ndf and set signature
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{Signature: expectedSig}
	testNdf := &NDF{
		RsaSignature: sig,
	}

	// Fetch signature
	receivedSig := testNdf.GetSig()

	// Compare fetched value to expected value
	if bytes.Compare(expectedSig, receivedSig) != 0 {
		t.Errorf("Signature does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedSig, receivedSig)
	}
}

// Error path (nil signature)
func TestNDF_GetSignature_NilCase(t *testing.T) {
	// Create ndf w/o signature object
	testNdf := &NDF{}

	// Attempt to get signature
	receivedSig := testNdf.GetSig()

	// Received sig should be nil
	if receivedSig != nil {
		t.Errorf("Signature should default to nil if not set!")
	}

}

// ----------------------------------------- Nonce tests ------------------------------------------------

// Happy path
func TestNDF_GetNonce(t *testing.T) {
	expectedNonce := []byte{1, 2, 45, 67, 42}

	// Create message with nonce value
	sig := &RSASignature{Nonce: expectedNonce}
	testNdf := &NDF{
		RsaSignature: sig,
	}

	// Retrieve the nonce
	receivedNonce := testNdf.GetNonce()

	// Compare to the value originally set
	if bytes.Compare(expectedNonce, receivedNonce) != 0 {
		t.Errorf("Nonce does not match one that was set!"+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, receivedNonce)

	}
}

// Error path (nil object)
func TestNDF_GetNonce_NilObject(t *testing.T) {
	// Create ndf w/o signature object
	testNdf := &NDF{}

	// Attempt to get nonce
	receivedSig := testNdf.GetNonce()

	// Received nonce should be nil
	if receivedSig != nil {
		t.Errorf("Nonce should default to nil if not set!")
	}

}

//
func TestNDF_SetNonce(t *testing.T) {
	// Create ndf message
	tempVal := []byte("fail Fail fail")
	tempSig := &RSASignature{Nonce: tempVal}
	testNdf := &NDF{RsaSignature: tempSig}

	// Set the sig
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testNdf.SetNonce(expectedNonce)

	// Check that the ndf's signature is identical to the one set
	if bytes.Compare(testNdf.RsaSignature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testNdf.RsaSignature.Nonce)
	}
}

// Happy path
func TestNDF_SetNonce_NilObject(t *testing.T) {
	testNdf := &NDF{}

	// Set the sig w/o signature being initialized
	expectedNonce := []byte{1, 2, 45, 67, 42}
	testNdf.SetNonce(expectedNonce)

	// Sig should be set to expected value
	if bytes.Compare(testNdf.RsaSignature.Nonce, expectedNonce) != 0 {
		t.Errorf("Signature should match value it was set to! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expectedNonce, testNdf.RsaSignature.Nonce)
	}
}

// Error path
func TestNDF_SetNonce_SetNil(t *testing.T) {
	// Create signature object
	expectedSig := []byte{1, 2, 45, 67, 42}
	sig := &RSASignature{
		Signature: expectedSig,
		Nonce:     expectedSig,
	}

	// Create ndf message
	testNdf := &NDF{RsaSignature: sig}

	// Set the sig to nil (error case)
	err := testNdf.SetNonce(nil)
	if err != nil {
		return
	}

	t.Errorf("Expected error path: Should not be able to set signature as nil")

}

// -------------------- Sign/Verify tests -------------------------------

// Happy path
func TestNdf_SignVerify(t *testing.T) {
	// Create ndf object
	ourNdf := []byte{25, 254, 123, 42}
	testNdf := &NDF{
		Ndf: ourNdf,
	}
	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Sign message
	err = signature.Sign(testNdf, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Verify signature
	err = signature.Verify(testNdf, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Failed to verify: %+v", err)
	}
}

// Error path
func TestNdf_SignVerify_Error(t *testing.T) {
	// Create ndf object
	ourNdf := []byte{25, 254, 123, 42}
	testNdf := &NDF{
		Ndf: ourNdf,
	}

	// Generate keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate key: %+v", err)
	}
	pubKey := privateKey.GetPublic()

	// Sign message
	err = signature.Sign(testNdf, privateKey)
	if err != nil {
		t.Errorf("Unable to sign message: %+v", err)
	}

	// Reset ndf value so verify()'s signature won't match
	testNdf.Ndf = []byte{1}

	// Verify signature
	err = signature.Verify(testNdf, pubKey)
	// Verify signature
	if err != nil {
		return
	}

	t.Error("Expected error path: Should not have verified!")

}
