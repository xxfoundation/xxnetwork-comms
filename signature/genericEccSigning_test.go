////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package signature

import (
	"crypto/rand"
	"gitlab.com/xx_network/crypto/signature/ec"
	"testing"
)

// Happy path / smoke test
func TestSignEddsa(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys
	privKey, err := ec.NewKeyPair(rand.Reader)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}

	// Sign message
	err = SignEddsa(testSig, privKey)
	if err != nil {
		t.Fatalf("Failed to sign message: %+v", err)
	}

}

// Happy path
func TestSignVerifyEddsa(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys
	privKey, err := ec.NewKeyPair(rand.Reader)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	// Sign object
	err = SignEddsa(testSig, privKey)
	if err != nil {
		t.Fatalf("Failed to sign: +%v", err)
	}
	// Verify the signature
	err = VerifyEddsa(testSig, pubKey)
	if err != nil {
		t.Errorf("Expected happy path! Verification resulted in: %+v", err)
	}

}

// Error path
func TestSignVerifyEddsa_Error(t *testing.T) {
	// Generate a test signable
	testSig := InitTestSignable()

	// Generate keys
	privKey, err := ec.NewKeyPair(rand.Reader)
	if err != nil {
		t.Errorf("Failed to generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	// Sign object
	err = SignEddsa(testSig, privKey)
	if err != nil {
		t.Fatalf("Failed to sign: +%v", err)
	}

	// Modify object post-signing
	testSig.id = []byte("i will fail")

	// Attempt to verify modified object
	err = VerifyEddsa(testSig, pubKey)
	if err != nil {
		return
	}
	t.Errorf("Expected error path: VerifyEddsa should not return true")

}
