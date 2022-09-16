////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package token

import (
	"bytes"
	nonce2 "gitlab.com/xx_network/crypto/nonce"
	"reflect"
	"testing"
)

// ##############
// # LIVE TESTS #
// ##############

// Unit test for NewLive
func TestNewLive(t *testing.T) {
	newToken := NewLive()

	// Test that token is empty on initialization
	if !reflect.DeepEqual(newToken.t, Token{}) {
		t.Errorf("New token's toke initialized incorrectly."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", nil, newToken.t)
	}
	if newToken.has {
		t.Errorf("New token'starts off not clear.")
	}
}

// Unit test for Set
func TestLive_SetToken(t *testing.T) {
	newToken := NewLive()
	expectedVal := []byte("testToken")
	tkn := Token{}
	copy(tkn[:], expectedVal)

	// Set token's value
	newToken.Set(tkn)

	// Check that the new value has been written to the token
	if !bytes.Equal(tkn[:], newToken.GetBytes()) {
		t.Errorf("Set did not write value as expected."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", tkn[:], newToken.GetBytes())
	}
}

// Unit test for Get
func TestLive_GetToken(t *testing.T) {
	newToken := NewLive()

	// Test Get on a newly initialized token object
	if newToken.Has() {
		t.Errorf("Get did not retrieve expected value on initialization.")
	}

	// Set a new value for token
	expectedVal := []byte("testToken")
	tkn := Token{}
	copy(tkn[:], expectedVal)
	newToken.Set(tkn)

	// Test that the new value is successfully retrieved by Get
	if !bytes.Equal(tkn[:], newToken.GetBytes()) {
		t.Errorf("Get did not retrieve expected value after a Set call."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedVal, newToken.GetBytes())
	}
}

// Unit test for Live.Clear
func TestLive_Clear(t *testing.T) {
	l := NewLive()
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)
	l.Set(t1)

	l.Clear()

	cleared := true
	for i := 0; i < len(l.t); i++ {
		if l.t[i] != 0 {
			cleared = false
		}
	}

	if l.has || !cleared {
		t.Errorf("Did not properly clear token")
	}
}

// Unit test for Live.Get
func TestLive_Get(t *testing.T) {
	l := NewLive()
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)
	l.Set(t1)

	t2, _ := l.Get()
	if !t2.Equals(t1) {
		t.Error("Did not get same token we set")
	}

	l.Clear()

	if l.t.Equals(t2) {
		t.Error("token from live was not deep copied")
	}
}

// Unit test for Live.GetBytes
func TestLive_GetBytes(t *testing.T) {
	l := NewLive()
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)
	l.Set(t1)

	b := l.GetBytes()
	if bytes.Compare(b, t1.Marshal()) != 0 {
		t.Error("Did not receive same token")
	}

	l.Clear()

	if bytes.Compare(b, l.GetBytes()) == 0 {
		t.Error("Clearing token also cleared received data")
	}
}

// Unit test for Live.Has
func TestLive_Has(t *testing.T) {
	l := NewLive()
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)
	if l.Has() {
		t.Error("Has was initialized as true")
	}

	l.Set(t1)
	if !l.Has() {
		t.Error("Has was not set properly")
	}
}

// Unit test for Live.Set
func TestLive_Set(t *testing.T) {
	l := NewLive()
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)

	l.Set(t1)

	if !l.t.Equals(t1) || !l.has {
		t.Error("Didn't properly set token")
	}
}

// ###############
// # TOKEN TESTS #
// ###############

// Unit test for GenerateToken
func TestGenerateToken(t *testing.T) {
	nonce, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce: %+v", err)
	}
	_ = GenerateToken(nonce)
}

// Unit test for Unmarshal
func TestUnmarshal(t *testing.T) {
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)
	m := t1.Marshal()

	t2, err := Unmarshal(m)
	if err != nil {
		t.Errorf("Error unmarshalling token data: %+v", err)
	}
	if !t2.Equals(t1) {
		t.Errorf("Unmarshalled and original tokens are different\nOriginal: %+v, Unmarshalled: %+v", t1, t2)
	}
}

// Error path test for Unmarshal
func TestUnmarshal_Error(t *testing.T) {
	badData := []byte("test")
	_, err := Unmarshal(badData)
	if err == nil {
		t.Error("Should have received an error for bad data")
	}
}

// Unit test for Token.Marshal
func TestToken_Marshal(t *testing.T) {
	nonce, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	newToken := GenerateToken(nonce)
	m := newToken.Marshal()

	if len(m) != len(newToken) {
		t.Error("Lengths are different, Did not properly marshal token")
	}

	m[0] = m[0] + 1
	if m[0] == newToken[0] {
		t.Error("Token was not properly copied to new location")
	}
}

// Unit test for Token.Equals
func TestToken_Equals(t *testing.T) {
	n1, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce")
	}
	t1 := GenerateToken(n1)
	t1_copy := GenerateToken(n1)
	if !t1.Equals(t1_copy) {
		t.Errorf("Tokens made with same data were not equal")
	}

	n2, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce 2")
	}
	t2 := GenerateToken(n2)
	if t1.Equals(t2) {
		t.Errorf("Tokens made with different data were equal")
	}
}
