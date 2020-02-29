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
