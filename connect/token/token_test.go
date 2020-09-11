///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package token

import (
	"bytes"
	"reflect"
	"testing"
)

// Unit test for NewLive
func TestNewToken(t *testing.T) {
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
func TestToken_SetToken(t *testing.T) {
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
func TestToken_GetToken(t *testing.T) {
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
