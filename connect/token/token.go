///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package token

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/crypto/nonce"
	"sync/atomic"
)

type Token [nonce.NonceLen]byte

// Generates a new token and adds it to internal state
func GenerateToken(newNonce nonce.Nonce) Token {
	var t Token
	copy(t[:], newNonce.Bytes())
	return t
}

func Unmarshal(newVal []byte) (Token, error) {
	if len(newVal) != nonce.NonceLen {
		return Token{}, errors.Errorf("New value is not of expected length. "+
			"Expected length of %d, received length: %d", nonce.NonceLen, len(newVal))
	}
	var t Token
	copy(t[:], newVal)
	return t, nil

}

func (t Token) Marshal() []byte {
	return t[:]
}

// Represents a reverse-authentication token
type Live struct {
	*atomic.Value
}

// Constructor which initializes a token for
// use by the associated host object
func NewLive() Live {
	return Live{
		Value: &atomic.Value{},
	}
}

// Set rewrites the token for negotiation or renegotiation
func (t *Live) Set(newToken []byte) {

	t.Store(newToken)
}

// Get reads and returns the token
func (t *Live) Get() []byte {
	retrievedVal := t.Load()
	if retrievedVal == nil {
		return nil
	}
	b := retrievedVal.([]byte)
	if len(b) == 0 {
		return nil
	}

	return b
}

// Clear is used to set token to a nil value
// as store will not let you do this explicitly
func (t *Live) Clear() {
	t.Store([]byte{})
}
