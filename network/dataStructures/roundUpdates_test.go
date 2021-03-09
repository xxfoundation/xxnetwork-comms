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

func TestUpdates_AddRound(t *testing.T) {
	u := NewUpdates()
	// Construct a mock round object
	ri := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	}
	rnd := NewRound(ri, testutils.LoadKeyTesting(t))
	err := u.AddRound(rnd)
	if err != nil {
		t.Errorf("Failed to add round: %+v", err)
	}
}

func TestUpdates_GetUpdate(t *testing.T) {
	u := NewUpdates()
	updateID := 3
	// Construct a mock round object
	ri := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: uint64(updateID),
	}
	if err := testutils.SignRoundInfo(ri, t); err != nil {
		t.Errorf("Failed to sign mock round info: %v", err)
	}
	rnd := NewRound(ri, testutils.LoadKeyTesting(t))
	_ = u.AddRound(rnd)
	_, err := u.GetUpdate(updateID)
	if err != nil {
		t.Errorf("Failed to get update: %+v", err)
	}
}

func TestUpdates_GetUpdates(t *testing.T) {
	u := NewUpdates()
	updateID := 3
	// Construct a mock round object
	roundInfoOne := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: uint64(updateID),
	}
	if err := testutils.SignRoundInfo(roundInfoOne, t); err != nil {
		t.Errorf("Failed to sign mock round info: %v", err)
	}
	roundOne := NewRound(roundInfoOne, testutils.LoadKeyTesting(t))

	// Construct a second eound
	roundInfoTwo := &mixmessages.RoundInfo{
		ID:       0,
		UpdateID: uint64(updateID + 1),
	}
	if err := testutils.SignRoundInfo(roundInfoTwo, t); err != nil {
		t.Errorf("Failed to sign mock round info: %v", err)
	}
	roundTwo := NewRound(roundInfoTwo, testutils.LoadKeyTesting(t))

	_ = u.AddRound(roundOne)
	// Add second round twice (shouldn't duplicate)
	_ = u.AddRound(roundTwo)
	_ = u.AddRound(roundTwo)
	l := u.GetUpdates(2)
	if len(l) != 2 {
		t.Error("Something went wrong, didn't get all results")
	}
}
