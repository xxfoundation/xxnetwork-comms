////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"bytes"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testutils"
	"testing"
)

// Smoke test for constructor
func TestNewRound(t *testing.T) {
	pubKey, _ := testutils.LoadPublicKeyTesting(t)
	ecKey, _ := testutils.LoadEllipticPublicKey(t)

	ri := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)}

	rnd := NewRound(ri, pubKey, ecKey.GetPublic())

	// Check that values in object match inputted values
	if rnd.info != ri || rnd.rsaPubKey != pubKey || !bytes.Equal(rnd.ecPubKey.Marshal(), ecKey.GetPublic().Marshal()) {
		t.Errorf("Initial round values from constructor are not expected."+
			"\n\tExpected round info: %v"+
			"\n\tReceived round info: %v"+
			"\n\tExpected rsa public key: %v"+
			"\n\tReceived rsa public key: %v"+
			"\n\tExpected EC public key: %v"+
			"\n\tReceived EC public key: %v",
			ri, rnd.info, pubKey, rnd.rsaPubKey, ecKey.GetPublic().Marshal(), rnd.ecPubKey.Marshal())
	}

}

// Smoke test for other constructor
func TestNewVerifiedRound(t *testing.T) {
	pubKey, _ := testutils.LoadPublicKeyTesting(t)
	ri := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)}

	rnd := NewVerifiedRound(ri, pubKey)

	// Check that values in object match inputted values
	if rnd.info != ri || rnd.rsaPubKey != pubKey || *rnd.needsValidation != 1 {
		t.Errorf("Initial round values from constructor are not expected."+
			"\n\tExpected round info: %v"+
			"\n\tReceived round info: %v"+
			"\n\tExpected public key: %v"+
			"\n\tReceived public key: %v"+
			"\n\tExpected needsValidation: %v"+
			"\n\tReceived needsValidation: %v",
			ri, rnd.info, pubKey, rnd.rsaPubKey, rnd.needsValidation, 1)
	}

}

// Unit test of Get()
func TestNewRound_Get(t *testing.T) {
	pubKey, _ := testutils.LoadPublicKeyTesting(t)
	ri := &mixmessages.RoundInfo{ID: uint64(1), UpdateID: uint64(1)}
	// Mock signature of roundInfo as it will be verified in codepath
	testutils.SignRoundInfoRsa(ri, t)
	ecPubKey, _ := testutils.LoadEllipticPublicKey(t)
	testutils.SignRoundInfoEddsa(ri, ecPubKey, t)

	rnd := NewRound(ri, pubKey, ecPubKey.GetPublic())

	// Check the initial value of the atomic value (lazily)
	if *rnd.needsValidation != 0 {
		t.Errorf("Validation value is not default value!")
	}

	// Check that the returned round info matches inputted value
	retrievedRI := rnd.Get()
	if retrievedRI != ri {
		t.Errorf("RoundInfo from Get() not expected."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", ri, retrievedRI)
	}

	// Check the atomic value has been incremented after a Get() call
	if *rnd.needsValidation != 1 {
		t.Errorf("Validation value is not set after Get() call!")
	}

	// Check the atomic value has not changed after a second Get() call
	_ = rnd.Get()
	if *rnd.needsValidation != 1 {
		t.Errorf("Validation value has been modified after a second Get() call!")
	}

}
