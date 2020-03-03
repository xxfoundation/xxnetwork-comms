////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package dataStructures

import (
	"gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

func TestData_UpsertRound(t *testing.T) {
	d := Data{}
	err := d.UpsertRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	})
	if err != nil {
		t.Errorf("Failed to upsert round: %+v", err)
	}
}

func TestData_GetRound(t *testing.T) {
	d := Data{}
	_ = d.UpsertRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 3,
	})
	_, err := d.GetRound(0)
	if err != nil {
		t.Errorf("Failed to get roundinfo with proper id")
	}
}

func TestData_ComparisonFunc(t *testing.T) {
	d := Data{}
	_ = d.UpsertRound(&mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 3,
	})
	_ = d.UpsertRound(&mixmessages.RoundInfo{
		ID:       2,
		UpdateID: 4,
	})
	r, err := d.GetRound(2)
	if err != nil {
		t.Errorf("Failed to get round: %+v", err)
	}
	if r.UpdateID != 4 {
		t.Error("Round did not properly upsert")
	}
}
