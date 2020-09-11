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
	"sync/atomic"
	"testing"
)

// Unit test for NewLive
func TestNewToken(t *testing.T) {
	newToken := NewLive()

	// Test that token is empty on initialization
	if !reflect.DeepEqual(newToken.Value, &atomic.Value{}) {
		t.Errorf("New token's toke initialized incorrectly."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", nil, newToken.Value)
	}
}

// Unit test for Set
func TestToken_SetToken(t *testing.T) {
	newToken := NewLive()
	expectedVal := []byte("testToken")

	// Set token's value
	newToken.Set(expectedVal)

	// Check that the new value has been written to the token
	if !bytes.Equal(expectedVal, newToken.Get()) {
		t.Errorf("Set did not write value as expected."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedVal, newToken.Get())
	}
}

// Unit test for Get
func TestToken_GetToken(t *testing.T) {
	newToken := NewLive()

	// Test Get on a newly initialized token object
	if newToken.Get() != nil {
		t.Errorf("Get did not retrieve expected value on initialization."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", nil, newToken.Get())
	}

	// Set a new value for token
	expectedVal := []byte("testToken")
	newToken.Set(expectedVal)

	// Test that the new value is successfully retrieved by Get
	if !bytes.Equal(expectedVal, newToken.Get()) {
		t.Errorf("Get did not retrieve expected value after a Set call."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", expectedVal, newToken.Get())
	}

}
