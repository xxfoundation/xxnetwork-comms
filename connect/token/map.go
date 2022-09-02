////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package token

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/crypto/nonce"
	"sync"
)

type Map struct {
	m   map[Token]nonce.Nonce
	mux sync.Mutex
}

func NewMap() *Map {
	return &Map{
		m:   make(map[Token]nonce.Nonce),
		mux: sync.Mutex{},
	}
}

func (m *Map) Generate() Token {
	newNonce, err := nonce.NewNonce(nonce.RegistrationTTL)
	if err != nil {
		jww.FATAL.Panicf("Failed to generate new Token/Nonce pair: %s", err)
	}

	newToken := GenerateToken(newNonce)

	m.mux.Lock()

	m.m[newToken] = newNonce

	m.mux.Unlock()

	jww.DEBUG.Printf("Token generated: %v", newToken)
	return newToken

}

func (m *Map) Validate(token Token) bool {
	m.mux.Lock()
	retrievedNonce, ok := m.m[token]
	delete(m.m, token)
	m.mux.Unlock()

	if !ok {
		return false
	}

	return retrievedNonce.IsValid()

}
