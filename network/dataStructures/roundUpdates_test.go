package dataStructures

import (
	"gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

func TestUpdates_AddRound(t *testing.T) {
	u := Updates{}
	err := u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 0,
	})
	if err != nil {
		t.Errorf("Failed to add round: %+v", err)
	}
}

func TestUpdates_GetUpdate(t *testing.T) {
	u := Updates{}
	_ = u.AddRound(&mixmessages.RoundInfo{
		ID:       0,
		UpdateID: 3,
	})
	_, err := u.GetUpdate(3)
	if err != nil {
		t.Errorf("Failed to get update: %+v", err)
	}
}
