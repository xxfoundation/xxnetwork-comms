///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package connect

import "sync"

// Represents a reverse-authentication token
type Token struct {
	token []byte
	lock  sync.RWMutex
}

// Constructor which initializes a token for
// use by the associated host object
func NewToken() Token {
	return Token{
		token: nil,
		lock:  sync.RWMutex{},
	}
}

// SetToken rewrites the token for negotiation or renegotiation
func (t *Token) SetToken(newToken []byte) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.token = newToken
}

// GetToken reads and returns the token
func (t *Token) GetToken() []byte {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.token
}
