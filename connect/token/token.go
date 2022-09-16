////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package token

import (
	"bytes"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/nonce"
	"sync"
)

type Token [nonce.NonceLen]byte

// Generates a new token and adds it to internal state
func GenerateToken(newNonce nonce.Nonce) Token {
	return Token(newNonce.Value)
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

func (t Token) Equals(u Token) bool {
	return bytes.Equal(t[:], u[:])
}

// Represents a reverse-authentication token
type Live struct {
	mux sync.RWMutex
	t   Token
	has bool
}

// Constructor which initializes a token for
// use by the associated host object
func NewLive() *Live {
	return &Live{
		has: false,
	}
}

// Get reads and returns the token
func (l *Live) Get() (Token, bool) {
	l.mux.RLock()
	defer l.mux.RUnlock()
	var tCopy Token
	copy(tCopy[:], l.t[:])
	return tCopy, l.has
}

// Get reads and returns the token
func (l *Live) GetBytes() []byte {
	t, ok := l.Get()
	if !ok {
		return nil
	} else {
		return t[:]
	}
}

//Returns true if a token is present
func (l *Live) Has() bool {
	l.mux.RLock()
	defer l.mux.RUnlock()
	return l.has
}

// Set rewrites the token for negotiation or renegotiation
func (l *Live) Set(newToken Token) {
	l.mux.Lock()
	copy(l.t[:], newToken[:])
	l.has = true
	l.mux.Unlock()
}

// Clear is used to set token to a nil value
// as store will not let you do this explicitly
func (l *Live) Clear() {
	l.mux.Lock()
	for i := 0; i < len(l.t); i++ {
		l.t[i] = 0
	}
	l.has = false
	l.mux.Unlock()
}
