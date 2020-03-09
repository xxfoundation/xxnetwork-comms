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

func TestUpdates_AddRound(t *testing.T) {
	u := NewUpdates()
	err := u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	})
	if err != nil {
		t.Errorf("Failed to add round: %+v", err)
	}
}

func TestUpdates_GetUpdate(t *testing.T) {
	u := NewUpdates()
	_ = u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 3,
	})
	_, err := u.GetUpdate(3)
	if err != nil {
		t.Errorf("Failed to get update: %+v", err)
	}
}

func TestUpdates_GetUpdates(t *testing.T) {
	u := NewUpdates()
	_ = u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 3,
	})
	_ = u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 4,
	})
	_ = u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 4,
	})
	l := u.GetUpdates(2)
	if len(l) != 2 {
		t.Error("Something went wrong, didn't get all results")
	}
}
