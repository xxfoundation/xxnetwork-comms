////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package token

import (
	nonce2 "gitlab.com/xx_network/crypto/nonce"
	"testing"
)

// Unit test for NewMap
func TestNewMap(t *testing.T) {
	m := NewMap()
	if m.m == nil {
		t.Errorf("Failed to initialize map")
	}
}

// Unit test for Map.Generate()
func TestMap_Generate(t *testing.T) {
	m := NewMap()
	_ = m.Generate()
}

// Unit test for Map.Validate()
func TestMap_Validate(t *testing.T) {
	m := NewMap()

	token := m.Generate()

	valid := m.Validate(token)
	if !valid {
		t.Error("Failed to validate token")
	}
}

// Test for invalid token
func TestMap_ValidateError(t *testing.T) {
	m := NewMap()

	nonce, err := nonce2.NewNonce(nonce2.NonceLen)
	if err != nil {
		t.Errorf("Failed to generate nonce for token")
	}
	token := GenerateToken(nonce)

	valid := m.Validate(token)
	if valid {
		t.Error("Token should not have validated")
	}
}
