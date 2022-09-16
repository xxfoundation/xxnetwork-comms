////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

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

// Interface for our GenericRsaSignable structure, to be used for protobuffs we want
// to be able to cryptographically sign using an RSA key
type GenericRsaSignable interface {
	// GetSig returns the RSA signature.
	// IF none exists, it creates it, adds it to the object, then returns it.
	GetSig() *messages.RSASignature
	// Digest hashes the contents of the message in a repeatable manner
	// using the provided cryptographic hash. It includes the nonce in the hash
	Digest(nonce []byte, h hash.Hash) []byte
}

// SignRsa takes a GenericRsaSignable object, marshals the data intended to be signed.
// It hashes that data and sets it as the signature of that object
func SignRsa(signable GenericRsaSignable, privKey *rsa.PrivateKey) error {
	// Create rand for signing and nonce generation
	rand := csprng.NewSystemRNG()

	// Generate nonce
	newNonce := make([]byte, 32)
	_, err := rand.Read(newNonce)
	if err != nil {
		return errors.Errorf("Failed to generate nonce: %+v", err)
	}

	// Prepare to hash the data
	sha := crypto.SHA256
	h := sha.New()

	// Generate the serialized data
	data := signable.Digest(newNonce, h)

	// SignRsa the message
	signature, err := rsa.Sign(rand, privKey, sha, data, nil)

	// Print results of signing
	jww.TRACE.Printf("RSA signature.Sign nonce: 0x%x", newNonce)
	jww.TRACE.Printf("RSA signature.Sign sig for nonce 0x%x 0x%x", newNonce[:8], signature)
	jww.TRACE.Printf("RSA signature.Sign digest for nonce 0x%x 0x%x", newNonce[:8], data)
	jww.TRACE.Printf("RSA signature.Sign data for nonce 0x%x: [%x]", newNonce[:8], data)
	jww.TRACE.Printf("RSA signature.Sign privKey for nonce 0x%x: N: 0x%v;; E: 0x%x;; D: 0x%v", newNonce[:8], privKey.N.Text(16), privKey.E, privKey.D.Text(16))
	jww.TRACE.Printf("RSA signature.Sign pubKey for nonce 0x%x: E: 0x%x;; V: 0x%v", newNonce[:8], privKey.PublicKey.E, privKey.PublicKey.N.Text(16))

	if err != nil {
		return errors.Errorf("Unable to sign message: %+v", err)
	}

	// Modify the signature for the new values
	// NOTE: This is the only way to change the internal of the interface object.
	// The code commented below would be cleaner, but the changes do not take
	signable.GetSig().Signature = signature
	signable.GetSig().Nonce = newNonce

	//ourSig := signable.GetSig()
	//ourSig = &messages.RSASignature{
	//	Nonce:     newNonce,
	//	Signature: signature,
	//}

	return nil
}

// VerifyRsa takes the signature from the verifiable message
// and verifies it on the public key. If
func VerifyRsa(verifiable GenericRsaSignable, pubKey *rsa.PublicKey) error {
	// Take the signature from the object
	sigMsg := verifiable.GetSig()
	nonce := sigMsg.Nonce
	sig := sigMsg.Signature

	// Prepare to hash the data
	sha := crypto.SHA256
	h := sha.New()

	// Get the data to replicate the signature
	data := verifiable.Digest(nonce, h)

	// Verify the signature using our implementation
	err := rsa.Verify(pubKey, sha, data, sig, nil)

	jww.TRACE.Printf("RSA signature.Verify nonce: 0x%x", nonce)
	jww.TRACE.Printf("RSA signature.Verify sig for nonce 0x%x: 0x%x", nonce[:8], sig)
	jww.TRACE.Printf("RSA signature.Verify digest for nonce 0x%x, 0x%x", nonce[:8], data)
	jww.TRACE.Printf("RSA signature.Verify data for nonce 0x%x: [%x]", nonce[:8], data)
	jww.TRACE.Printf("RSA signature.Verify pubKey for nonce 0x%x: E: 0x%x;; V: 0x%v", nonce[:8], pubKey.E, pubKey.N.Text(16))

	return err
}
