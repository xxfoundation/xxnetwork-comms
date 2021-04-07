////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////

package signature

import (
	"crypto"
	"github.com/katzenpost/core/crypto/eddsa"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/csprng"
	"hash"
)

// Interface for our GenericRsaSignable structure, to be used for protobuffs we want
// to be able to cryptographically sign using an RSA key
type GenericEccSignable interface {
	// GetSig returns the RSA signature.
	// IF none exists, it creates it, adds it to the object, then returns it.
	GetEccSig() *messages.ECCSignature

	// Digest hashes the contents of the message in a repeatable manner
	// using the provided cryptographic hash. It includes the nonce in the hash
	Digest(nonce []byte, h hash.Hash) []byte
}

// SignEddsa takes a GenericEccSignable object, marshals the data
// intended to be signed with a nonce.
func SignEddsa(signable GenericEccSignable, privKey *eddsa.PrivateKey) error {
	// Create rand for signing and nonce generation
	rand := csprng.NewSystemRNG()

	// Generate nonce
	newNonce := make([]byte, 32)
	_, err := rand.Read(newNonce)
	if err != nil {
		return errors.Errorf("Failed to generate nonce: %+v", err)
	}

	// Prepare to hash the data
	// fixme: change hash be faster for this interface?
	sha := crypto.SHA256
	h := sha.New()

	// Generate the serialized data
	data := signable.Digest(newNonce, h)

	// Sign the message
	signature := privKey.Sign(data)

	// Print results of signing
	jww.TRACE.Printf("signature.Sign nonce: 0x%x", newNonce)
	jww.TRACE.Printf("signature.Sign sig for nonce 0x%x 0x%x", newNonce[:8], signature)
	jww.TRACE.Printf("signature.Sign digest for nonce 0x%x 0x%x", newNonce[:8], data)
	jww.TRACE.Printf("signature.Sign data for nonce 0x%x: [%x]", newNonce[:8], data)
	jww.TRACE.Printf("signature.Sign privKey for nonce 0x%x: Type: %s;; Bytes: %x;; Identity: %s", newNonce[:8], privKey.KeyType(), privKey.Bytes(), privKey.Identity())
	jww.TRACE.Printf("signature.Sign pubKey for nonce 0x%x: pubKey: %s;", newNonce[:8], privKey.PublicKey().String())

	// Modify the signature for the new values
	// NOTE: This is the only way to change the internal of the interface object.
	// The code commented below would be cleaner, but the changes do not take
	signable.GetEccSig().Signature = signature
	signable.GetEccSig().Nonce = newNonce

	//ourSig := signable.GetEccSig()
	//ourSig = &messages.ECCSignature{
	//	Nonce:     newNonce,
	//	Signature: signature,
	//}
	return nil
}

// VerifyEddsa takes the signature from the verifiable message
// and verifies it on the public key. If
func VerifyEddsa(verifiable GenericEccSignable, pubKey *eddsa.PublicKey) error {
	sigMsg := verifiable.GetEccSig()
	nonce := sigMsg.Nonce

	sig := sigMsg.Signature
	// Prepare to hash the data
	// fixme: change hash be faster for this interface?
	sha := crypto.SHA256
	h := sha.New()

	// Generate the serialized data
	data := verifiable.Digest(nonce, h)

	if !pubKey.Verify(sig, data) {
		return errors.New("failed to verify EDDSA signature")
	}

	return nil
}
