////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package signature

import (
	"crypto/rand"
	"errors"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"hash"
	"testing"
)

func InitTestSignable() *TestSignable {
	// Arbitrary test values
	testId := []byte{1, 2, 3}
	// construct a TestSignable with arbitrary values
	return &TestSignable{
		id: testId,
	}

}

// Happy path / smoke test
func TestSign(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	// Sign message
	err = Sign(testSig, privKey)
	if err != nil {
		t.Errorf("Failed to sign message: %+v", err)
	}

	sigMsg := testSig.GetSignature()

	// Check if the signature is valid
	if !rsa.IsValidSignature(pubKey, sigMsg.Signature) {
		t.Errorf("Failed smoke test! Signature is not at least as long as the signer's public key."+
			"\n\tSignature: %+v"+
			"\n\tSigner's public key: %+v", len(sigMsg.Signature), pubKey.Size())
	}
}

// Error path
func TestSign_Error(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys for signing
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	// Sign object and fetch signature
	err = Sign(testSig, privKey)
	if err != nil {
		t.Errorf("Failed to sign: %+v", err)
	}
	ourSign := testSig.GetSignature()

	// Input a random set of bytes less than the signature
	randByte := make([]byte, len(ourSign.Signature)/2)
	rand.Read(randByte)

	// Compare signature to random set of bytes (expected to not match)
	// Test arbitrary slice with server's public key
	if rsa.IsValidSignature(pubKey, randByte) {
		t.Errorf("Invalid signature returned valid! "+
			"\n\t Signature: %+v "+
			"\n\t Signer's public key: %+v", len(randByte), pubKey.Size())
	}
}

// Happy path
func TestSignVerify(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	// Sign object
	err = Sign(testSig, privKey)
	if err != nil {
		t.Errorf("Failed to sign: +%v", err)
	}
	// Verify the signature
	err = Verify(testSig, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Verification resulted in: %+v", err)
	}

}

// Error path
func TestSignVerify_Error(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	// Sign object
	Sign(testSig, privKey)

	// Modify object post-signing
	testSig.id = []byte("i will fail")
	// Attempt to verify modified object
	err = Verify(testSig, pubKey)
	if err != nil {
		return
	}
	t.Errorf("Expected error path: Verify should not return true")

}

// --------- Create mock Signable object ------------------

// Test struct with arbitrary fields to be signed and verified
type TestSignable struct {
	id        []byte
	signature *messages.RSASignature
}

func (ts *TestSignable) Digest(nonce []byte, h hash.Hash) []byte {
	h.Write(nonce)
	h.Write(ts.id)
	return h.Sum(nil)

}

func (ts *TestSignable) GetSignature() *messages.RSASignature {
	return ts.signature
}

func (ts *TestSignable) SetSignature(signature, nonce []byte) error {
	if signature == nil {
		return errors.New("Cannot set signature to nil value")
	}
	ts.signature = &messages.RSASignature{
		Nonce:     nonce,
		Signature: signature,
	}
	return nil
}
