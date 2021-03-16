///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testutils"
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
	rnd := NewRound(ri, pubKey)
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
	testutils.SignRoundInfo(ri, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	rnd := NewRound(ri, pubKey)

	_ = d.UpsertRound(rnd)
	_, err = d.GetRound(0)
	if err != nil {
		t.Errorf("Failed to get roundinfo with proper id")
	}
}

func TestData_ComparisonFunc(t *testing.T) {
	d := NewData()

	// Construct a mock round object
	roundInfoOne := &mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 3,
	}
	testutils.SignRoundInfo(roundInfoOne, t)

	pubKey, err := testutils.LoadPublicKeyTesting(t)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		t.FailNow()
	}
	roundOne := NewRound(roundInfoOne, pubKey)
	_ = d.UpsertRound(roundOne)

	// Construct a mock round object
	roundInfoTwo := &mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 4,
	}
	testutils.SignRoundInfo(roundInfoTwo, t)

	roundTwo := NewRound(roundInfoTwo, pubKey)
	_ = d.UpsertRound(roundTwo)
	r, err := d.GetRound(2)
	if err != nil {
		t.Errorf("Failed to get round: %+v", err)
	}
	if r.UpdateID != 4 {
		t.Error("Round did not properly upsert")
	}
}
