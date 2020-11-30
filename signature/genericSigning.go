////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

// Contains a generic signing interface and implementations to sign the data
// as well as verify the signature

package signature

import (
	"crypto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"hash"
)

func init() {
	jww.ERROR.Println("Signature verification is curently disabled, if you " +
		"get this message STOP OPERATION and contact the developers at " +
		"nodes@xx.network")
}

// Interface for signing generically
type GenericSignable interface {
	// GetSignature returns the RSA signature.
	// IF none exists, it creates it, adds it to the object, then returns it.
	GetSignature() *messages.RSASignature
	// Digest hashes the contents of the message in a repeatable manner
	// using the provided cryptographic hash. It includes the nonce in the hash
	Digest(nonce []byte, h hash.Hash) []byte
	// SetSignature modifies the internal of the generic sign-able
	// in order to save the newly created signature.
	// It will generate an RSA signature message
	SetSignature(signature, nonce []byte) error
}

// Sign takes a genericSignable object, marshals the data intended to be signed.
// It hashes that data and sets it as the signature of that object
func Sign(signable GenericSignable, privKey *rsa.PrivateKey) error {

	// Create rand for signing and nonce generation
	rand := csprng.NewSystemRNG()

	// Generate nonce
	ourNonce := make([]byte, 32)
	_, err := rand.Read(ourNonce)
	if err != nil {
		return errors.Errorf("Failed to generate nonce: %+v", err)
	}

	// Prepare to hash the data
	sha := crypto.SHA256
	h := sha.New()

	// Generate the serialized data
	data := signable.Digest(ourNonce, h)

	// Get the data that is to be signed (including nonce)

	// Sign the message
	signature, err := rsa.Sign(rand, privKey, sha, data, nil)

	// Print results of signing
	jww.TRACE.Printf("signature.Sign nonce: 0x%x", ourNonce)
	jww.TRACE.Printf("signature.Sign sig for nonce 0x%x 0x%x", ourNonce[:8], signature)
	jww.TRACE.Printf("signature.Sign digest for nonce 0x%x 0x%x", ourNonce[:8], data)
	jww.TRACE.Printf("signature.Sign data for nonce 0x%x: [%x]", ourNonce[:8], data)
	jww.TRACE.Printf("signature.Sign privKey for nonce 0x%x: N: 0x%v;; E: 0x%x;; D: 0x%v", ourNonce[:8], privKey.N.Text(16), privKey.E, privKey.D.Text(16))
	jww.TRACE.Printf("signature.Sign pubKey for nonce 0x%x: E: 0x%x;; V: 0x%v", ourNonce[:8], privKey.PublicKey.E, privKey.PublicKey.N.Text(16))

	if err != nil {
		return errors.Errorf("Unable to sign message: %+v", err)
	}

	// And set the signature
	err = signable.SetSignature(signature, ourNonce)
	if err != nil {
		return errors.Errorf("Unable to finalize signature: %+v", err)
	}

	return nil
}

// Verify takes the signature from the object and clears it out.
// It then re-creates the signature and compares it to the original signature.
// If the recreation matches the original signature it returns true,
// else it returns false
func Verify(verifiable GenericSignable, pubKey *rsa.PublicKey) error {

	// Take the signature from the object
	sigMsg := verifiable.GetSignature()
	nonce := sigMsg.Nonce
	sig := sigMsg.Signature

	// Prepare to hash the data
	sha := crypto.SHA256
	h := sha.New()

	// Get the data to replicate the signature
	data := verifiable.Digest(nonce, h)

	// Verify the signature using our implementation
	err := rsa.Verify(pubKey, sha, data, sig, nil)

	jww.TRACE.Printf("signature.Verify nonce: 0x%x", nonce)
	jww.TRACE.Printf("signature.Verify sig for nonce 0x%x: 0x%x", nonce[:8], sig)
	jww.TRACE.Printf("signature.Verify digest for nonce 0x%x, 0x%x", nonce[:8], data)
	jww.TRACE.Printf("signature.Verify data for nonce 0x%x: [%x]", nonce[:8], data)
	jww.TRACE.Printf("signature.Verify pubKey for nonce 0x%x: E: 0x%x;; V: 0x%v", nonce[:8], pubKey.E, pubKey.N.Text(16))

	// And check for an error
	if err != nil {
		// If there is an error, then signature is invalid
		return err
	}

	// Otherwise it has been verified
	return nil

}
