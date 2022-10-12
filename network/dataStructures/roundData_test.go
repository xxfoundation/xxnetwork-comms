////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testutils"
	"reflect"
	"testing"
)

func TestData_UpsertRound(t *testing.T) {
	d := NewData()

	// Construct a mock round object
	ri := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	}

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := NewRound(ri, pubKey, nil)
	err = d.UpsertRound(rnd)
	if err != nil {
		t.Errorf("Failed to upsert round: %+v", err)
	}
}

func TestData_GetRound(t *testing.T) {
	d := NewData()

	// Construct a mock round object
	ri := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	}
	testutils.SignRoundInfoRsa(ri, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := NewRound(ri, pubKey, nil)

	_ = d.UpsertRound(rnd)
	_, err = d.GetRound(0)
	if err != nil {
		t.Errorf("Failed to get roundinfo with proper id")
	}
}

func TestData_GetWrappedRound(t *testing.T) {
	d := NewData()

	// Construct a mock round object
	ri := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	}
	testutils.SignRoundInfoRsa(ri, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := NewRound(ri, pubKey, nil)

	_ = d.UpsertRound(rnd)
	retrieved, err := d.GetWrappedRound(0)
	if err != nil {
		t.Errorf("Failed to get roundinfo with proper id")
	}

	if !reflect.DeepEqual(rnd, retrieved) {
		t.Errorf("Retrieved value did not match expected!"+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", rnd, retrieved)
	}
}

func TestData_ComparisonFunc(t *testing.T) {
	d := NewData()

	// Construct a mock round object
	roundInfoOne := &mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 3,
	}
	testutils.SignRoundInfoRsa(roundInfoOne, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	roundOne := NewRound(roundInfoOne, pubKey, nil)
	_ = d.UpsertRound(roundOne)

	// Construct a mock round object
	roundInfoTwo := &mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 4,
	}
	testutils.SignRoundInfoRsa(roundInfoTwo, t)

	roundTwo := NewRound(roundInfoTwo, pubKey, nil)
	_ = d.UpsertRound(roundTwo)
	r, err := d.GetRound(2)
	if err != nil {
		t.Errorf("Failed to get round: %+v", err)
	}
	if r.UpdateID != 4 {
		t.Error("Round did not properly upsert")
	}
}
