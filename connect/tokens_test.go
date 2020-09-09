///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package connect

import (
	"bytes"
	"reflect"
	"sync"
	"testing"
)

// Unit test for NewToken
func TestNewToken(t *testing.T) {
	newToken := NewToken()

	// Test that lock is empty on initialization
	if !reflect.DeepEqual(newToken.lock, sync.RWMutex{}) {
		t.Errorf("New token's lock initialized incorrectly."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", sync.RWMutex{}, newToken.lock)
	}

	// Test that token is empty on initialization
	if newToken.token != nil {
		t.Errorf("New token's toke initialized incorrectly."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", nil, newToken.token)
	}
}

// Unit test for SetToken
func TestToken_SetToken(t *testing.T) {
	newToken := NewToken()
	expectedVal := []byte("testToken")

	// Set token's value
	newToken.SetToken(expectedVal)

	// Check that the new value has been written to the token
	if !bytes.Equal(expectedVal, newToken.token) {
		t.Errorf("SetToken did not write value as expected."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedVal, newToken.token)
	}
}

// Unit test for GetToken
func TestToken_GetToken(t *testing.T) {
	newToken := NewToken()

	// Test GetToken on a newly initialized token object
	if newToken.GetToken() != nil {
		t.Errorf("GetToken did not retrieve expected value on initialization."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", nil, newToken.GetToken())
	}

	// Set a new value for token
	expectedVal := []byte("testToken")
	newToken.SetToken(expectedVal)

	// Test that the new value is successfully retrieved by GetToken
	if !bytes.Equal(expectedVal, newToken.GetToken()) {
		t.Errorf("GetToken did not retrieve expected value after a SetToken call."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedVal, newToken.GetToken())
	}

}
